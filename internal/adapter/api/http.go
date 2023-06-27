package api

import (
	ports "walls-notification-service/internal/port"
)

// Httphander for the api
type HTTPHandler struct {
	notificationService ports.NotificationService
}

func NewHTTPHandler(
	countryService ports.NotificationService) *HTTPHandler {
	return &HTTPHandler{
		notificationService: countryService,
	}
}
