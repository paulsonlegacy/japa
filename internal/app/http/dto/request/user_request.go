package request

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)


type CreateUserRequest struct {
	FullName string `json:"full_name" validate:"required,min=2,max=100"`
	Username  string `json:"username" validate:"required,min=2,max=100"`
	Email     string `json:"email" validate:"required,email"`
	Phone     string `json:"phone" validate:"required,e164"` // e164 complaint - +2348012345678, +447123456789 etc
	Password  string `json:"password" validate:"required,min=8"` // plain password; hash before saving
	Role      string `json:"role" validate:"required,oneof=admin user agent"` // customize roles as needed
}


// Bind parses and validates the request body and returns a User entity
func (req *CreateUserRequest) Bind(c *fiber.Ctx, v *validator.Validate) error {
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