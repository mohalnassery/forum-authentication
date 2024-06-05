package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"forum/database"
	"forum/models"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form data
	err := r.ParseMultipartForm(20 << 20) // 20 MB max file size
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the form values for title and body
	title := r.FormValue("title")
	body := r.FormValue("message")

	// Get the form values for categories
	categories := r.Form["options"]

	// Get the user from the session
	user, err := GetSessionUser(r)
	if err != nil {
		http.Error(w, "User not logged in", http.StatusUnauthorized)
		return
	}

	// Get the user's ID from the database
	userID, err := database.GetUserID(user.Username)
	if err != nil {
		http.Error(w, "Failed to retrieve user ID", http.StatusInternalServerError)
		return
	}

	// Get the image file from the form data
	file, header, err := r.FormFile("image")
	if err != nil && err != http.ErrMissingFile {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var imagePath string
	var filename string
	if file != nil {
		// Validate the file type
		allowedTypes := map[string]bool{
			"image/jpeg": true,
			"image/png":  true,
			"image/gif":  true,
		}
		if !allowedTypes[header.Header.Get("Content-Type")] {
			http.Error(w, "Invalid file type. Only JPEG, PNG, and GIF are allowed.", http.StatusBadRequest)
			return
		}

		// Generate a unique filename for the uploaded image
		ext := filepath.Ext(header.Filename)
		filename = fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		imagePath = filepath.Join("../client/uploads", filename)

		// Save the uploaded image file
		dst, err := os.Create(imagePath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Create a new post object with the form data and user ID
	post := &models.Post{
		Title:        title,
		Body:         body,
		CreationDate: time.Now().Format(time.RFC3339),
		AuthorID:     userID,
		Categories:   categories,
		Image:        filename,
	}

	err = database.InsertPost(post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to the main page after successfully creating a post
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	// Main selection
	query := `
        SELECT
            posts.post_id,
            posts.title,
            posts.body,
            posts.creation_date,
            users.username,
            GROUP_CONCAT(DISTINCT categories.name) AS categories,
            (SELECT COUNT(*) FROM "like-posts" WHERE "like-posts".postID = posts.post_id) AS likes,
            (SELECT COUNT(*) FROM "dislike-posts" WHERE "dislike-posts".postID = posts.post_id) AS dislikes,
            (SELECT COUNT(*) FROM comments WHERE post_id = posts.post_id) AS comment_count
        FROM posts
        JOIN users ON posts.author = users.id
        LEFT JOIN post_categories ON posts.post_id = post_categories.post_id
        LEFT JOIN categories ON post_categories.category_id = categories.id
    `

	var params []interface{}
	var whereClause []string

	// Filters
	fullQuery := r.URL.Query()
	if len(fullQuery) > 0 {
		// Get potential filter options
		categoryQuery := fullQuery["category"]
		filterLikes := fullQuery["liked"]
		filterCreated := fullQuery["created"]

		if len(categoryQuery) > 0 {
			categoryNames := strings.Split(categoryQuery[0], ",")
			placeholders := make([]string, len(categoryNames))
			for i := range categoryNames {
				placeholders[i] = "?"
				params = append(params, categoryNames[i])
			}
			whereClause = append(whereClause, "posts.post_id IN (SELECT post_id FROM post_categories WHERE category_id IN (SELECT id FROM categories WHERE name IN ("+strings.Join(placeholders, ",")+")))")
		}
		if len(filterCreated) > 0 {
			whereClause = append(whereClause, "users.username = ?")
			params = append(params, filterCreated[0])
		}
		if len(filterLikes) > 0 {
			whereClause = append(whereClause, "posts.post_id IN (SELECT postID FROM 'like-posts' WHERE user_id IN (SELECT id FROM users WHERE username = ?))")
			params = append(params, filterLikes[0])
		}
	}

	if len(whereClause) > 0 {
		query += " WHERE " + strings.Join(whereClause, " AND ")
	}

	// Sorting
	query += `
		GROUP BY posts.post_id
		ORDER BY posts.creation_date DESC;
	`

	// Execution
	rows, err := database.DB.Query(query, params...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []models.PostHome
	for rows.Next() {
		var post models.PostHome
		var categoriesString sql.NullString
		if err := rows.Scan(&post.PostID, &post.Title, &post.Body, &post.CreationDate, &post.Author, &categoriesString, &post.Likes, &post.Dislikes, &post.CommentCount); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if categoriesString.Valid {
			post.Categories = strings.Split(categoriesString.String, ",")
		} else {
			post.Categories = []string{}
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// setting the heading in the front-end to be of type application json
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(posts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	// Extract the postID from the URL path
	urlParts := strings.Split(r.URL.Path, "/")
	if len(urlParts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(urlParts[2])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	// Get the user from the session
	user, err := GetSessionUser(r)
	var userID int
	if err != nil {
		// User is not logged in, set userID to 0
		userID = 0
	} else {
		// Get the user's ID from the database
		userID, err = database.GetUserID(user.Username)
		if err != nil {
			http.Error(w, "Failed to retrieve user ID", http.StatusInternalServerError)
			return
		}
	}
	// Retrieve the post details from the database
	var post models.PostHome
	var categoriesString sql.NullString
	err = database.DB.QueryRow(`
		SELECT 
			posts.post_id, 
			posts.title, 
			posts.body, 
			posts.creation_date,
			posts.image,
			users.username,
			GROUP_CONCAT(DISTINCT categories.name) AS categories,
			(SELECT COUNT(*) FROM "like-posts" WHERE "like-posts".postID = posts.post_id) AS likes,
			(SELECT COUNT(*) FROM "dislike-posts" WHERE "dislike-posts".postID = posts.post_id) AS dislikes,
			(SELECT COUNT(*) FROM comments WHERE post_id = posts.post_id) AS comment_count,
			COALESCE((SELECT 1 FROM "like-posts" WHERE "like-posts".postID = posts.post_id AND user_id = ?), 0) AS user_liked,
			COALESCE((SELECT 1 FROM "dislike-posts" WHERE "dislike-posts".postID = posts.post_id AND user_id = ?), 0) AS user_disliked,
			CASE WHEN posts.author = ? THEN 1 ELSE 0 END AS is_author
		FROM posts
		JOIN users ON posts.author = users.id
		LEFT JOIN post_categories ON posts.post_id = post_categories.post_id
		LEFT JOIN categories ON post_categories.category_id = categories.id
		WHERE posts.post_id = ?
		GROUP BY posts.post_id
	`, userID, userID, userID, postID).Scan(
		&post.PostID,
		&post.Title,
		&post.Body,
		&post.CreationDate,
		&post.Image,
		&post.Author,
		&categoriesString,
		&post.Likes,
		&post.Dislikes,
		&post.CommentCount,
		&post.UserLiked,
		&post.UserDisliked,
		&post.IsAuthor,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert categories string to slice
	if categoriesString.Valid {
		post.Categories = strings.Split(categoriesString.String, ",")
	} else {
		post.Categories = []string{}
	}

	// Retrieve the comments for the post from the database
	rows, err := database.DB.Query(`
	SELECT 
		comments.id,
		comments.body,
		comments.creation_date,
		users.username,
		COALESCE(SUM(CASE WHEN 'like-comments'.id IS NOT NULL THEN 1 ELSE 0 END), 0) AS likes,
		COALESCE(SUM(CASE WHEN 'dislike-comments'.id IS NOT NULL THEN 1 ELSE 0 END), 0) AS dislikes,
		COALESCE((SELECT 1 FROM "like-comments" WHERE "like-comments".comment_id = comments.id AND user_id = ?), 0) AS user_liked,
		COALESCE((SELECT 1 FROM "dislike-comments" WHERE "dislike-comments".comment_id = comments.id AND user_id = ?), 0) AS user_disliked,
		CASE WHEN comments.author = ? THEN 1 ELSE 0 END AS is_author
	FROM comments
	JOIN users ON comments.author = users.id
	LEFT JOIN 'like-comments' ON comments.id = 'like-comments'.comment_id
	LEFT JOIN 'dislike-comments' ON comments.id = 'dislike-comments'.comment_id
	WHERE comments.post_id = ?
	GROUP BY comments.id
`, userID, userID, userID, postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var comments []*models.CommentHome
	for rows.Next() {
		var comment models.CommentHome
		err := rows.Scan(
			&comment.ID,
			&comment.Body,
			&comment.CreationDate,
			&comment.Author,
			&comment.Likes,
			&comment.Dislikes,
			&comment.UserLiked,
			&comment.UserDisliked,
			&comment.IsAuthor,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		comments = append(comments, &comment)
	}

	// Create a response object containing the post details and comments
	response := struct {
		Post     *models.PostHome      `json:"post"`
		Comments []*models.CommentHome `json:"comments"`
	}{
		Post:     &post,
		Comments: comments,
	}

	// Return the response as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	urlParts := strings.Split(r.URL.Path, "/")
	if len(urlParts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(urlParts[2])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
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

	err = database.DeletePost(postID, userID)
	if err != nil {
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}
