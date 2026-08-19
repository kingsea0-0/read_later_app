package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hoard-app/backend/internal/config"
)

var (
	// ErrMissingSecret is returned when the server has no JWT secret configured.
	ErrMissingSecret = errors.New("supabase jwt secret is not configured")
	// ErrInvalidToken is returned when a token fails signature or claim validation.
	ErrInvalidToken = errors.New("invalid token")
	// ErrMissingUserID is returned when a valid token carries no subject claim.
	ErrMissingUserID = errors.New("token is missing user id")
)

// ValidateToken parses and verifies a Supabase-issued JWT (HMAC/HS256) and
// returns the subject (user id) claim. It is the single source of truth for
// token validation, shared by the auth middleware and the session handler.
func ValidateToken(cfg *config.Config, tokenStr string) (string, error) {
	if cfg.SupabaseJWTSecret == "" {
		return "", ErrMissingSecret
	}
	if tokenStr == "" {
		return "", ErrInvalidToken
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(cfg.SupabaseJWTSecret), nil
	})
	if err != nil || !token.Valid {
		return "", ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrInvalidToken
	}

	userID, _ := claims["sub"].(string)
	if userID == "" {
		return "", ErrMissingUserID
	}

	return userID, nil
}

// AuthMiddleware rejects requests without a valid Bearer token and injects the
// authenticated user id into the gin context under "user_id".
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid authorization header"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		userID, err := ValidateToken(cfg, tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
