package routes

import (
	"encoding/json"
	"forum/database"
	"net/http"
)

func GetNotifications(w http.ResponseWriter, r *http.Request) {
	user, err := GetSessionUser(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := database.GetUserID(user.Username)
	if err != nil {
		http.Error(w, "Failed to retrieve user ID", http.StatusInternalServerError)
		return
	}

	notifications, err := database.GetNotifications(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve notifications", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}
