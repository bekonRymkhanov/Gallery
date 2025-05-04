package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	//photo routes
	router.HandlerFunc(http.MethodGet, "/photos", app.listPhotosHandler)
	router.HandlerFunc(http.MethodPost, "/photos", app.requirePermission("movies:write", app.createPhotoHandler))
	router.HandlerFunc(http.MethodGet, "/photos/:id", app.showPhotoHandler)
	router.HandlerFunc(http.MethodPatch, "/photo/:id", app.requirePermission("movies:write", app.updatePhotoHandler))
	router.HandlerFunc(http.MethodDelete, "/photo/:id", app.requirePermission("movies:write", app.deletePhotoHandler))

	// comment routes
	router.HandlerFunc(http.MethodGet, "/photos/:id/comments", app.listCommentsHandler)
	router.HandlerFunc(http.MethodGet, "/comments/:comment_id", app.showCommentHandler)
	router.HandlerFunc(http.MethodPost, "/comments", app.requireAuthenticatedUser(app.createCommentHandler))
	router.HandlerFunc(http.MethodPatch, "/comments/:comment_id", app.requireAuthenticatedUser(app.updateCommentHandler))
	router.HandlerFunc(http.MethodDelete, "/comments/:comment_id", app.requireAuthenticatedUser(app.deleteCommentHandler))

	// Rating routes
	router.HandlerFunc(http.MethodPost, "/ratings", app.requireAuthenticatedUser(app.createRatingHandler))
	router.HandlerFunc(http.MethodGet, "/ratings/:id", app.showRatingHandler)
	router.HandlerFunc(http.MethodPatch, "/ratings/:id", app.requireAuthenticatedUser(app.updateRatingHandler))
	router.HandlerFunc(http.MethodDelete, "/ratings/:id", app.requireAuthenticatedUser(app.deleteRatingHandler))
	router.HandlerFunc(http.MethodGet, "/photos/:id/ratings", app.listRatingsHandler)
	router.HandlerFunc(http.MethodGet, "/photo/:id/rating", app.requireAuthenticatedUser(app.listRatingOfUserHandler))

	// Like routes
	router.HandlerFunc(http.MethodPost, "/photos/:id/likes", app.requireAuthenticatedUser(app.likePhotoHandler))
	router.HandlerFunc(http.MethodDelete, "/photos/:id/likes", app.requireAuthenticatedUser(app.unlikePhotoHandler))
	router.HandlerFunc(http.MethodGet, "/users/likes", app.listLikesByUserHandler)
	router.HandlerFunc(http.MethodGet, "/photos/:id/likes/check", app.requireAuthenticatedUser(app.checkLikeHandler))

	// User routes
	router.HandlerFunc(http.MethodPost, "/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/tokens/authentication", app.createAuthenticationTokenHandler)

	//healthcheck route
	router.HandlerFunc(http.MethodGet, "/healthcheck", app.healthcheckHandler)

	return app.enableCORS(app.recoverPanic(app.rateLimit(app.authenticate(router))))

}
