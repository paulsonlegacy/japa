package handlers

import (
	"japa/internal/app/http/dto/apperror"
	"japa/internal/app/http/dto/request"
	"japa/internal/app/http/dto/response"
	//"japa/internal/domain/entity"
	"japa/internal/domain/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	//"go.uber.org/zap"
)

// TYPES

// VisaApplication handler
type VisaHandler struct {
	Validator *validator.Validate
	Usecase   *usecase.VisaUsecase
}

// METHODS

// Initialize VisaApplication handler
func NewVisaHandler(v *validator.Validate, uc *usecase.VisaUsecase) *VisaHandler {
	return &VisaHandler{v, uc}
}

// Handler for visa submission
func (vh *VisaHandler) SubmitVisaApplication(c *fiber.Ctx) error {
	var reqBody request.CreateVisaApplicationRequest
	if err := reqBody.Bind(c, vh.Validator); err != nil {
		return response.BadRequest(c, apperror.NewValidationErr(err.Error()))
	}

	// Pass to usecase layer
	if err := vh.Usecase.CreateVisaApplication(c.Context(), reqBody); err != nil {
		return response.InternalServerError(c, apperror.NewServerErr(err.Error()))
	}

	// If application successful
	return response.Success(c, "Visa application was successful")
}
