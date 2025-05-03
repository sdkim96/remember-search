package internal

import (
	"context"
	"log"
)

// Get remembers from the database.
func (handler *DBHandler) GetUsers() {
	rows, err := handler.db.Query("SELECT id, name, email FROM remeber")
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

func (handler *DBHandler) GetUsersWithCtx(ctx context.Context) error {
	rows, err := handler.db.QueryContext(ctx, "SELECT id, name, email FROM remeber")
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
