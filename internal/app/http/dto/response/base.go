package response

import (
	"github.com/gofiber/fiber/v2"
)


//  SUCCESS RESPONSES

// data ...map[string]any is to make data/payload optional
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
		"status":  "success",
		"data":    payload,
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
	})
}


//  ERROR RESPONSES

// message ...string is to make message optional
func BadRequest(c *fiber.Ctx, message ...string) error {
	var msg string
	if len(message) > 0 {
		msg = message[0]
	} else {
		msg = "bad request"
	}

	return c.Status(fiber.StatusBadRequest).JSON(map[string]any{
		"message": msg,
		"status":  "failed",
	})
}

func Unauthorized(c *fiber.Ctx, message ...string) error {
	var msg string
	if len(message) > 0 {
		msg = message[0]
	} else {
		msg = "unauthorized"
	}

	return c.Status(fiber.StatusUnauthorized).JSON(map[string]any{
		"message": msg,
		"status":  "failed",
	})
}

func Forbidden(c *fiber.Ctx, message ...string) error {
	var msg string
	if len(message) > 0 {
		msg = message[0]
	} else {
		msg = "forbidden"
	}

	return c.Status(fiber.StatusForbidden).JSON(map[string]any{
		"message": msg,
		"status":  "failed",
	})
}


func NotFound(c *fiber.Ctx, message ...string) error {
	var msg string
	if len(message) > 0 {
		msg = message[0]
	} else {
		msg = "not found"
	}

	return c.Status(fiber.StatusNotFound).JSON(map[string]any{
		"message": msg,
		"status":  "failed",
		"data": map[string]any{},
	})
}

func Conflict(c *fiber.Ctx, message ...string) error {
	var msg string
	if len(message) > 0 {
		msg = message[0]
	} else {
		msg = "conflict"
	}

	return c.Status(fiber.StatusConflict).JSON(map[string]any{
		"message": msg,
		"status":  "failed",
	})
}

func Unprocessable(c *fiber.Ctx, message ...string) error {
	var msg string
	if len(message) > 0 {
		msg = message[0]
	} else {
		msg = "unprocessable"
	}

	return c.Status(fiber.StatusUnprocessableEntity).JSON(map[string]any{
		"message": msg,
		"status":  "failed",
	})
}


func InternalServerError(c *fiber.Ctx, message ...string) error {
	var msg string
	if len(message) > 0 {
		msg = message[0]
	} else {
		msg = "internal server error"
	}

	return c.Status(fiber.StatusInternalServerError).JSON(
		map[string]any{
			"message": msg,
			"status":  "failed",
		},
	)
}