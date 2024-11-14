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
	Like_id       int
	Liked_Post_id int
	UserName      string
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
	fmt.Println("current user :", CurrentUser)
	fmt.Println("URL : ", r.URL.Query())
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
		var likes_Post []Likes
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
		likesRow, err := data.Db.Query("SELECT * FROM likes WHERE liked_post_id = ?;", id)
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
		defer likesRow.Close()
		for likesRow.Next() {
			var like_id, liked_post_id int
			var user_name string
			if err := likesRow.Scan(&like_id, &liked_post_id, &user_name); err != nil {
				fmt.Println("Scan ::::  ==> ", err)
				continue
			}
			likes_Post = append(likes_Post, Likes{
				Like_id:       like_id,
				Liked_Post_id: liked_post_id,
				UserName:      user_name,
			})
			fmt.Println("likes :", like_id, liked_post_id, user_name)
		}

		posts_toshow = append(posts_toshow, Post{
			Comments:          comments_toshow,
			Likes:             likes_Post,
			Postid:            id,
			Usernamepublished: usernamepublished,
			CurrentUsser:      CurrentUser,
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
	// if r.Method == "POST" {
	// 	fmt.Println("current user :", CurrentUser)
	// 	curr := r.FormValue("user")
	// 	reaction := r.FormValue("reaction")
	// 	postName := r.FormValue("postid")
	// 	fmt.Println("current :", curr)
	// 	fmt.Println("reaction :", reaction)
	// 	fmt.Println("post name :", postName)
	// }
	err = tmpl.Execute(w, struct {
		Currenuser string
		Posts      []Post
		// Likes      []Likes
	}{
		Currenuser: CurrentUser,
		Posts:      posts_toshow,
		// Likes:      likes_Post,
	})
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server", http.StatusInternalServerError)
		return
	}
}
