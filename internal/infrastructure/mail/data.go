package mailer

type MailData struct {
	Name       string
	Message    string
	Link       string
	SiteName   string
	SiteDomain string
	SiteEmail  string
}

func StandardMailData(name, message string, links ...string) MailData {
	link := ""
	if len(links) > 0 && links[0] != "" {
		link = links[0]
	}

	return MailData{
		Name:       name,
		Message:    message,
		Link:       link,
		SiteName:   "Legacy Technologies",
		SiteDomain: "www.legacywebhub.com",
		SiteEmail:  "legacywebhub@gmail.com",
	}
}
