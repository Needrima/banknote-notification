package api

import (
	ports "notification-service/internal/port"
)

type HTTPHandler struct {
	notificationService ports.NotificationService
}

func NewHTTPHandler(countryService ports.NotificationService) *HTTPHandler {
	return &HTTPHandler{
		notificationService: countryService,
	}
}
