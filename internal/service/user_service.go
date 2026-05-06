package service

import (
	"hello_world/internal/model"
	"hello_world/internal/repository"
)

type UserService struct {
	Repo *repository.UserRepo
}

func NewUserService(repo *repository.UserRepo) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) GetUsers() ([]model.User, error) {
	return s.Repo.FindAll()
}