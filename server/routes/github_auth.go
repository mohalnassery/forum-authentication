package routes

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"forum/database"
	"forum/models"
	"net/http"
	"regexp"
	"time"

	"github.com/google/go-github/v52/github"
	"golang.org/x/oauth2"
	githubOAuth "golang.org/x/oauth2/github"
)

var (
	githubOauthConfig = &oauth2.Config{
		ClientID:     "Ov23liAQ6ck4v8pHIbG0",
		ClientSecret: "7760bf27f11f1df3008dc854cd9385f8899e95b4",
		RedirectURL:  "https://localhost:8443/auth/github/callback",
		Scopes:       []string{"user:email"},
		Endpoint:     githubOAuth.Endpoint,
	}
)

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(20 * time.Minute)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration, HttpOnly: true, Secure: true}
	http.SetCookie(w, &cookie)

	return state
}

func GithubLogin(w http.ResponseWriter, r *http.Request) {
	oauthState := generateStateOauthCookie(w)
	url := githubOauthConfig.AuthCodeURL(oauthState, oauth2.SetAuthURLParam("prompt", "login"))
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GithubCallback(w http.ResponseWriter, r *http.Request) {
	oauthState, err := r.Cookie("oauthstate")
	if err != nil || r.FormValue("state") != oauthState.Value {
		fmt.Println("invalid oauth state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := githubOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Println("Failed to exchange token:", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	client := github.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token)))

	user, _, err := client.Users.Get(context.Background(), "")
	if err != nil {
		fmt.Println("Failed to get user info:", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Validate and sanitize user data
	userLogin := sanitizeInput(user.GetLogin())
	userEmail := sanitizeInput(user.GetEmail())

	existingUser, err := database.GetUserByEmail(userEmail)
	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if existingUser != nil {
		err = CreateSession(w, r, models.UserRegisteration{
			Username: existingUser.Username,
		})
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		newUser := &models.UserRegisteration{
			Username: userLogin,
			Email:    userEmail,
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

func sanitizeInput(input string) string {
	// Remove any characters that are not alphanumeric, @, ., or _
	re := regexp.MustCompile(`[^a-zA-Z0-9@._]`)
	return re.ReplaceAllString(input, "")
}
