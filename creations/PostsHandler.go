package creation

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"text/template"

	data "main/dataBase"
	"main/handler"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/CreatePost.html")
	if err != nil {
		http.Error(w, "Internal server error in Create page", http.StatusInternalServerError)
		return
	}
	// if handler.IsJavaScriptDisabled(r) {
	// 	fmt.Println("Error in javascript in post creation")
	// 	http.Redirect(w, r, "/", http.StatusSeeOther)
	// 	return
	// }
	// fmt.Println("PATH:", r.URL.Path)
	post_id := r.URL.Query().Get("postid")
	username := r.URL.Query().Get("user")

	tmpl.Execute(w, struct {
		Post_id  string
		Username string
	}{
		Post_id:  post_id,
		Username: username,
	})
}

func InsertPost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/InsertPost" {
		http.Error(w, "page not found", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Not Allowed Method", http.StatusMethodNotAllowed)
		return
	}

	CurrentUser := r.URL.Query().Get("user")
	Post_id, _ := strconv.Atoi(r.URL.Query().Get("postid"))
	fmt.Println(Post_id)

	title := r.FormValue("title")
	body := r.FormValue("body")
	// adding image
	var imageData []byte
	var ImageErr error
	image, _, err6 := r.FormFile("image")
	if err6 != nil {
		fmt.Println("Error getting image:", err6)
		imageData = nil
		// return
	} else {
		// fmt.Println("Image:", image)
		defer image.Close()
		imageData, ImageErr = io.ReadAll(image)
		if ImageErr != nil {
			fmt.Println("Error reading image:", ImageErr)
			imageData = nil
			return
		}
	}
	// end here
	categories := r.Form["categories"]
	if len(categories) == 0 {
		categories = append(categories, "All")
	}
	if title == "" || body == "" {
		handler.ChooseError(w, "Bad Request", 400)
		// http.Error(w, "bad request empty post", http.StatusBadRequest)
		return
	}
	row := data.Db.QueryRow("SELECT id FROM users WHERE username = ?", CurrentUser)
	var id int
	err := row.Scan(&id)
	if err != nil {
		fmt.Println("we can't retrive data")
		handler.ChooseError(w, "Inrternal Server Error", 500)
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Add The post to the posts table
	if imageData == nil {
		_, err = data.Db.Exec("INSERT INTO posts(post_creator, title, body, user_id) VALUES (?, ?, ?, ?)", CurrentUser, title, body, id)
		if err != nil {
			log.Println("Error inserting user:", err)
			http.Error(w, "Internal server error", 500)
			return
		}
	} else {
		_, err = data.Db.Exec("INSERT INTO posts(post_creator, title, body,image, user_id) VALUES (?, ?, ?, ?,?)", CurrentUser, title, body, imageData, id)
		if err != nil {
			log.Println("Error inserting user:", err)
			http.Error(w, "Internal server error", 500)
			return
		}
	}
	err = data.Db.QueryRow("SELECT id FROM posts WHERE post_creator = ? AND title = ? AND body = ? AND user_id = ?", CurrentUser, title, body, id).Scan(&Post_id)
	if err != nil {
		handler.ChooseError(w, "Internal Server Error", 500)
		return
	}
	for _, categorie := range categories {
		_, err = data.Db.Exec("INSERT INTO categories(post_id, categorie) VALUES (?, ?)", Post_id, categorie)
		if err != nil {
			log.Println("Error inserting user:", err)
			handler.ChooseError(w, "Internal Server Error", 500)
			// http.Error(w, "Internal server error", 500)
			return
		}
	}
	http.Redirect(w, r, "/forum", http.StatusSeeOther)
}
