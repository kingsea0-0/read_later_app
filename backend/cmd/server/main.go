package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hoard-app/backend/internal/config"
	"github.com/hoard-app/backend/internal/handler"
	"github.com/hoard-app/backend/internal/middleware"
	"github.com/hoard-app/backend/internal/repository"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("unable to create connection pool: %v", err)
	}
	defer pool.Close()

	bookmarkRepo := repository.NewBookmarkRepository(pool)
	tagRepo := repository.NewTagRepository(pool)

	r := gin.Default()
	r.Use(corsMiddleware())

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.POST("/api/v1/auth/session", handler.AuthSession(cfg))

	auth := r.Group("/api/v1", middleware.AuthMiddleware(cfg))
	{
		auth.DELETE("/auth/session", handler.AuthLogout())

		auth.GET("/bookmarks", handler.ListBookmarks(bookmarkRepo))
		auth.POST("/bookmarks", handler.CreateBookmark(bookmarkRepo))
		auth.GET("/bookmarks/:id", handler.GetBookmark(bookmarkRepo))
		auth.PUT("/bookmarks/:id", handler.UpdateBookmark(bookmarkRepo))
		auth.DELETE("/bookmarks/:id", handler.DeleteBookmark(bookmarkRepo))

		auth.GET("/tags", handler.ListTags(tagRepo))
		auth.POST("/tags", handler.CreateTag(tagRepo))
		auth.DELETE("/tags/:id", handler.DeleteTag(tagRepo))
	}

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("server starting on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(shutdownCtx)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "http://localhost:5173" || origin == "https://hoard.app" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
