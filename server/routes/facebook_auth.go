package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"forum/database"
	"forum/models"
	"net/http"

	"golang.org/x/oauth2"
	facebookOAuth "golang.org/x/oauth2/facebook"
)

var (
	facebookOauthConfig = &oauth2.Config{
		ClientID:     "1531896514060809",
		ClientSecret: "1022032e6c06c90bfa49f26307c476e4",
		RedirectURL:  "https://localhost:8443/auth/facebook/callback",
		Scopes:       []string{"email"},
		Endpoint:     facebookOAuth.Endpoint,
	}
)

func FacebookLogin(w http.ResponseWriter, r *http.Request) {
	oauthState := generateStateOauthCookie(w)
	url := facebookOauthConfig.AuthCodeURL(oauthState, oauth2.SetAuthURLParam("auth_type", "reauthenticate"))
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func FacebookCallback(w http.ResponseWriter, r *http.Request) {
	oauthState, err := r.Cookie("oauthstate")
	if err != nil || r.FormValue("state") != oauthState.Value {
		fmt.Println("invalid oauth state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	token, err := facebookOauthConfig.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		fmt.Println("Failed to exchange token:", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	client := &http.Client{}
	res, err := client.Get("https://graph.facebook.com/me?fields=id,name,email&access_token=" + token.AccessToken)
	if err != nil {
		fmt.Println("Failed to get user info:", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer res.Body.Close()

	var user struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	err = json.NewDecoder(res.Body).Decode(&user)
	if err != nil {
		fmt.Println("Failed to decode user info:", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Validate and sanitize user data
	user.Name = sanitizeInput(user.Name)
	user.Email = sanitizeInput(user.Email)

	existingUser, err := database.GetUserByEmail(user.Email)
	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if existingUser != nil {
<<<<<<< HEAD
=======
		// // Check if the existing user's auth type matches the current login platform
		// if existingUser.AuthType != "facebook" {
		// 	// Display an error message to the user
		// 	errorMessage := "Login failed. There is already an email registered with another platform."
		// 	http.Redirect(w, r, "/login?error="+url.QueryEscape(errorMessage), http.StatusTemporaryRedirect)
		// 	return
		// }


		// User already exists, perform login
>>>>>>> reboot/master
		err = CreateSession(w, r, models.UserRegisteration{
			Username: existingUser.Username,
		})
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		newUser := &models.UserRegisteration{
			Username: user.Name,
			Email:    user.Email,
			AuthType: "facebook",
		}
		err = database.InsertUser(newUser)
		if err != nil {
			fmt.Println(err.Error())
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		err = CreateSession(w, r, *newUser)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
