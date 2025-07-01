// Fiber handlers for user routes
package handlers

import (
	"time"
	"context"

	"japa/internal/app/http/dto/request"
	"japa/internal/app/http/dto/response"
	"japa/internal/domain/usecase"
	"japa/internal/pkg"
	//"japa/internal/domain/entity"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	//"go.uber.org/zap"
)

// TYPES

// User handler
type UserHandler struct {
	Validator *validator.Validate
	Usecase   *usecase.UserUsecase
}

// METHODS

// Initialize user handler
func NewUserHandler(v *validator.Validate, us *usecase.UserUsecase) *UserHandler {
	return &UserHandler{v, us}
}


// Register handler
func (uh *UserHandler) Register(c *fiber.Ctx) error {
	// Parse req body
	var reqBody request.CreateUserRequest
	if err := reqBody.Bind(c, uh.Validator); err != nil {
		return response.BadRequest(c)
	}

	// Registering user
	err := uh.Usecase.RegisterUser(c.Context(), reqBody)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	// If registeration succeeded
	return response.UserRegisteredOK(c)
}



// Login handler
func (uh *UserHandler) Login(c *fiber.Ctx) error {
	var reqBody struct {
		Account string `json:"account"` // Username or Email
		Password string `json:"password"`
	}

	// Parsing incoming payload into user object
	if err := c.BodyParser(&reqBody); err != nil {
		return  response.BadRequest(c, "invalid input")
	}

	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Confirming user and generate tokens
	accessToken, refreshToken, err := uh.Usecase.LoginUser(ctx, reqBody.Account, reqBody.Password)
	if err != nil || refreshToken == "" {
		return response.Unauthorized(c, err.Error())
	}

	// Store refresh token in secure HTTP-only cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour), // 7 days - adjust lifespan
		Secure:   true,
		HTTPOnly: true,
		SameSite: "Strict",
	})

	// Return JWT token
	return response.Success(c, "login successful", map[string]any{
		"token": accessToken,
	})
}



func (uh *UserHandler) Logout(c *fiber.Ctx) error {
	// Fetching refresh token
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return  response.BadRequest(c, "No refresh token")
	}

	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Delete from DB
	if err := uh.Usecase.Logout(ctx, refreshToken); err != nil {
		return response.InternalServerError(c, "Failed to revoke token")
	}

	// Clear the cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
	})

	return response.Success(c, "Logged out successfully")
}



func (uh *UserHandler) RefreshToken(c *fiber.Ctx) error {
	// get refresh token from cookie
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return response.Unauthorized(c, "Missing refresh token")
	}

	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Lookup token in DB
	token, err := uh.Usecase.GetRefreshToken(ctx, refreshToken)
	if err != nil || token == nil {
		return response.Unauthorized(c, "Invalid refresh token")
	}

	// Check if token is expired
	if time.Now().After(token.ExpiresAt) {
		return response.Unauthorized(c, "Refresh token expired")
	}

	// Fetch the user
	user, err := uh.Usecase.Repo.FindUserByID(ctx, token.UserID)
	if err != nil {
		return response.Unauthorized(c, "User not found")
	}

	// Issue new access token
	accessToken, err := pkg.GenerateJWT(user, uh.Usecase.JWTConfig)
	if err != nil {
		return response.InternalServerError(c, "Could not create token")
	}

	// Rotate refresh token (recommended)
	newRefreshToken, err := pkg.GenerateRefreshToken()
	if err != nil {
		return response.InternalServerError(c, "Could not create refresh token")
	}

	// Update the DB record
	token.Token = newRefreshToken
	token.ExpiresAt = time.Now().Add(7 * 24 * time.Hour)
	uh.Usecase.DB.Save(&token)

	// Set the new cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		Expires:  token.ExpiresAt,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
	})

	// Return new access token
	return response.Success(c, "", map[string]any{
		"token": accessToken,
	})
}