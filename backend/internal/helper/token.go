package helper

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
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

// GenerateOTPToken membuat kode OTP numerik 6 digit untuk verifikasi
func GenerateOTPToken(secret string, ttlMinutes int) (string, string, time.Time, error) {
	otp := ""
	for i := 0; i < 6; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", "", time.Time{}, err
		}
		otp += n.String()
	}

	expiry := time.Now().Add(time.Duration(ttlMinutes) * time.Minute)

	// Hash OTP untuk disimpan di DB
	hash := sha256.Sum256([]byte(otp + secret))
	tokenHash := hex.EncodeToString(hash[:])

	return otp, tokenHash, expiry, nil
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

// GenerateEmailChangeOTP membuat kode OTP numerik 6 digit, expires 5 menit
func GenerateEmailChangeOTP(secret string) (string, string, time.Time, error) {
	return GenerateOTPToken(secret, 5)
}

// ParseExpiresAt helper untuk parse expiry dari JWT claims
func ParseExpiresAt(claims *jwt.RegisteredClaims) time.Time {
	if claims.ExpiresAt != nil {
		return claims.ExpiresAt.Time
	}
	return time.Time{}
}
