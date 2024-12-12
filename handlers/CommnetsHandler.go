package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"text/template"

	"go.mod/dataBase"
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
		ChooseError(w, "Internal Server Error", 500)
		return
	}

	temp := r.URL.Query().Get("post_id")

	post_id, err := strconv.Atoi(temp)
	if err != nil {
		ChooseError(w, "Internal Server Error", 500)
		return
	}

	username := r.URL.Query().Get("writer")

	fmt.Println("url", r.URL.Path)
	var page int
	if r.URL.Query().Get("page") == "" {
		page = 1
	} else {
		var err error
		page, err = strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			ChooseError(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	fmt.Println("Page Query", page)
	var offset int
	var DBlength int
	err = dataBase.Db.QueryRow("SELECT COUNT(*) FROM comments").Scan(&DBlength)
	fmt.Println("Data Base Length", DBlength)
	if err != nil {
		ChooseError(w, "!Internal Server Error", http.StatusInternalServerError)
		return
	}
	if page < 1 || page > int(math.Ceil(float64(DBlength/5)+1)) {
		// ChooseError(w, "Bage Request You are in the last page", 400)
		// return
		if page < 1 {
			page = 1
		} else {
			page = int(math.Ceil(float64(DBlength/5) + 1))
		}
		offset = DBlength - (DBlength - (5 * (page - 1)))
	} else {
		offset = DBlength - (DBlength - (5 * (page - 1)))
	}

	var toshow []Comment
	comments_toshow, _, err := GetComments(tmpl, w, username, toshow, post_id, offset)
	if err != nil {
		ChooseError(w, err.Error(), 500)
		return
	}
	var title, body, post_creator string
	err = dataBase.Db.QueryRow(`SELECT post_creator, title, body FROM posts WHERE id = ?`, post_id).Scan(&post_creator, &title, &body)
	if err != nil {
		ChooseError(w, "There is no post to this comment", 500)
		return
	}

	tmpl.Execute(w, struct {
		Post_Id     int
		CurrentUser string
		Title       string
		Post_writer string
		Body        string
		PageIndex  int
		DataLength int
		Comments []Comment
	}{
		Post_Id:     post_id,
		CurrentUser: username,
		Title:       title,
		Post_writer: post_creator,
		Body:        body,
		PageIndex: page,
		DataLength: int(math.Ceil(float64(DBlength)/5)),
		Comments: comments_toshow,
	})
}

func CreatCommnet(w http.ResponseWriter, r *http.Request) {
	if IsJavaScriptDisabled(r) {
		http.Redirect(w, r, "/NoJs", http.StatusSeeOther)
		return
	}
	commentBody := r.FormValue("comments")
	commentWriter := r.URL.Query().Get("writer")

	var message string
	var errCode int

	var userID int
	err := dataBase.Db.QueryRow("SELECT id FROM users WHERE username = ?", commentWriter).Scan(&userID)
	if err != nil {

		message = "Unauthorized: Please log in to comment"
		errCode = http.StatusUnauthorized
		ResponseComments(message, w, errCode)
		return
	}

	postID, err := strconv.Atoi(r.URL.Query().Get("postid"))
	if err != nil {

		message = "Bad Request: Invalid Post ID"
		errCode = http.StatusBadRequest
		ResponseComments(message, w, errCode)
		return
	}

	if commentBody == "" || len(commentBody) > 500 {

		message = "Bad Request: Comment body cannot be empty"
		errCode = http.StatusBadRequest
		ResponseComments(message, w, errCode)
		return
	}

	_, err = dataBase.Db.Exec("INSERT INTO comments (comment_body, comment_writer, comment_writer_id, post_commented_id) VALUES (?, ?, ?, ?)", commentBody, commentWriter, userID, postID)
	if err != nil {

		message = "Internal Server Error: Failed to save comment"
		errCode = http.StatusInternalServerError
		ResponseComments(message, w, errCode)
		return
	}

	message = "Comment created successfully"
	errCode = http.StatusOK
	ResponseComments(message, w, errCode)
}

func ResponseComments(message string, w http.ResponseWriter, errCode int) {
	w.Header().Set("Content-Type", "application/json")

	response := struct {
		Message string `json:"message"`
	}{
		Message: message,
	}

	w.WriteHeader(errCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"check": false, "message": "Internal Server Error"}`))
	}
}

func GetComments(tmpl *template.Template, w http.ResponseWriter, CurrentUser string, comments_toshow []Comment, id int, offset int) ([]Comment, int, error) {
	comm_rows, err := dataBase.Db.Query("SELECT * FROM comments WHERE post_commented_id = ? ORDER BY comment_id DESC LIMIT 5 OFFSET ?", id, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("aNA hNA 1", err)
			return nil, 0, errors.New("internal server error")
		} else {
			fmt.Println("aNA hNA 2", err)

			return nil, 0, errors.New("internal Server Error You Droped The comment table restar the server")
		}
	}
	defer comm_rows.Close()
	var cid int
	var commented Comment

	for comm_rows.Next() {
		var comment_id, post_commented_id, userid int
		var comment_body, comment_writer string
		var time any
		if err := comm_rows.Scan(&comment_id, &comment_body, &comment_writer, &userid, &post_commented_id, &time); err != nil {
			return nil, 0, errors.New("internal server error")
		}
		err = dataBase.Db.QueryRow(`SELECT COUNT(*) FROM likes WHERE liked_comment_id = ?`, comment_id).Scan(&commented.Comment_likes_count)
		if err != nil {
			return nil, 0, errors.New("internal server error")
		}
		err = dataBase.Db.QueryRow(`SELECT COUNT(*) FROM dislikes WHERE disliked_comment_id = ?`, comment_id).Scan(&commented.Comment_dislikes_count)
		if err != nil {
			return nil, 0, errors.New("internal server error")
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
	// for i := 0; i < len(comments_toshow)-1; i++ {
	// 	for j := i + 1; j < len(comments_toshow); j++ {
	// 		comments_toshow[i], comments_toshow[j] = comments_toshow[j], comments_toshow[i]
	// 	}
	// }
	return comments_toshow, cid, nil
}
