package main

import (
	"fmt"
	"net/http"

	creations "main/creations"
	data "main/dataBase"
	handler "main/handler"
	reactions "main/reactions"
	userData "main/userData"
)

var port = "7089"

func middlewareAUTH(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userID int

		cookie, err := r.Cookie("session_token")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		err = data.Db.QueryRow("SELECT user_id FROM sessions WHERE session_id = ?;", cookie.Value).Scan(&userID)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		var token string
		err = data.Db.QueryRow("SELECT session_id FROM sessions WHERE user_id = ?;", userID).Scan(&token)
		if err == nil && token == cookie.Value {
			http.Redirect(w, r, "/forum", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func middlewareForum(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userID int

		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		if cookie.Value == "guest" {
			if r.URL.Path == "/forum" {
				next.ServeHTTP(w, r)
			} else {
				http.Redirect(w, r, "/login", http.StatusFound)
			}
			return
		}
		err = data.Db.QueryRow("SELECT user_id FROM sessions WHERE session_id = ?;", cookie.Value).Scan(&userID)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	http.HandleFunc("/", handler.Home)
	http.HandleFunc("/login", middlewareAUTH(userData.Login))
	http.HandleFunc("/guest", handler.Guest)
	http.HandleFunc("/register", middlewareAUTH(userData.HandleRegistration))
	http.HandleFunc("/forum", middlewareForum(handler.Forum))
	http.HandleFunc("/create_post", middlewareForum(creations.CreatePost))
	http.HandleFunc("/InsertPost", middlewareForum(creations.InsertPost))
	http.HandleFunc("/PostsLikes", middlewareForum(reactions.PostsLike))
	http.HandleFunc("/PostsDislikes", middlewareForum(reactions.PostsDislikes))
	http.HandleFunc("/CommentsLikes", middlewareForum(reactions.CommentsLike))
	http.HandleFunc("/CommentsDisLikes", middlewareForum(reactions.CommentsDislike))
	http.HandleFunc("/api/likes", middlewareForum(reactions.LikesCounterWithApi))
	http.HandleFunc("/logout", userData.Logout)
	http.HandleFunc("/style/", handler.Style)
	http.HandleFunc("/create_comment", creations.CreatCommnet)
	http.HandleFunc("/showcomments", middlewareForum(creations.ShowComments))
	fmt.Println("Server started on http://localhost:" + port)
	http.ListenAndServe(":"+port, nil)
}
