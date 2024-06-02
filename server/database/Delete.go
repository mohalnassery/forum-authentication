package database

func DeletePost(postID, userID int) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM posts WHERE post_id = ? AND author = ?", postID, userID)
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM `like-posts` WHERE postID = ?", postID)
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM `dislike-posts` WHERE postID = ?", postID)
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM comments WHERE post_id = ?", postID)
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM post_categories WHERE post_id = ?", postID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func DeleteComment(commentID, userID int) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete likes associated with the comment
	_, err = tx.Exec("DELETE FROM `like-comments` WHERE comment_id = ?", commentID)
	if err != nil {
		return err
	}

	// Delete dislikes associated with the comment
	_, err = tx.Exec("DELETE FROM `dislike-comments` WHERE comment_id = ?", commentID)
	if err != nil {
		return err
	}

	// Delete the comment itself
	_, err = tx.Exec("DELETE FROM comments WHERE id = ? AND author = ?", commentID, userID)
	if err != nil {
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
