package handler

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"net/http"

	data "main/dataBase"
	// handler "main/handler"
)

type Post struct {
	LikesCounter      int
	DislikeCounter    int
	Postid            int
	Usernamepublished string
	CurrentUsser      string
	CurrentUser_id    int
	Title             string
	Body              string
	Time              any
	Categorie         string
	Image             string
}

type Reactions struct {
	LikeCount    int
	DislikeCount int
}

var postt Post

func Forum(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/forum" {
        http.Error(w, "page not found", http.StatusNotFound)
        return
    }
    tmpl, err := template.ParseFiles("templates/forum.html")
    if err != nil {
        http.Error(w, "Internal Server Error with forum html page", http.StatusInternalServerError)
        return
    }

    var CurrentUser, CurrentSession string
    var session_id string
    cat_to_filter := r.FormValue("categories")
    cookie1, err:= r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Internal Server Error with forum html page", http.StatusInternalServerError)
		return
	}
    // isGuest := false
    if cookie1.Value == "guest" {
        CurrentUser = "guest"
        CurrentSession = "0"
        // isGuest = true
    } else {
        CurrentSession = cookie1.Value
        err = data.Db.QueryRow("SELECT user_id, session_id FROM sessions WHERE session_id = ?", CurrentSession).Scan(&postt.CurrentUser_id, &session_id)
        if err != nil {
            http.Redirect(w, r, "/login", http.StatusSeeOther)
            return
        }
        err = data.Db.QueryRow("SELECT username from users where id = ?", postt.CurrentUser_id).Scan(&CurrentUser)
        if err != nil {
            http.Redirect(w, r, "/login", http.StatusSeeOther)
            return
        }
    }

    posts_toshow, comment_id, post_id, _ := GetPosts(cat_to_filter, tmpl, w, CurrentUser)

    err = tmpl.Execute(w, struct {
        Currenuser string
        Curr_id    int
        comment_id int
        Post_id    int
        Posts      []Post
        IsGuest    bool 
    }{
        Currenuser: CurrentUser,
        Curr_id:    postt.CurrentUser_id,
        comment_id: comment_id,
        Post_id:    post_id,
        Posts:      posts_toshow,
        // IsGuest:    isGuest, 
    })

    if err != nil {
        fmt.Println(err)
        http.Error(w, "Internal Server", http.StatusInternalServerError)
        return
    }
}


func GetPosts(cat_to_filter string, tmpl *template.Template, w http.ResponseWriter, CurrentUser string) ([]Post, int, int, error) {
	var post_rows *sql.Rows
	var err error
	if cat_to_filter != "all" && cat_to_filter != "" {

		if cat_to_filter == "myposts" {
			post_rows, err = data.Db.Query(`SELECT * FROM posts WHERE post_creator = ?`, CurrentUser)
		} else if cat_to_filter == "likedposts" {
			// fmt.Println("liked posts")
			post_rows, err = data.Db.Query(`
				SELECT posts.* FROM posts 
				JOIN likes ON posts.id = likes.post_id 
				WHERE likes.username = ?`, CurrentUser)
		} else {
			post_rows, err = data.Db.Query(`
				SELECT p.* FROM posts p
				JOIN categories c ON p.id = c.post_id
				WHERE c.categorie = ?`, cat_to_filter)
		}
	} else {
		post_rows, err = data.Db.Query("SELECT * FROM posts;")
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, 0, 0, errors.New("no Feild in data base")
		} else {
			return nil, 0, 0, errors.New("internal server error you deleted a table from the data base shut down the server and restar it again")
		}
	}
	defer post_rows.Close()
	// post_rows, err = data.Db.Query("SELECT * FROM posts WHERE id = ?;", cat_to_filter)
	var posts_toshow []Post
	var comment_id, post_id int
	for post_rows.Next() {
		var id, user_id int
		var title, body, usernamepublished string
		var imageData []byte
		var time any
		if err := post_rows.Scan(&id, &user_id, &usernamepublished, &title, &body, &imageData, &time); err != nil {
			http.Error(w, "Error fetching post data", http.StatusInternalServerError)
			// fmt.Println(err)
			continue
		}
		// fmt.Println(usernamepublished, title, body)
		post_id = id
		var likee Reactions
		var dislikee Reactions
		err = data.Db.QueryRow(`SELECT COUNT(*) FROM likes WHERE post_id = ?`, post_id).Scan(&likee.LikeCount)
		if err != nil {
			// fmt.Println("Error fetching like count ==>", err)
			http.Error(w, "Error fetching like count", http.StatusInternalServerError)
			return nil, 0, 0, errors.New("internal server error")
		}
		err = data.Db.QueryRow(`SELECT COUNT(*) FROM dislikes WHERE post_id = ?`, post_id).Scan(&dislikee.DislikeCount)
		if err != nil {
			// fmt.Println("Error fetching like count ==>", err)
			http.Error(w, "Error fetching like count", http.StatusInternalServerError)
			return nil, 0, 0, errors.New("internal server error")
		}
		base64Image := base64.StdEncoding.EncodeToString(imageData)
		// fmt.Println("comments id= ", comment_id, "post id= ", post_id)
		posts_toshow = append(posts_toshow, Post{
			Postid:            id,
			LikesCounter:      likee.LikeCount,
			DislikeCounter:    dislikee.DislikeCount,
			Usernamepublished: usernamepublished,
			CurrentUsser:      CurrentUser,
			CurrentUser_id:    postt.CurrentUser_id,
			Title:             title,
			Body:              body,
			Time:              time,
			Image:             base64Image,
		})
	}
	if err := post_rows.Err(); err != nil {
		return nil, 0, 0, errors.New("error ma3reftch fin")
	}
	for i := 0; i < len(posts_toshow)-1; i++ {
		for j := i + 1; j < len(posts_toshow); j++ {
			posts_toshow[i], posts_toshow[j] = posts_toshow[j], posts_toshow[i]
		}
	}
	return posts_toshow, comment_id, post_id, nil
}
