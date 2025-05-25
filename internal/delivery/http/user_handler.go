// Package http implements the HTTP delivery layer of the application
package http

import (
	"encoding/json"
	"net/http"

	"github.com/tyobaskara/jeki-backend/internal/domain/entity"
	"github.com/tyobaskara/jeki-backend/internal/usecase"
)

// UserHandler handles HTTP requests related to user operations
// It translates HTTP requests into usecase calls and formats the responses
type UserHandler struct {
	userUsecase *usecase.UserUsecase // Business logic layer for user operations
}

// NewUserHandler creates a new instance of UserHandler
// Parameters:
//   - userUsecase: pointer to UserUsecase instance
// Returns:
//   - *UserHandler: new instance of UserHandler
func NewUserHandler(userUsecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

// CreateUser handles HTTP POST requests to create a new user
// It decodes the request body into a User entity and calls the usecase
// Parameters:
//   - w: http.ResponseWriter for sending the response
//   - r: *http.Request containing the request data
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	// Decode the request body into the user struct
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call the usecase to create the user
	if err := h.userUsecase.CreateUser(r.Context(), &user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GetUser handles HTTP GET requests to retrieve a user
// It extracts the user ID from the request and calls the usecase
// Parameters:
//   - w: http.ResponseWriter for sending the response
//   - r: *http.Request containing the request data
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// TODO: Extract user ID from request
	id := uint(1) // Placeholder

	// Call the usecase to get the user
	user, err := h.userUsecase.GetUserByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send success response
	json.NewEncoder(w).Encode(user)
} 