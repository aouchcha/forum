package main

import (
	"fmt"
	"net/http"

	handler "main/handler"
)

var port = "9090"

func main() {
	http.HandleFunc("/", handler.Home)
	http.HandleFunc("/login", handler.Login)
	http.HandleFunc("/register", handler.HandleRegistration)
	http.HandleFunc("/forum", handler.Forum) // this is where the forum would be handled after the login
	http.HandleFunc("/create_post", handler.CreatPost)
	http.HandleFunc("/logout", handler.Logout)
	http.HandleFunc("/style/", handler.Style)
	http.HandleFunc("/create_comment", handler.CreatCommnet)
	http.HandleFunc("/Likes", handler.Like)
	fmt.Println("Server started on http://localhost:" + port)
	http.ListenAndServe(":"+port, nil)
}
