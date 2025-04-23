package utils

import (
	"time"

	"japa/internal/config"
	"japa/internal/models"

	"github.com/golang-jwt/jwt/v5"
)


func GenerateJWT(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"username": user.Username,
		"email": user.Email,
		"exp": time.Now().Add(config.Settings.JWT.Expiry).Unix(),
		"iss": config.Settings.JWT.Issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Settings.JWT.JWTSecretKey))
}