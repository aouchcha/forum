package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"text/template"
)

func main() {
	http.HandleFunc("/", Feed)
	http.HandleFunc("/like", Like)
}

type Post struct {
	ID      int
	Title   string
	Content string
	Likes   int
}

func Feed(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query(`
		SELECT posts.id, posts.title, posts.content, COUNT(likes.post_id) AS like_count
		FROM posts
		LEFT JOIN likes ON posts.id = likes.post_id
		GROUP BY posts.id
	`)
	if err != nil {
		http.Error(w, "Error fetching posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Likes)
		if err != nil {
			http.Error(w, "Error scanning posts", http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	tmpl, err := template.ParseFiles("templates/feed.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, map[string]interface{}{"Posts": posts})
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func Like(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	postID, err := strconv.Atoi(r.URL.Query().Get("post_id"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	userID := 1 // Hardcoded user ID for simplicity; replace with session-based ID.

	db, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow("SELECT 1 FROM likes WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&exists)
	if err == sql.ErrNoRows {
		_, err = db.Exec("INSERT INTO likes (user_id, post_id) VALUES (?, ?)", userID, postID)
		if err != nil {
			http.Error(w, "Error liking post", http.StatusInternalServerError)
			return
		}
	} else {
		_, err = db.Exec("DELETE FROM likes WHERE user_id = ? AND post_id = ?", userID, postID)
		if err != nil {
			http.Error(w, "Error unliking post", http.StatusInternalServerError)
			return
		}
	}

	var likeCount int
	err = db.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = ?", postID).Scan(&likeCount)
	if err != nil {
		http.Error(w, "Error fetching like count", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"likes": likeCount})
}
