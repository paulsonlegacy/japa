
package request

import (
	"time"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type CreatePostRequest struct {
	AuthorID    string    `json:"author_id" validate:"required,ulid"`
	Title       string    `json:"title" validate:"required"`
	Content     string    `json:"content" validate:"required"`
	Excerpt     string    `json:"excerpt"`
	Tags        []string  `json:"tags" validate:"omitempty,dive,required"`
	AccessLevel *string      `json:"access_level"`
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


type PostWithAuthor struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	Excerpt     string     `json:"excerpt"`
	Tags        []string   `json:"tags"`
	TagsRaw     string     `json:"-"` // We store the raw DB value here before decoding
	CreatedAt   time.Time  `json:"created_at"`
	AccessLevel string     `json:"access_level"`
	AuthorName  string     `json:"author_name"`
}
