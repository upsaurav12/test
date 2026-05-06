package model

import (
    "gorm.io/gorm"
)

type User struct {
    gorm.Model
    ID int `json:"id"`
    Name string `json:"name"`
}

func init() {
    Register(&User{})
}