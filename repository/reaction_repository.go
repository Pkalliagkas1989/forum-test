package repository

import (
	"database/sql"
	"time"

	"forum/models"
)

type ReactionRepository struct {
	db *sql.DB
}

func NewReactionRepository(db *sql.DB) *ReactionRepository {
	return &ReactionRepository{db: db}
}

func (r *ReactionRepository) GetAllReactions() ([]models.Reaction, error) {
	rows, err := r.db.Query(`
		SELECT user_id, reaction_type, comment_id, post_id, created_at 
		FROM reactions`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reactions []models.Reaction
	for rows.Next() {
		var react models.Reaction
		err := rows.Scan(&react.UserID, &react.Type, &react.CommentID, &react.PostID, &react.CreatedAt)
		if err != nil {
			return nil, err
		}
		reactions = append(reactions, react)
	}

	return reactions, nil
}

func (r *ReactionRepository) AddReaction(userID string, reactionType int, postID, commentID *string) error {
	_, err := r.db.Exec(`INSERT OR REPLACE INTO reactions (user_id, reaction_type, comment_id, post_id, created_at) VALUES (?, ?, ?, ?, ?)`,
		userID, reactionType, commentID, postID, time.Now())
	return err
}
