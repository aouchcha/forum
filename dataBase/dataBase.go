package dataBase

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

func init() {
	var err error
	Db, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		fmt.Println("Error in open database")
		os.Exit(1)
		return
	}

	createUsersTable := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT UNIQUE NOT NULL,
        email TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

	_, err = Db.Exec(createUsersTable)
	if err != nil {
		fmt.Println("Error creating users table")
		os.Exit(1)
	}

	createSessionsTable := `
	CREATE TABLE IF NOT EXISTS sessions (
		user_id INTEGER NOT NULL PRIMARY KEY,
		session_id TEXT,
		expires_at DATETIME NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`

	_, err = Db.Exec(createSessionsTable)
	if err != nil {
		fmt.Println("Error creating sessions table")
		os.Exit(1)
	}

	const CreatPosts string = `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		post_creator TEXT NOT NULL,
		title TEXT NOT NULL,
		body TEXT NOT NULL,
		time DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`

	if _, err := Db.Exec(CreatPosts); err != nil {
		fmt.Println("Error creating posts table")
		os.Exit(1)
	}

	const CreatCategories string = `
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER NOT NULL,
		categorie TEXT NOT NULL,
		FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
	)`

	if _, err := Db.Exec(CreatCategories); err != nil {
		fmt.Println("Error creating categories table")
		os.Exit(1)
	}

	const CommentsTable string = `
	CREATE TABLE IF NOT EXISTS comments (
		comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
		comment_body TEXT NOT NULL,
		comment_writer TEXT NOT NULL,
		comment_writer_id INTEGER NOT NULL,
		post_commented_id INTEGER,
		time DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (comment_writer_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (post_commented_id) REFERENCES posts(id) ON DELETE CASCADE
	)`

	if _, err := Db.Exec(CommentsTable); err != nil {
		fmt.Println("Error creating comments table")
		os.Exit(1)
	}
	const LikesTable string = `
	CREATE TABLE IF NOT EXISTS likes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,             
    post_id INTEGER ,
	liked_comment_id INTEGER ,             
    username TEXT NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (liked_comment_id) REFERENCES comments(comment_id) ON DELETE CASCADE, 
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE  
);`
	if _, err := Db.Exec(LikesTable); err != nil {
		fmt.Println("Error creating likes table")
		os.Exit(1)
	}
	const DislikeTable string = `
	CREATE TABLE IF NOT EXISTS dislikes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER ,             
    post_id INTEGER ,
	disliked_comment_id INTEGER ,             
    username TEXT ,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (disliked_comment_id) REFERENCES comments(comment_id) ON DELETE CASCADE, 
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE  
);`
	if _, err := Db.Exec(DislikeTable); err != nil {
		fmt.Println("Error creating dislikes table")
		os.Exit(1)
	}
}
