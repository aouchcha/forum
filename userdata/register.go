package userdata

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"go.mod/dataBase"
	"go.mod/handlers"
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
		handlers.ChooseError(w, "Page Not Found", 404)
		return
	}
	if r.Method == "GET" {
		http.ServeFile(w, r, "./templates/register.html")
		return
	}
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {

			handlers.ChooseError(w, "BAd Request", 400)
			return
		}
		emailRegex := `^[a-zA-Z0-9]+@[a-zA-Z0-9]+\.[a-zA-Z]+$`

		re := regexp.MustCompile(emailRegex)
		email := strings.TrimSpace(r.FormValue("email"))
		username := strings.TrimSpace(r.FormValue("username"))
		password := strings.TrimSpace(r.FormValue("password"))

		if len(username) >= 10 || !re.MatchString(email) {
			handlers.ChooseError(w, "Bad Request", 400)
			return
		}

		if email == "" || username == "" || password == "" {
			handlers.ChooseError(w, "Bad Request", 400)
			return
		}

		var existing int
		// dataBase.Db.QueryRow("SELECT id FROM users WHERE username = ? ", username).Scan(&existingID)
		// if existingID {
		// 	handlers.ChooseError(w, "Username already taken", http.StatusConflict)
		// 	return
		// }
		// dataBase.Db.QueryRow("SELECT id FROM users WHERE email = ? ", email).Scan(&existingID)
		// if existingID {
		// 	handlers.ChooseError(w, "email already taken", http.StatusConflict)
		// 	return
		// }
		err = dataBase.Db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ? OR email = ?", username, email).Scan(&existing)
		if err != nil {
			handlers.ChooseError(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		fmt.Println("NUM OF ROWS :", existing)
		if existing != 0 {
			handlers.ChooseError(w, "The username or the passweord alreday taken", http.StatusConflict)
			return
		}

		hashedPassword, err := HashPassword(password)
		if err != nil {
			handlers.ChooseError(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		_, err = dataBase.Db.Exec(`
            INSERT INTO users (email, username, password)
            VALUES (?, ?, ?)`, email, username, hashedPassword)
		if err != nil {
			handlers.ChooseError(w, "The NAme or email already taken", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/forum", http.StatusSeeOther)
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
