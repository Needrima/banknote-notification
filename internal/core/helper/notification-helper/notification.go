package helper

import (
	"log"
	"net/smtp"
)

func SendNotification() {
	// Set up SMTP credentials
	smtpUsername := "nkenchor@osemeke.com"
	smtpPassword := "Chucky@2022"

	// Compose email
	from := "WallsPay"
	to := "noniegrmn@gmail.com"
	subject := "Your notification subject"
	body := "This is the content of your notification."

	message := "From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body

	// Set up SMTP authentication and connection
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, "smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587", auth, from, []string{to}, []byte(message))
	if err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}

}
