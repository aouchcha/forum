package reactions

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"go.mod/dataBase"
	"go.mod/helpers"
)

func LikesCounterWithApi(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/likes" {

		cookie, err := r.Cookie("session_token")
		if err != nil || cookie.Value == "guest" {
			ResponseReaction(w, 401, 0, 0)
		}

		var LikeCount, DislikeCount int
		var post_id, comment_id string
		if r.URL.Query().Get("postid") != "" {
			post_id = helpers.Unhash(r.URL.Query().Get("postid"))
		} else {
			comment_id = helpers.Unhash(r.URL.Query().Get("comment_id"))
		}
		var Check int

		if post_id == "" && comment_id != "" {
			err = dataBase.Db.QueryRow("SELECT COUNT(*) FROM comments WHERE comment_id = ?", comment_id).Scan(&Check)
		} else if comment_id == "" && post_id != "" {
			err = dataBase.Db.QueryRow("SELECT COUNT(*) FROM posts WHERE id = ?", post_id).Scan(&Check)
		} else {
			ResponseReaction(w, http.StatusBadRequest, 0, 0)
			return
		}
		if err != nil || Check == 0 {
			ResponseReaction(w, http.StatusNotFound, 0, 0)
			return
		}
		LikeCount, DislikeCount, err = getLikeAndDislikeCount(post_id, comment_id)
		if err != nil {
			ResponseReaction(w, 400, 0, 0)
			return
		}
		ResponseReaction(w, 200, LikeCount, DislikeCount)
	}
}

func ResponseReaction(w http.ResponseWriter, ErrCode, LikeCount, DislikeCount int) {
	response, err := json.Marshal(struct {
		LikeCount    string
		DislikeCount string
	}{
		LikeCount:    strconv.Itoa(LikeCount),
		DislikeCount: strconv.Itoa(DislikeCount),
	})
	if err != nil {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"check": false, "message": "Internal Server Error"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(ErrCode)
	w.Write(response)
}

func getLikeAndDislikeCount(post_id, comment_id string) (int, int, error) {
	var LikeCount, DislikeCount int
	if post_id != "" && comment_id == "" {
		err := dataBase.Db.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = ?", post_id).Scan(&LikeCount)
		if err != nil {
			return 0, 0, errors.New("Internal Server Error")
		}
		err = dataBase.Db.QueryRow("SELECT COUNT(*) FROM dislikes WHERE post_id = ?", post_id).Scan(&DislikeCount)
		if err != nil {
			return 0, 0, errors.New("Internal Server Error")
		}
	} else if post_id == "" && comment_id != "" {

		err := dataBase.Db.QueryRow("SELECT COUNT(*) FROM likes WHERE liked_comment_id = ?", comment_id).Scan(&LikeCount)
		if err != nil {
			return 0, 0, errors.New("Internal Server Error")
		}
		err = dataBase.Db.QueryRow("SELECT COUNT(*) FROM dislikes WHERE disliked_comment_id = ?", comment_id).Scan(&DislikeCount)
		if err != nil {
			return 0, 0, errors.New("Internal Server Error")
		}
	} else {
		return 0, 0, errors.New("bad dataBase from the front end")
	}
	return LikeCount, DislikeCount, nil
}
