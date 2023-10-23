package main

import (
	"database/sql"
	"github.com/joefazee/go-api-tdd/pkg/security"
	"github.com/joefazee/go-api-tdd/pkg/store/sqlstore/postgres"
	"log"

	_ "github.com/lib/pq"
)

const (
	postgresDNS = "postgres://root:secret@localhost:5455/blog?sslmode=disable"
	driver      = "postgres"
	key         = "0e8eef14e9677512b7fbc06989040096"
)

func main() {

	srv, err := setup()
	if err != nil {
		log.Fatal(err)
	}

	if err = srv.run(":8000"); err != nil {
		log.Fatal(err)
	}
}

func setup() (*server, error) {
	db, err := connectToDB(driver)
	if err != nil {
		return nil, err
	}

	pStore := postgres.NewPostgresStore(db)

	newJWT, err := security.NewJWT(key)
	if err != nil {
		return nil, err
	}

	srv := newServer(pStore, newJWT)

	srv.setupRoutes()

	return srv, nil
}

func connectToDB(driver string) (*sql.DB, error) {
	db, err := sql.Open(driver, postgresDNS)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
