package repository

import (
	"context"

	"hello_world/internal/model"
)

// UserRepository defines the data-access contract for User entities.
type UserRepository interface {
	FindAll(ctx context.Context, offset, limit int) ([]model.User, error)
	FindByID(ctx context.Context, id uint) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uint) error
}
