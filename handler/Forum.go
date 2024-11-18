package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	data "main/dataBase"
)

type Post struct {
	Comments          []Comment
	LikesCounter      int
	Postid            int
	Usernamepublished string
	CurrentUsser      string
	CurrentUser_id    int
	Title             string
	Body              string
	Time              any
	Categorie         string
}

type Comment struct {
	Liked_comment     []Likes
	Comment_id        int
	Comment_body      string
	Comment_writer    string
	Post_commented_id int
	Comment_time      any
}

type Likes struct {
	LikeCount       int
	Liked_Post_id   int
	Liked_User_id   int
	Liked_User_name string
}

func Forum(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/forum" {
		http.Error(w, "page not found", http.StatusNotFound)
		return
	}
	tmpl, err := template.ParseFiles("templates/forum.html")
	if err != nil {
		http.Error(w, "Internal Server Error with forum html page", http.StatusInternalServerError)
		return
	}
	cat_to_filter := r.FormValue("categories")
	CurrentUser := r.URL.Query().Get("user")
	// fmt.Println("current user :", CurrentUser)
	fmt.Println("URL : ", r.URL.Query())
	_, err = r.Cookie("session_token_" + CurrentUser)
	if err != nil {
		http.Error(w, "Your session is expaired login again", http.StatusNotFound)
		return
	}
	var session_id string
	var curr_user_id int
	err = data.Db.QueryRow("SELECT user_id, session_id FROM sessions WHERE session_id = ?", CurrentUser).Scan(&curr_user_id, &session_id)
	if err != nil {
		http.Error(w, "You need to log in", http.StatusNotFound)
		return
	}
	fmt.Println("cuurr user id :", curr_user_id, "userName : ", CurrentUser)
	// to get filtered posts
	var post_rows *sql.Rows
	if cat_to_filter != "all" && cat_to_filter != "" {
		post_rows, err = data.Db.Query("SELECT * FROM posts WHERE categorie = ?;", cat_to_filter)
	} else {
		post_rows, err = data.Db.Query("SELECT * FROM posts;")
	}
	if err != nil {
		if err == sql.ErrNoRows {
			if err := tmpl.Execute(w, CurrentUser); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Printf("Template execution error: %v", err)
			}
			return
		}
		fmt.Println(err)
		http.Error(w, "ana hna Internal server error", http.StatusInternalServerError)
		return
	}
	defer post_rows.Close()
	// to get comments
	var posts_toshow []Post
	for post_rows.Next() {
		var comments_toshow []Comment
		// var likes_Post []Likes
		var id int
		var title, body, usernamepublished, categorie string
		var time any
		if err := post_rows.Scan(&id, &usernamepublished, &title, &body, &categorie, &time); err != nil {
			fmt.Println(err)
			continue
		}
		// fmt.Println("POST id :", id)
		comm_rows, err := data.Db.Query("SELECT * FROM comments WHERE post_commented_id = ?;", id)
		if err != nil {
			if err == sql.ErrNoRows {
				if err := tmpl.Execute(w, CurrentUser); err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					log.Printf("Template execution error: %v", err)
					return
				}
			}
			http.Error(w, "ana hna Internal server error", http.StatusInternalServerError)
			return
		}
		defer comm_rows.Close()
		for comm_rows.Next() {
			var comment_id, post_commented_id int
			var comment_body, comment_writer string
			var time any
			if err := comm_rows.Scan(&comment_id, &comment_body, &comment_writer, &post_commented_id, &time); err != nil {
				fmt.Println(err)
				continue
			}
			comments_toshow = append(comments_toshow, Comment{
				Comment_id:     comment_id,
				Comment_body:   comment_body,
				Comment_writer: comment_writer,
				Comment_time:   time,
			})
		}
		var likeCount int
		err = data.Db.QueryRow(`SELECT COUNT(*) FROM likes WHERE post_id = ?`, id).Scan(&likeCount)
		if err != nil {
			fmt.Println("Error fetching like count ==>", err)
			http.Error(w, "Error fetching like count", http.StatusInternalServerError)
			return
		}
		// likeData := Likes{
		// 	LikeCount: likeCount,
		// }
		fmt.Println("likeCount :", likeCount)
		posts_toshow = append(posts_toshow, Post{
			Comments:          comments_toshow,
			LikesCounter:      likeCount,
			Postid:            id,
			Usernamepublished: usernamepublished,
			CurrentUsser:      CurrentUser,
			CurrentUser_id:    curr_user_id,
			Title:             title,
			Body:              body,
			Time:              time,
			Categorie:         categorie,
		})

	}
	if err := post_rows.Err(); err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(posts_toshow)-1; i++ {
		for j := i + 1; j < len(posts_toshow); j++ {
			posts_toshow[i], posts_toshow[j] = posts_toshow[j], posts_toshow[i]
		}
	}
	err = tmpl.Execute(w, struct {
		Currenuser string
		Posts      []Post
	}{
		Currenuser: CurrentUser,
		Posts:      posts_toshow,
	})
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server", http.StatusInternalServerError)
		return
	}
}
