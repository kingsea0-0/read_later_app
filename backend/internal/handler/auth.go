package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hoard-app/backend/internal/config"
	"github.com/hoard-app/backend/internal/middleware"
)

// AuthSession validates a Supabase access token supplied by the client and,
// on success, returns the authenticated user id. Clients (web / extension)
// call this once after OAuth to confirm the backend accepts their session.
func AuthSession(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			AccessToken string `json:"access_token" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "access_token is required"})
			return
		}

		userID, err := middleware.ValidateToken(cfg, req.AccessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid access token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": gin.H{
			"status":  "authenticated",
			"user_id": userID,
		}})
	}
}

func AuthLogout() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		// Invalidate session if needed
		c.JSON(http.StatusOK, gin.H{"data": gin.H{"user_id": userID}})
	}
}

func GetUserID(c *gin.Context) string {
	uid, _ := c.Get("user_id")
	return uid.(string)
}

func parseBearerToken(c *gin.Context) string {
	auth := c.GetHeader("Authorization")
	return strings.TrimPrefix(auth, "Bearer ")
}
