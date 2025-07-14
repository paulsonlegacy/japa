package response

import (
	"time"

	"japa/internal/app/http/dto/apperror"

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
		"error": map[string]any{},
	})
}

func Banned(c *fiber.Ctx, bannedUntil time.Time, banReason string) error {
	return c.Status(fiber.StatusForbidden).JSON(map[string]any{
		"message": "Account is banned until " + bannedUntil.Format(time.RFC3339),
		"status": "failed",
		"data": map[string]any{},
		"error": map[string]any{
			"code": apperror.ErrCodeUserBanned,
			"message": "Account is banned until " + bannedUntil.Format(time.RFC3339),
			"detail": map[string]string{
				"banned_until": bannedUntil.Format(time.RFC3339),
				"reason": banReason,
			},
		},
	})
}