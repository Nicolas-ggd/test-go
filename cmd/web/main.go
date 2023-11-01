package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/thedevsaddam/renderer"
	"test.nicolas.net/pkg/models"
)

var rnd *renderer.Render

type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	formDecoder    *form.Decoder
	users          models.UserModelInterface
	sessionManager *scs.SessionManager
}

func init() {
	opts := renderer.Options{
		ParseGlobPattern: "../../ui/html/*.html",
	}

	rnd = renderer.New(opts)
}

func main() {
	addr := flag.String("addr", ":8080", "HTTP network address")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(&sql.DB{})
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		formDecoder:    formDecoder,
		users:          &models.UserModel{DB: &sql.DB{}}, // maybe its not a good version how it write(users)
		sessionManager: sessionManager,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	log.Println("Server starting on: 8080")

	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)

}
