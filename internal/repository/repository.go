package repository

import (
	"context"
)

// Repository define el contrato abstracto de persistencia
type Repository interface {
	GetByID(ctx context.Context, id int64) (interface{}, error)
	Create(ctx context.Context, entity interface{}) error
	Update(ctx context.Context, entity interface{}) error
	Delete(ctx context.Context, id int64) error
}
