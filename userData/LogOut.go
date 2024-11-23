package handler

import (
	"fmt"
	"net/http"

	data "main/dataBase"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	CurrentSessionCookie, err := r.Cookie("session_token")
	_, err1 := r.Cookie("user_token")
	if err != nil || err1 != nil{
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	fmt.Println("Cookie", CurrentSessionCookie.Value)
	CSC := CurrentSessionCookie.Value
	cookie1 := &http.Cookie{
		Name:   "session_token",
		MaxAge: -1,
	}
	cookie2 := &http.Cookie{
		Name: "user_token",
		MaxAge: -1,
	}
	http.SetCookie(w,cookie1)
	http.SetCookie(w,cookie2)

	_, err = data.Db.Exec("DELETE FROM sessions WHERE session_id = ?", CSC)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
