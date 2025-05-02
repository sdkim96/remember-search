package internal

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const pgxDrverName string = "pgx"

// DBHandler is a struct that wraps the sql.DB object.
//
// It provides methods to interact with the database.
type DBHandler struct {
	db *sql.DB
}

// Initialize a new database handler.
func InitDB(pgURL string) *DBHandler {
	db, err := sql.Open(pgxDrverName, pgURL)
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}

	return &DBHandler{
		db: db,
	}

}

// Checks the health of the database connection.
func (handler *DBHandler) GetDBHealth() error {
	return handler.db.Ping()
}

// Defer the closing of the database connection.
// This should be called when the application is done using the database.
func (handler *DBHandler) Close() {
	if err := handler.db.Close(); err != nil {
		log.Fatal("Error closing database: ", err)
	}
}
