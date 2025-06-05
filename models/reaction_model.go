package models

import "time"

type Reaction struct {
	UserID    string    `json:"user_id"`
	Type      int       `json:"reaction_type"` // 1 = Like, 2 = Love, etc.
	CommentID *string   `json:"comment_id,omitempty"`
	PostID    *string   `json:"post_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// ReactionWithUser is a reaction along with the username of the user who reacted
type ReactionWithUser struct {
	UserID       string    `json:"user_id"`
	Username     string    `json:"username"`
	ReactionType int       `json:"reaction_type"`
	PostID       *string   `json:"post_id,omitempty"`
	CommentID    *string   `json:"comment_id,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// ReactionRequest represents a request to react to a post or comment
// TargetType must be either "post" or "comment".
type ReactionRequest struct {
	TargetID     string `json:"target_id"`
	TargetType   string `json:"target_type"`
	ReactionType int    `json:"reaction_type"`
}
