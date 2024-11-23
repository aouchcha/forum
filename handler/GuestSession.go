package handler

import (
	"net/http"
	"time"
)

func Guest(w http.ResponseWriter, r *http.Request) {
	username := "guest"
	cookie := &http.Cookie{
		Name:    "guest_token",
		Value:   username,
		Expires: time.Now().Add(4 * time.Minute),
		MaxAge:  60,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w,r,"forum",http.StatusSeeOther)
}
