package test

import (
    "testing"
    "time"

    "japa/internal/domain/entity"
    "japa/internal/app/http/dto/request" // Adjust if your module path differs

    "golang.org/x/crypto/bcrypt"
)

func TestToUser(t *testing.T) {
    req := &request.CreateUserRequest{
        FullName: "Jane Doe",
        Username: "janedoe",
        Email:    "jane@example.com",
        Phone:    "+12345678901",
        Password: "securepassword",
        Role:     "user",
    }

    user, err := request.ToUser(req)
    if err != nil {
        t.Fatalf("ToUser returned an error: %v", err)
    }

    // Check top-level fields
    if user.FullName != req.FullName {
        t.Errorf("expected FullName %q, got %q", req.FullName, user.FullName)
    }

    if user.Email != req.Email {
        t.Errorf("expected Email %q, got %q", req.Email, user.Email)
    }

    if user.Phone != req.Phone {
        t.Errorf("expected Phone %q, got %q", req.Phone, user.Phone)
    }

    if user.Role != req.Role {
        t.Errorf("expected Role %q, got %q", req.Role, user.Role)
    }

    // Check that ID is not zero (ULID)
    if user.ID.Compare(entity.User{}.ID) == 0 {
        t.Errorf("expected non-zero ULID")
    }

    // Check password is hashed
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("securepassword"))
    if err != nil {
        t.Errorf("expected hashed password to match original, but got error: %v", err)
    }

    // Check timestamps (not nil / not zero)
    if user.CreatedAt.IsZero() || user.UpdatedAt.IsZero() {
        t.Errorf("expected CreatedAt and UpdatedAt to be set")
    }

    // Optional: ensure CreatedAt ≈ UpdatedAt (small time delta)
    delta := user.UpdatedAt.Sub(user.CreatedAt)
    if delta > time.Second {
        t.Errorf("expected CreatedAt and UpdatedAt to be close, but delta was %v", delta)
    }
}
