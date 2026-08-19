package model

import "time"

type Bookmark struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	URL         string     `json:"url"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	FaviconURL  string     `json:"favicon_url,omitempty"`
	OGImageURL  string     `json:"og_image_url,omitempty"`
	SiteName    string     `json:"site_name,omitempty"`
	ContentType string     `json:"content_type"`
	IsArchived  bool       `json:"is_archived"`
	IsFavorite  bool       `json:"is_favorite"`
	ReadAt      *time.Time `json:"read_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateBookmarkRequest struct {
	URL         string `json:"url" binding:"required"`
	Title       string `json:"title" binding:"required"`
	ContentType string `json:"content_type,omitempty"`
}

type UpdateBookmarkRequest struct {
	Title       *string `json:"title,omitempty"`
	IsArchived  *bool   `json:"is_archived,omitempty"`
	IsFavorite  *bool   `json:"is_favorite,omitempty"`
}

type BookmarkListParams struct {
	Limit       int    `form:"limit" default:"20"`
	Offset      int    `form:"offset" default:"0"`
	ContentType string `form:"type"`
}

type BookmarkListResponse struct {
	Bookmarks []Bookmark `json:"bookmarks"`
	Total     int        `json:"total"`
}
