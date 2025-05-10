package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"path/filepath"

    "japa/internal/config"
)

type SMTPMailer struct {
	Host         string
	Port         int
	Username     string
	Password     string
	From         string
	TemplateDir  string // directory where templates are stored
}

// host string, port int, username, password, from, templateDir string
func NewSMTPMailer(config config.EmailConfig) *SMTPMailer {
	return &SMTPMailer{
		Host:        config.EMAIL_HOST,
		Port:        config.EMAIL_PORT,
		Username:    config.EMAIL_USERNAME,
		Password:    config.EMAIL_PASSWORD,
		From:        config.EMAIL_FROM,
		TemplateDir: config.TEMPLATE_DIR,
	}
}

// SendTemplate sends an email using an HTML template and dynamic data
func (s *SMTPMailer) SendTemplate(to, subject, templateFile string, data any) error {
	body, err := s.parseTemplate(templateFile, data)
	if err != nil {
		return err
	}
	return s.send(to, subject, body)
}

// send handles the low-level email delivery
func (s *SMTPMailer) send(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)

	msg := []byte(fmt.Sprintf(
		"To: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n%s",
		to, subject, body))

	return smtp.SendMail(addr, auth, s.From, []string{to}, msg)
}

// parseTemplate loads a template file and injects data
func (s *SMTPMailer) parseTemplate(templateFile string, data any) (string, error) {
	templatePath := filepath.Join(s.TemplateDir, templateFile)
	template, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := template.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}