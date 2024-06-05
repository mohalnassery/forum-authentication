package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"forum/database"
	"forum/models"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"io"
	"net/http"
	"net/url"
)

var (
	googleOauthConfig = &oauth2.Config{
		ClientID:     "632453699826-cahu1ldajpebroq35vb25v9v1n4r4vk0.apps.googleusercontent.com",
		ClientSecret: "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",
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	oauthStateString = "random_string"
)

func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	content, err := getUserInfo(r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	var user models.GoogleUser
	err = json.Unmarshal(content, &user)
	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Check if the user already exists in the database
	existingUser, err := database.GetUserByEmail(user.Email)
	if err != nil {
		fmt.Println(err.Error()) // Log the error
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if existingUser != nil {
		// // Check if the existing user's auth type matches the current login platform
		// if existingUser.AuthType != "google" {
		// 	// Display an error message to the user
		// 	errorMessage := "Login failed. There is already an email registered with another platform."
		// 	http.Redirect(w, r, "/login?error="+url.QueryEscape(errorMessage), http.StatusTemporaryRedirect)
		// 	return
		// }
		// User already exists, perform login
		err = CreateSession(w, r, models.UserRegisteration{
			Username: existingUser.Username,
		})
		if err != nil {
			fmt.Println(err.Error()) // Log the error
		}
	} else {
		// User doesn't exist, perform registration
		newUser := &models.UserRegisteration{
			Username: user.Name,
			Email:    user.Email,
			AuthType: "google",
		}
		err = database.InsertUser(newUser)
		if err != nil {
			fmt.Println(err.Error()) // Log the error
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		err = CreateSession(w, r, *newUser)
		if err != nil {
			fmt.Println(err.Error()) // Log the error
		}
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func getUserInfo(state string, code string) ([]byte, error) {
	if state != oauthStateString {
		return nil, fmt.Errorf("invalid oauth state")
	}

	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}

	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}

	return contents, nil
}
