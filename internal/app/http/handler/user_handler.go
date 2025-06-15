// Fiber handlers for user routes
package handlers

import (
	"japa/internal/app/http/dto/request"
	"japa/internal/app/http/dto/response"
	"japa/internal/domain/usecase"
	//"japa/internal/domain/entity"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	//"go.uber.org/zap"
)

// TYPES

// User handler
type UserHandler struct {
	Validator *validator.Validate
	Usecase   *usecase.UserUsecase
}

// METHODS

// Initialize user handler
func NewUserHandler(v *validator.Validate, us *usecase.UserUsecase) *UserHandler {
	return &UserHandler{v, us}
}

// Register handler
func (uh *UserHandler) Register(c *fiber.Ctx) error {
	// Parse req body
	var body request.CreateUserRequest
	if err := body.Bind(c, uh.Validator); err != nil {
		return response.BadRequest(c)
	}

	// Registering user
	err := uh.Usecase.RegisterUser(c.Context(), body)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	// If registeration succeeded
	return response.UserRegisteredOK(c)
}

// Login handler
func (uh *UserHandler) Login(c *fiber.Ctx) error {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parsing incoming payload into user object
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}

	// Confirming user
	token, err := uh.Usecase.Login(body.Email, body.Password)
	if err != nil {
		return response.Unauthorized(c, err.Error())
	}

	// Return JWT token
	return response.Success(c, "login successful", map[string]any{"token": token})
}
