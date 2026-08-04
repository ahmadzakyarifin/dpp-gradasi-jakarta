package helper

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Email    string `json:"email"`
	RoleName string `json:"role_name"`
	Name     string `json:"name"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID uint, email string, name string, role string, secret string, ttlMins int) (string, error) {
	claims := &JWTClaims{
		UserID:   userID,
		Email:    email,
		RoleName: role,
		Name:     name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(ttlMins) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func GenerateRefreshToken(ttlHours int) (string, time.Time, error) {
	expiry := time.Now().Add(time.Duration(ttlHours) * time.Hour)
	token := uuid.New().String()
	return token, expiry, nil
}

func ValidateToken(tokenString string, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token tidak valid")
}
