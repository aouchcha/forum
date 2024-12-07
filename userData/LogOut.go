package handler

import (
	"net/http"

	data "main/dataBase"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	CurrentSessionCookie, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	CSC := CurrentSessionCookie.Value
	cookie1 := &http.Cookie{
		Name:   "session_token",
		MaxAge: -1,
	}

	http.SetCookie(w, cookie1)
	_, err = data.Db.Exec("DELETE FROM sessions WHERE session_id = ?", CSC)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
