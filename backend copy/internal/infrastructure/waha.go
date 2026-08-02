package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/backend/config"
)

type WAHA struct {
	client  *http.Client
	baseURL string
	apiKey  string
	session string
}

type SendTextRequest struct {
	ChatID  string
	Session string
	Text    string
}

type SendMediaRequest struct {
	ChatID   string
	Session  string
	URL      string
	FileName string
	Caption  string
}

func NewWAHA(cfg *config.Config) *WAHA {
	timeout := time.Duration(cfg.WAHA.TimeoutSecs) * time.Second
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	return &WAHA{
		baseURL: strings.TrimRight(cfg.WAHA.URL, "/"),
		apiKey:  strings.TrimSpace(cfg.WAHA.APIKey),
		session: strings.TrimSpace(cfg.WAHA.Session),
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (w *WAHA) SendText(ctx context.Context, req SendTextRequest) error {
	req.ChatID = strings.TrimSpace(req.ChatID)
	req.Session = strings.TrimSpace(req.Session)
	req.Text = strings.TrimSpace(req.Text)

	if req.ChatID == "" {
		return errors.New("waha: chat id wajib diisi")
	}

	if req.Text == "" {
		return errors.New("waha: text wajib diisi")
	}

	payload := map[string]string{
		"chatId":  req.ChatID,
		"text":    req.Text,
		"session": w.sessionName(ctx, req.Session),
	}

	return w.request(ctx, http.MethodPost, "/api/sendText", payload)
}

func (w *WAHA) SendImage(ctx context.Context, req SendMediaRequest) error {
	req.ChatID = strings.TrimSpace(req.ChatID)
	req.Session = strings.TrimSpace(req.Session)
	req.URL = strings.TrimSpace(req.URL)
	req.Caption = strings.TrimSpace(req.Caption)

	if req.ChatID == "" {
		return errors.New("waha: chat id wajib diisi")
	}

	if req.URL == "" {
		return errors.New("waha: image url wajib diisi")
	}

	payload := map[string]any{
		"chatId":  req.ChatID,
		"session": w.sessionName(ctx, req.Session),
		"media": map[string]string{
			"url": req.URL,
		},
		"caption": req.Caption,
	}

	return w.request(ctx, http.MethodPost, "/api/sendImage", payload)
}

func (w *WAHA) SendDocument(ctx context.Context, req SendMediaRequest) error {
	req.ChatID = strings.TrimSpace(req.ChatID)
	req.Session = strings.TrimSpace(req.Session)
	req.URL = strings.TrimSpace(req.URL)
	req.FileName = strings.TrimSpace(req.FileName)
	req.Caption = strings.TrimSpace(req.Caption)

	if req.ChatID == "" {
		return errors.New("waha: chat id wajib diisi")
	}

	if req.URL == "" {
		return errors.New("waha: document url wajib diisi")
	}

	payload := map[string]any{
		"chatId":  req.ChatID,
		"session": w.sessionName(ctx, req.Session),
		"media": map[string]string{
			"url":      req.URL,
			"filename": req.FileName,
		},
		"caption": req.Caption,
	}

	return w.request(ctx, http.MethodPost, "/api/sendDocument", payload)
}

func (w *WAHA) sessionName(ctx context.Context, session string) string {
	session = strings.TrimSpace(session)
	if session != "" {
		return session
	}
	return w.session
}

func (w *WAHA) request(
	ctx context.Context,
	method string,
	endpoint string,
	payload any,
) error {

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("waha: marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		method,
		w.baseURL+endpoint,
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("waha: create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	if w.apiKey != "" {
		req.Header.Set("X-Api-Key", w.apiKey)
		req.Header.Set("Authorization", "Bearer "+w.apiKey)
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("waha: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		raw, _ := io.ReadAll(resp.Body)

		return fmt.Errorf(
			"waha: %s %s returned %d: %s",
			method,
			endpoint,
			resp.StatusCode,
			strings.TrimSpace(string(raw)),
		)
	}

	return nil
}

// GetSessions menampilkan daftar session WAHA.
func (w *WAHA) GetSessions(ctx context.Context) ([]map[string]any, error) {
	resp, err := w.requestRaw(ctx, http.MethodGet, "/api/sessions", nil)
	if err != nil {
		return nil, err
	}
	var sessions []map[string]any
	if err := json.Unmarshal(resp, &sessions); err != nil {
		return nil, fmt.Errorf("waha: unmarshal sessions: %w", err)
	}
	return sessions, nil
}

// GetSessionStatus mengambil status satu session.
func (w *WAHA) GetSessionStatus(ctx context.Context, session string) (map[string]any, error) {
	if session == "" {
		session = w.session
	}
	resp, err := w.requestRaw(ctx, http.MethodGet, "/api/sessions/"+session, nil)
	if err != nil {
		return nil, err
	}
	var status map[string]any
	if err := json.Unmarshal(resp, &status); err != nil {
		return nil, fmt.Errorf("waha: unmarshal session status: %w", err)
	}
	return status, nil
}

// StartSession memulai session WAHA.
func (w *WAHA) StartSession(ctx context.Context, session string) error {
	if session == "" {
		session = w.session
	}
	_, err := w.requestRaw(ctx, http.MethodPost, "/api/sessions/"+session+"/start", nil)
	return err
}

// StopSession menghentikan session WAHA.
func (w *WAHA) StopSession(ctx context.Context, session string) error {
	if session == "" {
		session = w.session
	}
	_, err := w.requestRaw(ctx, http.MethodDelete, "/api/sessions/"+session, nil)
	return err
}

// GetQR mengambil QR code untuk session (return base64 or raw).
func (w *WAHA) GetQR(ctx context.Context, session string) (string, error) {
	if session == "" {
		session = w.session
	}
	resp, err := w.requestRaw(ctx, http.MethodGet, "/api/sessions/"+session+"/qr", nil)
	if err != nil {
		return "", err
	}
	qrStr := strings.TrimSpace(string(resp))
	return qrStr, nil
}

// requestRaw adalah versi request tanpa json encode/decode payload.
func (w *WAHA) requestRaw(ctx context.Context, method, endpoint string, payload any) ([]byte, error) {
	var bodyReader io.Reader
	if payload != nil {
		body, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("waha: marshal payload: %w", err)
		}
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, w.baseURL+endpoint, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("waha: create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if w.apiKey != "" {
		req.Header.Set("X-Api-Key", w.apiKey)
		req.Header.Set("Authorization", "Bearer "+w.apiKey)
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("waha: request: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("waha: read body: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("waha: %s %s returned %d: %s", method, endpoint, resp.StatusCode, strings.TrimSpace(string(raw)))
	}

	return raw, nil
}
