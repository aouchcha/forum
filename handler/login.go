package handler

import (
	"fmt"
	data "main/dataBase"
	"net/http"
	"time"
)

//	func session(username string) []byte {
//		b := make([]byte, 32)
//		_, err := rand.Read(b)
//		if err != nil {
//			return nil
//		}
//		return b
//	}
func SessionCookie(w http.ResponseWriter, session string, expiration time.Time) {
	cookie := &http.Cookie{
		Name:    "session_token_" + session,
		Value:   session,
		Expires: expiration,
	}
	http.SetCookie(w, cookie)
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		http.Error(w, "page not found", http.StatusNotFound)
		return
	}
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var user int
		err := data.Db.QueryRow("SELECT id FROM users WHERE username = ? AND password = ?", username, password).Scan(&user)
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		// session := string(session(username))
		expiration := time.Now().Add(4 * time.Minute)
		SessionCookie(w, username, expiration)
		_, err = data.Db.Exec("INSERT INTO sessions (session_id, user_id, expires_at) VALUES (?, ?, ?)", username, user, expiration)
		if err != nil {
			fmt.Println("db error!", err)
		}
		fmt.Println("loging success")
		http.Redirect(w, r, "/forum?user="+username, http.StatusFound)
		return
	} else if r.Method == http.MethodGet {
		http.ServeFile(w, r, "./templates/login.html")
		return
	}
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}
