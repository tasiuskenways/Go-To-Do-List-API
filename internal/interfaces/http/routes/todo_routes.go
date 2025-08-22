package routes

import (
	"github.com/gofiber/fiber/v2"
	"tasius.my.id/todolistapi/internal/application/services"
	"tasius.my.id/todolistapi/internal/infrastructure/repositories"
	"tasius.my.id/todolistapi/internal/interfaces/http/handlers"
	"tasius.my.id/todolistapi/internal/interfaces/http/middleware"
)

func SetupTodoRoutes(app fiber.Router, deps RoutesDependencies) {

	todoRepo := repositories.NewTodoRepository(deps.Db)
	todoService := services.NewTodoService(todoRepo, deps.RedisClient)
	todoHandler := handlers.NewTodoHandler(todoService)

	todoGroup := app.Group("/todos", middleware.AuthMiddleware(deps.JWTManager))
	todoGroup.Post("", todoHandler.CreateTodo)
	todoGroup.Post("/:id", todoHandler.UpdateTodo)
	todoGroup.Delete("/:id", todoHandler.DeleteTodo)
	todoGroup.Get("", todoHandler.GetAllTodos)
	todoGroup.Get("/:id", todoHandler.GetTodoByID)
}
