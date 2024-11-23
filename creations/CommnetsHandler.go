package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"

	data "main/dataBase"
)

type Comment struct {
	Curr_commenter_id      int
	Curr_commenter         string
	Comment_id             int
	Comment_body           string
	Comment_writer         string
	Post_commented_id      int
	Comment_time           any
	Comment_likes_count    int
	Comment_dislikes_count int
}

var commented Comment

func ShowComments(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/ShowComment.html")
	if err != nil {
		http.Error(w, "Internal server error in parsing showcomments", http.StatusInternalServerError)
		return
	}
	temp := r.FormValue("postid")
	post_id, err := strconv.Atoi(temp)
	if err != nil {
		http.Error(w, "Internal server error in SHow coment castion post id", http.StatusInternalServerError)
		return
	}
	username := r.FormValue("writer")
	fmt.Println(post_id, username)
	var toshow []Comment
	comments_toshow, _ := GetComments(tmpl, w, username, toshow, post_id)
	fmt.Println(comments_toshow)
	tmpl.Execute(w, comments_toshow)
}

func CreatCommnet(w http.ResponseWriter, r *http.Request) {
	comment_body := r.FormValue("comments")
	comment_writer := r.URL.Query().Get("writer")

	post_id, err := strconv.Atoi(r.URL.Query().Get("postid"))
	if err != nil {
		http.Error(w, "Internal server Error", http.StatusInternalServerError)
		return
	}
	// err = data.Db.QueryRow("SELECT post_creator FROM users WHERE id=?", postid).Scan()
	if comment_body == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	_, err = data.Db.Exec("INSERT INTO comments (comment_body, comment_writer, post_commented_id) VALUES (?, ?, ?)", comment_body, comment_writer, post_id)
	if err != nil {
		log.Println("Error inserting user:", err)
		http.Error(w, "Internal server error", 500)
		return
	}
	http.Redirect(w, r, "/forum", http.StatusSeeOther)
}

func GetComments(tmpl *template.Template, w http.ResponseWriter, CurrentUser string, comments_toshow []Comment, id int) ([]Comment, int) {
	comm_rows, err := data.Db.Query("SELECT * FROM comments WHERE post_commented_id = ?;", id)
	if err != nil {
		if err == sql.ErrNoRows {
			if err := tmpl.Execute(w, CurrentUser); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Printf("Template execution error: %v", err)
				return nil, 0
			}
		}
		http.Error(w, "ana hna Internal server error", http.StatusInternalServerError)
		return nil, 0
	}
	defer comm_rows.Close()
	var cid int
	for comm_rows.Next() {
		var comment_id, post_commented_id int
		var comment_body, comment_writer string
		var time any
		if err := comm_rows.Scan(&comment_id, &comment_body, &comment_writer, &post_commented_id, &time); err != nil {
			fmt.Println(err)
			continue
		}
		err = data.Db.QueryRow(`SELECT COUNT(*) FROM likes WHERE liked_comment_id = ?`, comment_id).Scan(&commented.Comment_likes_count)
		if err != nil {
			fmt.Println("Error fetching like count in comments ==>", err)
			http.Error(w, "Error fetching like count", http.StatusInternalServerError)
			return nil, 0
		}
		fmt.Println("Comment_likes_count : ", commented.Comment_likes_count)
		err = data.Db.QueryRow(`SELECT COUNT(*) FROM dislikes WHERE disliked_comment_id = ?`, comment_id).Scan(&commented.Comment_dislikes_count)
		if err != nil {
			fmt.Println("Error fetching like count ==>", err)
			http.Error(w, "Error fetching like count", http.StatusInternalServerError)
			return nil, 0
		}
		cid = comment_id
		comments_toshow = append(comments_toshow, Comment{
			Comment_id:             comment_id,
			Curr_commenter:         CurrentUser,
			Curr_commenter_id:      commented.Curr_commenter_id,
			Comment_body:           comment_body,
			Comment_writer:         comment_writer,
			Comment_time:           time,
			Comment_likes_count:    commented.Comment_likes_count,
			Comment_dislikes_count: commented.Comment_dislikes_count,
		})
	}
	for i := 0; i < len(comments_toshow)-1; i++ {
		for j := i + 1; j < len(comments_toshow); j++ {
			comments_toshow[i], comments_toshow[j] = comments_toshow[j], comments_toshow[i]
		}
	}
	return comments_toshow, cid
}
