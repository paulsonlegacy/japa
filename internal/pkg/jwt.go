package pkg

import (
	"time"

	"japa/internal/config"
	"japa/internal/domain/entity"

	"github.com/golang-jwt/jwt/v5"
)


func GenerateJWT(user *entity.User) (string, error) {
	claims := jwt.MapClaims{
		"sub": user.ID,
		"fullname": user.FullName,
		"username": user.Username,
		"role": user.Role,
		"exp": time.Now().Add(config.Settings.JWT.Expiry).Unix(),
		"iss": config.Settings.JWT.Issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Settings.JWT.JWTSecretKey))
}