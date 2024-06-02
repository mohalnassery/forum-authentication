package database

import (
	"database/sql"
	"forum/models"

	"golang.org/x/crypto/bcrypt"
)

func IsUserTaken(user *models.UserRegisteration) (string, error) {
	// check if username exists in the database
	var countUser int
	var countEmail int
	err := DB.QueryRow(` SELECT  (SELECT COUNT(*) FROM users WHERE username = ?) AS countUsers, 
	(SELECT COUNT(*) FROM users WHERE email = ?) AS countEmail`, user.Username, user.Email).Scan(&countUser, &countEmail)
	if err != nil {
		return "", err
	}
	if countUser > 0 {
		return "Username already taken", nil
	}
	if countEmail > 0 {
		return "Email already taken", nil
	}
	return "", nil
}

func IsValidLogin(user *models.UserLogin) (bool, error) {
	var username, password string

	// Check if username exiists
	err := DB.QueryRow("SELECT username, password FROM users WHERE username = ?", user.Username).Scan(&username, &password)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	// check if password is correct
	if err = bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password)); err != nil {
		return false, nil
	}

	return true, nil
}
