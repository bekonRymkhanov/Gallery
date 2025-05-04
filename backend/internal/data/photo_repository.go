package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"ielts/internal/domain"
	"ielts/internal/filters"
	"ielts/internal/validator"
	"strings"
	"time"
)

type PhotoModel struct {
	DB *sql.DB
}

func (p PhotoModel) Insert(photo *domain.Photo) error {
	query := `INSERT INTO photos (title, description, author, category, tags, width, height, url, thumbnail_url, source, download_count, likes) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
				RETURNING id, version`

	args := []interface{}{
		photo.Title,
		photo.Description,
		photo.Author,
		photo.Category,
		photo.Tags,
		photo.Width,
		photo.Height,
		photo.URL,
		photo.ThumbnailURL,
		photo.Source,
		photo.DownloadCount,
		photo.Likes,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return p.DB.QueryRowContext(ctx, query, args...).Scan(&photo.ID, &photo.Version)
}

func (p PhotoModel) Get(id int64) (*domain.Photo, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `SELECT id, title, description, author, category, tags, width, height, url, thumbnail_url, source, download_count, likes, version
				FROM photos
				WHERE id = $1`
	var photo domain.Photo

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := p.DB.QueryRowContext(ctx, query, id).Scan(
		&photo.ID,
		&photo.Title,
		&photo.Description,
		&photo.Author,
		&photo.Category,
		&photo.Tags,
		&photo.Width,
		&photo.Height,
		&photo.URL,
		&photo.ThumbnailURL,
		&photo.Source,
		&photo.DownloadCount,
		&photo.Likes,
		&photo.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &photo, nil
}

func (p PhotoModel) Update(photo *domain.Photo) error {
	query := `UPDATE photos
				SET title = $1, description = $2, author = $3, category = $4, tags = $5, width = $6, height = $7, url = $8, thumbnail_url = $9, source = $10, download_count = $11, likes = $12, version = version + 1
				WHERE id = $13 AND version = $14
				RETURNING version`

	args := []interface{}{
		photo.Title,
		photo.Description,
		photo.Author,
		photo.Category,
		photo.Tags,
		photo.Width,
		photo.Height,
		photo.URL,
		photo.ThumbnailURL,
		photo.Source,
		photo.DownloadCount,
		photo.Likes,
		photo.ID,
		photo.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := p.DB.QueryRowContext(ctx, query, args...).Scan(&photo.Version)
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

func (p PhotoModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM photos
				WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := p.DB.ExecContext(ctx, query, id)
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

func (p PhotoModel) GetAll(mfilters filters.Filters, msearchOptions filters.PhotoSearch) ([]*domain.Photo, filters.Metadata, error) {
	where := []string{}
	args := []interface{}{}
	argPosition := 1

	if msearchOptions.Title != "" {
		where = append(where, fmt.Sprintf(
			"(to_tsvector('simple', title) @@ plainto_tsquery('simple', $%d) OR title ILIKE $%d)", argPosition, argPosition+1))
		args = append(args, msearchOptions.Title, "%"+msearchOptions.Title+"%")
		argPosition += 2
	}

	if msearchOptions.Author != "" {
		where = append(where, fmt.Sprintf(
			"(to_tsvector('simple', author) @@ plainto_tsquery('simple', $%d) OR author ILIKE $%d)", argPosition, argPosition+1))
		args = append(args, msearchOptions.Author, "%"+msearchOptions.Author+"%")
		argPosition += 2
	}

	if msearchOptions.Category != "" {
		where = append(where, fmt.Sprintf(
			"(to_tsvector('simple', category) @@ plainto_tsquery('simple', $%d) OR category ILIKE $%d)", argPosition, argPosition+1))
		args = append(args, msearchOptions.Category, "%"+msearchOptions.Category+"%")
		argPosition += 2
	}

	if msearchOptions.Tags != "" {
		where = append(where, fmt.Sprintf(
			"(to_tsvector('simple', tags) @@ plainto_tsquery('simple', $%d) OR tags ILIKE $%d)", argPosition, argPosition+1))
		args = append(args, msearchOptions.Tags, "%"+msearchOptions.Tags+"%")
		argPosition += 2
	}

	whereClause := ""
	if len(where) > 0 {
		whereClause = "WHERE " + strings.Join(where, " OR ")
	}

	query := fmt.Sprintf(`
	SELECT count(*) OVER(), id, title, description, author, category, tags, width, height, url, thumbnail_url, source, download_count, likes, version
	FROM photos
	%s
	ORDER BY %s %s, id ASC 
	LIMIT $%d OFFSET $%d`,
		whereClause,
		mfilters.SortColumn(),
		mfilters.SortDirection(),
		argPosition,
		argPosition+1)

	args = append(args, mfilters.Limit(), mfilters.Offset())
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := p.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, filters.Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	photos := []*domain.Photo{}

	for rows.Next() {
		var photo domain.Photo
		err := rows.Scan(
			&totalRecords,
			&photo.ID,
			&photo.Title,
			&photo.Description,
			&photo.Author,
			&photo.Category,
			&photo.Tags,
			&photo.Width,
			&photo.Height,
			&photo.URL,
			&photo.ThumbnailURL,
			&photo.Source,
			&photo.DownloadCount,
			&photo.Likes,
			&photo.Version,
		)
		if err != nil {
			return nil, filters.Metadata{}, err
		}
		photos = append(photos, &photo)
	}

	if err = rows.Err(); err != nil {
		return nil, filters.Metadata{}, err
	}

	metadata := filters.CalculateMetadata(totalRecords, mfilters.Page, mfilters.PageSize)

	return photos, metadata, nil
}

func (p PhotoModel) GetByCategory(category string) ([]*domain.Photo, error) {
	query := `
		SELECT id, title, description, author, category, tags, width, height, url, thumbnail_url, source, download_count, likes, version
		FROM photos
		WHERE category = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := p.DB.QueryContext(ctx, query, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []*domain.Photo

	for rows.Next() {
		var photo domain.Photo
		err := rows.Scan(
			&photo.ID,
			&photo.Title,
			&photo.Description,
			&photo.Author,
			&photo.Category,
			&photo.Tags,
			&photo.Width,
			&photo.Height,
			&photo.URL,
			&photo.ThumbnailURL,
			&photo.Source,
			&photo.DownloadCount,
			&photo.Likes,
			&photo.Version,
		)
		if err != nil {
			return nil, err
		}
		photos = append(photos, &photo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return photos, nil
}

func ValidatePhoto(v *validator.Validator, photo *domain.Photo) {
	v.Check(photo.Title != "", "title", "must be provided")
	v.Check(len(photo.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(photo.Author != "", "author", "must be provided")
	v.Check(photo.Category != "", "category", "must be provided")
	v.Check(photo.Width > 0, "width", "must be a positive integer")
	v.Check(photo.Height > 0, "height", "must be a positive integer")
	v.Check(photo.URL != "", "url", "must be provided")
	v.Check(photo.ThumbnailURL != "", "thumbnail_url", "must be provided")
	v.Check(photo.Source != "", "source", "must be provided")
}
