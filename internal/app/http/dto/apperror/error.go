package apperror


type AppError struct {
	Code    string `json:"code"`    // e.g., VALIDATION_ERROR, DB_ERROR
	Message string `json:"message"` // user-friendly message
	Details string `json:"details"` // actual error or dev-friendly info
}

func (e *AppError) Error() string {
	return e.Details
}

// Constructor
func New(code, message, details string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

func NewValidationErr(details string) *AppError {
	return New(ErrCodeValidation, "Invalid input", details)
}

func NewUnauthorizedErr(details string) *AppError {
	return New(ErrCodeUnauthorized, "Access denied", details)
}

func NewForbiddenErr(details string) *AppError {
	return New(ErrCodeForbidden, "Forbidden request", details)
}

func NewServerErr(details string) *AppError {
	return New(ErrCodeInternalServer, "Something went wrong", details)
}

