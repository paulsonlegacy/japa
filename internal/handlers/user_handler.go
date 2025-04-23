// Fiber handlers for user routes
package handlers

import (
	"japa/internal/models"
	"japa/internal/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	//"go.uber.org/zap"
)

// TYPES

// User handler
type UserHandler struct {
	Validator  *validator.Validate
	Service *services.UserService
}


// METHODS

// Initialize user handler
func NewUserHandler(v *validator.Validate, uc *services.UserService) *UserHandler {
	return &UserHandler{v, uc}
}


// Register handler
func (uh UserHandler) Register(c *fiber.Ctx) error {
	var user models.User

	// Parsing incoming payload into user object
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}

	// Registering user
	err := uh.Service.RegisterUser(c.Context(), &user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// If registeration succeeded
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success":"registeration successful"})
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
	token, err := uh.Service.Login(body.Email, body.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	// Return JWT token
	return c.JSON(fiber.Map{"token": token})
}