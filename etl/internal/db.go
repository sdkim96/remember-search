package internal

import (
	"database/sql"
)

func GetDB(connectionArgs string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionArgs)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
