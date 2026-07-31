package helper

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// GenerateResetToken membuat hash token reset password
func GenerateResetToken(secret string, ttlMinutes int) (string, string, time.Time, error) {
	raw := uuid.New().String()
	expiry := time.Now().Add(time.Duration(ttlMinutes) * time.Minute)

	// Hash token untuk disimpan di DB
	hash := sha256.Sum256([]byte(raw + secret))
	tokenHash := hex.EncodeToString(hash[:])

	return raw, tokenHash, expiry, nil
}

// ValidateResetToken memvalidasi raw token cocok dengan hash di DB
func HashToken(rawToken, secret string) string {
	hash := sha256.Sum256([]byte(rawToken + secret))
	return hex.EncodeToString(hash[:])
}

// GenerateActivationToken mirip GenerateResetToken
func GenerateActivationToken(secret string, ttlHours int) (string, string, time.Time, error) {
	return GenerateResetToken(secret, ttlHours*60)
}

// ParseExpiresAt helper untuk parse expiry dari JWT claims
func ParseExpiresAt(claims *jwt.RegisteredClaims) time.Time {
	if claims.ExpiresAt != nil {
		return claims.ExpiresAt.Time
	}
	return time.Time{}
}
