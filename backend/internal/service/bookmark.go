package service

import (
	"context"
	"fmt"

	"github.com/hoard-app/backend/internal/model"
	"github.com/hoard-app/backend/internal/repository"
)

type BookmarkService struct {
	repo *repository.BookmarkRepository
}

func NewBookmarkService(repo *repository.BookmarkRepository) *BookmarkService {
	return &BookmarkService{repo: repo}
}

func (s *BookmarkService) List(ctx context.Context, userID string, limit, offset int, contentType string) ([]model.Bookmark, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return s.repo.List(ctx, userID, limit, offset, contentType)
}

func (s *BookmarkService) Create(ctx context.Context, userID, url, title string) (*model.Bookmark, error) {
	if url == "" || title == "" {
		return nil, fmt.Errorf("url and title are required")
	}
	return s.repo.Create(ctx, userID, model.CreateBookmarkRequest{
		URL:   url,
		Title: title,
	})
}
