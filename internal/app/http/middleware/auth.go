package middleware

import (
	"github.com/gofiber/fiber/v2"
	"japa/internal/config"
	"japa/internal/app/http/dto/response"
	"japa/internal/pkg"
	"japa/internal/domain/entity"
	"gorm.io/gorm"
	"go.uber.org/zap"
	"regexp"
	"time"
)

// AuthMiddleware struct holds server configuration for authentication purposes
type AuthMiddleware struct {
	config.ServerConfig
	DB *gorm.DB
}

// NewAuthMiddleware initializes a new instance of AuthMiddleware
func NewAuthMiddleware(serverConfig config.ServerConfig, db *gorm.DB) *AuthMiddleware {
	return &AuthMiddleware{serverConfig, db}
}

// List of routes that do not require authentication
var unprotectedRoutes = []*regexp.Regexp{
	regexp.MustCompile("/api/v1/register"),
	regexp.MustCompile("/api/v1/login"),
	regexp.MustCompile(`/api/v1/updater/version(?:\\?[^/]*)?$`),
	regexp.MustCompile(`/api/v1/updater/download(?:\\?[^/]*)?$`),
}

// Handler returns a Fiber middleware that validates authentication tokens
func (middleware *AuthMiddleware) Handler() fiber.Handler {
	return func(c *fiber.Ctx) error {

		// Allow unprotected routes to bypass authentication
		for _, route := range unprotectedRoutes {
			if route.MatchString(c.OriginalURL()) {
				return c.Next()
			}
		}

		// Fetch token from Authorization header (or configured header)
		token := c.Get(middleware.AuthorizationHeaderPath)

		// Validate the token
		claims, err := pkg.ValidateJWT(token)
		if err != nil {
			// Log debug information if token validation fails
			zap.L().Debug(
				"Failed to validate token! Error: " + err.Error(),
				zap.String("token", token),
				zap.String("path", c.Path()),
				zap.String("method", c.Method()),
				zap.String("ip", c.IP()),
				zap.String("user_agent", c.Get(fiber.HeaderUserAgent)),
				zap.Any("headers", c.GetReqHeaders()),
			)
			return response.Unauthorized(c, "Unauthorized access")
		}

		sub := claims["sub"]

		// User placeholder 
		var user entity.User

		// Fetch user from DB
		if err := middleware.DB.
		Select("full_name", "username", "role", "banned_until", "ban_reason").
		First(&user, "id = ?", sub).Error; err != nil {
			return response.InternalServerError(c)
		}

		// Checking if user is banned
		if user.Role != "admin" && user.Role != "superadmin" && user.BannedUntil != nil {
			if time.Now().Before(*user.BannedUntil) { // If now is before the blocked datetime
				return response.Banned(c, *user.BannedUntil, *user.BanReason)
			}
		}

		// Save base user data to context
		c.Locals("user_id", sub)
		c.Locals("fullname", user.FullName)
		c.Locals("username", user.Username)
		c.Locals("role", user.Role)

		// Continue to next middleware/handler
		return c.Next()
	}
}


// RoleRequired returns a middleware that checks if the user has one of the allowed roles.
// Example usage: RoleRequired("admin", "superadmin")
func RoleRequired(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get role from previous auth middleware (from Locals)
		role, ok := c.Locals("role").(string)
		if !ok || role == "" {
			// Role not found in context, maybe user is not authenticated properly
			return response.Unauthorized(c)
		}

		// Check if user's role matches any allowed role
		for _, allowed := range allowedRoles {
			if role == allowed {
				return c.Next() // User has permission, continue
			}
		}

		// If no match found, deny access
		return response.Forbidden(c)
	}
}


// AdminOnly returns middleware for admin and superadmin only
func AdminOnly() fiber.Handler {
	return RoleRequired("admin", "superadmin")
}


// ModeratorOnly returns middleware for moderator and superadmin only
func ModeratorOnly() fiber.Handler {
	return RoleRequired("moderator", "superadmin")
}


// SuperadminOnly returns middleware for superadmin only
func SuperadminOnly() fiber.Handler {
	return RoleRequired("superadmin")
}
