package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"ielts/internal/domain"
	"ielts/internal/filters"
	"time"
)

type CommentModel struct {
	DB *sql.DB
}

func (c CommentModel) Insert(comment *domain.Comment) error {
	query := `
		INSERT INTO comments (photo_id, user_id, content, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, version
	`

	args := []interface{}{
		comment.PhotoID,
		comment.UserID,
		comment.Content,
		time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.DB.QueryRowContext(ctx, query, args...).Scan(&comment.ID, &comment.Version)
}

func (c CommentModel) Get(id int64) (*domain.Comment, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, photo_id, user_id, content, created_at, version
		FROM comments
		WHERE id = $1
	`

	var comment domain.Comment

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := c.DB.QueryRowContext(ctx, query, id).Scan(
		&comment.ID,
		&comment.PhotoID,
		&comment.UserID,
		&comment.Content,
		&comment.CreatedAt,
		&comment.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &comment, nil
}

func (c CommentModel) GetAllForPhoto(photoID int64, mfilters filters.Filters) ([]*domain.Comment, filters.Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, photo_id, user_id, content, created_at, version
		FROM comments
		WHERE photo_id = $1
		ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3`,
		mfilters.SortColumn(),
		mfilters.SortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{photoID, mfilters.Limit(), mfilters.Offset()}
	rows, err := c.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, filters.Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	comments := []*domain.Comment{}

	for rows.Next() {
		var comment domain.Comment
		err := rows.Scan(
			&totalRecords,
			&comment.ID,
			&comment.PhotoID,
			&comment.UserID,
			&comment.Content,
			&comment.CreatedAt,
			&comment.Version,
		)
		if err != nil {
			return nil, filters.Metadata{}, err
		}
		comments = append(comments, &comment)
	}

	if err = rows.Err(); err != nil {
		return nil, filters.Metadata{}, err
	}

	metadata := filters.CalculateMetadata(totalRecords, mfilters.Page, mfilters.PageSize)

	return comments, metadata, nil
}

func (c CommentModel) Update(comment *domain.Comment) error {
	query := `
		UPDATE comments
		SET content = $1, version = version + 1
		WHERE id = $2 AND user_id = $3 AND version = $4
		RETURNING version
	`

	args := []interface{}{
		comment.Content,
		comment.ID,
		comment.UserID,
		comment.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := c.DB.QueryRowContext(ctx, query, args...).Scan(&comment.Version)
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

func (c CommentModel) Delete(id, userID int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM comments
		WHERE id = $1 AND user_id = $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := c.DB.ExecContext(ctx, query, id, userID)
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
