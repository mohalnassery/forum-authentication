package database

import (
	"forum/models"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func InsertUser(user *models.UserRegisteration) error {
	// hashing the password
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return err
	}
	hashedPassword := string(bytes)

	// insert the user into the database. executes without returning any rows
	_, err = DB.Exec("INSERT INTO users (username, email, password, auth_type) VALUES (?, ?, ?, ?)", user.Username, user.Email, hashedPassword, user.AuthType)
	if err != nil {
		return err
	}

	return nil
}

func InsertPost(post *models.Post) error { // Insert the post into the database
	result, err := DB.Exec("INSERT INTO posts (title, body, creation_date, author) VALUES (?, ?, ?, ?)",
		post.Title, post.Body, post.CreationDate, post.AuthorID)
	if err != nil {
		return err
	}

	// Get the ID of the newly inserted post
	postID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// Set the ID of the post object
	post.PostID = int(postID)

	// Insert the categories for the post
	for _, category := range post.Categories {
		// Get the category ID from the database
		categoryID, err := GetCategoryID(category)
		if err != nil {
			return err
		}

		// Insert the post-category relationship into the database
		_, err = DB.Exec("INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)", post.PostID, categoryID)
		if err != nil {
			return err
		}
	}
	return nil
}

func InsertComment(username string, body string, postID int) error {
	// Get the user's ID from the database
	userID, err := GetUserID(username)
	if err != nil {
		return err
	}

	// Insert the comment into the database
	_, err = DB.Exec(`
        INSERT INTO comments (body, post_id, author)
        VALUES (?, ?, ?)
    `, body, postID, userID)
	if err != nil {
		return err
	}
	return nil
}

func InsertPostLike(postID int, userID int) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check if the user has already liked the post
	var count int
	err = tx.QueryRow(`SELECT COUNT(*) FROM "like-posts" WHERE postID = ? AND user_id = ?`, postID, userID).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		// User has already liked the post, so remove the like
		_, err = tx.Exec(`DELETE FROM "like-posts" WHERE postID = ? AND user_id = ?`, postID, userID)
		if err != nil {
			return err
		}
	} else {
		// Remove any existing dislike for the post by the user
		_, err = tx.Exec(`DELETE FROM "dislike-posts" WHERE postID = ? AND user_id = ?`, postID, userID)
		if err != nil {
			return err
		}

		// Insert the like for the post by the user
		_, err = tx.Exec(`INSERT INTO "like-posts" (postID, user_id) VALUES (?, ?)`, postID, userID)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func InsertPostDislike(postID int, userID int) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check if the user has already disliked the post
	var count int
	err = tx.QueryRow(`SELECT COUNT(*) FROM "dislike-posts" WHERE postID = ? AND user_id = ?`, postID, userID).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		// User has already disliked the post, so remove the dislike
		_, err = tx.Exec(`DELETE FROM "dislike-posts" WHERE postID = ? AND user_id = ?`, postID, userID)
		if err != nil {
			return err
		}
	} else {
		// Remove any existing like for the post by the user
		_, err = tx.Exec(`DELETE FROM "like-posts" WHERE postID = ? AND user_id = ?`, postID, userID)
		if err != nil {
			return err
		}

		// Insert the dislike for the post by the user
		_, err = tx.Exec(`INSERT INTO "dislike-posts" (postID, user_id) VALUES (?, ?)`, postID, userID)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func InsertCommentLike(commentID, userID int) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check if the user has already liked the comment
	var count int
	err = tx.QueryRow(`SELECT COUNT(*) FROM "like-comments" WHERE comment_id = ? AND user_id = ?`, commentID, userID).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		// User has already liked the comment, so remove the like
		_, err = tx.Exec(`DELETE FROM "like-comments" WHERE comment_id = ? AND user_id = ?`, commentID, userID)
		if err != nil {
			return err
		}
	} else {
		// Remove any existing dislike for the comment by the user
		_, err = tx.Exec(`DELETE FROM "dislike-comments" WHERE comment_id = ? AND user_id = ?`, commentID, userID)
		if err != nil {
			return err
		}

		// Insert the like for the comment by the user
		_, err = tx.Exec(`INSERT INTO "like-comments" (comment_id, user_id) VALUES (?, ?)`, commentID, userID)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func InsertCommentDislike(commentID, userID int) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check if the user has already disliked the comment
	var count int
	err = tx.QueryRow(`SELECT COUNT(*) FROM "dislike-comments" WHERE comment_id = ? AND user_id = ?`, commentID, userID).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		// User has already disliked the comment, so remove the dislike
		_, err = tx.Exec(`DELETE FROM "dislike-comments" WHERE comment_id = ? AND user_id = ?`, commentID, userID)
		if err != nil {
			return err
		}
	} else {
		// Remove any existing like for the comment by the user
		_, err = tx.Exec(`DELETE FROM "like-comments" WHERE comment_id = ? AND user_id = ?`, commentID, userID)
		if err != nil {
			return err
		}

		// Insert the dislike for the comment by the user
		_, err = tx.Exec(`INSERT INTO "dislike-comments" (comment_id, user_id) VALUES (?, ?)`, commentID, userID)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
