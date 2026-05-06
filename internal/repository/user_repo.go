package repository

import (
	"hello_world/internal/model"
	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (r *UserRepo) FindAll() ([]model.User, error) {
	var users []model.User
	err := r.DB.Find(&users).Error
	return users, err
}

func (r *UserRepo) Create(user*model.User) error {
	return r.DB.Create(user).Error
}
