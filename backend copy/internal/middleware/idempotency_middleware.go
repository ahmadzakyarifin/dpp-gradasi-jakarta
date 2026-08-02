package middleware

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	notificationentity "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/notification/entity"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
	"log"
)

const maxIdempotencyKeyLen = 255

type idempotencyResponse struct {
	StatusCode  int    `json:"status_code"`
	ContentType string `json:"content_type"`
	Body        string `json:"body"`
}

type responseCaptureWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseCaptureWriter) Write(data []byte) (int, error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}

func (w *responseCaptureWriter) WriteString(data string) (int, error) {
	w.body.WriteString(data)
	return w.ResponseWriter.WriteString(data)
}

// IdempotencyMiddleware caches successful unsafe responses for requests carrying
// X-Idempotency-Key in the MariaDB database.
func IdempotencyMiddleware(db *bun.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		if method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions {
			c.Next()
			return
		}

		key := strings.TrimSpace(c.GetHeader("X-Idempotency-Key"))
		if key == "" || db == nil {
			c.Next()
			return
		}
		if !isValidIdempotencyKey(key) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "X-Idempotency-Key tidak valid",
			})
			return
		}

		ctx := c.Request.Context()
		requestHash, err := idempotencyRequestHash(c.Request)
		if err != nil {
			log.Printf("ERROR: idempotency: failed to read request body. err: %v, key_fp: %s", err, idempotencyFingerprint(key))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "gagal membaca request body untuk idempotency",
			})
			return
		}

		claimed := &notificationentity.IdempotencyKey{
			Key:         key,
			Status:      notificationentity.IdempotencyStatusProcessing,
			RequestHash: requestHash,
		}
		if _, err := db.NewInsert().Model(claimed).Exec(ctx); err != nil {
			if !isDuplicateKeyError(err) {
				log.Printf("ERROR: idempotency: failed to claim key. err: %v, key_fp: %s", err, idempotencyFingerprint(key))
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"status":  false,
					"message": "gagal memproses idempotency key",
				})
				return
			}

			var ik notificationentity.IdempotencyKey
			err := db.NewSelect().Model(&ik).Where("`key` = ?", key).Scan(ctx)
			if errors.Is(err, sql.ErrNoRows) {
				log.Printf("WARN: idempotency: duplicate key vanished before select. key_fp: %s", idempotencyFingerprint(key))
				c.AbortWithStatusJSON(http.StatusConflict, gin.H{
					"status":  false,
					"message": "request dengan idempotency key ini sedang diproses",
				})
				return
			}
			if err != nil {
				log.Printf("ERROR: idempotency: failed to read existing key. err: %v, key_fp: %s", err, idempotencyFingerprint(key))
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"status":  false,
					"message": "gagal membaca idempotency key",
				})
				return
			}

			if ik.RequestHash != "" && ik.RequestHash != requestHash {
				c.AbortWithStatusJSON(http.StatusConflict, gin.H{
					"status":  false,
					"message": "idempotency key sudah digunakan untuk request yang berbeda",
				})
				return
			}

			switch ik.Status {
			case notificationentity.IdempotencyStatusCompleted:
				if ik.ResponsePayload != "" {
					var res idempotencyResponse
					if json.Unmarshal([]byte(ik.ResponsePayload), &res) == nil {
						contentType := res.ContentType
						if contentType == "" {
							contentType = "application/json; charset=utf-8"
						}
						c.Data(res.StatusCode, contentType, []byte(res.Body))
						c.Abort()
						return
					}
				}

				log.Printf("WARN: idempotency: completed response payload is invalid. key_fp: %s", idempotencyFingerprint(key))
				c.AbortWithStatusJSON(http.StatusConflict, gin.H{
					"status":  false,
					"message": "response idempotency tidak valid, silakan gunakan key baru",
				})
				return
			case notificationentity.IdempotencyStatusProcessing:
				c.Header("Retry-After", "2")
				c.AbortWithStatusJSON(http.StatusConflict, gin.H{
					"status":  false,
					"message": "request dengan idempotency key ini sedang diproses",
				})
				return
			default:
				c.AbortWithStatusJSON(http.StatusConflict, gin.H{
					"status":  false,
					"message": "status idempotency key tidak valid",
				})
				return
			}
		}

		markAsRetryable := func(reason string) {
			if _, err := db.NewDelete().
				Model((*notificationentity.IdempotencyKey)(nil)).
				Where("`key` = ?", key).
				Where("status = ?", notificationentity.IdempotencyStatusProcessing).
				Exec(context.WithoutCancel(ctx)); err != nil {
				log.Printf("ERROR: idempotency: failed to delete retryable key. err: %v, key_fp: %s, reason: %s", err, idempotencyFingerprint(key), reason)
			}
		}

		defer func() {
			if rec := recover(); rec != nil {
				markAsRetryable("panic")
				panic(rec)
			}
		}()

		capture := &responseCaptureWriter{ResponseWriter: c.Writer, body: bytes.NewBuffer(nil)}
		c.Writer = capture
		c.Next()

		status := c.Writer.Status()
		if status < http.StatusOK || status >= http.StatusMultipleChoices {
			markAsRetryable("non_2xx_response")
			return
		}

		res := idempotencyResponse{
			StatusCode:  status,
			ContentType: c.Writer.Header().Get("Content-Type"),
			Body:        capture.body.String(),
		}
		payload, err := json.Marshal(res)
		if err != nil {
			log.Printf("ERROR: idempotency: failed to marshal cached response. err: %v, key_fp: %s", err, idempotencyFingerprint(key))
			markAsRetryable("marshal_response_failed")
			return
		}

		if _, err := db.NewUpdate().
			Model((*notificationentity.IdempotencyKey)(nil)).
			Set("status = ?", notificationentity.IdempotencyStatusCompleted).
			Set("response_payload = ?", string(payload)).
			Set("updated_at = ?", time.Now()).
			Where("`key` = ?", key).
			Where("status = ?", notificationentity.IdempotencyStatusProcessing).
			Exec(context.WithoutCancel(ctx)); err != nil {
			log.Printf("ERROR: idempotency: failed to mark key completed. err: %v, key_fp: %s", err, idempotencyFingerprint(key))
			markAsRetryable("complete_update_failed")
		}
	}
}

