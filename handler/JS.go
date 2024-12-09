package handler

import (
	"net/http"
	"strings"
)

func IsJavaScriptDisabled(r *http.Request) bool {
	acceptHeader := r.Header.Get("Accept")
	if (!strings.Contains(acceptHeader, "application/javascript") || !strings.Contains(acceptHeader, "text/javascript")) && acceptHeader != "*/*" {
		return true
	}
	return false
}
