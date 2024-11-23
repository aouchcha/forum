package handler

import (
	"fmt"
	"net/http"
	"text/template"
	"time"

	data "main/dataBase"

	"github.com/gofrs/uuid"
)

type Exist struct {
	Exist bool
}

func session() string {
	id, _ := uuid.NewV7()
	// fmt.Println(id.String())
	return id.String()
}

func SessionCookie(w http.ResponseWriter, session string, expiration time.Time) {
	cookie := &http.Cookie{
		Name:    "session_token",
		Value:   session,
		Expires: expiration,
		MaxAge:  3600,
	}
	http.SetCookie(w, cookie)
}

func UserCookie(w http.ResponseWriter, session string, expiration time.Time) {
	cookie := &http.Cookie{
		Name:    "user_token",
		Value:   session,
		Expires: expiration,
		MaxAge:  3600,
	}
	http.SetCookie(w, cookie)
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		http.Error(w, "page not found", http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, "Internal server Error", http.StatusInternalServerError)
		return
	}
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")
		cookie := &http.Cookie{
			Name:   "guest_token",
			MaxAge: -1,
		}
		http.SetCookie(w, cookie)
		var userid int
		err := data.Db.QueryRow("SELECT id FROM users WHERE username = ? AND password = ?", username, password).Scan(&userid)
		if err != nil {
			// http.Error(w, "Invalid credentials", http.StatusUnauthorized)

			tmpl.Execute(w, Exist{Exist: false})
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		expiration := time.Now().Add(4 * time.Hour)
		session := session()
		SessionCookie(w, string(session), expiration)
		UserCookie(w, username, expiration)
		err = data.Db.QueryRow("SELECT user_id FROM sessions WHERE user_id = ? ", userid).Scan(&userid)
		if err != nil {
			fmt.Println(err)
			data.Db.Exec("INSERT INTO sessions (session_id, user_id, expires_at) VALUES (?, ?, ?)", session, userid, expiration)
		} else {
			fmt.Println("we find this session and we updat it")
			fmt.Println("USER ID", userid)
			_, err := data.Db.Exec("UPDATE sessions SET session_id = ?, expires_at = ?WHERE user_id = ?", session, expiration, userid)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		// fmt.Println("user id :",userid)

		http.Redirect(w, r, "/forum", http.StatusFound)
		return
	} else if r.Method == http.MethodGet {
		tmpl.Execute(w, Exist{Exist: true})
		return
	}
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}
