package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.homeDir))
	router.Handler(http.MethodGet, "/about", dynamic.ThenFunc(app.aboutDir))
	// signin and signup method need get and post method
	// get to render view, post to parse and insert data in database
	router.Handler(http.MethodGet, "/signin", dynamic.ThenFunc(app.signInDir))
	router.Handler(http.MethodGet, "/signup", dynamic.ThenFunc(app.signUpDir))
	router.Handler(http.MethodPost, "/signup", dynamic.ThenFunc(app.signUpDirPost))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}
