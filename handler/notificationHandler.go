package handler

import (
	"banknote-notification-service/errormsg"
	"banknote-notification-service/models"
	"banknote-notification-service/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NotificationHandler struct {
	Logger     *log.Logger
	Collection *mongo.Collection
}

// New return a new notification handler
func New(logger *log.Logger, collection *mongo.Collection) *NotificationHandler {
	return &NotificationHandler{
		Logger:     logger,
		Collection: collection,
	}
}

// SendNotification sends a notification based on the notification type (INSTANT or SCHEDULED).
// If the notification is SCHEDULED, its status will be "PENDING" until it is sent then SendNotification
// updates the status in the database to "SENT"
func (n *NotificationHandler) SendNotification() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get notification from client
		var notification models.Notification
		if err := c.BindJSON(&notification); err != nil {
			n.Logger.Println("error:", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"possible_errors": []string{
					errormsg.VALIDATION_ERROR.Error(),
					errormsg.DUPLICATE_ENTRY_ERROR.Error(),
				},
			})
			return
		}

		switch notification.Type {
		case "INSTANT":
			// send mail immediately
			if err := utils.SendMail(notification.From, notification.To, notification.Message); err != nil {
				n.Logger.Println("error sending notification mail:", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"possible_errors": []string{
						errormsg.INVALID_MESSAGE_ERROR.Error(),
						errormsg.MESSAGE_SERVICE_ERROR.Error(),
					},
				})
				return
			}

			// update reference, send time, status
			notification.Reference = uuid.NewV4().String()
			sendTime := time.Now()
			notification.SentAt = utils.ParseTimeToString(sendTime)
			notification.Status = "SENT"

			// insert notification into the database
			_, err := n.Collection.InsertOne(context.TODO(), notification)
			if err != nil {
				n.Logger.Println("inserting document into database:", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"possible_error": []string{
						errormsg.GENERIC_ERROR.Error(),
					},
				})
			}

			// send the notification reference as a response
			c.JSON(http.StatusOK, gin.H{
				"reference": notification.Reference,
			})

		case "SCHEDULED":
			// get the scheduled time for sending the notification and parse to time variable
			scheduledTime, err := utils.ParseTimeStringToTime(notification.SendAt)
			if err != nil {
				n.Logger.Println("could not parse time:", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"possible_errors": []string{errormsg.VALIDATION_ERROR.Error()},
				})
				return
			}

			// get time to scheduled time in seconds
			sendTime := utils.PeriodToScheduledTime(scheduledTime)

			// check if time to scheduled time is a future time (i.e at least 10 seconds later)
			if sendTime < 10 {
				n.Logger.Println("invalid schedule time")
				c.JSON(http.StatusInternalServerError, gin.H{
					"possible_errors": []string{errormsg.INVALID_SCHEDULE_DATE_ERROR.Error()},
				})
				return
			}

			// set notification reference and status to pending
			notification.Reference = uuid.NewV4().String()
			notification.Status = "PENDING"

			// insert notification to database
			_, err = n.Collection.InsertOne(context.TODO(), notification)
			if err != nil {
				n.Logger.Println("inserting document into database:", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"possible_error": []string{
						errormsg.GENERIC_ERROR.Error(),
					},
				})
			}

			// set a ticker to check if scheduled time has elapsed and launch a
			// function to update the database once the schedule time reaches
			ticker := time.NewTicker(time.Second * time.Duration(sendTime))
			go func(collection *mongo.Collection, not models.Notification) {
				select {
				case <-ticker.C:
					// send notification mail
					if err := utils.SendMail(notification.From, notification.To, notification.Message); err != nil {
						n.Logger.Println("error sending notification mail:", err)
						return
					}

					// update notification send_at time and status from "PENDING" to "SENT"
					not.SentAt = utils.ParseTimeToString(time.Now())
					not.Status = "SENT"

					// update notification in database
					_, err := n.Collection.UpdateOne(context.TODO(), bson.M{"reference": not.Reference}, bson.M{"$set": not})
					if err != nil {
						n.Logger.Printf("could not update notification with reference: %v\n", not.Reference)
						return
					}
				}
			}(n.Collection, notification)

			c.JSON(http.StatusOK, gin.H{
				"reference": notification.Reference,
			})
		}

	}
}

// GetNotificationStatus gets the status of a notification i.e "PENDING" or "SENT"
func (n *NotificationHandler) GetNotificationStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		reference := c.Param("reference")

		result := n.Collection.FindOne(context.TODO(), bson.M{"reference": reference})
		if result.Err() == mongo.ErrNoDocuments {
			n.Logger.Println("no document with provided reference:", result.Err())
			c.JSON(http.StatusBadRequest, gin.H{
				"possible_error": []string{
					errormsg.REQUEST_NOT_FOUND.Error(),
				},
			})
			return
		}

		var notification models.Notification
		if err := result.Decode(&notification); err != nil {
			n.Logger.Println("error getting notification:", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"possible_error": []string{
					errormsg.GENERIC_ERROR.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": notification.Status,
		})
	}
}

// GetNotificationByReference gets a notification based on the provided reference
func (n *NotificationHandler) GetNotificationByReference() gin.HandlerFunc {
	return func(c *gin.Context) {
		reference := c.Param("reference")
		result := n.Collection.FindOne(context.TODO(), bson.M{"reference": reference})
		if result.Err() == mongo.ErrNoDocuments {
			n.Logger.Println("no document with provided reference:", result.Err())
			c.JSON(http.StatusBadRequest, gin.H{
				"possible_error": []string{
					errormsg.REQUEST_NOT_FOUND.Error(),
				},
			})
			return
		}

		var notification models.Notification
		if err := result.Decode(&notification); err != nil {
			n.Logger.Println("error getting notification:", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"possible_error": []string{
					errormsg.GENERIC_ERROR.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusOK, notification)
	}
}

// GetScheduledNotificationList gets 10 notifications for a page based on the specified page number.
// Check the .env file to modify number of notifications returned
// If the page is less than one or incorrect page number, page number will be 1
// If the page number has no notification, nothing is returned
func (n *NotificationHandler) GetScheduledNotificationList() gin.HandlerFunc {
	// load .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("error loading environmental variables")
	}

	return func(c *gin.Context) {
		// get page number and convert to int
		page := c.Param("page")
		pageNum, err := strconv.Atoi(page)
		if err != nil || pageNum < 1 {
			pageNum = 1
		}

		notPerPage, _ := strconv.Atoi(os.Getenv("NOTIFICATIONS_PER_PAGE"))

		// set number of notifications to be skipped and number to be returned(i.e 10)
		skip := int64((pageNum - 1) * notPerPage)

		limit := int64(notPerPage)
		fmt.Println(skip, limit)

		findOptions := options.FindOptions{
			Skip:  &skip,
			Limit: &limit,
		}

		// query database for notification and return notifications
		var notifications []models.Notification
		cursor, err := n.Collection.Find(context.TODO(), bson.M{}, &findOptions)
		if err != nil {
			n.Logger.Println("finding document:", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"possible_error": []string{
					errormsg.GENERIC_ERROR.Error(),
				},
			})
			return
		}
		defer cursor.Close(context.TODO())

		if err := cursor.All(context.TODO(), &notifications); err != nil {
			n.Logger.Println("decoding document:", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"possible_error": []string{
					errormsg.GENERIC_ERROR.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusOK, notifications)
	}
}
