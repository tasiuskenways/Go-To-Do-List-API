package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"tasius.my.id/todolistapi/internal/application/dto"
	"tasius.my.id/todolistapi/internal/domain/entities"
	"tasius.my.id/todolistapi/internal/domain/repositories"
	"tasius.my.id/todolistapi/internal/domain/services"
	"tasius.my.id/todolistapi/internal/interfaces/http/middleware"
)

const (
	todoCacheKey     = "todo:%s"
	todoListCacheKey = "todos"
	cacheExpiration  = 5 * time.Minute
	errTodoNotFound   = "Not found todo with id: %s"
	errUnauthorized   = "Unauthorized"
)

type todoService struct {
	todoRepo    repositories.TodoRepository
	redisClient *redis.Client
}

// CreateTodo implements services.TodoService.
func (t *todoService) CreateTodo(ctx *fiber.Ctx, todo dto.TodoDTO) (*dto.TodoResponse, error) {
	todoEntity := &entities.Todo{
		Title:       todo.Title,
		Description: todo.Description,
		UserID:      todo.UserID,
	}

	err := t.todoRepo.Create(ctx.Context(), todoEntity)
	if err != nil {
		return nil, err
	}

	// Invalidate the todo list cache in background since we added a new todo
	go func() {
		t.invalidateListCache(context.Background())
	}()

	return t.generateTodoResponse(todoEntity), nil
}

// DeleteTodo implements services.TodoService.
func (t *todoService) DeleteTodo(ctx *fiber.Ctx, id string) error {

	existingTodo, err := t.todoRepo.GetByID(ctx.Context(), id)
	if err != nil {
		return err
	}

	if existingTodo == nil {
		return fmt.Errorf(errTodoNotFound, id)
	}

	userId, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		return err
	}

	if existingTodo.UserID != userId {
		return errors.New(errUnauthorized)
	}

	// Invalidate both the specific todo and the list cache in background
	go func() {
		bgCtx := context.Background()
		t.invalidateCache(bgCtx, id)
		t.invalidateListCache(bgCtx)
	}()

	return t.todoRepo.Delete(ctx.Context(), id)
}

// GetAllTodos implements services.TodoService.
func (t *todoService) GetAllTodos(ctx *fiber.Ctx) ([]dto.TodoResponse, error) {

	userId, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Try to get from cache first
	cached, err := t.redisClient.Get(ctx.Context(), todoListCacheKey).Result()
	if err == nil {
		var cachedTodos []dto.TodoResponse
		if err := json.Unmarshal([]byte(cached), &cachedTodos); err == nil {
			if cachedTodos[0].UserID != userId {
				return nil, errors.New(errUnauthorized)
			}
			return cachedTodos, nil
		}
	}

	// If not in cache, get from database
	todos, err := t.todoRepo.GetAll(ctx.Context(), userId)
	if err != nil {
		return nil, err
	}

	var result []dto.TodoResponse
	for _, todo := range todos {
		result = append(result, *t.generateTodoResponse(&todo))
	}

	// Cache the result
	if jsonData, err := json.Marshal(result); err == nil {
		go t.redisClient.Set(ctx.Context(), todoListCacheKey, jsonData, cacheExpiration)
	}

	return result, nil
}

// GetTodoByID implements services.TodoService.
func (t *todoService) GetTodoByID(ctx *fiber.Ctx, id string) (*dto.TodoResponse, error) {
	cacheKey := fmt.Sprintf(todoCacheKey, id)

	userId, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Try to get from cache first
	cached, err := t.redisClient.Get(ctx.Context(), cacheKey).Result()
	if err == nil {
		var cachedTodo dto.TodoResponse
		if err := json.Unmarshal([]byte(cached), &cachedTodo); err == nil {
			// Check if the cached todo belongs to the current user
			if cachedTodo.UserID == userId {
				return &cachedTodo, nil
			}
			return nil, errors.New(errUnauthorized)
		}
	}

	// If not in cache or invalid cache, get from database
	todo, err := t.todoRepo.GetByID(ctx.Context(), id)
	if err != nil {
		return nil, err
	}

	if todo == nil {
		return nil, fmt.Errorf(errTodoNotFound, id)
	}

	// Check if the todo belongs to the current user
	if todo.UserID != userId {
		return nil, errors.New(errUnauthorized)
	}

	response := t.generateTodoResponse(todo)

	// Cache the result
	if jsonData, err := json.Marshal(response); err == nil {
		go t.redisClient.Set(ctx.Context(), cacheKey, jsonData, cacheExpiration)
	}

	return response, nil
}

// UpdateTodo implements services.TodoService.
func (t *todoService) UpdateTodo(ctx *fiber.Ctx, id string, todo dto.UpdateTodoRequest) (*dto.TodoResponse, error) {
	existingTodo, err := t.todoRepo.GetByID(ctx.Context(), id)
	if err != nil {
		return nil, err
	}

	if existingTodo == nil {
		return nil, fmt.Errorf(errTodoNotFound, id)
	}

	userId, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if existingTodo.UserID != userId {
		return nil, errors.New(errUnauthorized)
	}

	existingTodo.Title = todo.Title
	existingTodo.Description = todo.Description

	err = t.todoRepo.Update(ctx.Context(), id, existingTodo)
	if err != nil {
		return nil, err
	}

	// Invalidate both the specific todo and the list cache in background
	go func() {
		bgCtx := context.Background()
		t.invalidateCache(bgCtx, id)
		t.invalidateListCache(bgCtx)
	}()

	return t.generateTodoResponse(existingTodo), nil
}

func NewTodoService(todoRepo repositories.TodoRepository, redisClient *redis.Client) services.TodoService {
	return &todoService{
		todoRepo:    todoRepo,
		redisClient: redisClient,
	}
}

func (t *todoService) invalidateCache(ctx context.Context, id string) {
	cacheKey := fmt.Sprintf(todoCacheKey, id)
	t.redisClient.Del(ctx, cacheKey)
}

func (t *todoService) invalidateListCache(ctx context.Context) {
	t.redisClient.Del(ctx, todoListCacheKey)
}

func (t *todoService) generateTodoResponse(todo *entities.Todo) *dto.TodoResponse {
	return &dto.TodoResponse{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		UserID:      todo.UserID,
	}
}