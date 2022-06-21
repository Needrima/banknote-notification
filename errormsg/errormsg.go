package errormsg

import (
	"errors"
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
