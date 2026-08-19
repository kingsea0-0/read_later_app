package middleware

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hoard-app/backend/internal/config"
)

const testSecret = "test-supabase-jwt-secret"

// signToken builds an HS256 JWT signed with the given secret and claims.
func signToken(t *testing.T, secret string, claims jwt.MapClaims) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}
	return signed
}

func TestValidateToken(t *testing.T) {
	cfg := &config.Config{SupabaseJWTSecret: testSecret}

	t.Run("valid token returns user id", func(t *testing.T) {
		tok := signToken(t, testSecret, jwt.MapClaims{
			"sub": "user-123",
			"exp": time.Now().Add(time.Hour).Unix(),
		})
		uid, err := ValidateToken(cfg, tok)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if uid != "user-123" {
			t.Errorf("got user id %q; want %q", uid, "user-123")
		}
	})

	t.Run("wrong secret is rejected", func(t *testing.T) {
		tok := signToken(t, "some-other-secret", jwt.MapClaims{
			"sub": "user-123",
			"exp": time.Now().Add(time.Hour).Unix(),
		})
		if _, err := ValidateToken(cfg, tok); err != ErrInvalidToken {
			t.Errorf("got err %v; want ErrInvalidToken", err)
		}
	})

	t.Run("expired token is rejected", func(t *testing.T) {
		tok := signToken(t, testSecret, jwt.MapClaims{
			"sub": "user-123",
			"exp": time.Now().Add(-time.Hour).Unix(),
		})
		if _, err := ValidateToken(cfg, tok); err != ErrInvalidToken {
			t.Errorf("got err %v; want ErrInvalidToken", err)
		}
	})

	t.Run("token without sub is rejected", func(t *testing.T) {
		tok := signToken(t, testSecret, jwt.MapClaims{
			"exp": time.Now().Add(time.Hour).Unix(),
		})
		if _, err := ValidateToken(cfg, tok); err != ErrMissingUserID {
			t.Errorf("got err %v; want ErrMissingUserID", err)
		}
	})

	t.Run("empty token is rejected", func(t *testing.T) {
		if _, err := ValidateToken(cfg, ""); err != ErrInvalidToken {
			t.Errorf("got err %v; want ErrInvalidToken", err)
		}
	})

	t.Run("garbage token is rejected", func(t *testing.T) {
		if _, err := ValidateToken(cfg, "not-a-jwt"); err != ErrInvalidToken {
			t.Errorf("got err %v; want ErrInvalidToken", err)
		}
	})

	t.Run("missing secret is rejected", func(t *testing.T) {
		emptyCfg := &config.Config{SupabaseJWTSecret: ""}
		tok := signToken(t, testSecret, jwt.MapClaims{"sub": "user-123"})
		if _, err := ValidateToken(emptyCfg, tok); err != ErrMissingSecret {
			t.Errorf("got err %v; want ErrMissingSecret", err)
		}
	})
}
