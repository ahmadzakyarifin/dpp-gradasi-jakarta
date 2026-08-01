package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// TurnstileVerifyResult adalah response dari Cloudflare Turnstile siteverify.
type TurnstileVerifyResult struct {
	Success    bool     `json:"success"`
	ErrorCodes []string `json:"error-codes"`
}

// VerifyTurnstile memvalidasi token Turnstile ke Cloudflare siteverify.
// Mengembalikan (true, nil) jika valid; (false, err) jika gagal/error.
func VerifyTurnstile(ctx context.Context, secret, token, remoteIP, verifyURL string) (bool, error) {
	// Dev bypass: jika secret kosong di development, lewati verifikasi
	// (berguna saat backend jalan tanpa key Turnstile sungguhan).
	if secret == "" && verifyURL == "" {
		return true, nil
	}

	if token == "" {
		return false, fmt.Errorf("captcha token kosong")
	}

	base := verifyURL
	if base == "" {
		base = "https://challenges.cloudflare.com/turnstile/v0/siteverify"
	}

	form := url.Values{}
	form.Set("secret", secret)
	form.Set("response", token)
	if remoteIP != "" {
		form.Set("remoteip", remoteIP)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, base, strings.NewReader(form.Encode()))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return false, err
	}

	var result TurnstileVerifyResult
	if err := json.Unmarshal(body, &result); err != nil {
		return false, fmt.Errorf("gagal parse response siteverify: %w", err)
	}

	if !result.Success {
		return false, fmt.Errorf("captcha invalid: %v", result.ErrorCodes)
	}

	return true, nil
}
