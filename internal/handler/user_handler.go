package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"hello_world/internal/repository"
	"hello_world/internal/service"
)

// UserHandler handles HTTP requests for user resources.
type UserHandler struct {
	svc service.UserService
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// GetUsers returns a paginated list of users.
//
//	GET /api/v1/users?page=1&limit=20
func (h *UserHandler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	users, err := h.svc.GetUsers(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to retrieve users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUser returns a single user by ID.
//
//	GET /api/v1/users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid user id"})
		return
	}

	user, err := h.svc.GetUserByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to retrieve user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func parseUintParam(c *gin.Context, key string) (uint, error) {
	v, err := strconv.ParseUint(c.Param(key), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(v), nil
}
