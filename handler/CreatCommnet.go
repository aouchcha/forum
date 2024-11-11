package handler

import (
	"fmt"
	"log"
	data "main/dataBase"
	"net/http"
	"strconv"
)

func CreatCommnet(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.FormValue("comments"))
	fmt.Println(r.URL.Query().Get("postid"))
	fmt.Println(r.URL.Query().Get("writer"))

	comment_body := r.FormValue("comments")
	comment_writer := r.URL.Query().Get("writer")

	post_id, err := strconv.Atoi(r.URL.Query().Get("postid"))
	if err != nil {
		http.Error(w, "Internal server Error", http.StatusInternalServerError)
		return
	}
	// err = data.Db.QueryRow("SELECT post_creator FROM users WHERE id=?", postid).Scan()
	if comment_body == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	_, err = data.Db.Exec("INSERT INTO comments (comment_body, comment_writer, post_commented_id) VALUES (?, ?, ?)", comment_body, comment_writer, post_id)
	if err != nil {
		log.Println("Error inserting user:", err)
		http.Error(w, "Internal server error", 500)
		return
	}
	http.Redirect(w, r, "/forum?user="+comment_writer, http.StatusSeeOther)
}
