package service

import (
	"context"

	"hello_world/internal/model"
)

// UserService defines the business-logic contract for user operations.
type UserService interface {
	// GetUsers returns a paginated list of users.
	GetUsers(ctx context.Context, page, limit int) ([]model.User, error)
	// GetUserByID returns a single user by primary key.
	GetUserByID(ctx context.Context, id uint) (*model.User, error)
	// Register creates a new user account with a hashed password.
	Register(ctx context.Context, name, email, password string) (*model.User, error)
	// Login validates credentials and returns a signed JWT token.
	Login(ctx context.Context, email, password string) (string, error)
}
