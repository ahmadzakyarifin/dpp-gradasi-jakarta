package otp

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	otpTTL        = 5 * time.Minute
	otpKeyPrefix  = "otp:change:"
	maxAttempts   = 5
	otpDigitCount = 6
)

type OTPResult struct {
	Code       string `json:"code"`
	ExpiresSec int    `json:"expires_sec"`
}

// Store generates a random OTP, stores it in Redis, and returns the code.
// Key pattern: otp:change:<user_id>:<email|phone>
func Store(ctx context.Context, rdb *redis.Client, userID uint, contactType string) (*OTPResult, error) {
	code := fmt.Sprintf("%06d", rand.Intn(999999))
	key := fmt.Sprintf("%s%d:%s", otpKeyPrefix, userID, contactType)

	if err := rdb.Set(ctx, key, code, otpTTL).Err(); err != nil {
		return nil, fmt.Errorf("otp: store: %w", err)
	}

	// Track attempts
	attemptsKey := key + ":attempts"
	rdb.Set(ctx, attemptsKey, 0, otpTTL)

	return &OTPResult{
		Code:       code,
		ExpiresSec: int(otpTTL.Seconds()),
	}, nil
}

// Verify checks the OTP code. Returns the stored code on success so caller
// can log it (for audit). Consumes the OTP on success (deletes key).
func Verify(ctx context.Context, rdb *redis.Client, userID uint, contactType, inputCode string) (string, error) {
	key := fmt.Sprintf("%s%d:%s", otpKeyPrefix, userID, contactType)
	attemptsKey := key + ":attempts"

	// Check attempts
	attempts, _ := rdb.Get(ctx, attemptsKey).Int()
	if attempts >= maxAttempts {
		return "", fmt.Errorf("terlalu banyak percobaan OTP, minta ulang")
	}

	stored, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("kode OTP tidak ditemukan atau sudah kedaluwarsa")
	}
	if err != nil {
		return "", fmt.Errorf("otp: verify: %w", err)
	}

	if stored != inputCode {
		rdb.Incr(ctx, attemptsKey)
		remaining := maxAttempts - attempts - 1
		if remaining <= 0 {
			rdb.Del(ctx, key, attemptsKey)
			return "", fmt.Errorf("OTP salah, percobaan habis, silakan minta ulang")
		}
		return "", fmt.Errorf("kode OTP salah, sisa %d percobaan", remaining)
	}

	// Success — consume OTP
	rdb.Del(ctx, key, attemptsKey)
	return stored, nil
}

// Invalidate removes any stored OTP for the user+type (e.g. on cancel/error).
func Invalidate(ctx context.Context, rdb *redis.Client, userID uint, contactType string) {
	key := fmt.Sprintf("%s%d:%s", otpKeyPrefix, userID, contactType)
	rdb.Del(ctx, key, key+":attempts")
}
