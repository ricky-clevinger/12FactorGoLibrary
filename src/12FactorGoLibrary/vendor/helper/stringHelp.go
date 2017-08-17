package helper

import (
	"html"
	"strings"
)

func HTMLClean(toClean string) string {

	toClean = strings.TrimSpace(toClean)
	cleaned := html.EscapeString(toClean)

	return cleaned
}
