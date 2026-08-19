package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/hoard-app/backend/internal/model"
	"github.com/hoard-app/backend/internal/repository"
)

func ListBookmarks(repo *repository.BookmarkRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetUserID(c)
		var params model.BookmarkListParams
		if err := c.ShouldBindQuery(&params); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query params"})
			return
		}
		if params.Limit <= 0 || params.Limit > 100 {
			params.Limit = 20
		}
		bookmarks, total, err := repo.List(c.Request.Context(), userID, params.Limit, params.Offset, params.ContentType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, model.BookmarkListResponse{
			Bookmarks: bookmarks,
			Total:     total,
		})
	}
}

func CreateBookmark(repo *repository.BookmarkRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetUserID(c)
		var req model.CreateBookmarkRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "url and title are required"})
			return
		}
		bookmark, err := repo.Create(c.Request.Context(), userID, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"data": bookmark})
	}
}

func GetBookmark(repo *repository.BookmarkRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetUserID(c)
		bookmarkID := c.Param("id")
		bookmark, err := repo.GetByID(c.Request.Context(), userID, bookmarkID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": bookmark})
	}
}

func UpdateBookmark(repo *repository.BookmarkRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetUserID(c)
		bookmarkID := c.Param("id")
		var req model.UpdateBookmarkRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		bookmark, err := repo.Update(c.Request.Context(), userID, bookmarkID, req)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": bookmark})
	}
}

func DeleteBookmark(repo *repository.BookmarkRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetUserID(c)
		bookmarkID := c.Param("id")
		if err := repo.Delete(c.Request.Context(), userID, bookmarkID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
