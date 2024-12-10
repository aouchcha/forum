package handler

import (
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ChooseError(w, "Page Not Found", 404)
		return
	}
	if r.Method == "Post" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "templates/homePage.html")
}
