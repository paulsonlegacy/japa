package request

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type CreateReplyRequest struct {
	CommentID string `json:"comment_id" validate:"required"`
	AuthorID  string `json:"author_id" validate:"required"`
	Content   string `json:"content" validate:"required"`
}

// Bind parses and validates the request body and returns an entity
func (req *CreateReplyRequest) Bind(c *fiber.Ctx, v *validator.Validate) error {
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
	