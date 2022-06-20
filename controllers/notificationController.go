package controllers

import (
	"banknote-notification-service/models"
	"banknote-notification-service/utils"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	INVALID_MESSAGE_ERROR       = errors.New("The message format read from the given topic is invalid")
	VALIDATION_ERROR            = errors.New("The request has validation errors")
	REQUEST_NOT_FOUND           = errors.New("The requested resource was NOT found")
	GENERIC_ERROR               = errors.New("Generic error occurred. See stacktrace for details")
	AUTHORIZATION_ERROR         = errors.New("You do NOT have adequate permission to access this resource")
	DUPLICATE_ENTRY_ERROR       = errors.New("Duplicate entry detected.")
	MESSAGE_SERVICE_ERROR       = errors.New("An error occurred while sending the message.")
	SMS_SERVICE_ERROR           = errors.New("An error occurred while sending SMS message.")
	INVALID_SCHEDULE_DATE_ERROR = errors.New("You cannot schedule a task in the past. You must provide a future date")
	NO_PRINCIPAL                = errors.New("Principal identifier NOT provided")
)

type NotificationController struct {
	Logger     *log.Logger
	Collection *mongo.Collection
}

func New(logger *log.Logger, collection *mongo.Collection) *NotificationController {
	return &NotificationController{
		Logger:     logger,
		Collection: collection,
	}
}

func (n *NotificationController) SendNotification() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data models.Notification
		if err := c.BindJSON(&data); err != nil {
			n.Logger.Println("error:", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": []string{VALIDATION_ERROR.Error(), DUPLICATE_ENTRY_ERROR.Error()},
			})
			return
		}

		switch data.Type {
		case "INSTANT":
			if err := utils.SendMail(data.From, data.To, data.Message); err != nil {
				n.Logger.Println("error sending notification mail:", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"errors": []string{INVALID_MESSAGE_ERROR.Error(), MESSAGE_SERVICE_ERROR.Error()},
				})
				return
			}
		case "SCHEDULED":
			scheduledTime, err := utils.ParseTime(data.SendAt)
			if err != nil {
				n.Logger.Println("could not parse time:", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"errors": []string{VALIDATION_ERROR.Error()},
				})
				return
			}

			if (scheduledTime.Hour() - time.Now().Hour()) < 24 {
				n.Logger.Println("invalid schedule time")
				c.JSON(http.StatusInternalServerError, gin.H{
					"errors": []string{VALIDATION_ERROR.Error()},
				})
				return
			}
		}

	}
}

func (n *NotificationController) GetNotificationStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		reference := c.Param("reference")
		c.JSON(http.StatusOK, gin.H{
			"status":    "PENDING",
			"reference": reference,
		})
	}
}

func (n *NotificationController) GetNotificationByReference() gin.HandlerFunc {
	return func(c *gin.Context) {
		reference := c.Param("reference")
		c.JSON(http.StatusOK, gin.H{
			"reference": reference,
		})
	}
}

func (n *NotificationController) GetScheduledNotificationList() gin.HandlerFunc {
	return func(c *gin.Context) {
		page := c.Param("page")
		c.JSON(http.StatusOK, gin.H{
			"page": page,
		})
	}
}
