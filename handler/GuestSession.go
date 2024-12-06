package handler

import (
	"fmt"
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
	fmt.Println("guest handle")
	http.Redirect(w,r,"/forum",http.StatusSeeOther)
}
