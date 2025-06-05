package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"forum/middleware"
	"forum/models"
	"forum/repository"
	"forum/utils"
)

// ActionHandler handles authenticated user actions like commenting and reacting

type ActionHandler struct {
	commentRepo  *repository.CommentRepository
	reactionRepo *repository.ReactionRepository
}

func NewActionHandler(commentRepo *repository.CommentRepository, reactionRepo *repository.ReactionRepository) *ActionHandler {
	return &ActionHandler{commentRepo: commentRepo, reactionRepo: reactionRepo}
}

// AddComment allows an authenticated user to add a comment to a post
func (h *ActionHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := middleware.GetCurrentUser(r)
	if user == nil {
		utils.ErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.CommentCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.PostID == "" || strings.TrimSpace(req.Content) == "" {
		utils.ErrorResponse(w, "Post ID and content required", http.StatusBadRequest)
		return
	}

	comment, err := h.commentRepo.AddComment(req.PostID, user.ID, req.Content)
	if err != nil {
		utils.ErrorResponse(w, "Failed to add comment", http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, comment, http.StatusCreated)
}

// React allows a user to react to a post or comment
func (h *ActionHandler) React(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := middleware.GetCurrentUser(r)
	if user == nil {
		utils.ErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.ReactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.TargetID == "" || (req.TargetType != "post" && req.TargetType != "comment") {
		utils.ErrorResponse(w, "Invalid target", http.StatusBadRequest)
		return
	}

	var postID, commentID *string
	if req.TargetType == "post" {
		postID = &req.TargetID
	} else {
		commentID = &req.TargetID
	}

	if err := h.reactionRepo.AddReaction(user.ID, req.ReactionType, postID, commentID); err != nil {
		utils.ErrorResponse(w, "Failed to react", http.StatusInternalServerError)
		return
	}

	var reactions []models.ReactionWithUser
	var err error
	if req.TargetType == "post" {
		reactions, err = h.reactionRepo.GetReactionsByPostWithUser(req.TargetID)
	} else {
		reactions, err = h.reactionRepo.GetReactionsByCommentWithUser(req.TargetID)
	}
	if err != nil {
		utils.ErrorResponse(w, "Failed to load reactions", http.StatusInternalServerError)
		return
	}
	utils.JSONResponse(w, reactions, http.StatusOK)
}
