package handlers

import (
	"net/http"
	"os"
	"strings"
)

func Style(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/style/" || strings.HasSuffix(r.URL.Path, "/") {
		ChooseError(w, "403 | Access to this resource is forbidden !", 403)
		return
	}
	filePath := "style" + strings.TrimPrefix(r.URL.Path, "/style")
	if !isFileExists(filePath) {
		ChooseError(w, "404 | Page Not Found", 404)
		return
	}
	http.ServeFile(w, r, filePath)
}

func isFileExists(filePath string) bool {
	if filePath == "" {
		return false
	}
	info, err := os.Stat(filePath)
	if err != nil || info.IsDir() {
		return false
	}
	return true
}
