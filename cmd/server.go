package main

import (
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
	// Initialize database
	db, err := internal.InitDB()
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize repository and service
	repo := internal.NewRepository(db)
	service := internal.NewService(repo).(*internal.Service)

	// Setup routes
	mux := http.NewServeMux()

	// GET /api/projects - Get all projects
	mux.HandleFunc("/api/projects", enableCORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			service.GetAllProjectsHandler(w, r)
		} else if r.Method == http.MethodPost {
			service.AddProjectHandler(w, r)
		} else {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	// Start server
	port := ":8080"
	log.Printf("Server starting on %s", port)

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
