package repository

import (
	"context"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/hoard-app/backend/internal/model"
)

type BookmarkRepository struct {
	pool *pgxpool.Pool
}

func NewBookmarkRepository(pool *pgxpool.Pool) *BookmarkRepository {
	return &BookmarkRepository{pool: pool}
}

func (r *BookmarkRepository) List(ctx context.Context, userID string, limit, offset int, contentType string) ([]model.Bookmark, int, error) {
	args := []interface{}{userID}
	whereClause := "user_id = $1 AND is_archived = false"

	if contentType != "" {
		args = append(args, contentType)
		whereClause += " AND content_type = $" + strconv.Itoa(len(args))
	}

	countQuery := "SELECT COUNT(*) FROM bookmarks WHERE " + whereClause
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, limit, offset)
	limitIdx := len(args) - 1
	offsetIdx := len(args)

	query := `SELECT id, user_id, url, title, COALESCE(description,''), COALESCE(favicon_url,''),
	       COALESCE(og_image_url,''), COALESCE(site_name,''), content_type, is_archived, is_favorite,
	       read_at, created_at, updated_at
	FROM bookmarks WHERE ` + whereClause + ` ORDER BY created_at DESC LIMIT $` + strconv.Itoa(limitIdx) + ` OFFSET $` + strconv.Itoa(offsetIdx)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var bookmarks []model.Bookmark
	for rows.Next() {
		var b model.Bookmark
		if err := rows.Scan(&b.ID, &b.UserID, &b.URL, &b.Title, &b.Description,
			&b.FaviconURL, &b.OGImageURL, &b.SiteName, &b.ContentType, &b.IsArchived, &b.IsFavorite,
			&b.ReadAt, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, 0, err
		}
		bookmarks = append(bookmarks, b)
	}

	return bookmarks, total, nil
}

func (r *BookmarkRepository) Create(ctx context.Context, userID string, req model.CreateBookmarkRequest) (*model.Bookmark, error) {
	var b model.Bookmark
	err := r.pool.QueryRow(ctx, `
		INSERT INTO bookmarks (user_id, url, title, content_type)
		VALUES ($1, $2, $3, COALESCE(NULLIF($4, ''), 'article'))
		RETURNING id, user_id, url, title, COALESCE(description,''), COALESCE(favicon_url,''),
		          COALESCE(og_image_url,''), COALESCE(site_name,''), content_type, is_archived, is_favorite,
		          read_at, created_at, updated_at
	`, userID, req.URL, req.Title, req.ContentType).Scan(
		&b.ID, &b.UserID, &b.URL, &b.Title, &b.Description,
		&b.FaviconURL, &b.OGImageURL, &b.SiteName, &b.ContentType, &b.IsArchived, &b.IsFavorite,
		&b.ReadAt, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BookmarkRepository) GetByID(ctx context.Context, userID, bookmarkID string) (*model.Bookmark, error) {
	var b model.Bookmark
	err := r.pool.QueryRow(ctx, `
		SELECT id, user_id, url, title, COALESCE(description,''), COALESCE(favicon_url,''),
		       COALESCE(og_image_url,''), COALESCE(site_name,''), content_type, is_archived, is_favorite,
		       read_at, created_at, updated_at
		FROM bookmarks
		WHERE id = $1 AND user_id = $2
	`, bookmarkID, userID).Scan(
		&b.ID, &b.UserID, &b.URL, &b.Title, &b.Description,
		&b.FaviconURL, &b.OGImageURL, &b.SiteName, &b.ContentType, &b.IsArchived, &b.IsFavorite,
		&b.ReadAt, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BookmarkRepository) Update(ctx context.Context, userID, bookmarkID string, req model.UpdateBookmarkRequest) (*model.Bookmark, error) {
	// Build dynamic update
	fields := []string{}
	args := []interface{}{}
	argIdx := 1

	if req.Title != nil {
		fields = append(fields, "title = $"+itoa(argIdx))
		args = append(args, *req.Title)
		argIdx++
	}
	if req.IsArchived != nil {
		fields = append(fields, "is_archived = $"+itoa(argIdx))
		args = append(args, *req.IsArchived)
		argIdx++
	}
	if req.IsFavorite != nil {
		fields = append(fields, "is_favorite = $"+itoa(argIdx))
		args = append(args, *req.IsFavorite)
		argIdx++
	}
	if req.IsArchived != nil && *req.IsArchived {
		fields = append(fields, "read_at = $"+itoa(argIdx))
		args = append(args, time.Now())
		argIdx++
	}
	fields = append(fields, "updated_at = now()")
	args = append(args, bookmarkID, userID)

	query := `UPDATE bookmarks SET ` + joinFields(fields) + ` WHERE id = $` + itoa(argIdx) + ` AND user_id = $` + itoa(argIdx+1) + `
		RETURNING id, user_id, url, title, COALESCE(description,''), COALESCE(favicon_url,''),
		          COALESCE(og_image_url,''), COALESCE(site_name,''), content_type, is_archived, is_favorite,
		          read_at, created_at, updated_at`

	var b model.Bookmark
	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&b.ID, &b.UserID, &b.URL, &b.Title, &b.Description,
		&b.FaviconURL, &b.OGImageURL, &b.SiteName, &b.ContentType, &b.IsArchived, &b.IsFavorite,
		&b.ReadAt, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BookmarkRepository) Delete(ctx context.Context, userID, bookmarkID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM bookmarks WHERE id = $1 AND user_id = $2`, bookmarkID, userID)
	return err
}

func itoa(i int) string {
	return strconv.Itoa(i)
}

func joinFields(fields []string) string {
	result := ""
	for i, f := range fields {
		if i > 0 {
			result += ", "
		}
		result += f
	}
	return result
}
