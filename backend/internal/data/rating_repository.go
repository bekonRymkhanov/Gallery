package data

import (
	"context"
	"database/sql"
	"errors"
	"ielts/internal/domain"
	"time"
)

type RatingModel struct {
	DB *sql.DB
}

func (r RatingModel) Insert(rating *domain.Rating) error {
	query := `
		INSERT INTO ratings (photo_id, user_id, score, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, version
	`

	args := []interface{}{
		rating.PhotoID,
		rating.UserID,
		rating.Score,
		time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return r.DB.QueryRowContext(ctx, query, args...).Scan(&rating.ID, &rating.Version)
}

func (r RatingModel) Get(id int64) (*domain.Rating, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, photo_id, user_id, score, created_at, version
		FROM ratings
		WHERE id = $1
	`

	var rating domain.Rating

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&rating.ID,
		&rating.PhotoID,
		&rating.UserID,
		&rating.Score,
		&rating.CreatedAt,
		&rating.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &rating, nil
}

func (r RatingModel) GetByPhotoAndUser(photoID, userID int64) (*domain.Rating, error) {
	query := `
		SELECT id, photo_id, user_id, score, created_at, version
		FROM ratings
		WHERE photo_id = $1 AND user_id = $2
	`

	var rating domain.Rating

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, photoID, userID).Scan(
		&rating.ID,
		&rating.PhotoID,
		&rating.UserID,
		&rating.Score,
		&rating.CreatedAt,
		&rating.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &rating, nil
}

func (r RatingModel) GetAllForPhoto(photoID int64) ([]*domain.Rating, error) {
	query := `
		SELECT id, photo_id, user_id, score, created_at, version
		FROM ratings
		WHERE photo_id = $1
		LIMIT 10
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.DB.QueryContext(ctx, query, photoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ratings []*domain.Rating

	for rows.Next() {
		var rating domain.Rating
		err := rows.Scan(
			&rating.ID,
			&rating.PhotoID,
			&rating.UserID,
			&rating.Score,
			&rating.CreatedAt,
			&rating.Version,
		)
		if err != nil {
			return nil, err
		}
		ratings = append(ratings, &rating)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ratings, nil
}

func (r RatingModel) Update(rating *domain.Rating) error {
	query := `
		UPDATE ratings
		SET score = $1, version = version + 1
		WHERE id = $2 AND user_id = $3 AND version = $4
		RETURNING version
	`

	args := []interface{}{
		rating.Score,
		rating.ID,
		rating.UserID,
		rating.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, args...).Scan(&rating.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (r RatingModel) Delete(id, userID int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM ratings
		WHERE id = $1 AND user_id = $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := r.DB.ExecContext(ctx, query, id, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (r RatingModel) GetAverageScoreForPhoto(photoID int64) (float64, error) {
	query := `
		SELECT COALESCE(AVG(score), 0)
		FROM ratings
		WHERE photo_id = $1
	`

	var avgScore float64

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, photoID).Scan(&avgScore)
	if err != nil {
		return 0, err
	}

	return avgScore, nil
}
