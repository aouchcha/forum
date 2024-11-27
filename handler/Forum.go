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
	LikesCounter      int
	DislikeCounter    int
	Postid            int
	Usernamepublished string
	CurrentUsser      string
	CurrentUser_id    int
	Title             string
	Body              string
	Time              any
	Categorie         string
}

type Reactions struct {
	LikeCount    int
	DislikeCount int
}

var postt Post

// likeee = data.Jaime{}

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

	var CurrentUser, CurrentSession string
	var session_id string
	cat_to_filter := r.FormValue("categories")
	cookie1, err1 := r.Cookie("session_token")
	cookie2, err2 := r.Cookie("user_token")

	if err1 != nil || err2 != nil {
		cookie3, err3 := r.Cookie("guest_token")
		if err3 != nil {
			fmt.Println("ana hna mouchkil fl cookie : ", err3)
			http.Error(w, "Bad Request if you want to continue as a guest choose it", http.StatusBadRequest)
			return
		}
		CurrentUser = cookie3.Value
		CurrentSession = "0"
	} else {
		CurrentUser = cookie2.Value
		CurrentSession = cookie1.Value
		err = data.Db.QueryRow("SELECT user_id, session_id FROM sessions WHERE session_id = ?", CurrentSession).Scan(&postt.CurrentUser_id, &session_id)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
	}

	posts_toshow, comment_id, post_id := GetPosts(cat_to_filter, tmpl, w, CurrentUser)
	// if posts_toshow == nil {
	// 	http.Error(w,"Internal Server You Droped a table while the code runing",http.StatusInternalServerError)
	// 	return
	// }
	err = tmpl.Execute(w, struct {
		Currenuser string
		comment_id int
		Post_id    int
		Posts      []Post
	}{
		Currenuser: CurrentUser,
		comment_id: comment_id,
		Post_id:    post_id,
		Posts:      posts_toshow,
	})
	fmt.Println("LikesCount :", postt.LikesCounter)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server", http.StatusInternalServerError)
		return
	}
}

func GetPosts(cat_to_filter string, tmpl *template.Template, w http.ResponseWriter, CurrentUser string) ([]Post, int, int) {
	var post_rows *sql.Rows
	var err error
	if cat_to_filter != "all" && cat_to_filter != "" {
		// post_rows, err = data.Db.Query("SELECT post_id FROM categories WHERE categorie = ?;", cat_to_filter)
		post_rows, err = data.Db.Query(`
			SELECT posts.* FROM posts
			JOIN categories ON posts.id = categories.post_id
			WHERE categories.categorie = ?`, cat_to_filter)
		if err != nil {
			fmt.Println("ERR", err)
			if err == sql.ErrNoRows {
				return nil, 0, 0
			}
		}
	} else {
		post_rows, err = data.Db.Query("SELECT * FROM posts;")
	}
	if err != nil {
		if err == sql.ErrNoRows {
			if err := tmpl.Execute(w, CurrentUser); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Printf("Template execution error: %v", err)
			}
		} else {

			fmt.Println("ERR==>", err)
			http.Error(w, "ana hna Internal server error", http.StatusInternalServerError)
			return nil, 0, 0
		}
	}
	defer post_rows.Close()
	// post_rows, err = data.Db.Query("SELECT * FROM posts WHERE id = ?;", cat_to_filter)
	var posts_toshow []Post
	var comment_id, post_id int
	for post_rows.Next() {
		var id int
		var title, body, usernamepublished string
		var time any
		if err := post_rows.Scan(&id, &usernamepublished, &title, &body, &time); err != nil {
			fmt.Println(err)
			continue
		}
		post_id = id
		var likee Reactions
		var dislikee Reactions
		err = data.Db.QueryRow(`SELECT COUNT(*) FROM likes WHERE post_id = ?`, post_id).Scan(&likee.LikeCount)
		if err != nil {
			fmt.Println("Error fetching like count ==>", err)
			http.Error(w, "Error fetching like count", http.StatusInternalServerError)
			return nil, 0, 0
		}
		err = data.Db.QueryRow(`SELECT COUNT(*) FROM dislikes WHERE post_id = ?`, post_id).Scan(&dislikee.DislikeCount)
		if err != nil {
			fmt.Println("Error fetching like count ==>", err)
			http.Error(w, "Error fetching like count", http.StatusInternalServerError)
			return nil, 0, 0
		}
		// fmt.Println("LikesCounter in posttt : ", postt.LikesCounter)
		fmt.Println("DislikeCounter in posttt ====> ", postt.DislikeCounter)
		// fmt.Println("comments id= ", comment_id, "post id= ", post_id)
		posts_toshow = append(posts_toshow, Post{
			Postid:            id,
			LikesCounter:      likee.LikeCount,
			DislikeCounter:    dislikee.DislikeCount,
			Usernamepublished: usernamepublished,
			CurrentUsser:      CurrentUser,
			CurrentUser_id:    postt.CurrentUser_id,
			Title:             title,
			Body:              body,
			Time:              time,
		})
	}
	// fmt.Println("postt.CurrentUser_id : ", postt.CurrentUser_id)
	// fmt.Println("LikesCounter in posttt : ", postt.LikesCounter)
	if err := post_rows.Err(); err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(posts_toshow)-1; i++ {
		for j := i + 1; j < len(posts_toshow); j++ {
			posts_toshow[i], posts_toshow[j] = posts_toshow[j], posts_toshow[i]
		}
	}
	return posts_toshow, comment_id, post_id
}
