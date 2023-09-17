package entities

import (
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	"walletban-api/internal/utils"
)

type User struct {
	gorm.Model
	Name        string
	Username    string `gorm:"unique;not null"`
	Email       string `gorm:"unique;not null"`
	PfpUrl      string
	Project     Project
	IsFirstTime bool `gorm:"default:true"`
}

func (u User) GetSignedJWT() string {
	claims := jwt.MapClaims{
		"username": u.Username,
		"uid":      u.ID,
		"pid":      u.Project.ID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, _ := token.SignedString([]byte(utils.JwtSecret))
	return t
}
