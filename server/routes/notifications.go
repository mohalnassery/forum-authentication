package routes

import (
	"encoding/json"
	"forum/database"
	"net/http"
	"strconv"
	"strings"
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

func ClearNotification(w http.ResponseWriter, r *http.Request) {
	urlParts := strings.Split(r.URL.Path, "/")
	if len(urlParts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	notificationID, err := strconv.Atoi(urlParts[2])
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	err = database.ClearNotification(notificationID)
	if err != nil {
		http.Error(w, "Failed to clear notification", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func MarkAllNotificationsAsRead(w http.ResponseWriter, r *http.Request) {
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

	err = database.MarkAllNotificationsAsRead(userID)
	if err != nil {
		http.Error(w, "Failed to mark all notifications as read", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
