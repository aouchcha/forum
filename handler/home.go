package handler

import (
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "templates/homePage.html")
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}