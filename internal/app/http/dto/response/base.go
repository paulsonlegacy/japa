package response

import (
	"github.com/gofiber/fiber/v2"
)


func Success(c *fiber.Ctx, messages ...string) error {
	msg := "success"
	if len(messages) > 0 && messages[0] != "" {
		msg = messages[0]
	}

	return c.Status(fiber.StatusOK).JSON(
		map[string]string{
			"message": msg,
		},
	)
}

func Created(c *fiber.Ctx, messages ...string) error {
	msg := "created"
	if len(messages) > 0 && messages[0] != "" {
		msg = messages[0]
	}

	return c.Status(fiber.StatusCreated).JSON(
		map[string]string{
			"message": msg,
		},
	)
}

func BadRequest(c *fiber.Ctx, messages ...string) error {
	msg := "bad request"
	if len(messages) > 0 && messages[0] != "" {
		msg = messages[0]
	}

	return c.Status(fiber.StatusBadRequest).JSON(
		map[string]string{
			"message": msg,
		},
	)
}

func Unauthorized(c *fiber.Ctx, messages ...string) error {
	msg := "unauthorized"
	if len(messages) > 0 && messages[0] != "" {
		msg = messages[0]
	}

	return c.Status(fiber.StatusUnauthorized).JSON(
		map[string]string{
			"message": msg,
		},
	)
}

func Forbidden(c *fiber.Ctx, messages ...string) error {
	msg := "forbidden"
	if len(messages) > 0 && messages[0] != "" {
		msg = messages[0]
	}
	
	return c.Status(fiber.StatusForbidden).JSON(
		map[string]string{
			"message": msg,
		},
	)
}

func NotFound(c *fiber.Ctx, messages ...string) error {
	msg := "not found"
	if len(messages) > 0 && messages[0] != "" {
		msg = messages[0]
	}
	
	return c.Status(fiber.StatusNotFound).JSON(
		map[string]string{
			"message": msg,
		},
	)
}

func Conflict(c *fiber.Ctx, messages ...string) error {
	msg := "conflict"
	if len(messages) > 0 && messages[0] != "" {
		msg = messages[0]
	}
	
	return c.Status(fiber.StatusConflict).JSON(
		map[string]string{
			"message": msg,
		},
	)
}

func Unprocessable(c *fiber.Ctx, messages ...string) error {
	msg := "unprocessable"
	if len(messages) > 0 && messages[0] != "" {
		msg = messages[0]
	}
	
	return c.Status(fiber.StatusUnprocessableEntity).JSON(
		map[string]string{
			"message": msg,
		},
	)
}

func InternalServerError(c *fiber.Ctx, messages ...string) error {
	msg := "internal server error"
	if len(messages) > 0 && messages[0] != "" {
		msg = messages[0]
	}
	
	return c.Status(fiber.StatusInternalServerError).JSON(
		map[string]string{
			"message": msg,
		},
	)
}