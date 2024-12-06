package main

import (
	"fmt"
	"net/http"

	creations "main/creations"
	handler "main/handler"
	reactions "main/reactions"
	userData "main/userData"
)

var port = "9089"

// func middlewareAUTH(next http.HandlerFunc) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		var userID int
// 		cookie, err := r.Cookie("session_token")
// 		if err != nil {
// 			next.ServeHTTP(w, r)
// 			return
// 		}
// 		err = data.Db.QueryRow("SELECT user_id FROM sessions WHERE session_id = ?;", cookie.Value).Scan(&userID)
// 		if err != nil {
// 			next.ServeHTTP(w, r)
// 			return
// 		}
// 		var token string
// 		err = data.Db.QueryRow("SELECT session_id FROM sessions WHERE user_id = ?;", userID).Scan(&token)
// 		if err == nil && token == cookie.Value {
// 			http.Redirect(w, r, "/forum", http.StatusFound)
// 			return
// 		}
// 		next.ServeHTTP(w, r)
// 	})
// }

// func middlewareForum(next http.HandlerFunc) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		var userID int
// 		cookie, err := r.Cookie("session_token")
// 		if err != nil {
// 			http.Redirect(w, r, "/login", http.StatusFound)
// 			return
// 		}
// 		err = data.Db.QueryRow("SELECT user_id FROM sessions WHERE session_id = ?;", cookie.Value).Scan(&userID)
// 		if err != nil {
// 			http.Redirect(w, r, "/login", http.StatusFound)
// 			return
// 		}
// 		next.ServeHTTP(w, r)
// 	})
// }

func main() {
	http.HandleFunc("/", handler.Home)
	http.HandleFunc("/login", (userData.Login))
	http.HandleFunc("/register", userData.HandleRegistration)
	http.HandleFunc("/guest", handler.Guest)
	http.HandleFunc("/forum", (handler.Forum))     // this is where the forum would be handled after the login
	http.HandleFunc("/noScript", handler.NoScript) // this is where the forum would be handled after the login
	http.HandleFunc("/create_post", creations.CreatePost)
	http.HandleFunc("/InsertPost", creations.InsertPost)
	http.HandleFunc("/PostsLikes", reactions.PostsLike)
	http.HandleFunc("/PostsDislikes", reactions.PostsDislikes)
	http.HandleFunc("/CommentsLikes", reactions.CommentsLike)
	http.HandleFunc("/CommentsDisLikes", reactions.CommentsDislike)
	http.HandleFunc("/api/likes", reactions.LikesCounterWithApi)
	http.HandleFunc("/logout", userData.Logout)
	http.HandleFunc("/style/", handler.Style)
	http.HandleFunc("/create_comment", creations.CreatCommnet)
	http.HandleFunc("/showcomments", creations.ShowComments)
	fmt.Println("Server started on http://localhost:" + port)
	http.ListenAndServe(":"+port, nil)
}
