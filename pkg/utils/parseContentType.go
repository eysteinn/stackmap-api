package utils

import "strings"

func ParseContentType(contentType string) string {
	return strings.Split(contentType, ";")[0]
}
