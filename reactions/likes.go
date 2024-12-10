package reactions

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	data "main/dataBase"
	handler "main/handler"
)

func PostsLike(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/PostsLikes" {
		if handler.IsJavaScriptDisabled(r) {
			fmt.Println("5555555")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		user := r.URL.Query().Get("user")
		postid := r.URL.Query().Get("Liked_Post_id")
		var user_id int
		err := data.Db.QueryRow("SELECT id FROM users WHERE username=?", user).Scan(&user_id)
		if err != nil {
			http.Error(w, "Inernal server error user not found in data base", http.StatusInternalServerError)
			return
		}
		var exists bool
		err = data.Db.QueryRow("SELECT id FROM likes WHERE user_id = ? AND post_id = ? AND username = ?", user_id, postid, user).Scan(&exists)
		if err == sql.ErrNoRows {
			// Delete from dislikes if exists
			err = data.Db.QueryRow("SELECT id FROM dislikes WHERE user_id = ? AND post_id = ? AND username = ?", user_id, postid, user).Scan(&exists)
			if err != sql.ErrNoRows {
				_, err = data.Db.Exec("DELETE FROM dislikes WHERE user_id = ? AND post_id = ? AND username = ?", user_id, postid, user)
				if err != nil {
					http.Error(w, "Error unliking post", http.StatusInternalServerError)
					return
				}
			}
			// Then Insert into likes if not exists
			_, err = data.Db.Exec("INSERT INTO likes (user_id, post_id, username) VALUES (?, ?, ?)", user_id, postid, user)
			if err != nil {
				fmt.Println("Hna error")
				fmt.Println("Error liking post", err)
				http.Error(w, "Error liking post", http.StatusInternalServerError)
				return
			}
		} else {
			// Delete from likes if exists
			_, err = data.Db.Exec("DELETE FROM likes WHERE user_id = ? AND post_id = ? AND username = ?", user_id, postid, user)
			if err != nil {
				http.Error(w, "Error unliking post", http.StatusInternalServerError)
				return
			}
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.Redirect(w, r, "/forum", http.StatusSeeOther)
}

func CommentsLike(w http.ResponseWriter, r *http.Request) {
	if handler.IsJavaScriptDisabled(r) {
		fmt.Println("5555555")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	Liked_comment_id, err := strconv.Atoi(r.URL.Query().Get("comment_id"))
	if err != nil {
		http.Error(w, "Internal server error in converting the comment id into int", http.StatusInternalServerError)
		return
	}
	// Post_id, err := strconv.Atoi(r.URL.Query().Get("post_id"))
	// if err =! nil {
	// 	http.Error(w, "Internal server error in converting the post id into int", http.StatusInternalServerError)
	// 	return
	// }
	User := r.URL.Query().Get("user")
	// fmt.Println(Liked_comment_id, Post_id, User)
	var User_id int
	//Check if the user is in our data base
	err = data.Db.QueryRow("SELECT id FROM users WHERE username = ?", User).Scan(&User_id)
	// fmt.Println("USER ID", User_id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "We can't find this user in our data base", http.StatusInternalServerError)
			return
		} else {
			http.Error(w, "We can't get the data", http.StatusInternalServerError)
			return
		}
	}
	var Exist bool
	//Check if the user allrady like the comment
	err = data.Db.QueryRow("SELECT id FROM likes WHERE user_id = ? AND username = ? AND liked_comment_id = ?", User_id, User, Liked_comment_id).Scan(&Exist)
	if err == sql.ErrNoRows {
		// The user didn't like the comment let's check if he already dislike the comment
		err = data.Db.QueryRow("SELECT id FROM dislikes WHERE user_id = ? AND username = ? AND disliked_comment_id = ?", User_id, User, Liked_comment_id).Scan(&Exist)
		if err != sql.ErrNoRows {
			//if he dislike the post we remove the dislike
			_, err = data.Db.Exec("DELETE FROM dislikes WHERE user_id = ? AND username = ? AND disliked_comment_id = ?", User_id, User, Liked_comment_id)
			if err != nil {
				http.Error(w, "Error unliking post", http.StatusInternalServerError)
				return
			}
		}
		//We add the like to our table
		_, err = data.Db.Exec("INSERT INTO likes (user_id, username, liked_comment_id) VALUES (?,?,?)", User_id, User, Liked_comment_id)
		if err != nil {
			http.Error(w, "Error liking post"+err.Error(), http.StatusInternalServerError)
			return
		}

	} else {
		//The user already like the comment now we remove the like from the table
		_, err = data.Db.Exec("DELETE FROM likes WHERE user_id = ? AND username = ? AND liked_comment_id = ?", User_id, User, Liked_comment_id)
		if err != nil {
			http.Error(w, "Error unliking post", http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, "/showcomments", http.StatusSeeOther)

}
