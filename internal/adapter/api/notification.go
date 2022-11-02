package api

import (
	"log"
	"net/http"
	"notification-service/internal/core/domain/entity"
	"notification-service/internal/core/helper"

	"github.com/gin-gonic/gin"
)

func (hdl *HTTPHandler) GetNotificationByRef(c *gin.Context) {
	notification, err := hdl.notificationService.GetNotificationByRef(c.Param("reference"))

	if err != nil {
		c.AbortWithStatusJSON(404, err)
		return
	}

	c.JSON(200, notification)
}

func (hdl *HTTPHandler) GetNotificationList(c *gin.Context) {
	notifications, err := hdl.notificationService.GetNotificationList(c.Param("page"))

	if err != nil {
		c.AbortWithStatusJSON(404, err)
		return
	}

	c.JSON(200, notifications)
}
func (hdl *HTTPHandler) CreateNotification(c *gin.Context) {
	body := entity.Notification{}
	if err := c.BindJSON(&body); err != nil {
		helper.LogEvent("INFO", "invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": []helper.ErrorBody{
				{
					Code:    helper.ValidationError,
					Source:  helper.Config.ServiceName,
					Message: "The request has validation errors",
				},
			},
		})
		return
	}

	reference, err := hdl.notificationService.CreateNotification(body)
	if err != nil {
		c.AbortWithStatusJSON(400, err)
		return
	}
	c.JSON(201, gin.H{"reference": reference})
}

func (hdl *HTTPHandler) GetNotificationStatus(c *gin.Context) {
	log.Println(c.Param("reference"))
	status, err := hdl.notificationService.GetNotificationStatus(c.Param("reference"))

	if err != nil {
		c.AbortWithStatusJSON(404, err)
		return
	}

	c.JSON(200, gin.H{
		"status": status,
	})
}
