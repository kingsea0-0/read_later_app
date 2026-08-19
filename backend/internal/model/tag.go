package model

import "time"

type Tag struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Color     string    `json:"color,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateTagRequest struct {
	Name  string `json:"name" binding:"required"`
	Color string `json:"color,omitempty"`
}

type TagListResponse struct {
	Tags  []Tag `json:"tags"`
	Total int   `json:"total"`
}
