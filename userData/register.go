package handler

import (
	"net/http"

	data "main/dataBase"
)

func HandleRegistration(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/register" {
		http.Error(w, "page not found", http.StatusNotFound)
		return
	}
	if r.Method == "GET" {
		http.ServeFile(w, r, "./templates/register.html")
		return
	}

	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		email := r.FormValue("email")
		username := r.FormValue("username")
		password := r.FormValue("password")

		if email == "" || username == "" || password == "" {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		data.Db.Exec(`
            INSERT INTO users (email, username, password)
            VALUES (?, ?, ?)`, email, username, password)

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
