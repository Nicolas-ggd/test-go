package models

import (
	"database/sql"
	"flag"
	"log"
)

var db *sql.DB

func init() {
	dsn := flag.String("dsn", "postgres://postgres:lothbrook11@localhost/postgres?sslmode=disable", "PostgreSQL data source name")

	var err error
	db, err = openDB(*dsn)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
