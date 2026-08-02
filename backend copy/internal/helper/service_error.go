package helper

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleServiceError menerjemahkan error dari service layer menjadi response
// HTTP yang konsisten di seluruh aplikasi:
//   - *ValidationError     -> 422 Unprocessable Entity
//   - *AuthenticationError -> 401 Unauthorized
//   - *NotFoundError       -> 404 Not Found
//   - error lain           -> 500 Internal Server Error
func HandleServiceError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	var valErr *ValidationError
	if errors.As(err, &valErr) {
		ErrorResponse(
			c,
			http.StatusUnprocessableEntity,
			"VALIDATION_ERROR",
			"Validasi gagal.",
			valErr.Errors,
		)
		return
	}

	var authErr *AuthenticationError
	if errors.As(err, &authErr) {
		code := authErr.Code
		if code == "" {
			code = "AUTH_ERROR"
		}

		ErrorResponse(
			c,
			http.StatusUnauthorized,
			code,
			authErr.Message,
			nil,
		)
		return
	}

	var notFoundErr *NotFoundError
	if errors.As(err, &notFoundErr) {
		ErrorResponse(
			c,
			http.StatusNotFound,
			"NOT_FOUND",
			notFoundErr.Message,
			nil,
		)
		return
	}

	ErrorResponse(
		c,
		http.StatusInternalServerError,
		"SERVER_ERROR",
		"Terjadi kesalahan pada server.",
		nil,
	)
}
