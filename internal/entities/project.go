package entities

import "gorm.io/gorm"

type Project struct {
	gorm.Model
	UserID       uint
	Name         string
	TokenName    string
	ClientId     string
	ClientSecret string
	ApiKey       string
	Consumers    []Consumer
}
