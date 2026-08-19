package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hoard-app/backend/internal/model"
	"github.com/hoard-app/backend/internal/repository"
)

func ListTags(repo *repository.TagRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetUserID(c)
		tags, total, err := repo.List(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, model.TagListResponse{
			Tags:  tags,
			Total: total,
		})
	}
}

func CreateTag(repo *repository.TagRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetUserID(c)
		var req model.CreateTagRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
			return
		}
		tag, err := repo.Create(c.Request.Context(), userID, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"data": tag})
	}
}

func DeleteTag(repo *repository.TagRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetUserID(c)
		tagID := c.Param("id")
		if err := repo.Delete(c.Request.Context(), userID, tagID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": "deleted"})
	}
}
