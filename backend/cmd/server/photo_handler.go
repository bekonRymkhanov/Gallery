package main

import (
	"ielts/internal/data"
	"ielts/internal/domain"
	"ielts/internal/filters"
	"ielts/internal/validator"

	"errors"
	"fmt"
	"net/http"
)

func (app *application) createPhotoHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title         string `json:"title"`
		Description   string `json:"description"`
		Author        string `json:"author"`
		Category      string `json:"category"`
		Tags          string `json:"tags"`
		Width         int    `json:"width"`
		Height        int    `json:"height"`
		URL           string `json:"url"`
		ThumbnailURL  string `json:"thumbnail_url"`
		Source        string `json:"source"`
		DownloadCount int64  `json:"download_count"`
		Likes         int64  `json:"likes"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	photo := &domain.Photo{
		Title:         input.Title,
		Description:   input.Description,
		Author:        input.Author,
		Category:      input.Category,
		Tags:          input.Tags,
		Width:         input.Width,
		Height:        input.Height,
		URL:           input.URL,
		ThumbnailURL:  input.ThumbnailURL,
		Source:        input.Source,
		DownloadCount: input.DownloadCount,
		Likes:         input.Likes,
	}

	if data.ValidatePhoto(v, photo); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Photo.Insert(photo)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/photos/%d", photo.ID))
	err = app.writeJSON(w, http.StatusCreated, envelope{"photo": photo}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showPhotoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(w, r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	photo, err := app.models.Photo.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"photo": photo}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updatePhotoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(w, r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	photo, err := app.models.Photo.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Title         *string `json:"title"`
		Description   *string `json:"description"`
		Author        *string `json:"author"`
		Category      *string `json:"category"`
		Tags          *string `json:"tags"`
		Width         *int    `json:"width"`
		Height        *int    `json:"height"`
		URL           *string `json:"url"`
		ThumbnailURL  *string `json:"thumbnail_url"`
		Source        *string `json:"source"`
		DownloadCount *int64  `json:"download_count"`
		Likes         *int64  `json:"likes"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Title != nil {
		photo.Title = *input.Title
	}
	if input.Description != nil {
		photo.Description = *input.Description
	}
	if input.Author != nil {
		photo.Author = *input.Author
	}
	if input.Category != nil {
		photo.Category = *input.Category
	}
	if input.Tags != nil {
		photo.Tags = *input.Tags
	}
	if input.Width != nil {
		photo.Width = *input.Width
	}
	if input.Height != nil {
		photo.Height = *input.Height
	}
	if input.URL != nil {
		photo.URL = *input.URL
	}
	if input.ThumbnailURL != nil {
		photo.ThumbnailURL = *input.ThumbnailURL
	}
	if input.Source != nil {
		photo.Source = *input.Source
	}
	if input.DownloadCount != nil {
		photo.DownloadCount = *input.DownloadCount
	}
	if input.Likes != nil {
		photo.Likes = *input.Likes
	}

	v := validator.New()

	if data.ValidatePhoto(v, photo); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Photo.Update(photo)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"photo": photo}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deletePhotoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(w, r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Photo.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "photo successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listPhotosHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		filters.PhotoSearch
		filters.Filters
	}
	v := validator.New()

	qs := r.URL.Query()

	input.PhotoSearch.Title = app.readString(qs, "title", "")
	input.PhotoSearch.Author = app.readString(qs, "author", "")
	input.PhotoSearch.Category = app.readString(qs, "category", "")
	input.PhotoSearch.Tags = app.readString(qs, "tags", "")

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")

	input.Filters.SortSafelist = []string{"id", "title", "author", "category", "width", "height", "download_count", "likes", "-id", "-title", "-author", "-category", "-width", "-height", "-download_count", "-likes"}

	if filters.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	photos, metadata, err := app.models.Photo.GetAll(input.Filters, input.PhotoSearch)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"photos": photos, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
