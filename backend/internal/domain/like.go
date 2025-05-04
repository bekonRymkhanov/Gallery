package domain

import (
	"time"
)

type Like struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	PhotoID   int64     `json:"photo_id"`
	CreatedAt time.Time `json:"created_at"`
	Version   int32     `json:"version"`
}
