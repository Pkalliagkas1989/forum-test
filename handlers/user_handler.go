package handlers

import (
	"net/http"

	"forum/middleware"
	"forum/models"
	"forum/repository"
	"forum/utils"
)

// UserHandler serves data for authenticated users.
type UserHandler struct {
	categoryRepo *repository.CategoryRepository
	postRepo     *repository.PostRepository
	commentRepo  *repository.CommentRepository
	reactionRepo *repository.ReactionRepository
}

// NewUserHandler creates a UserHandler.
func NewUserHandler(categoryRepo *repository.CategoryRepository, postRepo *repository.PostRepository, commentRepo *repository.CommentRepository, reactionRepo *repository.ReactionRepository) *UserHandler {
	return &UserHandler{
		categoryRepo: categoryRepo,
		postRepo:     postRepo,
		commentRepo:  commentRepo,
		reactionRepo: reactionRepo,
	}
}

// UserResponse embeds forum data and includes the authenticated user.
type UserResponse struct {
	User models.User `json:"user"`
	GuestResponse
}

// GetUserData returns forum data for an authenticated user, verifying the cookie.
func (h *UserHandler) GetUserData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := middleware.GetCurrentUser(r)
	if user == nil {
		utils.ErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	categories, err := h.categoryRepo.GetAll()
	if err != nil {
		utils.ErrorResponse(w, "Failed to load categories", http.StatusInternalServerError)
		return
	}

	var guestResp GuestResponse
	for _, cat := range categories {
		catResp := CategoryResponse{
			ID:    cat.ID,
			Name:  cat.Name,
			Posts: []PostResponse{},
		}

		posts, err := h.postRepo.GetPostsByCategoryWithUser(cat.ID)
		if err != nil {
			utils.ErrorResponse(w, "Failed to load posts", http.StatusInternalServerError)
			return
		}

		for _, post := range posts {
			postResp := PostResponse{
				ID:           post.ID,
				UserID:       post.UserID,
				Username:     post.Username,
				CategoryID:   post.CategoryID,
				CategoryName: cat.Name,
				Title:        post.Title,
				Content:      post.Content,
				CreatedAt:    post.CreatedAt,
				Comments:     []CommentResponse{},
				Reactions:    []ReactionResponse{},
			}

			comments, err := h.commentRepo.GetCommentsByPostWithUser(post.ID)
			if err != nil {
				utils.ErrorResponse(w, "Failed to load comments", http.StatusInternalServerError)
				return
			}

			for _, comment := range comments {
				commentResp := CommentResponse{
					ID:        comment.ID,
					UserID:    comment.UserID,
					Username:  comment.Username,
					Content:   comment.Content,
					CreatedAt: comment.CreatedAt,
					Reactions: []ReactionResponse{},
				}

				reactions, err := h.reactionRepo.GetReactionsByCommentWithUser(comment.ID)
				if err != nil {
					utils.ErrorResponse(w, "Failed to load reactions", http.StatusInternalServerError)
					return
				}
				for _, reaction := range reactions {
					commentResp.Reactions = append(commentResp.Reactions, ReactionResponse{
						UserID:       reaction.UserID,
						Username:     reaction.Username,
						ReactionType: reaction.ReactionType,
						CreatedAt:    reaction.CreatedAt,
					})
				}

				postResp.Comments = append(postResp.Comments, commentResp)
			}

			reactions, err := h.reactionRepo.GetReactionsByPostWithUser(post.ID)
			if err != nil {
				utils.ErrorResponse(w, "Failed to load reactions", http.StatusInternalServerError)
				return
			}
			for _, reaction := range reactions {
				postResp.Reactions = append(postResp.Reactions, ReactionResponse{
					UserID:       reaction.UserID,
					Username:     reaction.Username,
					ReactionType: reaction.ReactionType,
					CreatedAt:    reaction.CreatedAt,
				})
			}

			catResp.Posts = append(catResp.Posts, postResp)
		}

		guestResp.Categories = append(guestResp.Categories, catResp)
	}

	resp := UserResponse{User: *user, GuestResponse: guestResp}
	utils.JSONResponse(w, resp, http.StatusOK)
}
