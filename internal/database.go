package internal

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./my.db")
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to SQLite")

	// Test connection
	var version string
	if err := db.QueryRow("SELECT sqlite_version()").Scan(&version); err != nil {
		return nil, err
	}

	fmt.Println("SQLite version:", version)

	// Create tables
	if err := CreateTables(db); err != nil {
		return nil, err
	}

	return db, nil
}

func CreateTables(db *sql.DB) error {
	createProjectsTable := `
	CREATE TABLE IF NOT EXISTS projects (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT,
		github_url TEXT NOT NULL,
		owner_name TEXT NOT NULL,
		owner_avatar TEXT,
		language TEXT,
		stars INTEGER DEFAULT 0,
		category TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.Exec(createProjectsTable); err != nil {
		return err
	}
	return nil
}
