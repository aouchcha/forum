package handler

import (
	"fmt"
	"net/http"
	"text/template"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "page not found", http.StatusNotFound)
		return
	}
	templ, err := template.ParseFiles("templates/homePage.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	templ.Execute(w, nil)
}
