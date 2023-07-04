package repository

import (
	"crypto/tls"
	configuration "walls-notification-service/internal/core/helper/configuration-helper"

	"log"
	"net/smtp"
)

func ConnectToGoogleSmtp() (*smtp.Client, error) {
	configuration := configuration.ServiceConfiguration

	// Set up authentication
	auth := smtp.PlainAuth("", configuration.GoogleSmtpUsername, configuration.GoogleSmtpPassword, configuration.GoogleSmtpHost)

	// Connect to Google SMTP server
	client, err := smtp.Dial(configuration.GoogleSmtpHost + ":" + configuration.GoogleSmtpPort)
	if err != nil {
		log.Fatalf("Failed to connect to Google SMTP server: %v", err)
		return nil, err
	}
	// Start TLS encryption
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         configuration.GoogleSmtpHost,
	}
	err = client.StartTLS(tlsConfig)
	if err != nil {
		log.Fatalf("Failed to start TLS: %v", err)
		client.Close()
		return nil, err
	}
	// Authenticate with Google SMTP server
	err = client.Auth(auth)
	if err != nil {
		log.Fatalf("Failed to authenticate with Google SMTP server: %v", err)
		client.Close()
		return nil, err
	}

	return client, nil
}
