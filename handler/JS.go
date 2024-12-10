package handler

import (
	"fmt"
	"net/http"
	"strings"
)

func IsJavaScriptDisabled(r *http.Request) bool {
	acceptHeader := r.Header.Get("Accept")
	fmt.Println("Accept Header:", acceptHeader)
	if (!strings.Contains(acceptHeader, "application/javascript") || !strings.Contains(acceptHeader, "text/javascript")) && acceptHeader != "*/*" {
		return true
	}
	return false
}
