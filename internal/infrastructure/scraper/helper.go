package scraper

import "regexp"


var bgImageRe = regexp.MustCompile(`background-image:\s*url\(['"]?(.*?)['"]?\)`)

func extractBackgroundImageURL(style string) string {
	matches := bgImageRe.FindStringSubmatch(style)
	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}

// Helper function to avoid pointers to empty strings
func strPtrOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}