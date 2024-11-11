package handler

import (
	"fmt"
	data "main/dataBase"
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	fmt.Println("session_token_" + username)
	cookie, err := r.Cookie("session_token_" + username)
	if err != nil {
		fmt.Println("ERR", err)
		return
	}
	fmt.Println("Cookie", cookie.Value)
	_, err = data.Db.Exec("DELETE FROM sessions WHERE session_id = ?", username)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
