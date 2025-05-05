package main

import (
	"errors"
	"ielts/internal/data"
	"ielts/internal/domain"
	"ielts/internal/filters"
	"ielts/internal/validator"
	"net/http"
)

func (app *application) likePhotoHandler(w http.ResponseWriter, r *http.Request) {

	user := app.contextGetUser(r)
	if user.IsAnonymous() {
		app.authenticationRequiredResponse(w, r)
		return
	}

	photo_id, err := app.readIDParam(w, r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	like := &domain.Like{
		UserID:  user.ID,
		PhotoID: photo_id,
	}

	v := validator.New()
	filters.ValidateLike(v, like)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	_, err = app.models.Photo.Get(photo_id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	exists, err := app.models.Like.CheckLike(user.ID, photo_id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if exists {
		app.badRequestResponse(w, r, errors.New("user has already liked this photo"))
		return
	}

	err = app.models.Like.Insert(like)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"like": like}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) unlikePhotoHandler(w http.ResponseWriter, r *http.Request) {

	user := app.contextGetUser(r)

	photo_id, err := app.readIDParam(w, r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	_, err = app.models.Photo.Get(photo_id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	exists, err := app.models.Like.CheckLike(user.ID, photo_id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !exists {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Like.Delete(user.ID, photo_id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "like successfully removed"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listLikedPhotosByUserHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	if user.IsAnonymous() {
		app.authenticationRequiredResponse(w, r)
		return
	}

	photos, err := app.models.Photo.GetPhotosLikedByUser(user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"photos": photos}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) checkLikeHandler(w http.ResponseWriter, r *http.Request) {

	user := app.contextGetUser(r)
	if user.IsAnonymous() {
		app.authenticationRequiredResponse(w, r)
		return
	}
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

	liked, err := app.models.Like.CheckLike(user.ID, photoID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"liked": liked}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
