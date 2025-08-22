package repositories

import (
	"context"

	"tasius.my.id/todolistapi/internal/domain/entities"
)

type TodoRepository interface {
	Create(ctx context.Context, todo *entities.Todo) error
	GetAll(ctx context.Context, userID string) ([]entities.Todo, error)
	GetByID(ctx context.Context, id string) (*entities.Todo, error)
	Update(ctx context.Context, id string, todo *entities.Todo) error
	Delete(ctx context.Context, id string) error
}