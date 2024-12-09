package reactions

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	data "main/dataBase"
)

func LikesCounterWithApi(w http.ResponseWriter, r *http.Request) {
	fmt.Println("likes api")
	if r.URL.Path == "/api/likes" {

		cookie, err := r.Cookie("session_token")
		if err != nil || cookie.Value == "guest" {
			ResponseReaction(w, 401, 0, 0)
		}

		var LikeCount, DislikeCount int
		// var Check bool = true
		post_id := r.URL.Query().Get("postid")
		comment_id := r.URL.Query().Get("comment_id")
		fmt.Println(post_id, comment_id)
		// Fetch the Like and Dislike counts
		LikeCount, DislikeCount, err = getLikeAndDislikeCount(post_id, comment_id)
		if err != nil {
			ResponseReaction(w, 400, 0, 0)
			// http.Error(w, "Error fetching like and dislike count", http.StatusInternalServerError)
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
		// If there is an error in creating JSON, send an error response.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"check": false, "message": "Internal Server Error"}`))
		return
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(ErrCode)
	w.Write(response)
}

func getLikeAndDislikeCount(post_id, comment_id string) (int, int, error) {
	var LikeCount, DislikeCount int
	if post_id != "" && comment_id == "" {
		fmt.Println("ANA hna 1")
		err := data.Db.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = ?", post_id).Scan(&LikeCount)
		if err != nil {
			fmt.Println("Error fetching like count:", err)
			return 0, 0, err
		}
		err = data.Db.QueryRow("SELECT COUNT(*) FROM dislikes WHERE post_id = ?", post_id).Scan(&DislikeCount)
		if err != nil {
			fmt.Println("Error fetching dislike count:", err)
			return 0, 0, err
		}
	} else if post_id == "" && comment_id != "" {
		fmt.Println("ANA hna 2")

		err := data.Db.QueryRow("SELECT COUNT(*) FROM likes WHERE liked_comment_id = ?", comment_id).Scan(&LikeCount)
		if err != nil {
			fmt.Println("Error fetching like count:", err)
			return 0, 0, err
		}
		err = data.Db.QueryRow("SELECT COUNT(*) FROM dislikes WHERE disliked_comment_id = ?", comment_id).Scan(&DislikeCount)
		if err != nil {
			fmt.Println("Error fetching dislike count:", err)
			return 0, 0, err
		}
	} else {
		return 0, 0, errors.New("bad data from the front end")
	}
	return LikeCount, DislikeCount, nil
}
