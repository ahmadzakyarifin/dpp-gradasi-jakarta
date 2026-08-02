package validator

import "fmt"

func validationMessage(tag, field, param string) string {
	switch tag {
	case "required", "required_if", "required_with", "required_without":
		return fmt.Sprintf("%s wajib diisi", field)

	case "email":
		return fmt.Sprintf("Format %s tidak valid", field)

	case "min":
		return fmt.Sprintf("%s minimal %s karakter", field, param)

	case "eqfield":
		return fmt.Sprintf("%s harus sama", field)

	case "phone":
		return fmt.Sprintf("Format nomor telepon pada %s tidak valid", field)

	default:
		return fmt.Sprintf("%s tidak valid", field)
	}
}
