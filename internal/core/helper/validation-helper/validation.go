package helper

import (
	"regexp"
	"walls-notification-service/internal/core/domain/shared"

	"github.com/go-playground/validator/v10"
)

func ValidateValidChannel(fl validator.FieldLevel) bool {
	channel := fl.Field().Interface().(shared.Channel)

	// Check if the channel is either Phone or Email
	return channel == shared.Sms || channel == shared.Email
}

func ValidateValidContact(fl validator.FieldLevel) bool {
	contact := fl.Field().String()
	channel := fl.Parent().Elem().FieldByName("Channel").Interface().(shared.Channel)

	// Regular expression patterns for email and phone number
	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	phonePattern := `^\+\d{1,3}\d{4,}$`

	if channel == shared.Sms {
		// Check if the contact matches the phone number pattern
		match, _ := regexp.MatchString(phonePattern, contact)
		return match
	} else if channel == shared.Email {
		// Check if the contact matches the email pattern
		match, _ := regexp.MatchString(emailPattern, contact)
		return match
	}

	return false
}
