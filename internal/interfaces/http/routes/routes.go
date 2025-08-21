package routes

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"tasius.my.id/todolistapi/internal/config"
	"tasius.my.id/todolistapi/internal/interfaces/http/middleware"
	"tasius.my.id/todolistapi/internal/utils"
	"tasius.my.id/todolistapi/internal/utils/jwt"
)

type RoutesDependencies struct {
	Db          *gorm.DB
	RedisClient *redis.Client
	Config      *config.Config
	JWTManager  *jwt.TokenManager
}



func SetupRoutes(app *fiber.App, deps RoutesDependencies)  {
	api := app.Group("/api")

	api.Get("/health", func(c *fiber.Ctx) error {
		return utils.SuccessResponse(c, "OK", nil)
	})

	if deps.Config.AppEnv == "development" {
		api.Post("/decrypt", middleware.DecryptMiddleware(deps.Config.HybridEncryption.PrivateKeyPath), func(c *fiber.Ctx) error {
			var result map[string]interface{}
			if err := json.Unmarshal(c.Body(), &result); err != nil {
				return utils.ErrorResponse(c, fiber.StatusBadRequest, "failed to unmarshal JSON data")
			}

			return utils.SuccessResponse(c, "OK", result)
		})
	}

	SetupAuthRoutes(api, deps)
}