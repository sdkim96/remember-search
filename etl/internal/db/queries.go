package db

import (
	"context"
	"log"
	"time"
)

// Get remembers from the database.
func (h *DBHandler) GetUsers() {
	rows, err := h.conn.Query("SELECT id, name, email FROM remeber")
	if err != nil {
		log.Fatal("Error querying users: ", err)
	}

	for rows.Next() {
		var id int
		var name string
		var email string
		err := rows.Scan(&id, &name, &email)
		if err != nil {
			log.Fatal("Error scanning user: ", err)
		}
		log.Printf("User: %d, %s, %s\n", id, name, email)
	}
}

func (h *DBHandler) GetUsersWithCtx() error {

	// 10 seconds timeout
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second*10,
	)
	defer cancel()

	// Check if the context is done
	select {
	case <-ctx.Done():
		log.Println("Context done")
		return ctx.Err()
	default:
		log.Println("Context not done")
	}

	rows, err := h.conn.QueryContext(ctx, "SELECT id, name, email FROM remeber")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name, email string
		if err := rows.Scan(&id, &name, &email); err != nil {
			return err
		}
		log.Printf("User: %d, %s, %s\n", id, name, email)
	}

	if err := rows.Err(); err != nil {
		return err
	}
	return nil
}
