package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"path/filepath"

    "japa/internal/config"

	"go.uber.org/zap"
)

type SMTPMailer struct {
	ServerConfig    config.ServerConfig
	EmailConfig     config.SMTPConfig
	Logger          *zap.Logger
}

// host string, port int, username, password, from, templateDir string
func NewSMTPMailer(serverConfig config.ServerConfig, emailConfig config.SMTPConfig, logger *zap.Logger) *SMTPMailer {
	return &SMTPMailer{
		ServerConfig: serverConfig,
		EmailConfig:  emailConfig,
		Logger:       logger,
	}
}

// SendTemplate sends an email using an HTML template and dynamic data
func (s *SMTPMailer) SendViaSMTP(to string, subject string, data MailData) error {
	body, err := s.parseTemplate(data)
	if err != nil {
		return err
	}
	
	if err := s.send(to, subject, body); err != nil {
		s.Logger.Error("error sending mail via smtp", zap.Error(err))
		return err
	}

	return nil
}

// send handles the low-level email delivery
func (s *SMTPMailer) send(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", s.EmailConfig.EMAIL_HOST, s.EmailConfig.EMAIL_PORT)
	auth := smtp.PlainAuth("", s.EmailConfig.EMAIL_USERNAME, s.EmailConfig.EMAIL_PASSWORD, s.EmailConfig.EMAIL_HOST)

	msg := []byte(fmt.Sprintf(
		"To: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n%s",
		to, subject, body))

	return smtp.SendMail(addr, auth, s.ServerConfig.EmailFrom, []string{to}, msg)
}

// parseTemplate loads a template file and injects data
func (s *SMTPMailer) parseTemplate(data MailData) (string, error) {
	templatePath := filepath.Join(s.ServerConfig.TemplateDir, s.ServerConfig.EmailTemplate)
	tpl, err := template.ParseFiles(templatePath)
	if err != nil {
		s.Logger.Error("template parsing error", zap.Error(err)) // Log error
		return "", err
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		s.Logger.Error("template buffer error", zap.Error(err)) // Log error
		return "", err
	}

	return buf.String(), nil
}