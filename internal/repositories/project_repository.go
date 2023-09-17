package repositories

import (
	"gorm.io/gorm"
	"walletban-api/internal/entities"
)

type ProjectRepository interface {
	BaseRepository[entities.Project]
}

type projectRepository[Entity any] struct {
	Table *gorm.DB
	BaseRepository[Entity]
}

func NewProjectRepository(table *gorm.DB) ProjectRepository {
	br := NewBaseRepository[entities.Project](table)
	return &projectRepository[entities.Project]{Table: table,
		BaseRepository: br}
}
