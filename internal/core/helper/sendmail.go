package helper

import (
	"fmt"
	"net/smtp"
)

// SendMail sends notification name from "from" to "to" using google smtp API
func SendMail(to, from string, message string) error {

	smtpHost := Config.SMTPHost
	smtpPort := Config.SMTPPort
	password := Config.SMTPPassword

	headers := map[string]string{
		"From":                from,
		"To":                  to,
		"Subject":             "Mail from Amirdeen",
		"MIME-Version":        "1.0",
		"Content-Type":        "text/plain; charset=utf-8;",
		"Content-Disposition": "inline",
	}

	headerMessage := ""

	for header, value := range headers {
		headerMessage += fmt.Sprintf("%s: %s\r\n", header, value)
	}

	body := headerMessage + "\r\n" + message

	auth := smtp.PlainAuth("", Config.SMTPUsername, password, smtpHost)

	// Sending email.
	if err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(body)); err != nil {
		return err
	}

	return nil
}
