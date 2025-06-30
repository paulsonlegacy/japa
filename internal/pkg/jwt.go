package pkg

import (
	"time"
    "encoding/base64"

	"japa/internal/config"
	"japa/internal/domain/entity"

	"github.com/golang-jwt/jwt/v5"
)


// GenerateJWT creates and signs a JWT for the given user entity.
func GenerateJWT(user *entity.User, JWTConfig config.JWTConfig) (string, error) {
    // Define the claims payload for the token.
    claims := jwt.MapClaims{
        "sub":      user.ID,        // Standard claim: subject (unique user ID)
        //"fullname": user.FullName,  // Custom claim: user's full name
        //"username": user.Username,  // Custom claim: username
        //"role":     user.Role,      // Custom claim: user role (e.g., "admin")
        "exp":      time.Now().Add(JWTConfig.Expiry).Unix(), // Expiration time as UNIX timestamp
        "iss":      JWTConfig.Issuer,  // Issuer identifier
    }

    // Create a new JWT token object using HMAC SHA-256 signing method.
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    // Sign the token with your secret key (Base64-encoded from config).
    signedToken, err := token.SignedString([]byte(encodeBase64(JWTConfig.JWTSecretKey)))

    // Return the signed JWT string and any error encountered during signing.
    return signedToken, err
}



// ValidateJWT parses and validates a JWT token string.
// It returns all claims (both standard and custom) as a jwt.MapClaims map.
func ValidateJWT(tokenString string, JWTConfig config.JWTConfig) (jwt.MapClaims, error) {
    // Create an empty MapClaims to hold all token claims (custom + standard)
    claims := jwt.MapClaims{}

    // Parse the token string and populate the claims map.
    // The key function provides the signing key used to verify the token signature.
    token, err := jwt.ParseWithClaims(
        tokenString,
        claims,
        func(token *jwt.Token) (interface{}, error) {
            // Provide the signing secret key.
            return []byte(encodeBase64(JWTConfig.JWTSecretKey)), nil
        },
    )
    if err != nil {
        // Token parsing or signature verification failed.
        return nil, err
    }

    // If the token is valid, return the claims map.
    if token.Valid {
        return claims, nil
    }

    // If the token signature is invalid
    return nil, jwt.ErrTokenSignatureInvalid
}



// encodeBase64 encodes a string to standard Base64.
func encodeBase64(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}