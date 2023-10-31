package main

import (
	"log"
	"net/http"
)

func homeDir(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}

func aboutDir(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is about..."))
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", homeDir)
	mux.HandleFunc("/about", aboutDir)

	log.Println("Server starting on: 8080")

	err := http.ListenAndServe(":8080", mux)

	log.Fatal(err)
}
