package internal

import "log"

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
