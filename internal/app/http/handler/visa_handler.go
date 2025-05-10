package handlers

import (
	"japa/internal/app/http/dto/request"
	"japa/internal/domain/entity"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	//"go.uber.org/zap"
)

// TYPES

// VisaApplication handler
type VisaApplicationHandler struct {
	Validator *validator.Validate
	Usecase   *usecase.UserUsecase
}

// METHODS

// Initialize VisaApplication handler
func NewVisaApplicationHandler(v *validator.Validate, vas *usecase.VisaApplicationService) *VisaApplicationHandler {
	return &VisaApplicationHandler{v, vas}
}

// Handler for visa submission
func (vah *VisaApplicationHandler) SubmitVisaApplication(c *fiber.Ctx) error {
	var req request.CreateVisaApplicationRequest
	_ = c.BodyParser(&req)

	// Validate fields
	if err := vah.Validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Convert req to GORM model
	visaApp, err := entity.ToVisaApplication(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ULID",
		})
	}

	// Save to DB (assuming db is a *gorm.DB)
	if err := vah.Service.Repo.DB.Create(&visaApp).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save application",
		})
	}

	// If successful
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Visa application submitted",
		"data":    visaApp,
	})
}
