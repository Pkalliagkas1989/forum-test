package handlers

import (
	"encoding/json"
	"net/http"

	"forum/middleware"
	"forum/models"
	"forum/repository"
	"forum/utils"
)

// PostHandler handles post related endpoints
type PostHandler struct {
	PostRepo *repository.PostRepository
}

// NewPostHandler creates a new PostHandler
func NewPostHandler(repo *repository.PostRepository) *PostHandler {
	return &PostHandler{PostRepo: repo}
}

// CreatePost creates a new post for the authenticated user
func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := middleware.GetCurrentUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		CategoryID int    `json:"category_id"`
		Title      string `json:"title"`
		Content    string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.CategoryID == 0 || req.Title == "" || req.Content == "" {
		utils.ErrorResponse(w, "Category, title and content are required", http.StatusBadRequest)
		return
	}

	post := models.Post{
		UserID:     user.ID,
		CategoryID: req.CategoryID,
		Title:      req.Title,
		Content:    req.Content,
	}

	created, err := h.PostRepo.Create(post)
	if err != nil {
		utils.ErrorResponse(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, created, http.StatusCreated)
}
