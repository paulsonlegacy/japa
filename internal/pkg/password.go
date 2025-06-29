// password.go
package pkg

import (
    "strings"
    "crypto/subtle"
    "crypto/sha512"
    "crypto/rand"
    "encoding/hex"
    "encoding/base64"
    "golang.org/x/crypto/argon2"
)


// HashSHA512 computes the SHA-512 hash of a string and returns a hexadecimal string.
// This is a fast pre-hash step before applying Argon2 to normalize input length.
func HashSHA512(s string) string {
    h := sha512.Sum512([]byte(s))
    return hex.EncodeToString(h[:])
}


// HashAndEncodeArgon2 hashes the raw password (internally pre-hashing with SHA-512),
// generates a random salt, and returns salt$hash.
// saltLength don't need to be consistent to be verified but mostly 16 or 32
func HashAndEncodeArgon2(password string, saltLength int) string {
    // Pre-hash input to normalize
    preHashed := HashSHA512(password)

    saltBytes := make([]byte, saltLength)
    _, err := rand.Read(saltBytes)
    if err != nil {
        panic(err)
    }

    hash := argon2.IDKey([]byte(preHashed), saltBytes, 1, 64*1024, 4, 32)

    return base64.RawStdEncoding.EncodeToString(saltBytes) + "$" + base64.RawStdEncoding.EncodeToString(hash)
}
 

// VerifyArgon2 verifies the raw password against the stored salt$hash.
func VerifyArgon2(password, hashedPassword string) bool {
    parts := strings.Split(hashedPassword, "$")
    if len(parts) != 2 {
        return false
    }

    salt, err1 := base64.RawStdEncoding.DecodeString(parts[0])
    expectedHash, err2 := base64.RawStdEncoding.DecodeString(parts[1])
    if err1 != nil || err2 != nil {
        return false
    }

    preHashed := HashSHA512(password)

    computedHash := argon2.IDKey([]byte(preHashed), salt, 1, 64*1024, 4, 32)

    return subtle.ConstantTimeCompare(computedHash, expectedHash) == 1
}


// Compare takes a plaintext password and a stored encoded hash (salt$hash),
// pre-hashes the password with SHA-512, and verifies it with Argon2.
// Returns true if the password is correct.
func Compare(password string, hashedPassword string) bool {
    return VerifyArgon2(password, hashedPassword)
}
