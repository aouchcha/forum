package reactions

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	data "main/dataBase"
	handler "main/handler"
)

func PostsDislikes(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/PostsDislikes" {
		if handler.IsJavaScriptDisabled(r) {
			fmt.Println("5555555")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		user := r.URL.Query().Get("user")
		postid := r.URL.Query().Get("Disliked_Post_id")
		user_id, err := strconv.Atoi(r.URL.Query().Get("user_id"))
		if err != nil {
			http.Error(w, "Internal server Error", http.StatusInternalServerError)
			return
		}
		var exists bool
		err = data.Db.QueryRow("SELECT id FROM dislikes WHERE user_id = ? AND post_id = ? AND username = ?", user_id, postid, user).Scan(&exists)
		if err == sql.ErrNoRows {
			// Deleete from likes if exists first
			err = data.Db.QueryRow("SELECT id FROM likes WHERE user_id = ? AND post_id = ? AND username = ?", user_id, postid, user).Scan(&exists)
			if err != sql.ErrNoRows {
				_, err = data.Db.Exec("DELETE FROM likes WHERE user_id = ? AND post_id = ? AND username = ?", user_id, postid, user)
				if err != nil {
					http.Error(w, "Error unliking post", http.StatusInternalServerError)
					return
				}
			}
			// Then Insert into dislikes
			_, err = data.Db.Exec("INSERT INTO dislikes (user_id, post_id, username) VALUES (?, ?, ?)", user_id, postid, user)
			if err != nil {
				fmt.Println("Error liking post", err)
				http.Error(w, "Error liking post", http.StatusInternalServerError)
				return
			}

		} else {
			// Delete from dislikes if exists
			_, err = data.Db.Exec("DELETE FROM dislikes WHERE user_id = ? AND post_id = ? AND username = ?", user_id, postid, user)
			if err != nil {
				http.Error(w, "Error Disliked post", http.StatusInternalServerError)
				return
			}
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.Redirect(w, r, "forum", http.StatusSeeOther)
}

func CommentsDislike(w http.ResponseWriter, r *http.Request) {
	if handler.IsJavaScriptDisabled(r) {
		fmt.Println("5555555")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	Post_id, err := strconv.Atoi(r.URL.Query().Get("post_id"))
	if err != nil {
		http.Error(w, "Internal server error in converting the post id into int", http.StatusInternalServerError)
		return
	}
	User := r.URL.Query().Get("user")
	Disliked_comment_id, err := strconv.Atoi(r.URL.Query().Get("comment_id"))

	if err != nil {
		http.Error(w, "Internal server error in converting the comment id into int", http.StatusInternalServerError)
		return
	}
	fmt.Println(Post_id, User, Disliked_comment_id)
	var User_id int
	err = data.Db.QueryRow("SELECT id FROM users WHERE username = ?", User).Scan(&User_id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "We can't find this user in the data base", http.StatusInternalServerError)
			return
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
	var Exist bool
	//Check if The user already disliked the comment
	err = data.Db.QueryRow("SELECT id FROM dislikes WHERE user_id = ? AND username = ? AND disliked_comment_id = ?", User_id, User, Disliked_comment_id).Scan(&Exist)

	//if the user didn't disliked the comment
	if err == sql.ErrNoRows {
		//let's check if the user already like the commment
		err = data.Db.QueryRow("SELECT id FROM likes WHERE user_id = ? AND username = ? AND liked_comment_id = ?", User_id, User, Disliked_comment_id).Scan(&Exist)
		if err != sql.ErrNoRows {
			_, err = data.Db.Exec("DELETE FROM likes WHERE user_id = ? AND username = ? AND liked_comment_id = ?", User_id, User, Disliked_comment_id)
			if err != nil {
				http.Error(w, "Internal server Error", http.StatusInternalServerError)
				return
			}
			//if the user already liked the comment we should delete it
		}
		//the user didn't yet do any reaction to the comment so we add the dislike
		_, err = data.Db.Exec("INSERT INTO dislikes (user_id, disliked_comment_id, username) VALUES (?,?,?)", User_id, Disliked_comment_id, User)
		if err != nil {
			http.Error(w, "Error disliking post", http.StatusInternalServerError)
			return
		}

	} else {
		_, err = data.Db.Exec("DELETE FROM dislikes WHERE user_id = ? AND username = ? AND disliked_comment_id = ?", User_id, User, Disliked_comment_id)
		if err != nil {
			http.Error(w, "Internal server Error", http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, "forum", http.StatusSeeOther)

}
