package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"go.mod/dataBase"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/create_post" {
		ChooseError(w, "Page Not Found", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		ChooseError(w, "Mehtod Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles("templates/create_post.html")
	if err != nil {
		ChooseError(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	post_id := r.URL.Query().Get("postid")
	username := r.URL.Query().Get("user")

	err = tmpl.Execute(w, struct {
		Post_id  string
		Username string
	}{
		Post_id:  post_id,
		Username: username,
	})
	if err != nil {
		ChooseError(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func InsertPost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/InsertPost" {
		ChooseError(w, "Page Not Found", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodPost {
		ChooseError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	CurrentUser := r.URL.Query().Get("user")
	Post_id, _ := strconv.Atoi(r.URL.Query().Get("postid"))

	title := strings.TrimLeft(r.FormValue("title"), " ")
	body := strings.TrimLeft(r.FormValue("body"), " ")

	categories := r.Form["categories"]
	if len(categories) == 0 {
		categories = append(categories, "All")
	}

	Doubled, Check := CheckCategories(categories, "")
	if Doubled || !Check {
		ChooseError(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if title == "" || body == "" || ContentLength(title) > 500 || ContentLength(body) > 1000 {
		ChooseError(w, "you insered an empty field or more chars than the max", 400)
		return
	}
	row := dataBase.Db.QueryRow("SELECT id FROM users WHERE username = ?", CurrentUser)
	var id int
	err := row.Scan(&id)
	if err != nil {
		ChooseError(w, "You have chnaged the value of the query and this user didn't exist", http.StatusBadRequest)
		return
	}

	_, err = dataBase.Db.Exec("INSERT INTO posts(post_creator, title, body, user_id) VALUES (?, ?, ?, ?)", CurrentUser, title, body, id)
	if err != nil {
		ChooseError(w, "Inrternal Server Error", 500)
		return
	}

	err = dataBase.Db.QueryRow("SELECT id FROM posts WHERE post_creator = ? AND title = ? AND body = ? AND user_id = ?", CurrentUser, title, body, id).Scan(&Post_id)
	if err != nil {
		ChooseError(w, "Internal Server Error", 500)
		return
	}
	for _, categorie := range categories {
		_, err = dataBase.Db.Exec("INSERT INTO categories(post_id, categorie) VALUES (?, ?)", Post_id, categorie)
		if err != nil {

			ChooseError(w, "Internal Server Error", 500)
			return
		}
	}
	http.Redirect(w, r, "/forum", http.StatusSeeOther)
}

func CheckCategories(categories []string, categorie string) (bool, bool) {
	Doubled := false
	OurCat := map[string]bool{"it": true, "economie": true, "enteairtement": true, "politic": true, "sport": true, "All": true}
	Check := true
	if categories != nil {
		for i := 0; i < len(categories); i++ {
			for j := i + 1; j < len(categories); j++ {
				if categories[i] == categories[j] {
					Doubled = true
					break
				}
			}
			if Doubled {
				break
			}
		}

		for _, cat := range categories {
			if !OurCat[cat] {
				Check = false
				break
			}
		}
	} else {
		if !OurCat[categorie] {
			Check = false
		}
	}
	return Doubled, Check
}
