package services

import (
	"github.com/gofiber/fiber/v2"
	"tasius.my.id/todolistapi/internal/application/dto"
)

type TodoService interface {
	CreateTodo(ctx *fiber.Ctx, todo dto.TodoDTO) (*dto.TodoResponse, error)
	GetAllTodos(ctx *fiber.Ctx) ([]dto.TodoResponse, error)
	GetTodoByID(ctx *fiber.Ctx, id string) (*dto.TodoResponse, error)
	UpdateTodo(ctx *fiber.Ctx, id string, todo dto.UpdateTodoRequest) (*dto.TodoResponse, error)
	DeleteTodo(ctx *fiber.Ctx, id string) error
}