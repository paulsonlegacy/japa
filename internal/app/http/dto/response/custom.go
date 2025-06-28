package response

import (
	"time"
	"github.com/gofiber/fiber/v2"
)

func UserRegisteredOK(c *fiber.Ctx, data ...map[string]any) error {
	var payload map[string]any
	if len(data) > 0 {
		payload = data[0]
	} else {
		payload = map[string]any{}
	}

	return c.Status(fiber.StatusOK).JSON(map[string]any{
		"message": "user registered successfully",
		"status":  "success",
		"data":    payload,
	})
}

func Banned(c *fiber.Ctx, bannedUntil time.Time, banReason string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(map[string]any{
		"message": "Account is banned until " + bannedUntil.Format(time.RFC3339),
		"banned_until": bannedUntil.Format(time.RFC3339),
		"reason": banReason,
	})
}