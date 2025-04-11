// utils/jwt.go
package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWTClaim adalah struktur klaim dalam JWT token
type JWTClaim struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateJWT membuat JWT token baru untuk pengguna
func GenerateJWT(userID uint) (string, error) {
	// Mendapatkan secret key dari environment variable
	secretKey := getEnv("JWT_SECRET_KEY", "your-secret-key")

	claims := JWTClaim{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token berlaku 24 jam
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken memvalidasi token JWT
func ValidateToken(tokenString string) (*JWTClaim, error) {
	secretKey := getEnv("JWT_SECRET_KEY", "your-secret-key")

	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}

// getEnv mendapatkan nilai dari environment variable atau menggunakan nilai default
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
