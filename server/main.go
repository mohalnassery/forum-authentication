package main

import (
	"crypto/tls"
	"fmt"
	"forum/database"
	"forum/routes"
	"log"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

func main() {
	database.InitDatabaseTables()
	defer database.DB.Close()

	r := http.NewServeMux()

	// Create rate limiters for different routes
	generalLimiter := rate.NewLimiter(rate.Every(time.Minute), 500) // 500 requests per minute
	authLimiter := rate.NewLimiter(rate.Every(time.Second), 5)      // 5 requests per second
	postLimiter := rate.NewLimiter(rate.Every(time.Minute), 500)    // 500 requests per minute
	commentLimiter := rate.NewLimiter(rate.Every(time.Minute), 200) // 30 requests per minute

	// Displaying pages
	r.Handle("/", limitMiddleware(generalLimiter, http.HandlerFunc(routes.HandleGet)))
	r.Handle("/register", limitMiddleware(generalLimiter, http.HandlerFunc(routes.HandleGet)))
	r.Handle("/login", limitMiddleware(generalLimiter, http.HandlerFunc(routes.HandleGet)))
	r.Handle("/create", limitMiddleware(generalLimiter, http.HandlerFunc(routes.HandleGet)))
	r.Handle("/post-details/{postId}", limitMiddleware(generalLimiter, http.HandlerFunc(routes.HandleGet)))
	r.Handle("/favicon.ico", http.FileServer(http.Dir("../client/assets")))

	// Google auth endpoints
	r.Handle("/auth/google", limitMiddleware(authLimiter, http.HandlerFunc(routes.GoogleLogin)))
	r.Handle("/auth/google/callback", limitMiddleware(authLimiter, http.HandlerFunc(routes.GoogleCallback)))

	// GitHub auth endpoints
	r.Handle("/auth/github/login", limitMiddleware(authLimiter, http.HandlerFunc(routes.GithubLogin)))
	r.Handle("/auth/github/callback", limitMiddleware(authLimiter, http.HandlerFunc(routes.GithubCallback)))

	// Facebook auth endpoints
	r.Handle("/auth/facebook/login", limitMiddleware(authLimiter, http.HandlerFunc(routes.FacebookLogin)))
	r.Handle("/auth/facebook/callback", limitMiddleware(authLimiter, http.HandlerFunc(routes.FacebookCallback)))

	// Functionality endpoints
	r.Handle("/categories", limitMiddleware(generalLimiter, http.HandlerFunc(routes.GetCategories)))
	r.Handle("/auth/register", limitMiddleware(authLimiter, http.HandlerFunc(routes.UserRegister)))
	r.Handle("/auth/login", limitMiddleware(authLimiter, http.HandlerFunc(routes.UserLogin)))
	r.Handle("/posts", limitMiddleware(postLimiter, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			routes.CreatePost(w, r)
		case http.MethodGet:
			routes.GetPosts(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	r.Handle("/posts/{postId}", limitMiddleware(generalLimiter, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			routes.GetPost(w, r)
		case http.MethodDelete:
			routes.DeletePost(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	r.Handle("/posts/{postId}/like", limitMiddleware(postLimiter, http.HandlerFunc(routes.LikePost)))
	r.Handle("/posts/{postId}/dislike", limitMiddleware(postLimiter, http.HandlerFunc(routes.DislikePost)))
	r.Handle("/posts/{postId}/comments", limitMiddleware(commentLimiter, http.HandlerFunc(routes.CreateComment)))
	r.Handle("/comments/{commentId}/like", limitMiddleware(commentLimiter, http.HandlerFunc(routes.LikeComment)))
	r.Handle("/comments/{commentId}/dislike", limitMiddleware(commentLimiter, http.HandlerFunc(routes.DislikeComment)))
	r.Handle("/auth/is-logged-in", limitMiddleware(authLimiter, http.HandlerFunc(routes.IsLoggedIn)))
	r.Handle("/auth/logout", limitMiddleware(authLimiter, http.HandlerFunc(routes.Logout)))
	r.Handle("/user-stats", limitMiddleware(generalLimiter, http.HandlerFunc(routes.GetUserStats)))
	r.Handle("/all-stats", limitMiddleware(generalLimiter, http.HandlerFunc(routes.GetAllStats)))
	r.Handle("/leaderboard", limitMiddleware(generalLimiter, http.HandlerFunc(routes.GetLeaderboard)))
	r.Handle("/comments/{commentId}", limitMiddleware(commentLimiter, http.HandlerFunc(routes.DeleteComment)))

	// Serving static files
	r.Handle("/js/", http.FileServer(http.Dir("../client")))
	r.Handle("/css/", http.FileServer(http.Dir("../client")))
	r.Handle("/assets/", http.FileServer(http.Dir("../client")))
	r.Handle("/uploads/", http.FileServer(http.Dir("../client")))

	// Configure SSL/TLS
	tlsConfig := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	server := &http.Server{
		Addr:         ":8443",
		Handler:      r,
		TLSConfig:    tlsConfig,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	fmt.Println("Server is running on https://localhost:8443")
	if err := server.ListenAndServeTLS("../server/cert.pem", "../server/key.pem"); err != nil {
		log.Fatal(err)
	}
}

// limitMiddleware is a middleware that applies rate limiting
func limitMiddleware(limiter *rate.Limiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
