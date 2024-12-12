package handlers

import (
	"database/sql"
	"errors"
	"math"
	"net/http"
	"strconv"
	"text/template"

	"go.mod/dataBase"
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
	CommentsLength    int
}

type Reactions struct {
	LikeCount    int
	DislikeCount int
}

var postt Post

func Forum(w http.ResponseWriter, r *http.Request) {
	var page int
	if r.URL.Query().Get("page") == "" {
		page = 1
	} else {
		var err error
		page, err = strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			ChooseError(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	var offset int
	var DBlength int
	err := dataBase.Db.QueryRow("SELECT COUNT(*) FROM posts").Scan(&DBlength)
	if err != nil {
		ChooseError(w, "!Internal Server Error", http.StatusInternalServerError)
		return
	}
	if page < 1 || page > int(math.Ceil(float64(DBlength/10)+1)) {
		if page < 1 {
			page = 1
		} else {
			page = int(math.Ceil(float64(DBlength/10) + 1))
		}
		offset = DBlength - (DBlength - (10 * (page - 1)))
	} else {
		offset = DBlength - (DBlength - (10 * (page - 1)))
	}

	tmpl, err := template.ParseFiles("templates/forum.html")
	if err != nil {
		ChooseError(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var CurrentUser, CurrentSession string
	var session_id string
	cat_to_filter := r.FormValue("categories")

	cookie1, err := r.Cookie("session_token")
	if err != nil {
		ChooseError(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if cookie1.Value == "guest" {
		CurrentUser = "guest"
		CurrentSession = "0"
	} else {
		CurrentSession = cookie1.Value
		err = dataBase.Db.QueryRow("SELECT user_id, session_id FROM sessions WHERE session_id = ?", CurrentSession).Scan(&postt.CurrentUser_id, &session_id)
		if err != nil {
			ChooseError(w, "You don't have the right to be here sign in", http.StatusUnauthorized)
			return
		}
		err = dataBase.Db.QueryRow("SELECT username from users where id = ?", postt.CurrentUser_id).Scan(&CurrentUser)
		if err != nil {
			ChooseError(w, "You don't have the right to be here sign in", http.StatusUnauthorized)
			return
		}
	}

	posts_toshow, comment_id, post_id, err := GetPosts(cat_to_filter, tmpl, w, CurrentUser, offset)
	if err != nil {
		if err.Error() == "Bad Request" {
			ChooseError(w, err.Error(), http.StatusBadRequest)
		} else {
			ChooseError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err = tmpl.Execute(w, struct {
		Currenuser string
		Curr_id    int
		comment_id int
		Post_id    int
		PageIndex  int
		DataLength int
		Posts      []Post
	}{
		Currenuser: CurrentUser,
		Curr_id:    postt.CurrentUser_id,
		comment_id: comment_id,
		Post_id:    post_id,
		PageIndex:  page,
		DataLength: int(math.Ceil(float64(DBlength) / 10)),
		Posts:      posts_toshow,
	})
	if err != nil {
		ChooseError(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func GetPosts(cat_to_filter string, tmpl *template.Template, w http.ResponseWriter, CurrentUser string, offset int) ([]Post, int, int, error) {
	var post_rows *sql.Rows
	var err error
	if cat_to_filter != "all" && cat_to_filter != "" {
		if cat_to_filter == "myposts" {
			post_rows, err = dataBase.Db.Query(`SELECT * FROM posts WHERE post_creator = ? ORDER BY id DESC LIMIT 10 OFFSET ?`, CurrentUser, offset)
		} else if cat_to_filter == "likedposts" {
			post_rows, err = dataBase.Db.Query(`
				SELECT p.* FROM posts p
				JOIN likes li ON p.id = li.post_id 
				WHERE li.username = ? ORDER BY id DESC LIMIT 10 OFFSET ?`, CurrentUser, offset)
		} else {
			post_rows, err = dataBase.Db.Query(`
				SELECT p.* FROM posts p
				JOIN categories c ON p.id = c.post_id
				WHERE c.categorie = ? ORDER BY id DESC LIMIT 10 OFFSET ?`, cat_to_filter, offset)
		}
	} else {
		post_rows, err = dataBase.Db.Query("SELECT * FROM posts ORDER BY id DESC LIMIT 10 OFFSET ?", offset)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, 0, 0, errors.New("no Feild in dataBase base")
		} else {
			return nil, 0, 0, errors.New("internal server error you deleted a table from the dataBase base shut down the server and restar it again")
		}
	}
	defer post_rows.Close()

	var posts_toshow []Post
	var comment_id, post_id int
	for post_rows.Next() {
		var id, user_id int
		var title, body, usernamepublished string
		var time any
		if err := post_rows.Scan(&id, &user_id, &usernamepublished, &title, &body, &time); err != nil {
			return nil, 0, 0, errors.New("internal server error (we can't get the posts)")
		}

		post_id = id
		var likee Reactions
		var dislikee Reactions
		err = dataBase.Db.QueryRow(`SELECT COUNT(*) FROM likes WHERE post_id = ?`, post_id).Scan(&likee.LikeCount)
		if err != nil {
			return nil, 0, 0, errors.New("internal server error")
		}
		err = dataBase.Db.QueryRow(`SELECT COUNT(*) FROM dislikes WHERE post_id = ?`, post_id).Scan(&dislikee.DislikeCount)
		if err != nil {
			return nil, 0, 0, errors.New("internal server error")
		}

		var Length int
		err8 := dataBase.Db.QueryRow("SELECT COUNT(*) FROM comments WHERE post_commented_id = ?", id).Scan(&Length)
		if err8 != nil {
			return nil, 0, 0, errors.New("internal server error (we can't get the number of comments that we have)")
		}
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

			CommentsLength: Length,
		})
	}
	if err := post_rows.Err(); err != nil {
		return nil, 0, 0, errors.New("error during iteration on each row in the dataBasebase")
	}
	return posts_toshow, comment_id, post_id, nil
}
