package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	data "main/dataBase"
)

// var like = data.Jaime{}

func LikesCounterWithApi(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/likes" {
		var LikeCount, DislikeCount int
		post_id := r.URL.Query().Get("postid")
		comment_id := r.URL.Query().Get("comment_id")
		LikeCount, DislikeCount, LikeComment, DislikeComment, err := GetReactionsCounts(post_id, comment_id)
		if err != nil {
			http.Error(w, "Error fetching like and dislike count", http.StatusInternalServerError)
			return
		}
		response, err := json.Marshal(map[string]string{
			"post_id":        post_id,
			"comment_id":     comment_id,
			"LikeCount":      strconv.Itoa(LikeCount),
			"DislikeCount":   strconv.Itoa(DislikeCount),
			"LikeComment":    strconv.Itoa(LikeComment),
			"DislikeComment": strconv.Itoa(DislikeComment),
		})
		if err != nil {
			http.Error(w, "Error marshaling JSON response", http.StatusInternalServerError)
			return
		}
		fmt.Println("Response:", string(response))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

func GetReactionsCounts(post_id, comment_id string) (int, int, int, int, error) {
	var LikeCount, DislikeCount int
	var LikeComment, DislikeComment int
	if post_id != "" {
		err := data.Db.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = ?", post_id).Scan(&LikeCount)
		if err != nil {
			fmt.Println("Error fetching like count:", err)
			return 0, 0, 0, 0, err
		}
		err = data.Db.QueryRow("SELECT COUNT(*) FROM dislikes WHERE post_id = ?", post_id).Scan(&DislikeCount)
		if err != nil {
			fmt.Println("Error fetching dislike count:", err)
			return 0, 0, 0, 0, err
		}
		return LikeCount, DislikeCount, 0, 0, nil
	} else {
		err := data.Db.QueryRow("SELECT COUNT(*) FROM likes WHERE liked_comment_id = ?", comment_id).Scan(&LikeComment)
		if err != nil {
			fmt.Println("Error fetching like count:", err)
			return 0, 0, 0, 0, err
		}
		err = data.Db.QueryRow("SELECT COUNT(*) FROM dislikes WHERE disliked_comment_id = ?", comment_id).Scan(&DislikeComment)
		if err != nil {
			fmt.Println("Error fetching dislike count:", err)
			return 0, 0, 0, 0, err
		}
		return 0, 0, LikeComment, DislikeComment, nil
	}
	// return LikeCount, DislikeCount, LikeComment, DislikeComment, nil
}
