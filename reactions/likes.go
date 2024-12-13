package reactions

import (
	"database/sql"
	"net/http"
	"strconv"

	"go.mod/dataBase"
	"go.mod/handlers"
)

func PostsLike(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handlers.ChooseError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/PostsLikes" {
		handlers.ChooseError(w, "Page Not Found", http.StatusNotFound)
		return
	}

	if handlers.IsJavaScriptDisabled(r) {
		http.Redirect(w, r, "/NoJs", http.StatusSeeOther)
		return
	}
	user := r.URL.Query().Get("user")
	postid := r.URL.Query().Get("Liked_Post_id")
	var user_id int
	err := dataBase.Db.QueryRow("SELECT id FROM users WHERE username=?", user).Scan(&user_id)
	if err != nil {

		handlers.ChooseError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	var exists bool
	err = dataBase.Db.QueryRow("SELECT id FROM likes WHERE user_id = ? AND post_id = ? AND username = ?", user_id, postid, user).Scan(&exists)
	if err == sql.ErrNoRows {

		err = dataBase.Db.QueryRow("SELECT id FROM dislikes WHERE user_id = ? AND post_id = ? AND username = ?", user_id, postid, user).Scan(&exists)
		if err != sql.ErrNoRows {
			_, err = dataBase.Db.Exec("DELETE FROM dislikes WHERE user_id = ? AND post_id = ? AND username = ?", user_id, postid, user)
			if err != nil {

				handlers.ChooseError(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		_, err = dataBase.Db.Exec("INSERT INTO likes (user_id, post_id, username) VALUES (?, ?, ?)", user_id, postid, user)
		if err != nil {

			handlers.ChooseError(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	} else {

		_, err = dataBase.Db.Exec("DELETE FROM likes WHERE user_id = ? AND post_id = ? AND username = ?", user_id, postid, user)
		if err != nil {
			handlers.ChooseError(w, "Internal server error", http.StatusInternalServerError)

			return
		}
	}
}

func CommentsLike(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handlers.ChooseError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/CommentsLikes" {
		handlers.ChooseError(w, "Page Not Found", http.StatusNotFound)
		return
	}

	if handlers.IsJavaScriptDisabled(r) {
		http.Redirect(w, r, "/NoJs", http.StatusSeeOther)
		return
	}
	Liked_comment_id, err := strconv.Atoi(r.URL.Query().Get("comment_id"))
	if err != nil {
		handlers.ChooseError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	User := r.URL.Query().Get("user")
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

	err = dataBase.Db.QueryRow("SELECT id FROM likes WHERE user_id = ? AND username = ? AND liked_comment_id = ?", User_id, User, Liked_comment_id).Scan(&Exist)
	if err == sql.ErrNoRows {
		err = dataBase.Db.QueryRow("SELECT id FROM dislikes WHERE user_id = ? AND username = ? AND disliked_comment_id = ?", User_id, User, Liked_comment_id).Scan(&Exist)
		if err != sql.ErrNoRows {
			_, err = dataBase.Db.Exec("DELETE FROM dislikes WHERE user_id = ? AND username = ? AND disliked_comment_id = ?", User_id, User, Liked_comment_id)
			if err != nil {
				handlers.ChooseError(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		_, err = dataBase.Db.Exec("INSERT INTO likes (user_id, username, liked_comment_id) VALUES (?,?,?)", User_id, User, Liked_comment_id)
		if err != nil {
			handlers.ChooseError(w, "Internal server error", http.StatusInternalServerError)
			return
		}

	} else {
		_, err = dataBase.Db.Exec("DELETE FROM likes WHERE user_id = ? AND username = ? AND liked_comment_id = ?", User_id, User, Liked_comment_id)
		if err != nil {
			handlers.ChooseError(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}
