package response

import (
	"japa/internal/app/http/dto/apperror"

	"github.com/gofiber/fiber/v2"
)


//  SUCCESS RESPONSES

func Success(c *fiber.Ctx, message string, data ...map[string]any) error {
	if message == "" {
		message = "ok"
	}

	var payload map[string]any
	if len(data) > 0 {
		payload = data[0] // use firstpayload
	} else {
		payload = map[string]any{} // empty payload
	}

	return c.Status(fiber.StatusOK).JSON(map[string]any{
		"message": message,
		"status": "success",
		"data": payload,
		"error": map[string]any{},
	})
}

func Created(c *fiber.Ctx, data ...map[string]any) error {
	var payload map[string]any
	if len(data) > 0 {
		payload = data[0]
	} else {
		payload = map[string]any{}
	}

	return c.Status(fiber.StatusCreated).JSON(map[string]any{
		"message": "created",
		"status":  "success",
		"data":    payload,
		"error": map[string]any{},
	})
}


//  ERROR RESPONSES

func BadRequest(c *fiber.Ctx, appErr *apperror.AppError) error {
	if appErr == nil {
		appErr = apperror.New("BAD_REQUEST", "Bad request", "Bad request received")
	}

	return c.Status(fiber.StatusBadRequest).JSON(map[string]any{
		"message": appErr.Message,
		"status":  "failed",
		"error": fiber.Map{
			"message": appErr.Message,
			"details": appErr.Details,
			"code":    appErr.Code,
		},
	})
}

func Unauthorized(c *fiber.Ctx, appErr *apperror.AppError) error {
	if appErr == nil {
		appErr = apperror.New("UNAUTHORIZED", "Unauthorized", "Unauthorized request")
	}

	return c.Status(fiber.StatusUnauthorized).JSON(map[string]any{
		"message": appErr.Message,
		"status":  "failed",
		"error": fiber.Map{
			"message": appErr.Message,
			"details": appErr.Details,
			"code":    appErr.Code,
		},
	})
}

func Forbidden(c *fiber.Ctx, appErr *apperror.AppError) error {
	if appErr == nil {
		appErr = apperror.New("FORBIDDEN", "Forbidden", "Forbidden request")
	}

	return c.Status(fiber.StatusForbidden).JSON(map[string]any{
		"message": appErr.Message,
		"status":  "failed",
		"error": fiber.Map{
			"message": appErr.Message,
			"details": appErr.Details,
			"code":    appErr.Code,
		},
	})
}

func NotFound(c *fiber.Ctx, appErr *apperror.AppError) error {
	if appErr == nil {
		appErr = apperror.New("NOT_FOUND", "Not found", "Resource not found")
	}

	return c.Status(fiber.StatusNotFound).JSON(map[string]any{
		"message": appErr.Message,
		"status":  "failed",
		"error": fiber.Map{
			"message": appErr.Message,
			"details": appErr.Details,
			"code":    appErr.Code,
		},
	})
}

func Conflict(c *fiber.Ctx, appErr *apperror.AppError) error {
	if appErr == nil {
		appErr = apperror.New("CONFLICT", "Conflict", "Conflicted request")
	}

	return c.Status(fiber.StatusConflict).JSON(map[string]any{
		"message": appErr.Message,
		"status":  "failed",
		"error": fiber.Map{
			"message": appErr.Message,
			"details": appErr.Details,
			"code":    appErr.Code,
		},
		
	})
}

func Unprocessable(c *fiber.Ctx, appErr *apperror.AppError) error {
	if appErr == nil {
		appErr = apperror.New("UNPROCESSABLE", "Unprocessable", "Unprocessable request")
	}

	return c.Status(fiber.StatusUnprocessableEntity).JSON(map[string]any{
		"message": appErr.Message,
		"status":  "failed",
		"error": fiber.Map{
			"message": appErr.Message,
			"details": appErr.Details,
			"code":    appErr.Code,
		},
	})
}


func InternalServerError(c *fiber.Ctx, appErr *apperror.AppError) error {
	if appErr == nil {
		appErr = apperror.New("INTERNAL_SERVER_ERROR", "Internal Server Error", "An error occured while processing request")
	}

	return c.Status(fiber.StatusInternalServerError).JSON(
		map[string]any{
			"message": appErr.Message,
			"status":  "failed",
			"error": fiber.Map{
			"message": appErr.Message,
			"details": appErr.Details,
			"code":    appErr.Code,
		},
		},
	)
}