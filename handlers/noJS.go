package handlers

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

func NoJs(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/NoJs" {
		ChooseError(w, "Page Not Found", 404)
		return
	}

	if r.Method != http.MethodGet {
		ChooseError(w, "Method Not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "templates/NoJs.html")
}
