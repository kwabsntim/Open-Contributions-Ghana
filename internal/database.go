package internal

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func InitDB(databaseURL string, useLocal bool) (*sql.DB, error) {
	var db *sql.DB
	var err error

	if useLocal {
		// Use local SQLite for development
		db, err = sql.Open("sqlite3", "./my.db")
		if err != nil {
			return nil, fmt.Errorf("failed to open local database: %w", err)
		}
		fmt.Println("Connected to local SQLite database")
	} else {
		// Use Turso for production
		db, err = sql.Open("libsql", databaseURL)
		if err != nil {
			return nil, fmt.Errorf("failed to open Turso database: %w", err)
		}
		fmt.Println("Connected to Turso database")
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("Database connection verified")

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
