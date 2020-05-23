package utils

import "fmt"

type PostInfo struct {
	SourceTitle string
	PostTitle   string
	URL         string
	Description string
}

func FormatPost(pi PostInfo) string {
	return fmt.Sprintf("<b>%s</b>\n<a href='%s'>%s</a>\n\n%s",
		pi.SourceTitle, pi.URL, pi.PostTitle, pi.Description)
}
