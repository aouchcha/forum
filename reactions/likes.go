package reactions

import (
	"database/sql"
	"net/http"
	"strconv"

	"go.mod/dataBase"
	"go.mod/handlers"
	"go.mod/helpers"
)

func PostsLike(w http.ResponseWriter, r *http.Request) {
	//////////////////////////////////////Check if the request is good////////////////////////////////////////////////////////////
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
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	//////////////////////////////////////////// Get the userid and the username values from the session /////////////////////////////////////////////////////////
	cookie, _ := r.Cookie("session_token")
	var user string
	var user_id int
	dataBase.Db.QueryRow("SELECT id, username FROM users INNER JOIN sessions ON users.id = sessions.user_id WHERE sessions.session_id = ?", cookie.Value).Scan(&user_id, &user)
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	/////////////////////////////////////////////////// Get The Post id from the url query ////////////////////////////////////////////////////////////////
	temp := helpers.Unhash(r.URL.Query().Get("Liked_Post_id"))
	postid, err := strconv.Atoi(temp)
	if err != nil {
		handlers.ChooseError(w, "Bad Request You Chnaged In The Url Query", http.StatusBadRequest)
		return
	}
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	//////////////////////////////// check if the post id is valid (The user didn't cange it from inspect)////////////////////////////////////////////////
	var Check int
	err = dataBase.Db.QueryRow("SELECT COUNT(*) FROM posts WHERE id = ?", postid).Scan(&Check)
	if err != nil || Check == 0 {
		return
	}
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	//////////////////////////////////// add the reaction to the DB ////////////////////////////////////////////////////////////////////
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
	/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
}

func CommentsLike(w http.ResponseWriter, r *http.Request) {
	//////////////////////////////////////Check if the request is good////////////////////////////////////////////////////////////
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
	/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	///////////////////////////////////////// Get The comment id from query ////////////////////////////////////////////
	temp := helpers.Unhash(r.URL.Query().Get("comment_id"))
	Liked_comment_id, err := strconv.Atoi(temp)
	if err != nil {
		handlers.ChooseError(w, "Bad Request You Chnaged In The Url Query", http.StatusBadRequest)
		return
	}
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	////////////////////////////////// Check if the comment id is valid (already we have it in the data base in case the user change it from the inspect)///////////////////////////////////
	var Check int
	err = dataBase.Db.QueryRow("SELECT COUNT(*) FROM comments WHERE comment_id = ?", Liked_comment_id).Scan(&Check)
	if err != nil || Check == 0 {
		return
	}
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	//////////////////////////////////////////// Get the userid and the username values from the session /////////////////////////////////////////////////////////
	cookie, _ := r.Cookie("session_token")
	var User_id int
	var User string
	dataBase.Db.QueryRow("SELECT id, username FROM users INNER JOIN sessions ON users.id = sessions.user_id WHERE sessions.session_id = ?", cookie.Value).Scan(&User_id, &User)
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	//////////////////////////////////// add the reaction to the DB ////////////////////////////////////////////////////////////////////
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
	/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
}
