package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"hello_world/internal/model"
	"hello_world/internal/repository"
)

const bcryptCost = 12

type userService struct {
	repo      repository.UserRepository
	jwtSecret []byte
	jwtExpiry time.Duration
}

// NewUserService creates a UserService with the provided repository and JWT settings.
func NewUserService(
	repo repository.UserRepository,
	jwtSecret string,
	jwtExpiryHours int,
) UserService {
	return &userService{
		repo:      repo,
		jwtSecret: []byte(jwtSecret),
		jwtExpiry: time.Duration(jwtExpiryHours) * time.Hour,
	}
}

func (s *userService) GetUsers(ctx context.Context, page, limit int) ([]model.User, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit
	return s.repo.FindAll(ctx, offset, limit)
}

func (s *userService) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *userService) Register(ctx context.Context, name, email, password string) (*model.User, error) {
	_, err := s.repo.FindByEmail(ctx, email)
	if err == nil {
		// FindByEmail returned nil error → user already exists.
		return nil, ErrEmailAlreadyTaken
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return nil, fmt.Errorf("checking existing email: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}

	user := &model.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
	}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("registering user: %w", err)
	}
	return user, nil
}

func (s *userService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", fmt.Errorf("finding user for login: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"sub": user.ID,
		"iat": now.Unix(),
		"exp": now.Add(s.jwtExpiry).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("signing JWT: %w", err)
	}
	return signed, nil
}