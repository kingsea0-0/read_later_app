package service

import (
	"context"

	"github.com/hoard-app/backend/internal/model"
	"github.com/hoard-app/backend/internal/repository"
)

type TagService struct {
	repo *repository.TagRepository
}

func NewTagService(repo *repository.TagRepository) *TagService {
	return &TagService{repo: repo}
}

func (s *TagService) List(ctx context.Context, userID string) ([]model.Tag, int, error) {
	return s.repo.List(ctx, userID)
}
