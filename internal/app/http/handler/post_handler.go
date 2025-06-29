package handlers

import (
	"time"
	"context"
	"strconv"

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
type PostHandler struct {
	Validator *validator.Validate
	Usecase   *usecase.PostUsecase
}


// METHODS

// Initialize VisaApplication handler
func NewPostHandler(v *validator.Validate, uc *usecase.PostUsecase) *PostHandler {
	return &PostHandler{v, uc}
}

// Handler to create post
func (ph *PostHandler) CreatePost(c *fiber.Ctx) error {
	// Payload binding
	var reqBody request.CreatePostRequest
	if err := reqBody.Bind(c, ph.Validator); err != nil {
		return response.BadRequest(c)
	}

	// Contexts and timeouts
	ctx, cancel := context.WithTimeout(context.Background(), 20 * time.Second)
    defer cancel()

	// Pass to usecase layer
	if err := ph.Usecase.CreatePost(ctx, reqBody); err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, "Post created")
}

// Handler to fetch posts (/api/v1/posts?page=2&limit=20)
func (ph *PostHandler) FetchPosts(c *fiber.Ctx) error {
	// Parse query params
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit", "20"))
	if err != nil || limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call usecase
	posts, totalPosts, err := ph.Usecase.FetchPosts(ctx, limit, offset)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, "", map[string]any{
		"items": posts,
		"total": totalPosts,
	})
}
