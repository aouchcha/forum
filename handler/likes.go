package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	data "main/dataBase"
)

func Like(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Query().Get("Liked_Post_id"))
	fmt.Println(r.URL.Query().Get("user"))
	user := r.URL.Query().Get("user")
	liked_post_id, err := strconv.Atoi(r.URL.Query().Get("Liked_Post_id"))
	if err != nil {
		http.Error(w, "Internal server Error", http.StatusInternalServerError)
		return
	}
	_, err = data.Db.Exec("INSERT INTO likes (liked_post_id, user_name_like) VALUES (?, ?)", liked_post_id, user)
	if err != nil {
		log.Println("Error inserting user:", err)
		http.Error(w, "Internal server error", 500)
		return
	}
	http.Redirect(w, r, "/forum?user="+user, http.StatusSeeOther)
}
