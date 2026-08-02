package validator

import (
	"github.com/go-playground/validator/v10"
)

// RegisterRules mendaftarkan custom rule/tag ke validator engine.
func RegisterRules(v *validator.Validate) {
	// Mendaftarkan tag "phone"
	_ = v.RegisterValidation("phone", validatePhone)
}

func validatePhone(fl validator.FieldLevel) bool {
	val := fl.Field().String()

	// Gunakan "ID" sebagai default region.
	// Jika user input awalan +62 atau +65, libphonenumber akan otomatis
	// menggunakan kode negara dari input tersebut.
	_, err := Normalize("ID", val)

	// Jika tidak ada error dari fungsi Normalize, berarti nomor valid
	return err == nil
}
