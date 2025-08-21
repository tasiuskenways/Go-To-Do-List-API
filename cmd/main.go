package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"tasius.my.id/todolistapi/internal/config"
	"tasius.my.id/todolistapi/internal/infrastructure/db"
	"tasius.my.id/todolistapi/internal/interfaces/routes"
	"tasius.my.id/todolistapi/internal/utils/jwt"
)

func main()  {

	cfg := config.Load()

	postgres, err := db.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	redis, err := db.NewRedisConnection(cfg)
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	// Initialize JWT manager
	jwtManager, err := jwt.NewTokenManager(&cfg.JWT, redis)
	if err != nil {
		log.Fatal("Failed to initialize JWT manager:", err)
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error":   true,
				"message": err.Error(),
			})
		},
	})

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	routes.SetupRoutes(app, routes.RoutesDependencies{
		Db:         postgres,
		RedisClient: redis,
		Config:     cfg,
		JWTManager: jwtManager,
	})
	
	log.Printf("Server starting on port %s", cfg.AppPort)
	if err := app.Listen(":" + cfg.AppPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
	
}