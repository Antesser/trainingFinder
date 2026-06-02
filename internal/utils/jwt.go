package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(id string, secretKey []byte, duration time.Duration) (string, error) {
	iat := time.Now().UTC()

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": id,
		//"role": role,                     // Service role
		//"iat":  iat.Unix(),               // Issued at
		"exp": iat.Add(duration).Unix(), // Expiration time
	}).SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}

	return token, nil
}
