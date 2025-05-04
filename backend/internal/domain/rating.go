package domain

import "time"

type Rating struct {
	ID        int64     `json:"id"`
	PhotoID   int64     `json:"photo_id"`
	UserID    int64     `json:"user_id"`
	Score     int       `json:"score"` //1-5
	CreatedAt time.Time `json:"created_at"`
	Version   int32     `json:"version"`
}
