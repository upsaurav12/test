package handler

import (
	"hello_world/internal/service"

	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{Service: s}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	users, _ := h.Service.GetUsers()
	c.JSON(http.StatusOK, users)
}
