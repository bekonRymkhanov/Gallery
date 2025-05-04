package main

import (
	"errors"
	"fmt"
	"ielts/internal/data"
	"ielts/internal/domain"
	"ielts/internal/filters"
	"ielts/internal/validator"
	"net/http"
)

func (app *application) createRatingHandler(w http.ResponseWriter, r *http.Request) {

	user := app.contextGetUser(r)
	if user.IsAnonymous() {
		app.authenticationRequiredResponse(w, r)
		return
	}

	var input struct {
		PhotoID int64 `json:"photo_id"`
		Score   int   `json:"score"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	rating := &domain.Rating{
		PhotoID: input.PhotoID,
		UserID:  user.ID,
		Score:   input.Score,
	}

	v := validator.New()
	filters.ValidateRating(v, rating)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	_, err = app.models.Photo.Get(input.PhotoID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	existingRating, err := app.models.Rating.GetByPhotoAndUser(input.PhotoID, user.ID)
	if err != nil && !errors.Is(err, data.ErrRecordNotFound) {
		app.serverErrorResponse(w, r, err)
		return
	}

	if existingRating != nil {
		app.failedValidationResponse(w, r, map[string]string{
			"rating": "you have already rated this photo",
		})
		return
	}

	err = app.models.Rating.Insert(rating)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	photoObj, err := app.models.Photo.Get(input.PhotoID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	photoObj.Likes++
	err = app.models.Photo.Update(photoObj)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/ratings/%d", rating.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"rating": rating}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showRatingHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(w, r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	rating, err := app.models.Rating.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"rating": rating}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateRatingHandler(w http.ResponseWriter, r *http.Request) {

	user := app.contextGetUser(r)
	if user.IsAnonymous() {
		app.authenticationRequiredResponse(w, r)
		return
	}

	id, err := app.readIDParam(w, r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	rating, err := app.models.Rating.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if rating.UserID != user.ID {
		app.notPermittedResponse(w, r)
		return
	}

	var input struct {
		Score *int `json:"score"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Score != nil {
		rating.Score = *input.Score
	}

	v := validator.New()
	filters.ValidateRating(v, rating)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Rating.Update(rating)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"rating": rating}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteRatingHandler(w http.ResponseWriter, r *http.Request) {

	user := app.contextGetUser(r)
	if user.IsAnonymous() {
		app.authenticationRequiredResponse(w, r)
		return
	}

	id, err := app.readIDParam(w, r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	rating, err := app.models.Rating.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if rating.UserID != user.ID {
		app.notPermittedResponse(w, r)
		return
	}

	err = app.models.Rating.Delete(id, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	photoObj, err := app.models.Photo.Get(rating.PhotoID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if photoObj.Likes > 0 {
		photoObj.Likes--
		err = app.models.Photo.Update(photoObj)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "rating successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listRatingsHandler(w http.ResponseWriter, r *http.Request) {

	photoID, err := app.readIDParam(w, r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	_, err = app.models.Photo.Get(photoID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	ratings, err := app.models.Rating.GetAllForPhoto(photoID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	avgScore, err := app.models.Rating.GetAverageScoreForPhoto(photoID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{
		"ratings":       ratings,
		"count":         len(ratings),
		"average_score": avgScore,
	}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listRatingOfUserHandler(w http.ResponseWriter, r *http.Request) {

	user := app.contextGetUser(r)
	if user.IsAnonymous() {
		app.authenticationRequiredResponse(w, r)
		return
	}
	id, err := app.readIDParam(w, r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	rating, err := app.models.Rating.GetByPhotoAndUser(id, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"rating": rating}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
