package models

type CommentHome struct {
	ID           int    `json:"id"`
	Body         string `json:"body"`
	CreationDate string `json:"creationDate"`
	Author       string `json:"author"`
	Likes        int    `json:"likes"`
	Dislikes     int    `json:"dislikes"`
	UserLiked    bool   `json:"userLiked"`
	UserDisliked bool   `json:"userDisliked"`
	IsAuthor     bool   `json:"isAuthor"`
	PostID       int    `json:"postID"`
}

type CommentRequest struct {
	Body     string `json:"body"`
	Likes    int    `json:"likes"`
	Dislikes int    `json:"dislikes"`
}
