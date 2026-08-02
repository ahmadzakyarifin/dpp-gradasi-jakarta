package helper

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func doHandle(err error) (int, map[string]any) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	HandleServiceError(c, err)

	var body map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	return w.Code, body
}

func TestServiceErrorToHTTP(t *testing.T) {
	cases := []struct {
		code string
		want int
	}{
		{"AUTH_INVALID_CREDENTIALS", http.StatusUnauthorized},
		{"AUTH_SESSION_EXPIRED", http.StatusUnauthorized},
		{"AUTH_TOKEN_INVALID_OR_EXPIRED", http.StatusUnauthorized},
		{"AUTH_ACCOUNT_INACTIVE", http.StatusForbidden},
		{"AUTH_ACCOUNT_PENDING", http.StatusForbidden},
		{"FORBIDDEN", http.StatusForbidden},
		{"NOT_FOUND", http.StatusNotFound},
		{"VALIDATION_ERROR", http.StatusUnprocessableEntity},
		{"AUTH_EMAIL_EXISTS", http.StatusConflict},
		{"DUPLICATE_ENTRY", http.StatusConflict},
		{"UNKNOWN_CODE", http.StatusInternalServerError},
		{"", http.StatusInternalServerError},
	}

	for _, tc := range cases {
		got := ServiceErrorToHTTP(tc.code)
		if got != tc.want {
			t.Errorf("ServiceErrorToHTTP(%q) = %d, want %d", tc.code, got, tc.want)
		}
	}
}

func TestHandleServiceError(t *testing.T) {
	// ServiceError NOT_FOUND -> 404 + code di body
	code, body := doHandle(NewServiceError("NOT_FOUND", "Berita tidak ditemukan", nil))
	if code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", code)
	}
	if body["code"] != "NOT_FOUND" || body["success"] != false {
		t.Errorf("unexpected body: %v", body)
	}

	// ServiceError AUTH_INVALID_CREDENTIALS -> 401
	code, _ = doHandle(NewServiceError("AUTH_INVALID_CREDENTIALS", "Email atau password salah.", nil))
	if code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", code)
	}

	// Error non-ServiceError -> 500 SERVER_ERROR
	code, body = doHandle(errors.New("something broke"))
	if code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", code)
	}
	if body["code"] != "SERVER_ERROR" {
		t.Errorf("expected SERVER_ERROR code, got %v", body["code"])
	}

	// nil error -> tidak panic, response 200 kosong (handler seharusnya tidak memanggil dengan nil)
	code, _ = doHandle(nil)
	if code != http.StatusOK {
		t.Errorf("nil error: expected 200 (no-op), got %d", code)
	}
}
