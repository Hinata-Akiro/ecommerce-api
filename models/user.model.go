package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"unique" json:"email"`
	Name     string `json:"name"`
	Password string `json:"-"`
	IsAdmin  bool   `json:"is_admin"`
}
