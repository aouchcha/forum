package handler

import (
	"fmt"
	"html/template"
	"net/http"
)

func NoScript(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/noScript" {
		fmt.Println("rani hna no script")
		// http.Error(w, "", http.StatusNotFound)
		template, err := template.ParseFiles("templates/noScript.html")
		if err != nil {
			http.Error(w, "page not found", http.StatusNotFound)
			return
		}
		template.Execute(w, nil)
		return
	}
}
