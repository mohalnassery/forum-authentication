package database

import (
	"database/sql"
	"fmt"
	"forum/models"
)

func GetUserID(username string) (int, error) {
	var userID int
	err := DB.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func GetCategoryID(categoryName string) (int, error) {
	var categoryID int
	row := DB.QueryRow("SELECT id FROM categories WHERE name = ?", categoryName)
	err := row.Scan(&categoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("category not found: %s", categoryName)
		}
		return 0, fmt.Errorf("failed to retrieve category ID: %v", err)
	}
	return categoryID, nil
}

func GetUserByEmail(email string) (*models.UserRegisteration, error) {
	var user models.UserRegisteration
	err := DB.QueryRow("SELECT username, email, auth_type FROM users WHERE email = ?", email).Scan(&user.Username, &user.Email, &user.AuthType)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func GetUsernameByID(userID int) (string, error) {
	var username string
	err := DB.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}

func GetNotifications(userID int) ([]models.Notification, error) {
	rows, err := DB.Query("SELECT id, message, created_at, is_read FROM notifications WHERE user_id = ? ORDER BY created_at DESC", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var notification models.Notification
		err := rows.Scan(&notification.ID, &notification.Message, &notification.CreatedAt, &notification.IsRead)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}
	return notifications, nil
}

func ClearNotification(notificationID int) (int, error) {
	var postID int
	err := DB.QueryRow("SELECT post_id FROM notifications WHERE id = ?", notificationID).Scan(&postID)
	if err != nil {
		return 0, err
	}

	_, err = DB.Exec("UPDATE notifications SET is_read = 1 WHERE id = ?", notificationID)
	if err != nil {
		return 0, err
	}

	return postID, nil
}

func MarkAllNotificationsAsRead(userID int) error {
	_, err := DB.Exec("UPDATE notifications SET is_read = TRUE WHERE user_id = ?", userID)
	return err
}

func GetPostIDFromNotification(notificationID int) (int, error) {
	var postID int
	err := DB.QueryRow("SELECT post_id FROM notifications WHERE id = ?", notificationID).Scan(&postID)
	if err != nil {
		return 0, err
	}
	return postID, nil
}

func GetPostIDByCommentID(commentID int) (int, error) {
	var postID int
	err := DB.QueryRow("SELECT post_id FROM comments WHERE id = ?", commentID).Scan(&postID)
	if err != nil {
		return 0, err
	}
	return postID, nil
}

func GetUserActivity(userID int) (models.UserActivity, error) {
	var activity models.UserActivity

	// Retrieve user-created posts
	rows, err := DB.Query("SELECT post_id, title, body FROM posts WHERE author = ?", userID)
	if err != nil {
		return activity, err
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.PostID, &post.Title, &post.Body)
		if err != nil {
			return activity, err
		}
		activity.CreatedPosts = append(activity.CreatedPosts, post)
	}

	// Retrieve posts where the user left a like or dislike
	rows, err = DB.Query(`
		SELECT p.post_id, p.title, p.body, l.liked
		FROM posts p
		JOIN "like-posts" l ON p.post_id = l.postID
		WHERE l.user_id = ?`, userID)
	if err != nil {
		return activity, err
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		var liked sql.NullBool
		err := rows.Scan(&post.PostID, &post.Title, &post.Body, &liked)
		if err != nil {
			return activity, err
		}
		if liked.Valid && liked.Bool {
			activity.LikedPosts = append(activity.LikedPosts, post)
		} else if liked.Valid && !liked.Bool {
			activity.DislikedPosts = append(activity.DislikedPosts, post)
		}
	}

	// Retrieve comments made by the user, along with the corresponding post information
	rows, err = DB.Query(`
		SELECT c.id, c.body, p.post_id, p.title, p.body
		FROM comments c
		JOIN posts p ON c.post_id = p.post_id
		WHERE c.author = ?`, userID)
	if err != nil {
		return activity, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment models.CommentHome
		var post models.Post
		err := rows.Scan(&comment.ID, &comment.Body, &post.PostID, &post.Title, &post.Body)
		if err != nil {
			return activity, err
		}
		comment.PostID = post.PostID
		activity.Comments = append(activity.Comments, comment)
	}

	// Retrieve comments liked by the user
	rows, err = DB.Query(`
		SELECT c.id, c.body, p.post_id, p.title, p.body
		FROM comments c
		JOIN posts p ON c.post_id = p.post_id
		JOIN "like-comments" lc ON c.id = lc.comment_id
		WHERE lc.user_id = ?`, userID)
	if err != nil {
		return activity, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment models.CommentHome
		var post models.Post
		err := rows.Scan(&comment.ID, &comment.Body, &post.PostID, &post.Title, &post.Body)
		if err != nil {
			return activity, err
		}
		comment.PostID = post.PostID
		activity.LikedComments = append(activity.LikedComments, comment)
	}

	// Retrieve comments disliked by the user
	rows, err = DB.Query(`
		SELECT c.id, c.body, p.post_id, p.title, p.body
		FROM comments c
		JOIN posts p ON c.post_id = p.post_id
		JOIN "dislike-comments" dc ON c.id = dc.comment_id
		WHERE dc.user_id = ?`, userID)
	if err != nil {
		return activity, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment models.CommentHome
		var post models.Post
		err := rows.Scan(&comment.ID, &comment.Body, &post.PostID, &post.Title, &post.Body)
		if err != nil {
			return activity, err
		}
		comment.PostID = post.PostID
		activity.DislikedComments = append(activity.DislikedComments, comment)
	}

	return activity, nil
}
