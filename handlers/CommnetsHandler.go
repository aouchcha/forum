package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"go.mod/dataBase"
	"go.mod/helpers"
)

type Comment struct {
	Curr_commenter_id      int
	Curr_commenter         string
	Comment_id             string
	Comment_body           string
	Comment_writer         string
	Post_commented_id      int
	Comment_time           any
	Comment_likes_count    int
	Comment_dislikes_count int
}

func ShowComments(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/showcomments" {
		ChooseError(w, "Page Not Found", http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("templates/showcomments.html")
	if err != nil {
		ChooseError(w, "Internal Server Error", 500)
		return
	}

	temp := helpers.Unhash(r.URL.Query().Get("post_id"))
	post_id, err := strconv.Atoi(temp)
	if err != nil {
		ChooseError(w, "Bad Request You change in the post id", 500)
		return
	}

	// Get the userid and the username values from the session
	cookie, _ := r.Cookie("session_token")
	var username string
	dataBase.Db.QueryRow("SELECT username FROM users INNER JOIN sessions ON users.id = sessions.user_id WHERE sessions.session_id = ?", cookie.Value).Scan(&username)

	// Fix The number of pages
	var page int
	if r.URL.Query().Get("page") == "" {
		page = 1
	} else {
		var err error
		page, err = strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			ChooseError(w, "Page Not found", 404)
			return
		}
	}

	var offset int
	var DBlength int
	err = dataBase.Db.QueryRow("SELECT COUNT(*) FROM comments").Scan(&DBlength)
	if err != nil {
		ChooseError(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if page < 1 || page > int(math.Ceil(float64(DBlength/5)+1)) {
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
		ChooseError(w, "There is no post to this comment", http.StatusBadRequest)
		return
	}
	tmpl.Execute(w, struct {
		Post_Id     string
		CurrentUser string
		Title       string
		Post_writer string
		Body        string
		PageIndex   int
		DataLength  int
		Comments    []Comment
	}{
		Post_Id:     helpers.Hash(post_id),
		CurrentUser: username,
		Title:       title,
		Post_writer: post_creator,
		Body:        body,
		PageIndex:   page,
		DataLength:  int(math.Ceil(float64(DBlength) / 5)),
		Comments:    comments_toshow,
	})
}

func CreatCommnet(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/create_comment" {
		ChooseError(w, "Page Not Found", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodPost {
		ChooseError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if IsJavaScriptDisabled(r) {
		http.Redirect(w, r, "/NoJs", http.StatusSeeOther)
		return
	}
	commentBody := strings.TrimLeft(r.FormValue("comments"), " ")
	// Get the userid and the username values from the session
	cookie, _ := r.Cookie("session_token")
	var userID int
	var commentWriter string
	err9 := dataBase.Db.QueryRow("SELECT id, username FROM users INNER JOIN sessions ON users.id = sessions.user_id WHERE sessions.session_id = ?", cookie.Value).Scan(&userID, &commentWriter)
	if err9 != nil {
		ChooseError(w, "You must login first", http.StatusUnauthorized)
		return
	}
	var message string
	var errCode, Check int
	temp := helpers.Unhash(r.URL.Query().Get("postid"))
	postID, err := strconv.Atoi(temp)
	if err != nil {
		message = "Bad Request: Invalid Post ID"
		errCode = http.StatusBadRequest
		ResponseComments(message, w, errCode)
		return
	}

	err = dataBase.Db.QueryRow("SELECT COUNT(*) FROM posts WHERE id = ?", postID).Scan(&Check)
	if err != nil || Check == 0 {
		message = "Bad Request: Invalid Post ID"
		errCode = http.StatusBadRequest
		ResponseComments(message, w, errCode)
		return
	}

	if commentBody == "" || ContentLength(commentBody) > 1000 {
		message = "Bad Request: Comment body cannot be empty and can't depasse 500 char"
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
		fmt.Println("Error encoding response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"check": false, "message": "Internal Server Error"}`))
	}
}

func GetComments(tmpl *template.Template, w http.ResponseWriter, CurrentUser string, comments_toshow []Comment, id int, offset int) ([]Comment, int, error) {
	comm_rows, err := dataBase.Db.Query("SELECT * FROM comments WHERE post_commented_id = ? ORDER BY comment_id DESC LIMIT 5 OFFSET ?", id, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, 0, errors.New("internal server error")
		} else {
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
		// if comment_writer != "" {
		// 	ChooseError(w, "You are not allowed to see this comment", http.StatusForbidden)
		// 	return nil, 0, errors.New("internal server error")
		// }
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
			Comment_id:             helpers.Hash(comment_id),
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
	return comments_toshow, cid, nil
}

func ContentLength(s string) int {
	var content_length int
	for _, char := range s {
		if char != '\r' {
			content_length++
		}
	}
	return content_length
}
