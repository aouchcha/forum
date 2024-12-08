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
		var LikeCount, DislikeCount int
		post_id := r.URL.Query().Get("postid")
		comment_id := r.URL.Query().Get("comment_id")
		fmt.Println(post_id, comment_id)
		// Fetch the Like and Dislike counts
		LikeCount, DislikeCount, err := getLikeAndDislikeCount(post_id, comment_id)
		if err != nil {
			http.Error(w, "Error fetching like and dislike count", http.StatusInternalServerError)
			return
		}

		response, err := json.Marshal(map[string]string{
			"LikeCount":    strconv.Itoa(LikeCount),
			"DislikeCount": strconv.Itoa(DislikeCount),
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
