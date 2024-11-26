package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	// "strconv"

	data "main/dataBase"
)

// var post []Post

func Like(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/Likes" {
		user := r.URL.Query().Get("user")
		postid := r.URL.Query().Get("Liked_Post_id")
		comment_id_str := r.URL.Query().Get("comment_id")
		fmt.Println(user, postid, comment_id_str, "TTTTTTTTTTT")
		var user_id int
		err := data.Db.QueryRow("SELECT id FROM users WHERE username=?", user).Scan(&user_id)
		if err != nil {
			http.Error(w, "Inernal server error first err in like func", http.StatusInternalServerError)
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
		// if comment_id_str != "" {
		// 	liked_comment_id, err := strconv.Atoi(comment_id_str)
		// 	if err != nil {
		// 		fmt.Println("Error converting comment_id to int")
		// 		http.Error(w, "Internal server Error", http.StatusInternalServerError)
		// 		return
		// 	}
		// 	var exists bool
		// 	err = data.Db.QueryRow("SELECT id FROM likes WHERE liked_comment_id = ? ", liked_comment_id).Scan(&exists)
		// 	if err == sql.ErrNoRows {
		// 		// Delete from dislikes if exists
		// 		err = data.Db.QueryRow("SELECT id FROM dislikes WHERE user_id = ? AND disliked_comment_id = ? AND username = ?", user_id, liked_comment_id, user).Scan(&exists)
		// 		if err != sql.ErrNoRows {
		// 			_, err = data.Db.Exec("DELETE FROM dislikes WHERE user_id = ? AND disliked_comment_id = ? AND username = ?", user_id, liked_comment_id, user)
		// 			if err != nil {
		// 				http.Error(w, "Error unliking post", http.StatusInternalServerError)
		// 				return
		// 			}
		// 		}
		// 		// Then Insert into likes if not exists
		// 		_, err = data.Db.Exec("INSERT INTO likes (user_id, liked_comment_id, username) VALUES (?, ?, ?)", user_id, liked_comment_id, user)
		// 		if err != nil {
		// 			fmt.Println("Error liking post", err)
		// 			http.Error(w, "Error liking post", http.StatusInternalServerError)
		// 			return
		// 		}
		// 	} else {
		// 		// Delete from likes if exists
		// 		_, err = data.Db.Exec("DELETE FROM likes WHERE liked_comment_id = ? ", liked_comment_id)
		// 		if err != nil {
		// 			http.Error(w, "Error unliking post", http.StatusInternalServerError)
		// 			return
		// 		}
		// 	}
		// }
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
