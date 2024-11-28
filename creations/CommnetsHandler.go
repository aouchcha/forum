package creation

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	data "main/dataBase"
	handler "main/handler"
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

func ShowComments(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/ShowComment.html")
	if err != nil {
		handler.ChooseError(w, "Internal Server Error", 500)
		return
	}
	temp := r.FormValue("postid")
	post_id, err := strconv.Atoi(temp)
	if err != nil {
		handler.ChooseError(w, "Internal Server Error", 500)
		return
	}
	username := r.FormValue("writer")
	fmt.Println(post_id, username)
	var toshow []Comment
	comments_toshow, _, err := GetComments(tmpl, w, username, toshow, post_id)
	if err != nil {
		handler.ChooseError(w, err.Error(), 500)
		return
	}
	var title, body, post_creator string
	err = data.Db.QueryRow(`SELECT post_creator, title, body FROM posts WHERE id = ?`, post_id).Scan(&post_creator, &title, &body)
	if err != nil {
		handler.ChooseError(w, "There is no post to this comment", 500)
		return
	}

	// fmt.Println(comments_toshow)
	tmpl.Execute(w, struct {
		Post_Id     int
		CurrentUser string
		Title       string
		Post_writer string
		Body        string
		// Length      int
		Comments []Comment
	}{
		Post_Id:     post_id,
		CurrentUser: username,
		Title:       title,
		Post_writer: post_creator,
		Body:        body,
		// Length:      len(comments_toshow),
		Comments: comments_toshow,
	})
}

func CreatCommnet(w http.ResponseWriter, r *http.Request) {
	comment_body := r.FormValue("comments")
	comment_writer := r.URL.Query().Get("writer")

	_, err := r.Cookie("session_token")
	_, err2 := r.Cookie("user_token")

	if err != nil || err2 != nil {
		fmt.Println("ana hna akhay")
		handler.ChooseError(w, "Bad Request if you want to continue as a guest choose it", 400)
		return
	}

	post_id, err := strconv.Atoi(r.URL.Query().Get("postid"))
	if err != nil {
		handler.ChooseError(w, "Inernal Server Error", 500)
		return
	}
	// err = data.Db.QueryRow("SELECT post_creator FROM users WHERE id=?", postid).Scan()
	if comment_body == "" {
		handler.ChooseError(w, "Bad Request", 400)
		return
	}

	_, err = data.Db.Exec("INSERT INTO comments (comment_body, comment_writer, post_commented_id) VALUES (?, ?, ?)", comment_body, comment_writer, post_id)
	if err != nil {
		handler.ChooseError(w, "Inernal Server Error", 500)
		return
	}
}

func GetComments(tmpl *template.Template, w http.ResponseWriter, CurrentUser string, comments_toshow []Comment, id int) ([]Comment, int, error) {
	comm_rows, err := data.Db.Query("SELECT * FROM comments WHERE post_commented_id = ?;", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, 0, errors.New("no Feild in data base")
		} else {
			return nil, 0, errors.New("internal Server Error You Droped T-he comment table restar the server")
		}
	}
	defer comm_rows.Close()
	var cid int
	var commented Comment

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
			return nil, 0, errors.New("Error fetching like count")
		}
		err = data.Db.QueryRow(`SELECT COUNT(*) FROM dislikes WHERE disliked_comment_id = ?`, comment_id).Scan(&commented.Comment_dislikes_count)
		if err != nil {
			return nil, 0, errors.New("Error fetching like count")
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
			Post_commented_id:      id,
		})
	}
	for i := 0; i < len(comments_toshow)-1; i++ {
		for j := i + 1; j < len(comments_toshow); j++ {
			comments_toshow[i], comments_toshow[j] = comments_toshow[j], comments_toshow[i]
		}
	}
	return comments_toshow, cid, nil
}
