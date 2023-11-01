package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// The main purpose of CSP is to mitigate and detect XSS attacks. XSS attacks exploit the browser's trust in the content received from the server.
		// The victim's browser is exposed to execution of malicious scripts, because the browser trusts the source of the content.
		// -------WARNING---------
		// It's not a safe way to use 'unsafe-inlane' and 'unsafe-eval' methods, but for now, it's work. for future when application size increase
		// then its become a problem
		res.Header().Set("Content-Security-Policy",
			"default-src 'self' https://getbootstrap.com; style-src 'self' fonts.googleapis.com 'unsafe-inline' https://getbootstrap.com; font-src 'self' fonts.gstatic.com https://getbootstrap.com; script-src 'self' 'unsafe-inline' 'unsafe-eval' http://www.google.com https://getbootstrap.com https://ajax.googleapis.com;")

		// A policy that controls how much information is shared through the HTTP Referer header. Helps to protect user privacy.
		res.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		// This header is used to block browsers' MIME type sniffing,
		// which can transform non-executable MIME types into executable MIME types (MIME Confusion Attacks).
		res.Header().Set("X-Content-Type-Options", "nosniff")
		// X-Frame-Options allows content publishers to prevent their own content from being used in an invisible frame by attackers.
		res.Header().Set("X-Frame-Options", "deny")
		// The HTTP X-XSS-Protection response header is a feature of Internet Explorer,
		// Chrome, and Safari that stops pages from loading when they detect reflected cross-site scripting (XSS) attacks.
		res.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(res, req)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", req.RemoteAddr, req.Proto, req.Method, req.URL.RequestURI())

		next.ServeHTTP(res, req)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				res.Header().Set("Connection", "close")
				app.serverError(res, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(res, req)
	})
}

func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

		if id == 0 {
			next.ServeHTTP(w, r)

			return
		}

		exists, err := app.users.UserExists(id)
		if err != nil {
			app.serverError(w, err)

			return
		}

		if exists {
			ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check if user is authenticated or not, return from the middleware
		// chain so that no subsequent handlers in the cain are executed.
		if !app.isAuthenticated(r) {
			app.sessionManager.Put(r.Context(), "redirectPathAfterLogin", r.URL.Path)
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		// Otherwise set the "Cache-Control: no-store" header so that pages
		// require authentication are not stored in the users browser cache (or
		// other intermediary cache).
		w.Header().Add("Cache-Control", "no-store")

		// call next hanlder in the chain
		next.ServeHTTP(w, r)
	})
}