func isValidIdempotencyKey(key string) bool {
	if key == "" || len(key) > maxIdempotencyKeyLen || !utf8.ValidString(key) {
		return false
	}
	for _, r := range key {
		switch {
		case r >= 'a' && r <= 'z':
		case r >= 'A' && r <= 'Z':
		case r >= '0' && r <= '9':
		case r == '-' || r == '_' || r == '.' || r == ':' || r == '~':
		default:
			return false
		}
	}
	return true
}

func idempotencyRequestHash(r *http.Request) (string, error) {
	var body []byte
	if r.Body != nil {
		var err error
		body, err = io.ReadAll(r.Body)
		if err != nil {
			return "", err
		}
		r.Body = io.NopCloser(bytes.NewReader(body))
	}

	sum := sha256.Sum256([]byte(fmt.Sprintf("%s\n%s\n%s\n%s", r.Method, r.URL.Path, r.URL.RawQuery, string(body))))
	return fmt.Sprintf("%x", sum), nil
}

func isDuplicateKeyError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062
	}
	return strings.Contains(strings.ToLower(err.Error()), "duplicate")
}

func idempotencyFingerprint(value string) string {
	sum := sha256.Sum256([]byte(value))
	return fmt.Sprintf("%x", sum[:8])
}

func CleanupIdempotencyKeys(ctx context.Context, db *bun.DB, completedTTL time.Duration, processingTTL time.Duration) (int64, error) {
	if db == nil {
		return 0, nil
	}
	if completedTTL <= 0 {
		completedTTL = 24 * time.Hour
	}
	if processingTTL <= 0 {
		processingTTL = 15 * time.Minute
	}

	now := time.Now()
	res, err := db.NewDelete().
		Model((*notificationentity.IdempotencyKey)(nil)).
		WhereGroup(" AND ", func(q *bun.DeleteQuery) *bun.DeleteQuery {
			return q.
				WhereOr("status = ? AND updated_at < ?", notificationentity.IdempotencyStatusCompleted, now.Add(-completedTTL)).
				WhereOr("status = ? AND updated_at < ?", notificationentity.IdempotencyStatusProcessing, now.Add(-processingTTL))
		}).
		Exec(ctx)
	if err != nil {
		log.Printf("ERROR: idempotency: cleanup failed. err: %v", err)
		return 0, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("ERROR: idempotency: failed to read cleanup rows affected. err: %v", err)
		return 0, err
	}
	if rows > 0 {
		log.Printf("INFO: idempotency: cleaned expired keys. rows: %d", rows)
	}
	return rows, nil
}

func StartIdempotencyCleanupJob(ctx context.Context, db *bun.DB, interval time.Duration, completedTTL time.Duration, processingTTL time.Duration) {
	if db == nil {
		return
	}
	if interval <= 0 {
		interval = time.Hour
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				_, _ = CleanupIdempotencyKeys(ctx, db, completedTTL, processingTTL)
			}
		}
	}()
}
