package repositories

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BaseRepository[Entity any] interface {
	Create(ctx context.Context, entity Entity) (*Entity, error)
	Update(ctx context.Context, entity Entity) (*Entity, error)
	FindOne(ctx context.Context, entity Entity) (*Entity, error)
	FindAll() (*[]Entity, error)
	FindWithFilter(filter map[string]interface{}) (*[]Entity, error)
	UpdateWithFilter(entity Entity, filter map[string]interface{}) (*Entity, error)
	Delete(ctx context.Context, entity Entity, id uint) error
}

type baseRepository[Entity any] struct {
	Table *gorm.DB
}

func (b baseRepository[Entity]) UpdateWithFilter(entity Entity, filter map[string]interface{}) (*Entity, error) {
	result := b.Table.Model(&entity).Updates(filter)
	if result.Error != nil {
		return nil, result.Error
	}
	return &entity, nil
}

func (b baseRepository[Entity]) Update(ctx context.Context, entity Entity) (*Entity, error) {
	result := b.Table.WithContext(ctx).Model(&entity).Updates(&entity)
	if result.Error != nil {
		return nil, result.Error
	}
	return &entity, nil
}

func (b baseRepository[Entity]) FindWithFilter(filter map[string]interface{}) (*[]Entity, error) {
	var EntityHolder []Entity
	result := b.Table.Where(filter).Find(&EntityHolder)
	if result.Error != nil {
		return nil, result.Error
	}
	return &EntityHolder, nil
}

func (b baseRepository[Entity]) Create(ctx context.Context,
	entity Entity) (*Entity,
	error) {
	result := b.Table.WithContext(ctx).Create(&entity)
	if result.Error != nil {
		return nil, result.Error
	}
	return &entity, nil
}

func (b baseRepository[Entity]) FindOne(ctx context.Context,
	entity Entity) (*Entity, error) {
	result := b.Table.WithContext(ctx).Preload(clause.Associations).Where(&entity).First(&entity)
	if result.Error != nil {
		return nil, result.Error
	}
	return &entity, nil
}

func (b baseRepository[Entity]) FindAll() (*[]Entity,
	error) {
	var EntityHolder []Entity
	result := b.Table.Preload(clause.Associations).Find(&EntityHolder)
	if result.Error != nil {
		return nil, result.Error
	}
	return &EntityHolder, nil
}

func (b baseRepository[Entity]) Delete(ctx context.Context, entity Entity, id uint) error {
	result := b.Table.WithContext(ctx).Delete(&entity, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func NewBaseRepository[Entity any](collection *gorm.DB) BaseRepository[Entity] {
	return &baseRepository[Entity]{Table: collection}
}
