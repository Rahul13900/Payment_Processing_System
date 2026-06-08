package handler

import (
	"errors"
	"net/http"
	"user_service/internal/models"
	"user_service/internal/service"
	"user_service/internal/usererrors"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// Register handles POST /auth/register
func (h *UserHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	// parse and validate the request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Validation Error",
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	// Call Service
	user, err := h.service.Register(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, usererrors.ErrInvalidEmail):
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   "Invalid Email",
				Message: "Email is invalid",
			})
			return
		case errors.Is(err, usererrors.ErrWeakPassword):
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   "WEAK_PASSWORD",
				Message: "Password must be at least 8 characters",
			})
			return
		case errors.Is(err, usererrors.ErrEmailExists):
			c.JSON(http.StatusConflict, models.ErrorResponse{
				Error:   "EMAIL_EXISTS",
				Message: "Email already registered",
			})
			return
		default:
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "INTERNAL_ERROR",
				Message: "Failed to register user",
			})
			return
		}
	}
	// Success Response
	resp := models.AuthResponse{
		UserID: user.ID,
		Email:  user.Email,
		Name:   user.Name,
	}
	c.JSON(http.StatusCreated, resp)
}

// Login handles POST /auth/login
func (h *UserHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	// Parse and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "VALIDATION_ERROR",
			Message: "Invalid request",
		})
		return
	}

	// Call service
	user, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		// Don't reveal which field is wrong
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "INVALID_CREDENTIALS",
			Message: "Invalid email or password",
		})
		return
	}

	// Success response
	resp := models.AuthResponse{
		UserID: user.ID,
		Email:  user.Email,
		Name:   user.Name,
	}

	c.JSON(http.StatusOK, resp)
}

// GetProfile handles GET /users/me
func (h *UserHandler) GetProfile(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "UNAUTHORIZED",
			Message: "User ID not found in context",
		})
		return
	}

	userID, ok := userIDInterface.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "INTERNAL_ERROR",
			Message: "Invalid user ID type",
		})
		return
	}

	// Get user
	user, err := h.service.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "NOT_FOUND",
			Message: "User not found",
		})
		return
	}

	// Response
	resp := models.UserProfileResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
	}

	c.JSON(http.StatusOK, resp)
}
