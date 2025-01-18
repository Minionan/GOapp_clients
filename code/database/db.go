// code/database/db.go
package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// Initialize sets up the database connection
func Initialize(dbPath string) error {
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	log.Println("Database connection successful")
	return nil
}

// Close closes the database connection
func Close() error {
	return db.Close()
}

// GetDB returns the database instance
func GetDB() *sql.DB {
	return db
}
