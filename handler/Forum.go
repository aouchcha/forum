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
	Likes             []Likes
	Postid            int
	Usernamepublished string
	CurrentUsser      string
	Title             string
	Body              string
	Time              any
	Categorie         string
}

type Comment struct {
	Comment_id        int
	Comment_body      string
	Comment_writer    string
	Post_commented_id int
	Comment_time      any
}
type Likes struct {
	Post_id  int
	UserName string
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
	_, err = r.Cookie("session_token_" + CurrentUser)
	if err != nil {
		http.Error(w, "Your session is expaired login again", http.StatusNotFound)
		return
	}
	var session_id string
	err = data.Db.QueryRow("SELECT session_id FROM sessions WHERE session_id = ?", CurrentUser).Scan(&session_id)
	if err != nil {
		http.Error(w, "You need to log in", http.StatusNotFound)
		return
	}
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
			// fmt.Println("comment Data :", comment_id, comment_body, comment_writer, post_commented_id)
			comments_toshow = append(comments_toshow, Comment{
				Comment_id:     comment_id,
				Comment_body:   comment_body,
				Comment_writer: comment_writer,
				Comment_time:   time,
			})
		}
		// like := r.FormValue("like")
		// fmt.Println("like :", like)
		posts_toshow = append(posts_toshow, Post{
			Comments:          comments_toshow,
			Postid:            id,
			Usernamepublished: usernamepublished,
			CurrentUsser:      CurrentUser,
			Title:             title,
			Body:              body,
			Time:              time,
			Categorie:         categorie,
		})
	}
	// like := r.FormValue("like")
	// dislike := r.FormValue("dislike")
	// fmt.Println("like :", like)
	// fmt.Println("dislike :", dislike)
	if err := post_rows.Err(); err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(posts_toshow)-1; i++ {
		for j := i + 1; j < len(posts_toshow); j++ {
			posts_toshow[i], posts_toshow[j] = posts_toshow[j], posts_toshow[i]
		}
	}
	// if r.Method == "POST" {
	likes := r.PostFormValue("like")
	dislikes := r.PostFormValue("dislike")
	if likes == "KAFKA_like" {
		fmt.Println("something in here")
	}
	fmt.Println("PostFormValue like :", likes)
	fmt.Println("PostFormValue dislike :", dislikes)
	fmt.Println("current user :", CurrentUser)
	like := r.FormValue("like")
	// dislike := r.FormValue("dislike")
	fmt.Println("like :", like)
	// }
	err = tmpl.Execute(w, struct {
		Currenuser string
		Posts      []Post
		// Likes      []Likes
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
