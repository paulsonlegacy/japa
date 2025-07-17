package handlers

import (
	"time"
	"context"
	"strconv"
	"errors"
	"encoding/json"

	"japa/internal/app/http/dto/apperror"
	"japa/internal/app/http/dto/request"
	"japa/internal/app/http/dto/response"
	//"japa/internal/domain/entity"
	"japa/internal/domain/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
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
		return response.BadRequest(c, apperror.New(
			apperror.ErrCodeValidation, 
			"Invalid request", 
			err.Error(),
		))
	}

	// Contexts and timeouts
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
    defer cancel()

	// Pass to usecase layer
	if err := ph.Usecase.CreatePost(ctx, reqBody); err != nil {
		return response.InternalServerError(c, apperror.New(
			apperror.ErrCodeDatabase, 
			"Something went wrong while creating post", 
			err.Error(),
		))
	}

	return response.Success(c, "Post created")
}

// Handler to update post
func (ph *PostHandler) UpdatePost(c *fiber.Ctx) error {
	// Parse request body
	var reqBody request.UpdatePostRequest
	if err := reqBody.Bind(c, ph.Validator); err != nil {
		return response.BadRequest(c, apperror.New(
			apperror.ErrCodeValidation,
			"Invalid request",
			err.Error(),
		))
	}

	// Check if user is allowed to update author
	if reqBody.AuthorID != nil && !isAdminRole(c) {
		return response.Forbidden(c, apperror.New(
			apperror.ErrCodeValidation,
			"Invalid request",
			"Not permitted to update post author",
		))
	}

	// Context for DB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Fetch existing post
	post, err := ph.Usecase.FetchPost(ctx, reqBody.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NotFound(c, apperror.New(
				apperror.ErrCodePostNotFound,
				"Post not found",
				err.Error(),
			))
		}
		return response.InternalServerError(c, apperror.New(
			apperror.ErrCodeDatabase,
			"Failed to fetch post",
			err.Error(),
		))
	}

	// Fill in missing fields from existing post
	if reqBody.AuthorID == nil {
		reqBody.AuthorID = post.AuthorID
	}
	if reqBody.Title == nil {
		reqBody.Title = &post.Title
	}
	if reqBody.Content == nil {
		reqBody.Content = &post.Content
	}
	if reqBody.Excerpt == nil {
		reqBody.Excerpt = post.Excerpt
	}
	if reqBody.Tags == nil || (post.Tags != nil && len(*post.Tags) > 0) {
		tagsJSON, _ := json.Marshal(*post.Tags)
		reqBody.Tags = &tagsJSON
	}
	if reqBody.Source == nil {
		reqBody.Source = post.Source
	}
	if reqBody.AccessLevel == nil {
		reqBody.AccessLevel = &post.AccessLevel
	}

	// Pass to usecase
	if err := ph.Usecase.UpdatePost(ctx, reqBody); err != nil {
		return response.InternalServerError(c, apperror.New(
			apperror.ErrCodeDatabase,
			"Failed to update post",
			err.Error(),
		))
	}

	return response.Success(c, "Post updated successfully")
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
		return response.InternalServerError(c, apperror.New(
			apperror.ErrCodeFetchPosts,
			"Error fetching posts",
			err.Error(),
		))
	}

	return response.Success(c, "", map[string]any{
		"items": posts,
		"total": totalPosts,
	})
}

func (ph *PostHandler) FetchPost(c *fiber.Ctx) error {
	postID := c.Params("post_id")

	if _, err := ulid.Parse(postID); err != nil {
		return response.BadRequest(c, apperror.New(
			apperror.ErrCodeInvalidID,
			"invalid post id format",
			err.Error(),
		))
	}

	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call usecase
	post, err := ph.Usecase.FetchPost(ctx, postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NotFound(c, apperror.New(
				apperror.ErrCodePostNotFound,
				"Post not found",
				err.Error(),
			))
		}
		return response.InternalServerError(c, apperror.New(
			apperror.ErrCodeInternalServer,
			"Error occured while fetching post",
			err.Error(),
		))
	}

	return response.Success(c, "", map[string]any{
		"items": post,
	})
}