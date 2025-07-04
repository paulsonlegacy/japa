package request

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ScrapedPostRequest struct {
	Category     string    `json:"category" validate:"default='news'"`
    Title        string    `json:"title" validate:"required,max=255"`
    Excerpt      *string   `json:"excerpt" validate:"omitempty"`
    PostImg      *string   `json:"post_img" validate:"omitempty"`
    ContentHTML  string    `json:"content_html" validate:"required"`
    ContentText  *string   `json:"content_text" validate:"omitempty"`
	Source       string    `json:"source" validate:"required"`
}


// Bind parses and validates the request body and returns a User entity
func (req *ScrapedPostRequest) Bind(c *fiber.Ctx, v *validator.Validate) error {
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