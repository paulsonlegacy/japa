package mailer



func StandardMailData(name string, message string, links ...string) map[string]string {
	link := ""
	if (len(links) > 0 && links[0] != "") {
		link = links[0]
	}

	return map[string]string{
		"Name":    name,
		"Message": message,
		"Link": link,
		"SiteName": "Legacy Technologies",
		"SiteDomain": "www.legacywebhub.com",
		"SiteEmail": "legacywebhun@gmail.com",
	}
}