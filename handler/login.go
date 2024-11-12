package handler

import (
	"fmt"
	"net/http"
	"text/template"
	"time"

	data "main/dataBase"
)

//	func session(username string) []byte {
//		b := make([]byte, 32)
//		_, err := rand.Read(b)
//		if err != nil {
//			return nil
//		}
//		return b
//	}
type Exist struct {
	Exist bool
}

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
	// dataE := Exist{Exist: true}
	template, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var user int
		err := data.Db.QueryRow("SELECT id FROM users WHERE username = ? AND password = ?", username, password).Scan(&user)
		if err != nil {
			// http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			// Exist.Exist = false
			fmt.Println("dataE : ", Exist{Exist: false})
			template.Execute(w, Exist{Exist: false})
			// http.Error(w,template.Execute(w, Exist{Exist: false}), http.StatusUnauthorized)
			return
		}
		// session := string(session(username))
		expiration := time.Now().Add(4 * time.Hour)
		SessionCookie(w, username, expiration)
		_, err = data.Db.Exec("INSERT INTO sessions (session_id, user_id, expires_at) VALUES (?, ?, ?)", username, user, expiration)
		if err != nil {
			fmt.Println("db error!", err)
		}
		// template.Execute(w, Exist{Exist: false}))
		fmt.Println("loging success")
		http.Redirect(w, r, "/forum?user="+username, http.StatusFound)
		return
	} else if r.Method == http.MethodGet {
		fmt.Println("Exist{Exist: false})", Exist{Exist: false})
		template.Execute(w, Exist{Exist: true})
		return
	}
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}
