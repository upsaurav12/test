package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"hello_world/internal/service"
)

// AuthHandler handles authentication-related HTTP requests.
type AuthHandler struct {
	svc service.UserService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(svc service.UserService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// RegisterRequest is the request body for user registration.
type RegisterRequest struct {
	Name     string `json:"name"     binding:"required,min=2,max=100"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginRequest is the request body for user login.
type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// TokenResponse is returned on successful authentication.
type TokenResponse struct {
	Token string `json:"token"`
}

// Register creates a new user account.
//
//	POST /auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	user, err := h.svc.Register(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyTaken) {
			c.JSON(http.StatusConflict, ErrorResponse{Error: "email already registered"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "registration failed"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Login authenticates a user and returns a signed JWT token.
//
//	POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	token, err := h.svc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "login failed"})
		return
	}

	c.JSON(http.StatusOK, TokenResponse{Token: token})
}
