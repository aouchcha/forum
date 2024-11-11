package handler

import (
	"database/sql"
	"fmt"
	"log"
	data "main/dataBase"
	"net/http"
)

func CreatPost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/create_post" {
		http.Error(w, "page not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Not Allowed Method", http.StatusMethodNotAllowed)
		return
	}
	CurrentUser := r.URL.Query().Get("user")
	title := r.FormValue("title")
	body := r.FormValue("body")
	categorie := r.FormValue("categories")
	if title == "" || body == "" {
		http.Error(w, "bad request empty post", http.StatusBadRequest)
		return
	}
	row := data.Db.QueryRow("SELECT username FROM users WHERE username = ?", CurrentUser)
	var username string
	err := row.Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("we this user don't exist")
			http.Error(w, "you are in the guest session", http.StatusInternalServerError)
			return
		} else {
			fmt.Println("we can't retrive data")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
	// fmt.Println(username)
	_, err = data.Db.Exec("INSERT INTO posts(post_creator, title, body, categorie) VALUES (?, ?, ?, ?)", CurrentUser, title, body, categorie)
	if err != nil {
		log.Println("Error inserting user:", err)
		http.Error(w, "Internal server error", 500)
		return
	}
	// fmt.Println(title, body, time)
	http.Redirect(w, r, "/forum?user="+CurrentUser, http.StatusSeeOther)
}
