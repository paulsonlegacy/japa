package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func isAdminRole(c *fiber.Ctx) bool {
	role := c.Locals("role")
	return role == "admin" || role == "superadmin"
}
