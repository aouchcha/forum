package handlers

import (
	"net/http"
	"time"
)

func Guest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/guest" {
		ChooseError(w,"Page Not Found", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		ChooseError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	cookie := &http.Cookie{
		Name:    "session_token",
		Value:   "guest",
		Expires: time.Now().Add(4 * time.Minute),
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/forum", http.StatusSeeOther)
}
