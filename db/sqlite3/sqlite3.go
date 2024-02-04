package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var Conn *sql.DB = NewClient()

func Execute(query string) (any, error) {
	res, err := Conn.Exec(query)
	return res, err
}

func NewClient() *sql.DB {
	db, err := sql.Open("sqlite3", "app.db")
	if err != nil {
		panic("Could not open database")
	}
	return db
}
