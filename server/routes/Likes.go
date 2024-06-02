package routes

import (
	"forum/database"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func LikePost(w http.ResponseWriter, r *http.Request) {
	postID, userID, errStr, errCode := getIDs(r)
	if errStr != "" {
		http.Error(w, errStr, errCode)
		return
	}

	err := database.InsertPostLike(postID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DislikePost(w http.ResponseWriter, r *http.Request) {
	postID, userID, errStr, errCode := getIDs(r)
	if errStr != "" {
		http.Error(w, errStr, errCode)
		return
	}

	err := database.InsertPostDislike(postID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func LikeComment(w http.ResponseWriter, r *http.Request) {
	commentID, userID, errStr, errCode := getIDs(r)
	if errStr != "" {
		http.Error(w, errStr, errCode)
		return
	}

	err := database.InsertCommentLike(commentID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DislikeComment(w http.ResponseWriter, r *http.Request) {
	commentID, userID, errStr, errCode := getIDs(r)
	if errStr != "" {
		http.Error(w, errStr, errCode)
		return
	}

	err := database.InsertCommentDislike(commentID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getIDs(r *http.Request) (int, int, string, int) {
	// Extract the object from the URL path
	urlParts := strings.Split(r.URL.Path, "/")
	if len(urlParts) < 4 {
		return 0, 0, "Invalid URL", http.StatusBadRequest
	}

	objID, err := strconv.Atoi(urlParts[2])
	if err != nil || objID < 1 {
		return 0, 0, "Invalid object ID", http.StatusBadRequest
	}

	user, err := GetSessionUser(r)
	if err != nil {
		return 0, 0, "User not logged in", http.StatusUnauthorized
	}

	userID, err := database.GetUserID(user.Username)
	if err != nil {
		return 0, 0, "User not found", http.StatusInternalServerError
	}
	return objID, userID, "", 0
}
