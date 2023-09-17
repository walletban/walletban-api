package entities

import "gorm.io/gorm"

type Project struct {
	gorm.Model
	UserID    uint
	Name      string
	TokenName string
	Consumers []Consumer
}
