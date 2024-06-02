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
