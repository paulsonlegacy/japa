package mailer

import (
	"fmt"

	"japa/internal/config"

	"go.uber.org/zap"
)

type SendgridMailer struct {
	SiteConfig  config.SiteConfig
	EmailConfig config.SMTPConfig
	Logger      *zap.Logger
}

func (s *SendgridMailer) Send(to any, emailData *EmailData) error {
	// Here you do your Sendgrid API call.
	// Return error if it fails.
	return fmt.Errorf("SendgridMailer not implemented") // stub for now
}
