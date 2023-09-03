package requests

import (
	"net/http"
	"strings"
)

func GetContent(request *http.Request) string {

	contenttype := request.Header.Get("Content-Type")
	contenttype = strings.Split(contenttype, ";")[0]
	return contenttype
}
