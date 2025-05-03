package db

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const pgxDrverName string = "pgx"

// DBHandler is a struct that wraps the sql.DB object.
//
// It provides methods to interact with the database.
type DBHandler struct {
	conn *sql.DB
}

// Initialize a new database handler.
func InitDB(pgURL string) *DBHandler {
	db, err := sql.Open(pgxDrverName, pgURL)
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}

	return &DBHandler{
		conn: db,
	}

}

// Checks the health of the database connection.
func (handler *DBHandler) GetDBHealth() {

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := handler.conn.PingContext(ctx)
	if err != nil {
		log.Fatal("Error pinging database. Check connectivity of connection: ", err)
	} else {
		log.Println("Database connection is healthy")
	}
}

// Defer the closing of the database connection.
// This should be called when the application is done using the database.
func (handler *DBHandler) Close() {
	if err := handler.conn.Close(); err != nil {
		log.Fatal("Error closing database: ", err)
	}
}
