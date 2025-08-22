package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"tasius.my.id/todolistapi/internal/domain/entities"
	"tasius.my.id/todolistapi/internal/domain/repositories"
)

const (
	errIDRequired      = "id is required"
	errTodoNil         = "todo cannot be nil"
	errUserIDRequired  = "user_id is required"
	errInvalidIDFormat = "invalid id format: %v"
)

type todoRepository struct {
	db *gorm.DB
}

// Create implements repositories.TodoRepository.
func (t *todoRepository) Create(ctx context.Context, todo *entities.Todo) error {
	if todo == nil {
		return errors.New(errTodoNil)
	}

	if todo.UserID == "" {
		return errors.New(errUserIDRequired)
	}

	if _, err := uuid.Parse(todo.UserID); err != nil {
		return fmt.Errorf("invalid user_id format: %v", err)
	}

	if err := t.db.WithContext(ctx).Create(todo).Error; err != nil {
		return fmt.Errorf("failed to create todo: %w", err)
	}

	return nil
}

// Delete implements repositories.TodoRepository.
func (t *todoRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New(errIDRequired)
	}

	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf(errInvalidIDFormat, err)
	}

	result := t.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.Todo{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete todo: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// GetAll implements repositories.TodoRepository.
func (t *todoRepository) GetAll(ctx context.Context, userID string) ([]entities.Todo, error) {
	var todos []entities.Todo
	if err := t.db.WithContext(ctx).Where("user_id = ?", userID).Find(&todos).Error; err != nil {
		return nil, err
	}
	return todos, nil
}

// GetByID implements repositories.TodoRepository.
func (t *todoRepository) GetByID(ctx context.Context, id string) (*entities.Todo, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("invalid id format: %v", err)
	}

	var todo entities.Todo
	if err := t.db.WithContext(ctx).Where("id = ?", id).First(&todo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}

	return &todo, nil
}

// Update implements repositories.TodoRepository.
func (t *todoRepository) Update(ctx context.Context, id string, todo *entities.Todo) error {
	if id == "" {
		return errors.New(errIDRequired)
	}

	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf(errInvalidIDFormat, err)
	}

	if todo == nil {
		return errors.New(errTodoNil)
	}

	result := t.db.WithContext(ctx).Model(&entities.Todo{}).Where("id = ?", id).Updates(todo)
	if result.Error != nil {
		return fmt.Errorf("failed to update todo: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func NewTodoRepository(db *gorm.DB) repositories.TodoRepository {
	return &todoRepository{
		db: db,
	}
}
