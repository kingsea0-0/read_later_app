package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/hoard-app/backend/internal/model"
)

type TagRepository struct {
	pool *pgxpool.Pool
}

func NewTagRepository(pool *pgxpool.Pool) *TagRepository {
	return &TagRepository{pool: pool}
}

func (r *TagRepository) List(ctx context.Context, userID string) ([]model.Tag, int, error) {
	countQuery := `SELECT COUNT(*) FROM tags WHERE user_id = $1`
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, name, COALESCE(color,''), created_at
		FROM tags WHERE user_id = $1
		ORDER BY name ASC
	`, userID)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tags []model.Tag
	for rows.Next() {
		var t model.Tag
		if err := rows.Scan(&t.ID, &t.UserID, &t.Name, &t.Color, &t.CreatedAt); err != nil {
			return nil, 0, err
		}
		tags = append(tags, t)
	}

	return tags, total, nil
}

func (r *TagRepository) Create(ctx context.Context, userID string, req model.CreateTagRequest) (*model.Tag, error) {
	var t model.Tag
	err := r.pool.QueryRow(ctx, `
		INSERT INTO tags (user_id, name, color)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, name, COALESCE(color,''), created_at
	`, userID, req.Name, req.Color).Scan(&t.ID, &t.UserID, &t.Name, &t.Color, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TagRepository) Delete(ctx context.Context, userID, tagID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM tags WHERE id = $1 AND user_id = $2`, tagID, userID)
	return err
}
