package main

import (
	"fmt"
	"forum/database"
	"forum/routes"
	"log"
	"net/http"
)

func main() {
	database.InitDatabaseTables()
	defer database.DB.Close()

	r := http.NewServeMux()

	// displaying pages
	r.HandleFunc("/", routes.HandleGet)
	r.HandleFunc("/register", routes.HandleGet)
	r.HandleFunc("/login", routes.HandleGet)
	r.HandleFunc("/create", routes.HandleGet)
	r.HandleFunc("/post-details/{postId}", routes.HandleGet)
	r.Handle("/favicon.ico", http.FileServer(http.Dir("../client/assets")))
	// google auth end points
	r.HandleFunc("/auth/google", routes.GoogleLogin)
	r.HandleFunc("/auth/google/callback", routes.GoogleCallback)
	// github auth end points
	r.HandleFunc("/auth/github/login", routes.GithubLogin)
	r.HandleFunc("/auth/github/callback", routes.GithubCallback)
	// facebook auth end points
	r.HandleFunc("/auth/facebook/login", routes.FacebookLogin)
	r.HandleFunc("/auth/facebook/callback", routes.FacebookCallback)

	// functionality end points (NOTE: USING GO VERSION 1.22 FOR BETTER ROUTING)
	r.HandleFunc("/categories", routes.GetCategories)
	r.HandleFunc("POST /auth/register", routes.UserRegister)
	r.HandleFunc("POST /auth/login", routes.UserLogin)
	r.HandleFunc("POST /posts", routes.CreatePost)
	r.HandleFunc("GET /posts", routes.GetPosts)
	r.HandleFunc("GET /posts/{postId}", routes.GetPost)

	r.HandleFunc("POST /posts/{postId}/like", routes.LikePost)
	r.HandleFunc("POST /posts/{postId}/dislike", routes.DislikePost)
	r.HandleFunc("POST /posts/{postId}/comments", routes.CreateComment)
	r.HandleFunc("POST /comments/{commentId}/like", routes.LikeComment)
	r.HandleFunc("POST /comments/{commentId}/dislike", routes.DislikeComment)

	r.HandleFunc("/auth/is-logged-in", routes.IsLoggedIn)
	r.HandleFunc("POST /auth/logout", routes.Logout)
	r.HandleFunc("GET /user-stats", routes.GetUserStats)
	r.HandleFunc("GET /all-stats", routes.GetAllStats)
	r.HandleFunc("GET /leaderboard", routes.GetLeaderboard)
	r.HandleFunc("DELETE /posts/{postId}", routes.DeletePost)
	r.HandleFunc("DELETE /comments/{commentId}", routes.DeleteComment)
	// serving static files
	r.Handle("/js/", http.FileServer(http.Dir("../client")))
	r.Handle("/css/", http.FileServer(http.Dir("../client")))
	r.Handle("/assets/", http.FileServer(http.Dir("../client")))

	fmt.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
