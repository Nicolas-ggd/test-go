package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.homeDir))
	router.Handler(http.MethodGet, "/about", dynamic.ThenFunc(app.aboutDir))
	// signin and signup method need get and post method
	// get to render view, post to parse and insert data in database

	// use secure routes
	protectedRoutes := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodGet, "/signin", protectedRoutes.ThenFunc(app.signInDir))
	router.Handler(http.MethodGet, "/signup", protectedRoutes.ThenFunc(app.signUpDir))
	router.Handler(http.MethodPost, "/signup", protectedRoutes.ThenFunc(app.signUpDirPost))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}
