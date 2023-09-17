package entities

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string
	Username string
	Email    string
	PfpUrl   string
	Project  Project
}
