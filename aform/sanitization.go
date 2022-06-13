package aform

import (
	"github.com/microcosm-cc/bluemonday"
	"html"
	"regexp"
	"strings"
)

func sanitizeToPlainText(value string) string {
	return sanitizeToNoHTML(trimSpace(value))
}

func sanitizeToOneLinePlainText(value string) string {
	return sanitizeToNoHTML(removeNewlineAndTrimSpace(value))
}

func sanitizeToNoHTML(value string) string {
	return html.UnescapeString(bluemonday.StrictPolicy().Sanitize(html.UnescapeString(value)))
}

func removeNewlineAndTrimSpace(value string) string {
	regexString := `\r?\n`
	reg := regexp.MustCompile(regexString)
	return trimSpace(reg.ReplaceAllString(value, " "))
}

func trimSpace(value string) string {
	return strings.TrimSpace(value)
}
