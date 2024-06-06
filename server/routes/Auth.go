package routes

import (
	"encoding/json"
	"fmt"
	"forum/database"
	"forum/models"
	"net/http"
	"regexp"

	_ "github.com/mattn/go-sqlite3"
)

func UserRegister(w http.ResponseWriter, r *http.Request) {
	var user *models.UserRegisteration
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	taken, err := database.IsUserTaken(user)
	if err != nil {
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}
	if taken != "" {
		http.Error(w, taken, http.StatusBadRequest)
		return
	}

	// regex for valid email
	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{1,}$`

	match, err := regexp.MatchString(emailRegex, user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !match {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	err = database.InsertUser(user)
	if err != nil {
		http.Error(w, "Failed to insert into database", http.StatusInternalServerError)
		return
	}
	// Create a session for the user
	err = CreateSession(w, r, *user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Return a success message
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User registered successfully")
}
func UserLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "You shouldn't be here", http.StatusBadRequest)
		return
	}
	var user *models.UserLogin
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	valid, err := database.IsValidLogin(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !valid {
		http.Error(w, "Username or password is incorrect", http.StatusBadRequest)
		return
	}

	// Create a new session for the user
	err = CreateSession(w, r, models.UserRegisteration{
		Username: user.Username,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Send a response to the client indicating successful login
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User logged in successfully")
}

func IsLoggedIn(w http.ResponseWriter, r *http.Request) {
	user, err := GetSessionUser(r)
	if err != nil {
		if err.Error() == "newer session found" {
			DestroySession(w, r)
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "not_logged_in"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":    "logged_in",
		"username":  user.Username,
		"sessionID": UserSessions[user.Username],
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	user, err := GetSessionUser(r)
	if err != nil {
		delete(UserSessions, user.Username)
	}
	DestroySession(w, r)
	json.NewEncoder(w).Encode(map[string]string{"status": "logged_out"})
}
