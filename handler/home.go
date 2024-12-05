package handler

import (
	"fmt"
	"net/http"
	"text/template"
	// data "main/dataBase"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "page not found", http.StatusNotFound)
		return
	}
	cookies := r.Cookies()
	if len(cookies) != 0 {

		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		CurrentUser := cookie.Value
		fmt.Println(CurrentUser)
		http.Redirect(w, r, "/forum", http.StatusSeeOther)
		return
	}

	templ, err := template.ParseFiles("templates/homePage.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	templ.Execute(w, nil)
}
