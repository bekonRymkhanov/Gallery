package domain

import "time"

type Comment struct {
	ID        int64     `json:"id"`
	PhotoID   int64     `json:"photo_id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Version   int32     `json:"version"`
}
