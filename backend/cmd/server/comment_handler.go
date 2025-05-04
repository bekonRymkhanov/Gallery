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

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {

	user := app.contextGetUser(r)
	if user.IsAnonymous() {
		app.authenticationRequiredResponse(w, r)
		return
	}

	var input struct {
		PhotoID int64  `json:"photo_id"`
		Content string `json:"content"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	comment := &domain.Comment{
		PhotoID: input.PhotoID,
		UserID:  user.ID,
		Content: input.Content,
	}

	v := validator.New()
	filters.ValidateComment(v, comment)

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

	err = app.models.Comment.Insert(comment)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/comments/%d", comment.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"comment": comment}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showCommentHandler(w http.ResponseWriter, r *http.Request) {

	comment_id, err := app.readCommentIDParam(w, r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	comment, err := app.models.Comment.Get(comment_id)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"comment": comment}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateCommentHandler(w http.ResponseWriter, r *http.Request) {

	user := app.contextGetUser(r)
	if user.IsAnonymous() {
		app.authenticationRequiredResponse(w, r)
		return
	}

	comment_id, err := app.readCommentIDParam(w, r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	comment, err := app.models.Comment.Get(comment_id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if comment.UserID != user.ID {
		app.notPermittedResponse(w, r)
		return
	}

	var input struct {
		Content *string `json:"content"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Content != nil {
		comment.Content = *input.Content
	}

	v := validator.New()
	filters.ValidateComment(v, comment)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Comment.Update(comment)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"comment": comment}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteCommentHandler(w http.ResponseWriter, r *http.Request) {

	user := app.contextGetUser(r)
	if user.IsAnonymous() {
		app.authenticationRequiredResponse(w, r)
		return
	}

	comment_id, err := app.readCommentIDParam(w, r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	comment, err := app.models.Comment.Get(comment_id)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if comment.UserID != user.ID {
		app.notPermittedResponse(w, r)
		return
	}

	err = app.models.Comment.Delete(comment_id, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "comment successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listCommentsHandler(w http.ResponseWriter, r *http.Request) {

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

	var input struct {
		filters.Filters
	}
	v := validator.New()

	qs := r.URL.Query()

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "-created_at")

	input.Filters.SortSafelist = []string{"created_at", "-created_at", "id", "-id"}

	if filters.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	comments, metadata, err := app.models.Comment.GetAllForPhoto(photoID, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"comments": comments, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
