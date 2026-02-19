package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"open-source-ghana/internal"
)

// CORS middleware
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func main() {
	// Load configuration from environment variables
	config := internal.LoadConfig()

	// Initialize database (local SQLite or Turso)
	var db *sql.DB
	var err error

	if config.UseLocalDB {
		db, err = internal.InitDB("", true)
		log.Println("Using local SQLite database for development")
	} else {
		db, err = internal.InitDB(config.GetDatabaseURL(), false)
		log.Printf("Using Turso database: %s", config.TursoDatabaseURL)
	}

	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize repository and service
	repo := internal.NewRepository(db)
	service := internal.NewService(repo).(*internal.Service)

	// Setup routes
	mux := http.NewServeMux()

	// API endpoints
	mux.HandleFunc("/api/projects", enableCORS(func(w http.ResponseWriter, r *http.Request) {
		// OPTIONS is already handled by CORS middleware
		if r.Method == http.MethodGet {
			service.GetAllProjectsHandler(w, r)
		} else if r.Method == http.MethodPost {
			service.AddProjectHandler(w, r)
		} else if r.Method != http.MethodOptions {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	// Start server
	port := fmt.Sprintf(":%s", config.Port)
	log.Printf("Server starting on port %s", config.Port)

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
