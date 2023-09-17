package repositories

import (
	"gorm.io/gorm"
	"walletban-api/internal/entities"
)

type UserRepository interface {
	BaseRepository[entities.User]
}

type userRepository[Entity any] struct {
	Table *gorm.DB
	BaseRepository[Entity]
}

func NewUserRepository(table *gorm.DB) UserRepository {
	br := NewBaseRepository[entities.User](table)
	return &userRepository[entities.User]{Table: table,
		BaseRepository: br}
}
