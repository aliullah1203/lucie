package helpers

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Secret key for signing JWT (must be loaded from env)
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// ValidateToken validates a JWT and returns claims
func ValidateToken(tokenString string) (*Claims, error) {
	// Re-load secret for safety if it wasn't loaded on startup (or use an init func)
	if len(jwtSecret) == 0 {
		jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	}
	if len(jwtSecret) == 0 {
		return nil, errors.New("JWT_SECRET not configured")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// GenerateToken creates a JWT token
func GenerateToken(userID, role string) (string, error) {
	// Re-load secret for safety
	if len(jwtSecret) == 0 {
		jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	}
	if len(jwtSecret) == 0 {
		return "", errors.New("JWT_SECRET not configured")
	}

	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
