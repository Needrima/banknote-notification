package helper

import (
	"fmt"
	"net/smtp"
)

// SendMail sends notification name from "from" to "to" using twilio's sendgrid API
func SendMail(from, to, message string) error {
	
	// Sender data.
	from := "oyebodeamirdeen@gmail.com"
	password := "blqgjjmsewlqzylb"
  
	// Receiver email address.
	to := []string{
	  "oyebodeamirdeen@example.com",
	}
  
	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
  
	// Message.
	message := []byte("This is a test email message.")
	
	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)
	
	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		return err
	  
	}

	return nil
}
