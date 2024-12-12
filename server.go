package main

import (
	"fmt"
	"net/http"
	"time"

	"go.mod/dataBase"
	"go.mod/handlers"
	"go.mod/reactions"
	"go.mod/userdata"
)

var port = "9845"

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
		err = dataBase.Db.QueryRow(
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
		if err != nil || cookie.Value == "" || cookie.Value == "guest" {
			next(w, r)
			return
		}
		err = dataBase.Db.QueryRow(
			"SELECT expires_at FROM sessions WHERE session_id = ?;",
			cookie.Value,
		).Scan(&expiresAt)
		fmt.Println(expiresAt)
		if err != nil || time.Now().After(expiresAt) {
			next(w, r)
			return
		}

		http.Redirect(w, r, "/forum", http.StatusFound)
	}
}

func main() {
	http.HandleFunc("/", auth(handlers.Home))
	http.HandleFunc("/login", auth(userdata.Login))
	http.HandleFunc("/guest", handlers.Guest)
	http.HandleFunc("/register", auth(userdata.HandleRegistration))
	http.HandleFunc("/logout", userdata.Logout)
	http.HandleFunc("/style/", handlers.Style)
	http.HandleFunc("/forum", middleware(handlers.Forum, true))
	http.HandleFunc("/showcomments", middleware(handlers.ShowComments, true))
	http.HandleFunc("/create_post", middleware(handlers.CreatePost, false))
	http.HandleFunc("/InsertPost", middleware(handlers.InsertPost, false))
	http.HandleFunc("/PostsLikes", middleware(reactions.PostsLike, false))
	http.HandleFunc("/PostsDislikes", middleware(reactions.PostsDislikes, false))
	http.HandleFunc("/CommentsLikes", middleware(reactions.CommentsLike, false))
	http.HandleFunc("/CommentsDisLikes", middleware(reactions.CommentsDislike, false))
	http.HandleFunc("/api/likes", reactions.LikesCounterWithApi)
	http.HandleFunc("/create_comment", handlers.CreatCommnet)
	http.HandleFunc("/NoJs", handlers.NoJs)
	fmt.Println("Server started on http://localhost:" + port)
	http.ListenAndServe(":"+port, nil)
}
