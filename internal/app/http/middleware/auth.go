package middleware

import (
	"github.com/gofiber/fiber/v2"
	"japa/internal/config"
	"japa/internal/app/http/dto/apperror"
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
	config.JWTConfig
	DB *gorm.DB
}

// NewAuthMiddleware initializes a new instance of AuthMiddleware
func NewAuthMiddleware(serverConfig config.ServerConfig, JWTConfig config.JWTConfig,  db *gorm.DB) *AuthMiddleware {
	return &AuthMiddleware{serverConfig, JWTConfig, db}
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
		claims, err := pkg.ValidateJWT(token, middleware.JWTConfig)
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
			return response.Unauthorized(c, apperror.NewUnauthorizedErr(err.Error()))
		}

		////// USER LOGIC & AUTHORIZATION //////

		userID := claims["sub"]
		var user entity.User
		
		// Fetch user basic info
		if err := middleware.DB.
			Select("id", "full_name", "username", "role", "banned_until", "ban_reason").
			First(&user, "id = ?", userID).Error; err != nil {
			return response.InternalServerError(c, apperror.NewServerErr(err.Error()))
		}

		// Checking if user is banned
		if user.Role != "admin" && user.Role != "superadmin" && user.BannedUntil != nil {
			if time.Now().Before(*user.BannedUntil) { // If now is before the blocked datetime
				return response.Banned(c, *user.BannedUntil, *user.BanReason)
			}
		}

		////// SUBSCRIPTION LOGIC & CANCELATION ///////

		var userSubscriptions []entity.Subscription

		// Fetch all subscriptions for this user
		if err := middleware.DB.
			Where("user_id = ?", user.ID).
			//Where("status = ?", "active").
			//Where("expires_at > ?", time.Now()).
			Order("started_at DESC").
			Find(&userSubscriptions).Error; err != nil {
			return response.InternalServerError(c, apperror.NewServerErr(err.Error()))
		}

		// User active subscription
		var userActiveSubscription *entity.Subscription

		// Loop through subscriptions to expire outdated ones
		// and pick the latest active
		for i := range userSubscriptions {
			s := &userSubscriptions[i]

			// Check if expired and status still active
			if s.ExpiresAt.Before(time.Now()) && s.Status == "active" {
				s.Status = "expired"
				if err := middleware.DB.Model(s).Update("status", "expired").Error; err != nil {
					// Log error
					zap.L().Error(
						"Failure updating subscription status",
						zap.String("subscriptionID", s.ID),
					)
				}
			}

			// If still active, pick as userActiveSubscription
			if s.Status == "active" && s.ExpiresAt.After(time.Now()) && userActiveSubscription == nil {
				userActiveSubscription = s
			}
		}

		// Save user data to context
		c.Locals("user_id", userID)
		c.Locals("full_name", user.FullName)
		c.Locals("username", user.Username)
		c.Locals("role", user.Role)
		c.Locals("subscription", userActiveSubscription)

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
			return response.Unauthorized(c, apperror.NewUnauthorizedErr("User role cannot be identified"))
		}

		// Check if user's role matches any allowed role
		for _, allowed := range allowedRoles {
			if role == allowed {
				return c.Next() // User has permission, continue
			}
		}

		// If no match found, deny access
		return response.Forbidden(c, apperror.NewForbiddenErr("Access denied for user"))
	}
}


// AgentOnly returns middleware for agents and admins only
func AgentOnly() fiber.Handler {
	return RoleRequired("agent", "admin", "superadmin")
}


// AdminOnly returns middleware for admins only
func AdminOnly() fiber.Handler {
	return RoleRequired("admin", "superadmin")
}


// SuperadminOnly returns middleware for superadmin only
func SuperadminOnly() fiber.Handler {
	return RoleRequired("superadmin")
}
