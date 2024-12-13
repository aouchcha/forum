package userdata

import (
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"go.mod/dataBase"
	"go.mod/handlers"

	"golang.org/x/crypto/bcrypt"
)

func CheckPassword(hashed, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}

func SessionCookie(w http.ResponseWriter, session_id string, expiration time.Time) {
	cookie := &http.Cookie{
		Name:    "session_token",
		Value:   session_id,
		Expires: expiration,
	}
	http.SetCookie(w, cookie)
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		handlers.ChooseError(w, "Page Not Found", 404)
		return
	}
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var userID int
		var hashed string
		err := dataBase.Db.QueryRow("SELECT id, password FROM users WHERE username = ?", username).Scan(&userID, &hashed)
		if err != nil {
			handlers.ChooseError(w, "Invalid credentials (Wrong username)", http.StatusUnauthorized)
			return
		}

		err = CheckPassword(hashed, password)
		if err != nil {
			handlers.ChooseError(w, "Invalid Credentials (Wrong password)", http.StatusUnauthorized)

			return
		}

		_, err = dataBase.Db.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
		if err != nil {
			handlers.ChooseError(w, "Internal Server Error", 500)
			return
		}
		uuid, err := uuid.NewV4()
		if err != nil {
			handlers.ChooseError(w, "Internal Server Error", 500)
			return
		}
		session := uuid.String()
		expiration := time.Now().Add(5 * time.Minute)
		_, err = dataBase.Db.Exec("INSERT INTO sessions (session_id, user_id, expires_at) VALUES (?, ?, ?)", session, userID, expiration)
		if err != nil {
			handlers.ChooseError(w, "Internal Server Error", 500)
			return
		}

		SessionCookie(w, session, expiration)

		http.Redirect(w, r, "/forum", http.StatusFound)
		return
	} else if r.Method == http.MethodGet {
		http.ServeFile(w, r, "./templates/login.html")
		return
	}
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}
