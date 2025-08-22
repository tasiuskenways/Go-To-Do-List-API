package handlers

import (
	"github.com/gofiber/fiber/v2"
	"tasius.my.id/todolistapi/internal/application/dto"
	"tasius.my.id/todolistapi/internal/domain/services"
	"tasius.my.id/todolistapi/internal/interfaces/http/middleware"
	validators "tasius.my.id/todolistapi/internal/interfaces/validator"
	"tasius.my.id/todolistapi/internal/utils"
)

type TodoHandler struct {
	todoService services.TodoService
}

func NewTodoHandler(todoService services.TodoService) *TodoHandler {
	return &TodoHandler{
		todoService: todoService,
	}
}

func (h *TodoHandler) CreateTodo(c *fiber.Ctx) error {
	var req dto.TodoDTO
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, INVALID_REQUEST_BODY)
	}

	// Add user id from middleware
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, err.Error())
	}
	req.UserID = userID

	// Validate request
	if errors := validators.ValidateCreateTodo(&req); len(errors) > 0 {
		return utils.ValidationErrorResponse(c, errors)
	}

	// Create todo
	response, err := h.todoService.CreateTodo(c, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.CreatedResponse(c, "Todo created successfully", response)
}

func (h *TodoHandler) UpdateTodo(c *fiber.Ctx) error {
	var req dto.UpdateTodoRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, INVALID_REQUEST_BODY)
	}

	id := c.Params("id")
	if id == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "id is required")
	}

	// Validate request
	if errors := validators.ValidateUpdateTodo(&req); len(errors) > 0 {
		return utils.ValidationErrorResponse(c, errors)
	}

	// Update todo
	response, err := h.todoService.UpdateTodo(c, id, req)
	if err != nil {
		switch err.Error() {
		case "Unauthorized":
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, err.Error())
		case "Not found todo with id: " + id:
			return utils.ErrorResponse(c, fiber.StatusNotFound, err.Error())
		default:
			return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
		}
	}

	return utils.SuccessResponse(c, "Todo updated successfully", response)
}

func (h *TodoHandler) DeleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "id is required")
	}

	// Delete todo
	if err := h.todoService.DeleteTodo(c, id); err != nil {
		switch err.Error() {
		case "Unauthorized":
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, err.Error())
		case "Not found todo with id: " + id:
			return utils.ErrorResponse(c, fiber.StatusNotFound, err.Error())
		default:
			return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
		}
	}

	return utils.SuccessResponse(c, "Todo deleted successfully", nil)
}

func (h *TodoHandler) GetAllTodos(c *fiber.Ctx) error {
	todos, err := h.todoService.GetAllTodos(c)
	if err != nil {
		switch err.Error() {
		case "Unauthorized":
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, err.Error())
		default:
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
		}
	}

	return utils.SuccessResponse(c, "Todos fetched successfully", todos)
}

func (h *TodoHandler) GetTodoByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "id is required")
	}

	// Get todo
	todo, err := h.todoService.GetTodoByID(c, id)
	if err != nil {
		switch err.Error() {
		case "Unauthorized":
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, err.Error())
		case "Not found todo with id: " + id:
			return utils.ErrorResponse(c, fiber.StatusNotFound, err.Error())
		default:
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
		}
	}

	return utils.SuccessResponse(c, "Todo fetched successfully", todo)
}



