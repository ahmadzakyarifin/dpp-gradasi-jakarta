package validator

import (
	"errors"
	"reflect"
	"strings"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var registerOnce sync.Once

// RegisterValidator mendaftarkan validator global Gin.
// Panggil sekali saat aplikasi start.
func RegisterValidator() {
	registerOnce.Do(func() {
		v, ok := binding.Validator.Engine().(*validator.Validate)
		if !ok {
			return
		}

		v.RegisterTagNameFunc(func(f reflect.StructField) string {
			if label := f.Tag.Get("label"); label != "" {
				return label
			}

			name := strings.SplitN(f.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}

			return name
		})

		// Mendaftarkan custom rules
		RegisterRules(v)
	})
}

// ValidationErrorItem adalah format error 422 sesuai kontrak API.
type ValidationErrorItem struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Param   string `json:"param"`
	Message string `json:"message"`
}

// ErrorArray mengubah validator.ValidationErrors menjadi []ValidationErrorItem.
// Format: [{field, tag, param, message}] - sesuai standar contract API.
func Errors(err error) []ValidationErrorItem {
	if err == nil {
		return nil
	}

	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return nil
	}

	result := make([]ValidationErrorItem, 0, len(validationErrors))
	for _, e := range validationErrors {
		result = append(result, ValidationErrorItem{
			Field:   e.Field(),
			Tag:     strings.ToLower(e.Tag()),
			Param:   e.Param(),
			Message: validationMessage(e.Tag(), e.Field(), e.Param()),
		})
	}
	return result
}
