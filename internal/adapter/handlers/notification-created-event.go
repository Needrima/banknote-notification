package handlers

import (
	"encoding/json"
	"fmt"
	"walls-notification-service/internal/core/domain/dto"
	"walls-notification-service/internal/core/domain/shared"
	"walls-notification-service/internal/core/services"
	eto "walls-notification-service/internal/core/helper/event-helper/eto"
	ports "walls-notification-service/internal/port"
)

// Event handler function
func OtpCreatedEventHandler(event interface{}, notificationRepository ports.NotificationRepository) {

	// Convert interface{} to byte array
	jsonBytes, err := json.Marshal(event)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	// Deserialize JSON to MyStruct
	var otpCreatedEvent eto.Event
	err = json.Unmarshal(jsonBytes, &otpCreatedEvent)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	fmt.Println("Received data:", otpCreatedEvent.EventData)
	// Access the field
	device := otpCreatedEvent.EventData.(map[string]interface{})["Device"].(map[string]interface{})
	device_reference := device["Reference"]


	createNotificationDto := dto.CreateNotification{
			UserReference: otpCreatedEvent.EventData.(map[string]interface{})["Reference"].(string),
			DeviceReference:   device_reference.(string),
			Contact:         otpCreatedEvent.EventData.(map[string]interface{})["Contact"].(string),
			Channel:         otpCreatedEvent.EventData.(map[string]interface{})["Channel"].(shared.Channel),
			Subject:         otpCreatedEvent.EventData.(map[string]interface{})["Subject"].(string),
			MessageBody:     otpCreatedEvent.EventData.(map[string]interface{})["MessageBody"].(string),
	}

	// Create an instance of the NotificationService

	services.NotificationService.CreateNotification(createNotificationDto)
	

}
