package routes

import (
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"tasius.my.id/todolistapi/internal/domain/entities"
	"tasius.my.id/todolistapi/internal/middleware"
	"tasius.my.id/todolistapi/internal/utils/jwt"
	"tasius.my.id/todolistapi/internal/utils/password"
	"tasius.my.id/todolistapi/internal/utils/response"
)

var validate = validator.New()

type AuthHandler struct {
	DB          *gorm.DB
	JWTManager  *jwt.TokenManager
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

func SetupAuthRoutes(api fiber.Router, deps RoutesDependencies) {
	handler := &AuthHandler{
		DB:          deps.Db,
		JWTManager:  deps.JWTManager,
	}

	auth := api.Group("/auth")
	{
		auth.Post("/register", handler.Register)
		auth.Post("/login", handler.Login)
		auth.Post("/refresh", handler.RefreshToken)
		auth.Post("/logout", middleware.AuthMiddleware(handler.JWTManager), handler.Logout)
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	req := new(RegisterRequest)
	if err := c.BodyParser(req); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		return response.ValidationErrorResponse(c, err)
	}

	// Check if user already exists
	var existingUser entities.User
	if err := h.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return response.ErrorResponse(c, http.StatusConflict, "Email already registered")
	}

	// Hash password
	hashedPassword, err := password.HashPassword(req.Password)
	if err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user")
	}

	// Create user
	user := &entities.User{
		Email:    req.Email,
		Password: hashedPassword,
		Name:     req.Name,
	}

	if err := h.DB.Create(user).Error; err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user")
	}

	// Generate tokens
	tokens, err := h.JWTManager.GenerateTokenPair(user)
	if err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate tokens")
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"data": AuthResponse{
			AccessToken:  tokens[jwt.AccessToken],
			RefreshToken: tokens[jwt.RefreshToken],
			ExpiresIn:    15 * 60, // 15 minutes in seconds
		},
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	req := new(LoginRequest)
	if err := c.BodyParser(req); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		return response.ValidationErrorResponse(c, err)
	}

	// Find user by email
	var user entities.User
	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrorResponse(c, http.StatusUnauthorized, "Invalid credentials")
		}
		return response.ErrorResponse(c, http.StatusInternalServerError, "Failed to login")
	}

	// Verify password
	if err := password.CheckPassword(req.Password, user.Password); err != nil {
		return response.ErrorResponse(c, http.StatusUnauthorized, "Invalid credentials")
	}

	// Generate tokens
	tokens, err := h.JWTManager.GenerateTokenPair(&user)
	if err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate tokens")
	}

	return response.SuccessResponse(c, "Login successful", AuthResponse{
		AccessToken:  tokens[jwt.AccessToken],
		RefreshToken: tokens[jwt.RefreshToken],
		ExpiresIn:    15 * 60, // 15 minutes in seconds
	})
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return response.ErrorResponse(c, http.StatusBadRequest, "Refresh token is required")
	}

	// Remove "Bearer " prefix if present
	token = strings.TrimPrefix(token, "Bearer ")

	tokens, err := h.JWTManager.RefreshToken(token)
	if err != nil {
		return response.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired refresh token")
	}

	return response.SuccessResponse(c, "Token refreshed successfully", AuthResponse{
		AccessToken:  tokens[jwt.AccessToken],
		RefreshToken: tokens[jwt.RefreshToken],
		ExpiresIn:    15 * 60, // 15 minutes in seconds
	})
}

func (h *AuthHandler) GetCurrentUser(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return response.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
	}

	var user entities.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		return response.ErrorResponse(c, http.StatusNotFound, "User not found")
	}

	// Don't return password hash
	user.Password = ""

	return response.SuccessResponse(c, "User retrieved successfully", user)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return response.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
	}

	// Invalidate all tokens for this user
	err = h.JWTManager.Logout(userID)
	if err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, "Failed to logout")
	}

	return response.SuccessResponse(c, "Successfully logged out", nil)
}
