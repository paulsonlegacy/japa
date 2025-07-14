package apperror

const (
	// General
	ErrCodeUnknown               = "UNKNOWN_ERROR"
	ErrCodeInternalServer        = "INTERNAL_SERVER_ERROR"
	ErrCodeBadRequest            = "BAD_REQUEST"
	ErrCodeTimeout               = "TIMEOUT_ERROR"
	ErrCodeRateLimited           = "RATE_LIMITED"

	// Validation
	ErrCodeValidation            = "VALIDATION_ERROR"
	ErrCodeMissingField          = "MISSING_FIELD"
	ErrCodeInvalidFormat         = "INVALID_FORMAT"
	ErrCodeOutOfRange            = "OUT_OF_RANGE"
	ErrCodeAlreadyExists         = "ALREADY_EXISTS"

	// Authentication / Authorization
	ErrCodeUserNotFound          = "USER_NOT_FOUND"
	ErrCodeUserBlocked           = "USER_BLOCKED"
	ErrCodeUserBanned            = "USER_BANNED"
	ErrCodeUnauthorized          = "UNAUTHORIZED"         // Not logged in
	ErrCodeForbidden             = "FORBIDDEN"            // Logged in, but not allowed
	ErrCodeInvalidCredentials    = "INVALID_CREDENTIALS"
	ErrCodeTokenExpired          = "TOKEN_EXPIRED"
	ErrCodeTokenInvalid          = "TOKEN_INVALID"

	// Database
	ErrCodeDatabase              = "DATABASE_ERROR"
	ErrCodeDuplicateKey          = "DUPLICATE_KEY"
	ErrCodeRecordNotFound        = "RECORD_NOT_FOUND"
	ErrCodeForeignKeyFailed      = "FOREIGN_KEY_ERROR"
	ErrCodeConstraintFailed      = "CONSTRAINT_ERROR"

	// Network / 3rd party
	ErrCodeNetworkError          = "NETWORK_ERROR"
	ErrCodeServiceUnavailable    = "SERVICE_UNAVAILABLE"
	ErrCodeExternalAPIError      = "EXTERNAL_API_ERROR"

	// File uploads
	ErrCodeInvalidFile           = "INVALID_FILE_FORMAT"
	ErrCodeFileTooLarge          = "FILE_TOO_LARGE"
	ErrCodeFileTypeUnsupported   = "FILE_TYPE_UNSUPPORTED"
	ErrCodeUploadFailed          = "UPLOAD_FAILED"

	// Business logic
	ErrCodeInvalidID             = "INVALID_ID"
	ErrCodePostNotFound          = "POST_NOT_FOUND"
	ErrCodeFetchPost             = "ERROR_FETCHING_POST"
	ErrCodeFetchPosts            = "ERROR_FETCHING_POSTS"
	ErrCodeInvalidApplication    = "INVALID_APPLICATION"
	ErrCodeDocNotFound           = "DOCUMENT_NOT_FOUND"
	ErrCodeInvalidDocFormat      = "INVALID_DOC_FORMAT"
	ErrCodePaymentFailed         = "PAYMENT_FAILED"
	ErrCodePostLimitReached      = "POST_LIMIT_REACHED"
	ErrCodePlanExpired           = "PLAN_EXPIRED"
)