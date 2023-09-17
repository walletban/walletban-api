package services

import (
	"gorm.io/gorm"
	"walletban-api/internal/repositories"
)

type ApplicationService struct {
	db                 *gorm.DB
	UserRepository     repositories.UserRepository
	ProjectRepository  repositories.ProjectRepository
	ConsumerRepository repositories.ConsumerRepository
}

func NewService(
	db *gorm.DB,
	UserRepository repositories.UserRepository,
	ProjectRepository repositories.ProjectRepository,
	ConsumerRepository repositories.ConsumerRepository,
) ApplicationService {
	return ApplicationService{
		db,
		UserRepository,
		ProjectRepository,
		ConsumerRepository,
	}
}
