// init_db.go
package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./db/data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createClientsTable := `
    CREATE TABLE IF NOT EXISTS clients (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        clientName TEXT NOT NULL,
        parentName TEXT NOT NULL,
        address1 TEXT NOT NULL,
        address2 TEXT NOT NULL,
        phone TEXT NOT NULL,
        email TEXT NOT NULL,
        abbreviation TEXT NOT NULL,
        active BOOLEAN NOT NULL DEFAULT 1,
        invoice_lock BOOLEAN NOT NULL DEFAULT 0
    );
    `
	_, err = db.Exec(createClientsTable)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database initialized and users and clients tables created.")
}
