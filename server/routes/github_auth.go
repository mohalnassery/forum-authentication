package routes

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"forum/database"
	"forum/models"
	"net/http"
	"net/url"
	"time"

	"github.com/google/go-github/v52/github"
	"golang.org/x/oauth2"
	githubOAuth "golang.org/x/oauth2/github"
)

var (
	githubOauthConfig = &oauth2.Config{
		ClientID:     "Ov23liAQ6ck4v8pHIbG0",
		ClientSecret: "7760bf27f11f1df3008dc854cd9385f8899e95b4",
		RedirectURL:  "http://localhost:8080/auth/github/callback",
		Scopes:       []string{"user:email"},
		Endpoint:     githubOAuth.Endpoint,
	}
)

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(20 * time.Minute)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

func GithubLogin(w http.ResponseWriter, r *http.Request) {
	oauthState := generateStateOauthCookie(w)
	url := githubOauthConfig.AuthCodeURL(oauthState, oauth2.SetAuthURLParam("prompt", "login"))
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GithubCallback(w http.ResponseWriter, r *http.Request) {
	// Get the OAuth state from the cookie
	oauthState, err := r.Cookie("oauthstate")
	if err != nil {
		fmt.Println("Missing oauth state cookie")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if r.FormValue("state") != oauthState.Value {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthState.Value, r.FormValue("state"))
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Exchange the authorization code for an access token
	code := r.FormValue("code")
	token, err := githubOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Println("Failed to exchange token:", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Create a new GitHub client using the access token
	client := github.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token)))

	// Get the authenticated user's information
	user, _, err := client.Users.Get(context.Background(), "")
	if err != nil {
		fmt.Println("Failed to get user info:", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Check if the user already exists in the database
	existingUser, err := database.GetUserByEmail(user.GetEmail())
	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if existingUser != nil {
		// Check if the existing user's auth type matches the current login platform
		if existingUser.AuthType != "github" {
			// Display an error message to the user
			errorMessage := "Login failed. There is already an email registered with another platform."
			http.Redirect(w, r, "/login?error="+url.QueryEscape(errorMessage), http.StatusTemporaryRedirect)
			return
		}

		// User already exists, perform login
		err = CreateSession(w, r, models.UserRegisteration{
			Username: existingUser.Username,
		})
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		// User doesn't exist, perform registration
		newUser := &models.UserRegisteration{
			Username: user.GetLogin(),
			Email:    user.GetEmail(),
			AuthType: "github",
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
