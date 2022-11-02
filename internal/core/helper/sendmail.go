package helper

import (
	"net/smtp"
)

// SendMail sends notification name from "from" to "to" using google smtp API
func SendMail(to string, message string) error {

	// Sender data.
	from := "oyebodeamirdeen@gmail.com"
	password := "blqgjjmsewlqzylb"

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	if err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(message)); err != nil {
		return err
	}

	return nil
}
