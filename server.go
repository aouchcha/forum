package main

import (
	"fmt"
	handler "main/handler"
	"net/http"
)

var port = "3434"

func main() {
	http.HandleFunc("/", handler.Home)
	http.HandleFunc("/login", handler.Login)
	http.HandleFunc("/register", handler.HandleRegistration)
	http.HandleFunc("/forum", handler.Forum) //this is where the forum would be handled after the login
	http.HandleFunc("/create_post", handler.CreatPost)
	http.HandleFunc("/logout", handler.Logout)
	http.HandleFunc("/style/", handler.Style)
	http.HandleFunc("/create_comment", handler.CreatCommnet)
	fmt.Println("Server started on http://localhost:" + port)
	http.ListenAndServe(":"+port, nil)
}
