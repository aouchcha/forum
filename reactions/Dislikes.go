package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	data "main/dataBase"
)

func Dislikes(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/Dislikes" {
		user := r.URL.Query().Get("user")
		postid := r.URL.Query().Get("Disliked_Post_id")
		user_id, err := strconv.Atoi(r.URL.Query().Get("user_id"))
		if err != nil {
			http.Error(w, "Internal server Error", http.StatusInternalServerError)
			return
		}
		Disliked_Comment_id_str := r.URL.Query().Get("comment_id")
		// if err != nil {
		// 	fmt.Println("Error converting liked_post_id to int")
		// 	http.Error(w, "Internal server Error", http.StatusInternalServerError)
		// 	return
		// }
		fmt.Println("Disliked_Post_id :", Disliked_Comment_id_str)
		fmt.Println("user_id :", user_id)
		if postid != "" {
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
		} else if Disliked_Comment_id_str != "" {
			Disliked_comment_id, err := strconv.Atoi(Disliked_Comment_id_str)
			if err != nil {
				fmt.Println("Error converting comment_id to int")
				http.Error(w, "Internal server Error", http.StatusInternalServerError)
				return
			}
			var exists bool
			err = data.Db.QueryRow("SELECT id FROM dislikes WHERE disliked_comment_id = ? ", Disliked_comment_id).Scan(&exists)
			if err == sql.ErrNoRows {
				// Delete from dislikes if exists
				err = data.Db.QueryRow("SELECT id FROM likes WHERE user_id = ? AND liked_comment_id = ? AND username = ?", user_id, Disliked_comment_id, user).Scan(&exists)
				if err != sql.ErrNoRows {
					_, err = data.Db.Exec("DELETE FROM likes WHERE user_id = ? AND liked_comment_id = ? AND username = ?", user_id, Disliked_comment_id, user)
					if err != nil {
						http.Error(w, "Error unliking post", http.StatusInternalServerError)
						return
					}
				}
				// Then Insert into likes if not exists
				_, err = data.Db.Exec("INSERT INTO dislikes (user_id, disliked_comment_id, username) VALUES (?, ?, ?)", user_id, Disliked_comment_id, user)
				if err != nil {
					fmt.Println("Error disliking post", err)
					http.Error(w, "Error liking post", http.StatusInternalServerError)
					return
				}
			} else {
				// Delete from likes if exists
				_, err = data.Db.Exec("DELETE FROM dislikes WHERE disliked_comment_id = ? ", Disliked_comment_id)
				if err != nil {
					http.Error(w, "Error undisliking post", http.StatusInternalServerError)
					return
				}
			}

		}
		http.Redirect(w, r, "/forum", http.StatusSeeOther)
		return
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
