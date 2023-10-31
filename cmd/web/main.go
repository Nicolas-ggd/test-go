package main

import (
	"log"
	"net/http"

	"github.com/thedevsaddam/renderer"
)

var rnd *renderer.Render

func init() {
	opts := renderer.Options{
		ParseGlobPattern: "../../ui/html/*.html",
	}

	rnd = renderer.New(opts)
}

func homeDir(w http.ResponseWriter, r *http.Request) {
	rnd.HTML(w, http.StatusOK, "home", nil)
}

func aboutDir(w http.ResponseWriter, r *http.Request) {
	rnd.HTML(w, http.StatusOK, "about", nil)
}

func signInDir(w http.ResponseWriter, r *http.Request) {
	rnd.HTML(w, http.StatusOK, "signin", nil)
}

func signUpDir(w http.ResponseWriter, r *http.Request) {
	rnd.HTML(w, http.StatusOK, "signup", nil)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", homeDir)
	mux.HandleFunc("/about", aboutDir)
	mux.HandleFunc("/signin", signInDir)
	mux.HandleFunc("/signup", signUpDir)

	log.Println("Server starting on: 8080")

	err := http.ListenAndServe(":8080", mux)

	log.Fatal(err)
}
