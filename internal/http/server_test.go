package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/peterkuchinov/The-link-shortener-on-Golang/internal/apperror"
	"github.com/peterkuchinov/The-link-shortener-on-Golang/internal/http/mocks"
	"go.uber.org/zap"
)

func TestServer_HTTP_Endpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	tests := []struct {
		name           string
		method         string
		target         string
		body           interface{}
		setupMock      func(m *mocks.MockLinkServiceShortener)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "201 Created on successful shorten",
			method: "POST",
			target: "/shorten",
			body:   gin.H{"url": "https://google.com", "custom_code": "googl"},
			setupMock: func(m *mocks.MockLinkServiceShortener) {
				m.ShortenFunc = func(ctx context.Context, url, customCode string) (string, error) {
					return "googl", nil
				}
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `"code":"googl"`,
		},
		{
			name:   "404 Not Found on missing code",
			method: "GET",
			target: "/fake-code",
			body:   nil,
			setupMock: func(m *mocks.MockLinkServiceShortener) {
				m.GetOriginalURLFunc = func(ctx context.Context, code string) (string, error) {
					return "", apperror.ErrNotFound
				}
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"Resource not found"}`,
		},
		{
			name:   "409 Conflict on busy custom code",
			method: "POST",
			target: "/shorten",
			body:   gin.H{"url": "https://google.com", "custom_code": "busy"},
			setupMock: func(m *mocks.MockLinkServiceShortener) {
				m.ShortenFunc = func(ctx context.Context, url, customCode string) (string, error) {
					return "", apperror.ErrCodeAlreadyExists
				}
			},
			expectedStatus: http.StatusConflict,
			expectedBody:   `{"error":"This short code is already taken"}`,
		},
		{
			name:           "400 Bad Request on broken JSON payload",
			method:         "POST",
			target:         "/shorten",
			body:           "this is definitely not a json string {",
			setupMock:      func(m *mocks.MockLinkServiceShortener) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid request body"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mocks.MockLinkServiceShortener{}
			tt.setupMock(mockSvc)

			serverWrapper := NewServer(":8080", "http://localhost:8080", logger, mockSvc)
			router := serverWrapper.srv.Handler

			var req *http.Request
			if strBody, ok := tt.body.(string); ok {
				req = httptest.NewRequest(tt.method, tt.target, bytes.NewBufferString(strBody))
			} else if tt.body != nil {
				jsonBytes, _ := json.Marshal(tt.body)
				req = httptest.NewRequest(tt.method, tt.target, bytes.NewBuffer(jsonBytes))
			} else {
				req = httptest.NewRequest(tt.method, tt.target, nil)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedBody != "" && !bytes.Contains(w.Body.Bytes(), []byte(tt.expectedBody)) {
				t.Errorf("Expected body to contain %s, got %s", tt.expectedBody, w.Body.String())
			}
		})
	}
}
