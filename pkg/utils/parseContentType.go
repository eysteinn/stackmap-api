package utils

import "mime"

/*
func ParseContentType(contentType string) string {
	return strings.Split(contentType, ";")[0]
}
*/

func ParseContentType(contentType string) string {
	if contentType == "" {
		return ""
	}
	if ct, _, err := mime.ParseMediaType(contentType); err == nil {
		return ct
	}
	return ""
}
