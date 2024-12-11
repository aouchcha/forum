package handlers

import (
	"net/http"
	"time"
)

func Guest(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:    "session_token",
		Value:   "guest",
		Expires: time.Now().Add(4 * time.Minute),
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/forum", http.StatusSeeOther)
}
