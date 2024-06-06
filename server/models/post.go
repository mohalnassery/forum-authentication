package models

import "database/sql"

type Post struct {
	Title        string   `json:"title"`
	Body         string   `json:"body"`
	CreationDate string   `json:"creationdate"`
	Categories   []string `json:"categories"`
	PostID       int      `json:"postID"`
	AuthorID     int      `json:"authorID"`
	Image        string   `json:"image"`
}

type PostHome struct {
	PostID       int            `json:"post_id"`
	Title        string         `json:"title"`
	Body         string         `json:"body"`
	CreationDate string         `json:"creationdate"`
	Image        sql.NullString `json:"image"`
	Author       string         `json:"author"`
	Categories   []string       `json:"categories"`
	Likes        int            `json:"likes"`
	Dislikes     int            `json:"dislikes"`
	CommentCount int            `json:"commentCount"`
	UserLiked    bool           `json:"userLiked"`
	UserDisliked bool           `json:"userDisliked"`
	IsAuthor     bool           `json:"isAuthor"`
}

type PostCreate struct {
	Title      string   `json:"title"`
	Body       string   `json:"body"`
	Image      string   `json:"image"`
	Categories []string `json:"categories"`
	ErrorMsg   string   `json:"errormsg"`
}
