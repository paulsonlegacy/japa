
package request

import (
	"time"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)


type PostWithAuthor struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Content     string     `json:"content"`
	Excerpt     *string    `json:"excerpt"`
	Tags        *[]byte    `json:"tags"`
	TagsRaw     *[]string  `json:"-"` // We store the raw tags value here
	Source      *string    `json:"source"`
	CreatedAt   time.Time  `json:"created_at"`
	AccessLevel string     `json:"access_level"`
	AuthorID    *string    `json:"author_id"`
	AuthorName  *string    `json:"author_name"`
}


type CreatePostRequest struct {
	AuthorID    *string    `json:"author_id" validate:"omitempty,ulid"`
	Title       string    `json:"title" validate:"required"`
	Slug        string    `json:"slug" validate:"required"`
	Content     string    `json:"content" validate:"required"`
	Excerpt     *string    `json:"excerpt"`
	Tags        *[]string `json:"tags" validate:"omitempty,dive,required"`
	Source      *string   `json:"source" validate:"omitempty,url"`
	AccessLevel *string   `json:"access_level"`
}

// Bind parses and validates the request body and returns a User entity
func (req *CreatePostRequest) Bind(c *fiber.Ctx, v *validator.Validate) error {
	// Parse request body into req
	if err := c.BodyParser(req); err != nil {
		return err
	}

	// Validate request struct
	if err := v.Struct(req); err != nil {
		return err
	}

	return nil
}

type UpdatePostRequest struct {
	ID          string    `json:"id" validate:"required"`
	AuthorID    *string    `json:"author_id" validate:"omitempty,ulid"`
	Title       *string    `json:"title" validate:"omitempty"`
	Content     *string    `json:"content" validate:"omitempty"`
	Excerpt     *string   `json:"excerpt" validate:"omitempty"`
	Tags        *[]byte   `json:"tags" validate:"omitempty,dive,required"`
	Source      *string   `json:"source" validate:"omitempty,url"`
	AccessLevel *string   `json:"access_level" validate:"omitempty"`
}

// Bind parses and validates the request body and returns a User entity
func (req *UpdatePostRequest) Bind(c *fiber.Ctx, v *validator.Validate) error {
	// Parse request body into req
	if err := c.BodyParser(req); err != nil {
		return err
	}

	// Validate request struct
	if err := v.Struct(req); err != nil {
		return err
	}

	return nil
}
