package data

import (
	"context"
	"database/sql"
	"errors"
	"ielts/internal/domain"
	"time"
)

type LikeModel struct {
	DB *sql.DB
}

func (l LikeModel) Insert(like *domain.Like) error {
	query := `
		INSERT INTO likes (user_id, photo_id, created_at)
		VALUES ($1, $2, $3)
		RETURNING id, version
	`

	args := []interface{}{
		like.UserID,
		like.PhotoID,
		time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := l.DB.QueryRowContext(ctx, query, args...).Scan(&like.ID, &like.Version)
	if err != nil {
		return err
	}

	updateQuery := `
		UPDATE photos
		SET likes = likes + 1, version = version + 1
		WHERE id = $1
		RETURNING likes
	`

	var newLikeCount int64
	err = l.DB.QueryRowContext(ctx, updateQuery, like.PhotoID).Scan(&newLikeCount)
	if err != nil {
		return err
	}

	return nil
}

func (l LikeModel) Delete(userID int64, photoID int64) error {
	query := `
		DELETE FROM likes
		WHERE user_id = $1 AND photo_id = $2
		RETURNING id
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int64
	err := l.DB.QueryRowContext(ctx, query, userID, photoID).Scan(&id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	updateQuery := `
		UPDATE photos
		SET likes = GREATEST(likes - 1, 0), version = version + 1
		WHERE id = $1
		RETURNING likes
	`

	var newLikeCount int64
	err = l.DB.QueryRowContext(ctx, updateQuery, photoID).Scan(&newLikeCount)
	if err != nil {
		return err
	}

	return nil
}

func (l LikeModel) CheckLike(userID int64, photoID int64) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM likes
			WHERE user_id = $1 AND photo_id = $2
		)
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var exists bool
	err := l.DB.QueryRowContext(ctx, query, userID, photoID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (l LikeModel) GetLikesByPhotoID(photoID int64) ([]*domain.Like, error) {
	query := `
		SELECT id, user_id, photo_id, created_at, version
		FROM likes
		WHERE photo_id = $1
		ORDER BY created_at DESC
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := l.DB.QueryContext(ctx, query, photoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var likes []*domain.Like

	for rows.Next() {
		var like domain.Like
		err := rows.Scan(
			&like.ID,
			&like.UserID,
			&like.PhotoID,
			&like.CreatedAt,
			&like.Version,
		)
		if err != nil {
			return nil, err
		}
		likes = append(likes, &like)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return likes, nil
}

func (l LikeModel) GetLikesByUserID(userID int64) ([]*domain.Like, error) {
	query := `
		SELECT id, user_id, photo_id, created_at, version
		FROM likes
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := l.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var likes []*domain.Like

	for rows.Next() {
		var like domain.Like
		err := rows.Scan(
			&like.ID,
			&like.UserID,
			&like.PhotoID,
			&like.CreatedAt,
			&like.Version,
		)
		if err != nil {
			return nil, err
		}
		likes = append(likes, &like)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return likes, nil
}
