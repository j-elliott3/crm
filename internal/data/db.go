package data

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

const dbPath = "crm.db"


func OpenDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}