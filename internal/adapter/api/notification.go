package api

import (
	"walls-notification-service/internal/core/domain/dto"
	errorhelper "walls-notification-service/internal/core/helper/error-helper"
	validation "walls-notification-service/internal/core/helper/validation-helper"

	"github.com/gin-gonic/gin"
)

// @Summary Create Notification
// @Description Create an Notification
// @Tags Notification
// @Accept json
// @Produce json
// @Success 200 {string} entity.Notification.Reference "Success"
// @Failure 500 {object} helper.ErrorResponse
// @Param requestBody body dto.CreateNotification true "Notification request body"
// @Router /api/notification [post]
func (hdl *HTTPHandler) CreateNotification(c *gin.Context) {
	body := dto.CreateNotification{}
	_ = c.BindJSON(&body)
	err := validation.Validate(&body)
	if err != nil {
		c.AbortWithStatusJSON(400, err)
		return
	}

	notification, err := hdl.notificationService.CreateNotification(body)
	if err != nil {
		c.AbortWithStatusJSON(500, errorhelper.ErrorMessage(errorhelper.MongoDBError, err.Error()))
		return
	}
	c.JSON(201, gin.H{"reference": notification})
}



// @Summary Get Notification list by user device reference
// @Description Retrieve the Notification list for a user's device
// @Tags Notification
// @Accept json
// @Produce json
// @Success 200 {string} entity.Notification.Notification "Success"
// @Failure 500 {object} helper.ErrorResponse
// @Param device_reference path string true "Device Reference"
// @Param page path string true "Page"
// @Router /api/notification/device/{device_reference}/{page} [get]
func (hdl *HTTPHandler) GetNotificationByDeviceReference(c *gin.Context) {
	notifications, err := hdl.notificationService.GetNotificationByDeviceReference(c.Param("device_reference"), c.Param("page"))

	if err != nil {
		c.AbortWithStatusJSON(500, errorhelper.ErrorMessage(errorhelper.NoRecordError, err.Error()))
		return
	}

	c.JSON(200, notifications)
}


// @Summary Get Notification list by user reference
// @Description Retrieve the Notification list for a user
// @Tags Notification
// @Accept json
// @Produce json
// @Success 200 {string} entity.Notification.Notification "Success"
// @Failure 500 {object} helper.ErrorResponse
// @Param user_reference path string true "User Reference"
// @Param device_reference path string true "Page"
// @Router /api/notification/user/{user_reference}/{page} [get]
func (hdl *HTTPHandler) GetNotificationByUserReference(c *gin.Context) {
	notifications, err := hdl.notificationService.GetNotificationByUserReference(c.Param("user_reference"), c.Param("page"))

	if err != nil {
		c.AbortWithStatusJSON(500, errorhelper.ErrorMessage(errorhelper.NoRecordError, err.Error()))
		return
	}

	c.JSON(200, notifications)
}

// @Summary Get Notification by reference
// @Description Retrieve the notification by reference
// @Tags Notification
// @Accept json
// @Produce json
// @Success 200 {string} entity.Notification.Notification "Success"
// @Failure 500 {object} helper.ErrorResponse
// @Param Reference path string true "Reference"
// @Router /api/notification/{reference} [get]
func (hdl *HTTPHandler) GetNotificationByReference(c *gin.Context) {
	notifications, err := hdl.notificationService.GetNotificationByReference(c.Param("reference"))

	if err != nil {
		c.AbortWithStatusJSON(500, errorhelper.ErrorMessage(errorhelper.NoRecordError, err.Error()))
		return
	}

	c.JSON(200, notifications)
}