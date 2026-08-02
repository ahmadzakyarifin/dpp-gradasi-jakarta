package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ServiceError struct {
	Code    string
	Message string
	Err     error
}

func (e *ServiceError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func NewServiceError(code, message string, err error) *ServiceError {
	return &ServiceError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func (e *ServiceError) Unwrap() error {
	return e.Err
}

// ServiceErrorToHTTP memetakan kode error ke HTTP status code.
// Satu sumber kebenaran untuk semua modul (sebelumnya diduplikasi per handler).
func ServiceErrorToHTTP(code string) int {
	switch code {
	case "AUTH_INVALID_CREDENTIALS", "AUTH_SESSION_EXPIRED", "AUTH_TOKEN_INVALID_OR_EXPIRED":
		return http.StatusUnauthorized
	case "AUTH_ACCOUNT_INACTIVE", "AUTH_ACCOUNT_PENDING", "FORBIDDEN":
		return http.StatusForbidden
	case "NOT_FOUND":
		return http.StatusNotFound
	case "VALIDATION_ERROR":
		return http.StatusUnprocessableEntity
	case "AUTH_EMAIL_EXISTS", "DUPLICATE_ENTRY":
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

// HandleServiceError menulis response error secara konsisten dari *ServiceError.
// Handler cukup memanggil ini — tidak perlu mapping manual lagi.
func HandleServiceError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	svc, ok := err.(*ServiceError)
	if !ok {
		ErrorResponse(c, http.StatusInternalServerError, "SERVER_ERROR", "Terjadi kesalahan pada server.", nil)
		return
	}

	ErrorResponse(c, ServiceErrorToHTTP(svc.Code), svc.Code, svc.Message, nil)
}
