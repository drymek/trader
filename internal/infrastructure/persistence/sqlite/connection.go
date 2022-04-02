package sqlite

import (
	"database/sql"
	"errors"
)

var ErrCannotOpenDatabase = errors.New("Cannot open database")

func NewConnection(filename string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, ErrCannotOpenDatabase
	}

	return db, nil
}
