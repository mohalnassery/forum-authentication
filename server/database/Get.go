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
