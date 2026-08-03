package helper

import (
	"net/http"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/validator"
	"github.com/gin-gonic/gin"
)

// ServiceError — error standar service layer, membawa kode untuk kontrol HTTP status.
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

// ============================================================
// Typed errors — dipakai service untuk kontrol HTTP status
// ============================================================

// ValidationError — validasi bisnis (bukan binding), HTTP 422.
type ValidationError struct {
	Fields map[string]string
	Errors map[string]string
}

func NewValidationError() *ValidationError {
	return &ValidationError{Fields: map[string]string{}, Errors: map[string]string{}}
}

func (e *ValidationError) Add(field, message string) {
	e.Fields[field] = message
	e.Errors[field] = message
}

func (e *ValidationError) Error() string {
	return "validasi gagal"
}

// AuthenticationError — kredensial/sesi bermasalah, HTTP 401/403.
type AuthenticationError struct {
	Message string
	Code    string
}

func (e *AuthenticationError) Error() string {
	return e.Message
}

// NotFoundError — data tidak ditemukan, HTTP 404.
type NotFoundError struct {
	Message string
}

func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{Message: message}
}

func (e *NotFoundError) Error() string {
	return e.Message
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
	case "AUTH_EMAIL_EXISTS", "DUPLICATE_ENTRY", "DUPLICATE_TITLE":
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

// HandleServiceError menulis response error secara konsisten.
// Handler cukup memanggil ini — tidak perlu mapping manual lagi.
func HandleServiceError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	switch e := err.(type) {
	case *ServiceError:
		ErrorResponse(c, ServiceErrorToHTTP(e.Code), e.Code, e.Message, nil)
	case *ValidationError:
		items := make([]validator.ValidationErrorItem, 0, len(e.Fields))
		for field, msg := range e.Fields {
			items = append(items, validator.ValidationErrorItem{Field: field, Message: msg})
		}
		ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "validasi gagal", items)
	case *AuthenticationError:
		code := e.Code
		if code == "" {
			code = "AUTHENTICATION_ERROR"
		}
		ErrorResponse(c, ServiceErrorToHTTP(code), code, e.Message, nil)
	case *NotFoundError:
		ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", e.Message, nil)
	default:
		ErrorResponse(c, http.StatusInternalServerError, "SERVER_ERROR", "Terjadi kesalahan pada server.", nil)
	}
}
