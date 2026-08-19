package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hoard-app/backend/internal/config"
	"github.com/hoard-app/backend/internal/middleware"
)

const testSecret = "test-supabase-jwt-secret"

func init() {
	gin.SetMode(gin.TestMode)
}

func signToken(t *testing.T, secret, sub string, exp time.Time) string {
	t.Helper()
	claims := jwt.MapClaims{"exp": exp.Unix()}
	if sub != "" {
		claims["sub"] = sub
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}
	return signed
}

func TestAuthSession(t *testing.T) {
	cfg := &config.Config{SupabaseJWTSecret: testSecret}
	router := gin.New()
	router.POST("/api/v1/auth/session", AuthSession(cfg))

	doRequest := func(body string) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/session", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		return w
	}

	t.Run("valid access token authenticates", func(t *testing.T) {
		tok := signToken(t, testSecret, "user-abc", time.Now().Add(time.Hour))
		w := doRequest(`{"access_token":"` + tok + `"}`)

		if w.Code != http.StatusOK {
			t.Fatalf("got status %d; want 200. body=%s", w.Code, w.Body.String())
		}
		var resp struct {
			Data struct {
				Status string `json:"status"`
				UserID string `json:"user_id"`
			} `json:"data"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to parse response: %v", err)
		}
		if resp.Data.Status != "authenticated" {
			t.Errorf("got status %q; want authenticated", resp.Data.Status)
		}
		if resp.Data.UserID != "user-abc" {
			t.Errorf("got user_id %q; want user-abc", resp.Data.UserID)
		}
	})

	t.Run("missing access_token returns 400", func(t *testing.T) {
		w := doRequest(`{}`)
		if w.Code != http.StatusBadRequest {
			t.Errorf("got status %d; want 400", w.Code)
		}
	})

	t.Run("invalid token returns 401", func(t *testing.T) {
		tok := signToken(t, "wrong-secret", "user-abc", time.Now().Add(time.Hour))
		w := doRequest(`{"access_token":"` + tok + `"}`)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("got status %d; want 401", w.Code)
		}
	})
}

func TestAuthMiddlewareGuardsProtectedRoutes(t *testing.T) {
	cfg := &config.Config{SupabaseJWTSecret: testSecret}
	router := gin.New()
	protected := router.Group("/api/v1", middleware.AuthMiddleware(cfg))
	protected.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"user_id": GetUserID(c)})
	})

	t.Run("missing authorization header returns 401", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/ping", nil)
		router.ServeHTTP(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("got status %d; want 401", w.Code)
		}
	})

	t.Run("valid bearer token passes and injects user id", func(t *testing.T) {
		tok := signToken(t, testSecret, "user-xyz", time.Now().Add(time.Hour))
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/ping", nil)
		req.Header.Set("Authorization", "Bearer "+tok)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("got status %d; want 200", w.Code)
		}
		var resp struct {
			UserID string `json:"user_id"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to parse response: %v", err)
		}
		if resp.UserID != "user-xyz" {
			t.Errorf("got user_id %q; want user-xyz", resp.UserID)
		}
	})
}
