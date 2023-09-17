package repositories

import (
	"gorm.io/gorm"
	"walletban-api/internal/entities"
)

type ConsumerRepository interface {
	BaseRepository[entities.Consumer]
}

type consumerRepository[Entity any] struct {
	Table *gorm.DB
	BaseRepository[Entity]
}

func NewConsumerRepository(table *gorm.DB) ConsumerRepository {
	br := NewBaseRepository[entities.Consumer](table)
	return &consumerRepository[entities.Consumer]{Table: table,
		BaseRepository: br}
}
