package validator

import (
	"fmt"
	"strings"

	"github.com/nyaruka/phonenumbers"
)

// Normalize mengubah nomor telepon menjadi format E.164.
// Contoh:
// ID + 081234567890 -> +6281234567890
func Normalize(countryISO, phone string) (string, error) {
	phone = strings.TrimSpace(phone)

	if phone == "" {
		return "", fmt.Errorf("nomor telepon wajib diisi")
	}

	num, err := phonenumbers.Parse(
		phone,
		strings.ToUpper(countryISO),
	)
	if err != nil {
		return "", fmt.Errorf("format nomor telepon tidak valid")
	}

	if !phonenumbers.IsValidNumber(num) {
		return "", fmt.Errorf("nomor telepon tidak valid")
	}

	return phonenumbers.Format(
		num,
		phonenumbers.E164,
	), nil
}

// WhatsAppChatID mengubah nomor menjadi ChatID WAHA.
//
// Contoh:
// +6281234567890
// menjadi
// 6281234567890@c.us
func WhatsAppChatID(countryISO, phone string) (string, error) {
	number, err := Normalize(countryISO, phone)
	if err != nil {
		return "", err
	}

	return strings.TrimPrefix(number, "+") + "@c.us", nil
}

// NormalizePhoneNumber mengubah nomor menjadi format E.164 (+62...).
func NormalizePhoneNumber(phone string) string {
	normalized, err := Normalize("ID", phone)
	if err != nil {
		return ""
	}
	return normalized
}

// ValidatePhoneNumber memvalidasi format nomor telepon Indonesia.
func ValidatePhoneNumber(phone string) error {
	if _, err := Normalize("ID", phone); err != nil {
		return fmt.Errorf("nomor telepon tidak valid")
	}
	return nil
}
