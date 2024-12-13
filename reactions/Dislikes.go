package reactions

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"go.mod/dataBase"
	"go.mod/handlers"
)

func PostsDislikes(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/PostsDislikes" {
		handlers.ChooseError(w, "Page Not Found", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodPost {
		handlers.ChooseError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	if handlers.IsJavaScriptDisabled(r) {
		http.Redirect(w, r, "/NoJs", http.StatusSeeOther)
		return
	}
	user := r.URL.Query().Get("user")
	postid := r.URL.Query().Get("Disliked_Post_id")
	user_id, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		handlers.ChooseError(w, "Internal server error", http.StatusInternalServerError)

		return
	}
	var exists bool
	err = dataBase.Db.QueryRow("SELECT id FROM dislikes WHERE user_id = ? AND post_id = ? AND username = ?", user_id, postid, user).Scan(&exists)
	if err == sql.ErrNoRows {

		err = dataBase.Db.QueryRow("SELECT id FROM likes WHERE user_id = ? AND post_id = ? AND username = ?", user_id, postid, user).Scan(&exists)
		if err != sql.ErrNoRows {
			_, err = dataBase.Db.Exec("DELETE FROM likes WHERE user_id = ? AND post_id = ? AND username = ?", user_id, postid, user)
			if err != nil {

				handlers.ChooseError(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		_, err = dataBase.Db.Exec("INSERT INTO dislikes (user_id, post_id, username) VALUES (?, ?, ?)", user_id, postid, user)
		if err != nil {

			handlers.ChooseError(w, "Internal server error", http.StatusInternalServerError)
			return
		}

	} else {

		_, err = dataBase.Db.Exec("DELETE FROM dislikes WHERE user_id = ? AND post_id = ? AND username = ?", user_id, postid, user)
		if err != nil {

			handlers.ChooseError(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func CommentsDislike(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/CommentsDisLikes" {
		handlers.ChooseError(w, "Page Not Found", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodPost {
		handlers.ChooseError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	if handlers.IsJavaScriptDisabled(r) {
		http.Redirect(w, r, "/NoJs", http.StatusSeeOther)
		return
	}
	User := r.URL.Query().Get("user")
	fmt.Println("USER", User)
	Disliked_comment_id, err := strconv.Atoi(r.URL.Query().Get("comment_id"))
	fmt.Println("Disliked_comment_id", )
	if err != nil {
		handlers.ChooseError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	var User_id int
	err = dataBase.Db.QueryRow("SELECT id FROM users WHERE username = ?", User).Scan(&User_id)
	if err != nil {
		if err == sql.ErrNoRows {

			handlers.ChooseError(w, "Internal server error", http.StatusInternalServerError)
			return
		} else {

			handlers.ChooseError(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
	var Exist bool

	err = dataBase.Db.QueryRow("SELECT id FROM dislikes WHERE user_id = ? AND username = ? AND disliked_comment_id = ?", User_id, User, Disliked_comment_id).Scan(&Exist)

	if err == sql.ErrNoRows {

		err = dataBase.Db.QueryRow("SELECT id FROM likes WHERE user_id = ? AND username = ? AND liked_comment_id = ?", User_id, User, Disliked_comment_id).Scan(&Exist)
		if err != sql.ErrNoRows {
			_, err = dataBase.Db.Exec("DELETE FROM likes WHERE user_id = ? AND username = ? AND liked_comment_id = ?", User_id, User, Disliked_comment_id)
			if err != nil {

				handlers.ChooseError(w, "Internal server error", http.StatusInternalServerError)
				return
			}

		}

		_, err = dataBase.Db.Exec("INSERT INTO dislikes (user_id, disliked_comment_id, username) VALUES (?,?,?)", User_id, Disliked_comment_id, User)
		if err != nil {

			handlers.ChooseError(w, "Internal server error", http.StatusInternalServerError)
			return
		}

	} else {
		_, err = dataBase.Db.Exec("DELETE FROM dislikes WHERE user_id = ? AND username = ? AND disliked_comment_id = ?", User_id, User, Disliked_comment_id)
		if err != nil {

			handlers.ChooseError(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}
