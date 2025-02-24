package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

var server *Server = NewServer()

type Server struct {
	db *sql.DB
}

func NewServer() *Server {
	db, err := sql.Open("postgres", fmt.Sprintf(
		"database=task_tracker_db sslmode=disable user=%s password=%s host=postgres port=5432",
		os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD")))
	if err != nil {
		log.Fatal(err)
	}
	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(50)

	return &Server{db: db}
}
