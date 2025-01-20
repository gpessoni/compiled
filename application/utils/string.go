package utils

import (
	"regexp"
	"strings"
)

func RemoveHTMLTags(input string) string {
	breakTags := map[string]bool{
		"br": true,
		"p":  true, "/p": true,
		"div": true, "/div": true,
		"li": true, "/li": true,
		"h1": true, "/h1": true,
		"h2": true, "/h2": true,
		"h3": true, "/h3": true,
		"h4": true, "/h4": true,
		"h5": true, "/h5": true,
		"h6": true, "/h6": true,
	}

	re := regexp.MustCompile(`<[^>]+>`)

	output := re.ReplaceAllStringFunc(input, func(tag string) string {
		cleanTag := strings.Trim(tag, "<>/")
		cleanTag = strings.ToLower(strings.Split(cleanTag, " ")[0])

		if breakTags[cleanTag] {
			return "\n"
		}
		return ""
	})

	output = strings.ReplaceAll(output, "\n ", "\n")
	output = strings.ReplaceAll(output, " \n", "\n")

	output = regexp.MustCompile(`\n{3,}`).ReplaceAllString(output, "\n\n")

	output = strings.TrimSpace(output)

	return output
}
