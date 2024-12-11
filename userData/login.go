package handler

import (
	"fmt"
	"net/http"
	"time"

	data "main/dataBase"
	"main/handler"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

//	type userExist struct {
//		Yes bool
//	}
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
		handler.ChooseError(w, "Page Not Found", 404)
		return
	}
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var userID int
		var hashed string
		err := data.Db.QueryRow("SELECT id, password FROM users WHERE username = ?", username).Scan(&userID, &hashed)
		if err != nil {
			handler.ChooseError(w, "Invalid credentials (Wrong username)", http.StatusUnauthorized)
			// http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		err = CheckPassword(hashed, password)
		if err != nil {
			handler.ChooseError(w, "Invalid Credentials (Wrong password)", http.StatusUnauthorized)
			// http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		_, err = data.Db.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
		if err != nil {
			fmt.Println("Error deleting old sessions:", err)
			handler.ChooseError(w, "Internal Server Error", 500)
			return
		}

		session := uuid.New().String()
		expiration := time.Now().Add(5 * time.Minute)
		_, err = data.Db.Exec("INSERT INTO sessions (session_id, user_id, expires_at) VALUES (?, ?, ?)", session, userID, expiration)
		if err != nil {
			handler.ChooseError(w, "Internal Server Error", 500)
			// http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		SessionCookie(w, session, expiration)
		// fmt.Println("login success")
		http.Redirect(w, r, "/forum", http.StatusFound)
		return
	} else if r.Method == http.MethodGet {
		http.ServeFile(w, r, "./templates/login.html")
		return
	}
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}
