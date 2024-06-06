package routes

import (
	"encoding/json"
	"forum/database"
	"forum/models"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func CreateComment(w http.ResponseWriter, r *http.Request) {
	// Extract the postID from the URL path
	urlParts := strings.Split(r.URL.Path, "/")
	if len(urlParts) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(urlParts[2])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Parse the request body
	var commentRequest models.CommentRequest
	err = json.NewDecoder(r.Body).Decode(&commentRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate comment length
	maxCharacters := 500
	if len(commentRequest.Body) > maxCharacters {
		http.Error(w, "Comment exceeds the maximum character limit", http.StatusBadRequest)
		return
	}

	// Get the user from the session
	user, err := GetSessionUser(r)
	if err != nil {
		http.Error(w, "User not logged in", http.StatusUnauthorized)
		return
	}

	err = database.InsertComment(user.Username, commentRequest.Body, postID)
	if err != nil {
		http.Error(w, "Failed to insert comment", http.StatusInternalServerError)
		return
	}

	// Create a notification for the post author
	postAuthorID, err := database.GetPostAuthorID(postID)
	if err != nil {
		http.Error(w, "Failed to retrieve post author ID", http.StatusInternalServerError)
		return
	}
	message := user.Username + " commented on your post"
	err = database.InsertNotification(postAuthorID, message)
	if err != nil {
		http.Error(w, "Failed to create notification", http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Comment posted successfully"})
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	urlParts := strings.Split(r.URL.Path, "/")
	if len(urlParts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	commentID, err := strconv.Atoi(urlParts[2])
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	user, err := GetSessionUser(r)
	if err != nil {
		http.Error(w, "User not logged in", http.StatusUnauthorized)
		return
	}

	userID, err := database.GetUserID(user.Username)
	if err != nil {
		http.Error(w, "Failed to retrieve user ID", http.StatusInternalServerError)
		return
	}

	err = database.DeleteComment(commentID, userID)
	if err != nil {
		http.Error(w, "Failed to delete comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
