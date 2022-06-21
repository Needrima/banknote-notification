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

func New(logger *log.Logger, collection *mongo.Collection) *NotificationHandler {
	return &NotificationHandler{
		Logger:     logger,
		Collection: collection,
	}
}

func (n *NotificationHandler) SendNotification() gin.HandlerFunc {
	return func(c *gin.Context) {
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

			notification.Reference = uuid.NewV4().String()
			sendtAt := time.Now()
			notification.SentAt = utils.ParseTimeToString(sendtAt)
			notification.Status = "SENT"

			result, err := n.Collection.InsertOne(context.TODO(), notification)
			if err != nil {
				n.Logger.Println("inserting document into database:", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"possible_error": []string{
						errormsg.GENERIC_ERROR.Error(),
					},
				})
			}
			n.Logger.Println("inserted id:", result.InsertedID)

			c.JSON(http.StatusOK, gin.H{
				"reference": notification.Reference,
			})

		case "SCHEDULED":
			scheduledTime, err := utils.ParseTimeStringToTime(notification.SendAt)
			if err != nil {
				n.Logger.Println("could not parse time:", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"possible_errors": []string{errormsg.VALIDATION_ERROR.Error()},
				})
				return
			}

			sendTime := utils.PeriodToScheduledTime(scheduledTime)
			if sendTime < 10 {
				n.Logger.Println("invalid schedule time")
				c.JSON(http.StatusInternalServerError, gin.H{
					"possible_errors": []string{errormsg.INVALID_SCHEDULE_DATE_ERROR.Error()},
				})
				return
			}

			notification.Reference = uuid.NewV4().String()
			notification.Status = "PENDING"

			_, err = n.Collection.InsertOne(context.TODO(), notification)
			if err != nil {
				n.Logger.Println("inserting document into database:", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"possible_error": []string{
						errormsg.GENERIC_ERROR.Error(),
					},
				})
			}

			ticker := time.NewTicker(time.Second * time.Duration(sendTime))
			go func(collection *mongo.Collection, not models.Notification) {
				select {
				case <-ticker.C:
					if err := utils.SendMail(notification.From, notification.To, notification.Message); err != nil {
						n.Logger.Println("error sending notification mail:", err)
						return
					}
					n.Logger.Println("sending scheduled mail successful")

					not.SentAt = utils.ParseTimeToString(time.Now())
					not.Status = "SENT"

					_, err := n.Collection.UpdateOne(context.TODO(), bson.M{"reference": not.Reference}, bson.M{"$set": not})
					if err != nil {
						n.Logger.Printf("could not update notification with reference: %v\n", not.Reference)
						return
					}

					n.Logger.Printf("updated notificantion with reference: %v\n", not.Reference)
				}
			}(n.Collection, notification)

			c.JSON(http.StatusOK, gin.H{
				"reference": notification.Reference,
			})
		}

	}
}

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

func (n *NotificationHandler) GetScheduledNotificationList() gin.HandlerFunc {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("error loading environmental variables")
	}

	return func(c *gin.Context) {
		page := c.Param("page")
		pageNum, err := strconv.Atoi(page)
		if err != nil || pageNum < 1 {
			pageNum = 1
		}

		notPerPage, _ := strconv.Atoi(os.Getenv("NOTIFICATIONS_PER_PAGE"))

		skip := int64((pageNum - 1) * notPerPage)

		limit := int64(notPerPage)
		fmt.Println(skip, limit)

		findOptions := options.FindOptions{
			Skip:  &skip,
			Limit: &limit,
		}

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
