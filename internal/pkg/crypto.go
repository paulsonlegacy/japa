// crypto.go
package pkg

import (
    "crypto/sha512"
    "encoding/hex"
    "crypto/rand"
    "encoding/base64"
    "golang.org/x/crypto/argon2"
)

// HashSHA512 computes the SHA-512 hash of a string and returns a hexadecimal string.
// This is a fast pre-hash step before applying Argon2 to normalize input length.
func HashSHA512(s string) string {
    h := sha512.Sum512([]byte(s))
    return hex.EncodeToString(h[:])
}


// HashAndEncodeArgon2 generates a random salt, hashes the password using Argon2id,
// and returns a single string in the format: salt$hash.
// Store this string in your database for later verification.
func HashAndEncodeArgon2(password string, saltLength int) string {
    saltBytes := make([]byte, saltLength)
    _, err := rand.Read(saltBytes)
    if err != nil {
        panic(err)
    }
    salt := base64.RawStdEncoding.EncodeToString(saltBytes)

    hash := argon2.IDKey([]byte(password), saltBytes, 1, 64*1024, 4, 32)
    hashEncoded := base64.RawStdEncoding.EncodeToString(hash)

    return salt + "$" + hashEncoded
}


func encodeBase64(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
} 