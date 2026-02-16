package internal

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TursoDatabaseURL string
	TursoAuthToken   string
	Port             string
	GitHubToken      string
	UseLocalDB       bool
}

// LoadConfig loads environment variables from .env file (if exists) and environment
func LoadConfig() *Config {
	// Load .env file if it exists (for local development)
	// In production (Render), environment variables are set directly
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	config := &Config{
		TursoDatabaseURL: getEnv("TURSO_DATABASE_URL", ""),
		TursoAuthToken:   getEnv("TURSO_AUTH_TOKEN", ""),
		Port:             getEnv("PORT", "8080"),
		GitHubToken:      getEnv("GITHUB_TOKEN", ""),
		UseLocalDB:       getEnv("USE_LOCAL_DB", "true") == "true",
	}

	// If not using local DB, validate Turso credentials
	if !config.UseLocalDB {
		if config.TursoDatabaseURL == "" {
			log.Fatal("TURSO_DATABASE_URL environment variable is required when USE_LOCAL_DB=false")
		}

		if config.TursoAuthToken == "" {
			log.Fatal("TURSO_AUTH_TOKEN environment variable is required when USE_LOCAL_DB=false")
		}
	}

	return config
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetDatabaseURL returns the formatted database URL for Turso
func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf("%s?authToken=%s", c.TursoDatabaseURL, c.TursoAuthToken)
}
