package test

import (
	//"fmt"
	"os"
	"path/filepath"
	"testing"

	"japa/internal/config"
	"japa/internal/infrastructure/mail"

	"go.uber.org/zap/zaptest"
)

func setupTestMailer(t *testing.T) *mailer.SMTPMailer {
	logger := zaptest.NewLogger(t)
	var TEMPLATES_PATH string
	BASE_PATH, err := os.Getwd()

	if err == nil {
		TEMPLATES_PATH = filepath.Join(BASE_PATH, "../../templates")
	} else {
		TEMPLATES_PATH = "../../templates"
	}

	return mailer.NewSMTPMailer(
		config.ServerConfig{TemplatesDir: TEMPLATES_PATH},
		config.SiteConfig{
			SiteName:   "TestSite",
			SiteEmail:  "noreply@test.com",
			SiteDomain: "https://test.com",
		},
		config.SMTPConfig{
			EMAIL_HOST:     "smtp.gmail.com", // 👈 set to dummy to force failure
			EMAIL_PORT:     587,
			EMAIL_USERNAME: "legacywebtechnologies@gmail.com",
			EMAIL_PASSWORD: "qwzwsjxljxgabecx",
		},
		logger,
	)
}

func TestSendWelcomeEmail_Fail(t *testing.T) {
	testMailer := setupTestMailer(t)
	data := mailer.WelcomeMail("Paulson")

	err := testMailer.Send("rqwyeyehdirikf", data)
	if err == nil {
		t.Errorf("expected error but got nil")
	}
}

func TestSendWelcomeEmail_Success(t *testing.T) {
	testMailer := setupTestMailer(t)
	data := mailer.WelcomeMail("Paulson")

	err := testMailer.Send("paulsonbosah@gmail.com", data)
	if err != nil {
		t.Errorf("expected success, got error: %v", err)
	}
}

func TestSendVisaApplicationEmail_Success(t *testing.T) {
	testMailer := setupTestMailer(t)
	data := mailer.VisaApplicationSuccessMail("Paulson", "APP123")

	err := testMailer.Send("paulsonbosah@gmail.com", data)
	if err != nil {
		t.Errorf("expected success, got error: %v", err)
	}
}

func TestSendEmail_InvalidRecipientType(t *testing.T) {
	testMailer := setupTestMailer(t)
	data := mailer.WelcomeMail("Paulson")

	err := testMailer.Send(12345, data)
	if err == nil {
		t.Errorf("expected error with invalid recipient type, got nil")
	}
}
