package models

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"regexp"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique" json:"email"`
	Name     string `json:"name"`
	Password string `json:"-"`
	IsAdmin  bool   `json:"is_admin"`
}


var PasswordRegex = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[!@#\$%\^&\*])[A-Za-z\d!@#\$%\^&\*]{8,}$`

func PasswordValidation(fl validator.FieldLevel) bool {
    password := fl.Field().String()
    valid, _ := ValidatePassword(password)
    return valid
}

func ValidatePassword(password string) (bool, string) {
    if len(password) < 8 {
        return false, "Password must be at least 8 characters long"
    }
    if match, _ := regexp.MatchString(`[A-Za-z]`, password); !match {
        return false, "Password must contain at least one letter"
    }
    if match, _ := regexp.MatchString(`\d`, password); !match {
        return false, "Password must contain at least one digit"
    }
    if match, _ := regexp.MatchString(`[!@#\$%\^&\*]`, password); !match {
        return false, "Password must contain at least one special character (e.g., !@#$%^&*)"
    }
    return true, ""
}