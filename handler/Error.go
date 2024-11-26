package handler

import (
	"net/http"
	"text/template"
)

func ChooseError(w http.ResponseWriter, ErrMessage string, ErrCode int) {
	tmpl2, err2 := template.ParseFiles("templates/errors.html")
	if err2 != nil {
		http.Error(w, "Internal server error in parsing error page", http.StatusInternalServerError)
		return
	}
	err := tmpl2.Execute(w, struct {
		ErrMessage string
		ErrCode    int
	}{
		ErrMessage: ErrMessage,
		ErrCode:    ErrCode,
	})
	if err != nil {
		http.Error(w, "Internal server error in exuting error page", http.StatusInternalServerError)
		return
	}
}
