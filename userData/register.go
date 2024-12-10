package handler

import (
	data "main/dataBase"
	"main/handler"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
func HandleRegistration(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/register" {
		handler.ChooseError(w, "Page Not Found", 404)
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
		// Regular expression for validating email addresses
		emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

		// Compiling the regex
		re := regexp.MustCompile(emailRegex)
		email := strings.TrimSpace(r.FormValue("email"))
		username := strings.TrimSpace(r.FormValue("username"))
		password := strings.TrimSpace(r.FormValue("password"))

		if len(username) >= 8 || !re.MatchString(email) {
			handler.ChooseError(w, "Bad Request", 400)
			return
		}

		if email == "" || username == "" || password == "" {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		var existingID int
		data.Db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&existingID)
		if existingID > 0 {
			http.Error(w, "Username already taken", http.StatusConflict)
			return
		}

		hashedPassword, err := HashPassword(password)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		data.Db.Exec(`
            INSERT INTO users (email, username, password)
            VALUES (?, ?, ?)`, email, username, hashedPassword)

		http.Redirect(w, r, "/forum", http.StatusSeeOther)
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
