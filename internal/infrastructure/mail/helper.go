package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"

	"go.uber.org/zap"
)

// ParseEmailTemplate loads the base and content templates and renders the HTML.
// - templatesDir: The root folder for templates.
// - mailData: Contains the email template and data injected into the template.
// - logger: Optional logger (can be nil)
func ParseEmailTemplate(templatesDir string, mailData *EmailData, logger *zap.Logger) (string, error) {
	baseTemplatePath := filepath.Join(templatesDir, "email/base.html") // contains {{define "email_layout"}}
	contentTemplatePath := filepath.Join(templatesDir, "email", mailData.EmailTemplate) // contains {{define "content"}}

	// Parse both files
	parsedTemplate, err := template.ParseFiles(baseTemplatePath, contentTemplatePath)
	if err != nil {
		if logger != nil {
			logger.Error("template parsing error", zap.Error(err))
		}
		return "", fmt.Errorf("parsing templates: %w", err)
	}

	var buf bytes.Buffer
	if err := parsedTemplate.ExecuteTemplate(&buf, "email_layout", mailData); err != nil {
		if logger != nil {
			logger.Error("template execution error", zap.Error(err))
		}
		return "", fmt.Errorf("executing template: %w", err)
	}

	return buf.String(), nil
}