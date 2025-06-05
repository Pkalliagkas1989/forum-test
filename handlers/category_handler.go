package handlers

import (
	"net/http"

	"forum/repository"
	"forum/utils"
)

// CategoryHandler serves category data.
type CategoryHandler struct {
	categoryRepo *repository.CategoryRepository
}

// NewCategoryHandler creates a CategoryHandler.
func NewCategoryHandler(categoryRepo *repository.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{categoryRepo: categoryRepo}
}

// GetCategories returns all forum categories.
func (h *CategoryHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	categories, err := h.categoryRepo.GetAll()
	if err != nil {
		utils.ErrorResponse(w, "Failed to load categories", http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, categories, http.StatusOK)
}
