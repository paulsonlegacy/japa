// password.go
package pkg

import (
    "strings"
    "crypto/subtle"
    "encoding/base64"
    "golang.org/x/crypto/argon2"
)

// VerifyArgon2 verifies a password against a stored Argon2 hash string in the format salt$hash.
// It decodes the salt, recomputes the hash with the same parameters, and compares them in constant time.
func VerifyArgon2(password, encodedHash string) bool {
    parts := strings.Split(encodedHash, "$")
    if len(parts) != 2 {
        return false
    }

    salt, err := base64.RawStdEncoding.DecodeString(parts[0])
    if err != nil {
        return false
    }

    expectedHash, err := base64.RawStdEncoding.DecodeString(parts[1])
    if err != nil {
        return false
    }

    computedHash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

    return subtle.ConstantTimeCompare(computedHash, expectedHash) == 1
}

// Compare takes a plaintext password and a stored encoded hash (salt$hash),
// pre-hashes the password with SHA-512, and verifies it with Argon2.
// Returns true if the password is correct.
func Compare(pw string, hashedPw string) bool {
    pw = HashSHA512(pw)
    return VerifyArgon2(pw, hashedPw)
}
