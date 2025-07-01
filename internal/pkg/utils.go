package pkg

import (
	"time"
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
)


// ContextWithTimeoutAndCancelFromFiber combines context.WithTimeout and the Fiber cancellation
// Fiber cancellation automatically terminates request when client disconnects or gives up
// Improved with context.Timeout to also force timeouts
//
// Use:
// ctx := ContextWithTimeoutAndCancelFromFiber(c, 5*time.Second)
// defer cancel()
func ContextWithTimeoutAndCancelFromFiber(
	c *fiber.Ctx, 
	timeout time.Duration,
) (
	context.Context, 
	context.CancelFunc,
) {
    // Create a base context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), timeout)

    // Create a goroutine that will cancel if the client disconnects
    go func() {
        <-c.Context().Done()
        cancel()
    }()

    return ctx, cancel
}

// ULIDToString returns the canonical string representation of a ULID.
func ULIDToString(id ulid.ULID) string {
	return id.String()
}

// StringToULID parses a ULID string and returns the ULID object.
// Returns an error if the string is invalid.
func StringToULID(s string) (ulid.ULID, error) {
	return ulid.Parse(s)
}