package main

import (
	"fmt"
	"net/http"
	"time"

	creations "main/creations"
	data "main/dataBase"
	handler "main/handler"
	reactions "main/reactions"
	userData "main/userData"
)

var port = "8080"

func middleware(next http.HandlerFunc, allowGuest bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userID int

		cookie, err := r.Cookie("session_token")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		if cookie.Value == "guest" {
			if allowGuest {
				next(w, r)
			} else {
				http.Redirect(w, r, "/login", http.StatusFound)
			}
			return
		}
		err = data.Db.QueryRow(
			"SELECT user_id FROM sessions WHERE session_id = ?;",
			cookie.Value,
		).Scan(&userID)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		next(w, r)
	}
}
func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var expiresAt time.Time

		cookie, err := r.Cookie("session_token")
		if err != nil || cookie.Value == "" {
			next(w, r)
			return
		}
		err = data.Db.QueryRow(
			"SELECT expires_at FROM sessions WHERE session_id = ?;",
			cookie.Value,
		).Scan(&expiresAt)
		if err != nil || time.Now().After(expiresAt) {
			next(w, r)
			return
		}
		http.Redirect(w, r, "/forum", http.StatusFound)
	}
}

func main() {

	http.HandleFunc("/", handler.Home)
	http.HandleFunc("/login", auth(userData.Login))
	http.HandleFunc("/guest", handler.Guest)
	http.HandleFunc("/register", userData.HandleRegistration)
	http.HandleFunc("/logout", userData.Logout)
	http.HandleFunc("/style/", handler.Style)
	http.HandleFunc("/forum", middleware(handler.Forum, true))
	http.HandleFunc("/showcomments", middleware(creations.ShowComments, true))
	http.HandleFunc("/create_post", middleware(creations.CreatePost, false))
	http.HandleFunc("/InsertPost", middleware(creations.InsertPost, false))
	http.HandleFunc("/PostsLikes", middleware(reactions.PostsLike, false))
	http.HandleFunc("/PostsDislikes", middleware(reactions.PostsDislikes, false))
	http.HandleFunc("/CommentsLikes", middleware(reactions.CommentsLike, false))
	http.HandleFunc("/CommentsDisLikes", middleware(reactions.CommentsDislike, false))
	http.HandleFunc("/api/likes", reactions.LikesCounterWithApi)
	http.HandleFunc("/create_comment", creations.CreatCommnet)
	fmt.Println("Server started on http://localhost:" + port)
	http.ListenAndServe(":"+port, nil)
}
