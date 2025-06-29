package mailer

import (
	"fmt"
	"time"
	
	"japa/internal/config"
)

var (
	SiteName   string = config.Settings.SiteConfig.SiteName
	SiteDomain string = config.Settings.SiteConfig.SiteDomain
	SiteEmail  string = config.Settings.SiteConfig.SiteEmail
	LogoURL    string = config.Settings.SiteConfig.LogoURL
	Year       int    = time.Now().Year()
)

type EmailData struct {
	Name          string
	Subject       string
	Heading       string
	Message       string
	LinkURL       string
	LinkText      string
	LogoURL       string // optional
	ApplicationID string // optional
	SiteName      string
	SiteDomain    string
	SiteEmail     string
	Year          int
	EmailTemplate string
}

func WelcomeMail(name string) *EmailData {
	return &EmailData{
		Name:          name,
		Subject:       fmt.Sprintf("Welcome %s", name),
		SiteName:      SiteName,
		SiteEmail:     SiteEmail,
		SiteDomain:    SiteDomain,
		EmailTemplate: "welcome.html",
		Year:          Year,
	}
}

func VisaApplicationSuccessMail(name string, applicationID string) *EmailData {
	return &EmailData{
		Name:          name,
		Subject:       "Visa Application Request Successful",
		ApplicationID: applicationID,
		SiteName:      SiteName,
		SiteEmail:     SiteEmail,
		SiteDomain:    SiteDomain,
		EmailTemplate: "visa_application_success.html",
		Year:          Year,
	}
}
