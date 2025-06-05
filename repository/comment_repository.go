package repository

import (
	"database/sql"
	"time"

	"forum/models"
	"forum/utils"
)

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) GetAllComments() ([]models.Comment, error) {
	rows, err := r.db.Query(`
		SELECT comment_id, post_id, user_id, content, created_at, updated_at 
		FROM comments ORDER BY created_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}

func (r *CommentRepository) AddComment(postID, userID, content string) (*models.Comment, error) {
	id := utils.GenerateUUID()
	created := time.Now()
	_, err := r.db.Exec(`INSERT INTO comments (comment_id, post_id, user_id, content, created_at) VALUES (?, ?, ?, ?, ?)`,
		id, postID, userID, content, created)
	if err != nil {
		return nil, err
	}
	return &models.Comment{ID: id, PostID: postID, UserID: userID, Content: content, CreatedAt: created}, nil
}
