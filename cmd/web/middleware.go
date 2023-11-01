package main

import (
	"fmt"
	"net/http"
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
