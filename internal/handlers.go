package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// parseGitHubURL extracts owner and repo name from a GitHub URL
// Supports formats:
// - https://github.com/owner/repo
// - https://github.com/owner/repo.git
// - github.com/owner/repo
func parseGitHubURL(url string) (owner, repo string, err error) {
	// Remove trailing .git if present
	url = strings.TrimSuffix(url, ".git")

	// Match github.com/owner/repo pattern
	re := regexp.MustCompile(`github\.com/([^/]+)/([^/]+)`)
	matches := re.FindStringSubmatch(url)

	if len(matches) != 3 {
		return "", "", fmt.Errorf("URL must be in format: github.com/owner/repo")
	}

	return matches[1], matches[2], nil
}

// GetAllProjectsHandler returns all projects from the database
func (s *Service) GetAllProjectsHandler(w http.ResponseWriter, r *http.Request) {
	projects, err := s.GetAllProjects(r.Context())
	if err != nil {
		log.Printf("GetAllProjectsHandler error: %v", err)
		http.Error(w, fmt.Sprintf("failed to get projects: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(projects); err != nil {
		http.Error(w, fmt.Sprintf("failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}

// AddProjectHandler fetches a GitHub repository and saves it to the database
func (s *Service) AddProjectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read raw body for logging and parse JSON
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("AddProjectHandler failed to read body: %v", err)
		http.Error(w, fmt.Sprintf("failed to read request body: %v", err), http.StatusBadRequest)
		return
	}

	// Log a concise representation of the incoming request for debugging
	log.Printf("AddProjectHandler: headers: User-Agent=%s, Content-Type=%s", r.Header.Get("User-Agent"), r.Header.Get("Content-Type"))
	log.Printf("AddProjectHandler: raw body: %s", string(bodyBytes))

	// Reject excessively large payloads
	const maxBodySize = 100 * 1024 // 100 KB
	if len(bodyBytes) > maxBodySize {
		log.Printf("AddProjectHandler: request body too large: %d bytes", len(bodyBytes))
		http.Error(w, "request body too large", http.StatusBadRequest)
		return
	}

	var req struct {
		GithubURL   string `json:"github_url"`
		Name        string `json:"name,omitempty"`
		OwnerName   string `json:"owner_name,omitempty"`
		OwnerAvatar string `json:"owner_avatar,omitempty"`
		Description string `json:"description,omitempty"`
		Language    string `json:"language,omitempty"`
		Stars       int    `json:"stars,omitempty"`
		CreatedAt   string `json:"created_at,omitempty"`
	}

	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		log.Printf("AddProjectHandler invalid body JSON: %v", err)
		http.Error(w, fmt.Sprintf("invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	if req.GithubURL == "" {
		log.Printf("AddProjectHandler missing github_url in request")
		http.Error(w, "github_url is required", http.StatusBadRequest)
		return
	}

	// Extract owner and repo name from GitHub URL
	owner, repoName, err := parseGitHubURL(req.GithubURL)
	if err != nil {
		log.Printf("AddProjectHandler invalid GitHub URL: %v", err)
		http.Error(w, fmt.Sprintf("invalid GitHub URL: %v", err), http.StatusBadRequest)
		return
	}

	// Helper to truncate long strings
	truncate := func(s string, n int) string {
		if len(s) <= n {
			return s
		}
		return s[:n]
	}

	// If the client provided full repo metadata in the request (from preview), use it and avoid calling GitHub again.
	if req.Name != "" && req.OwnerName != "" {
		// Basic validation and sanitization
		name := strings.TrimSpace(req.Name)
		ownerName := strings.TrimSpace(req.OwnerName)
		if name == "" || ownerName == "" {
			http.Error(w, "name and owner_name must be provided", http.StatusBadRequest)
			return
		}

		// Cap lengths
		name = truncate(name, 200)
		ownerName = truncate(ownerName, 200)
		description := truncate(strings.TrimSpace(req.Description), 2000)
		ownerAvatar := strings.TrimSpace(req.OwnerAvatar)
		if ownerAvatar != "" && !strings.HasPrefix(ownerAvatar, "https://") {
			// only allow https avatar URLs; drop otherwise
			ownerAvatar = ""
		}
		language := truncate(strings.TrimSpace(req.Language), 100)
		stars := req.Stars
		if stars < 0 {
			stars = 0
		}

		// Duplicate check by github_url
		existing, err := s.repo.GetProjectByGithubURL(r.Context(), req.GithubURL)
		if err != nil {
			log.Printf("AddProjectHandler duplicate check error: %v", err)
			http.Error(w, fmt.Sprintf("failed to check existing project: %v", err), http.StatusInternalServerError)
			return
		}
		if existing != nil {
			// Idempotent: return existing project
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(existing); err != nil {
				http.Error(w, fmt.Sprintf("failed to encode response: %v", err), http.StatusInternalServerError)
			}
			return
		}

		// Build project from provided fields
		proj := &Project{
			Name:        name,
			Description: description,
			GithubURL:   req.GithubURL,
			OwnerName:   ownerName,
			OwnerAvatar: ownerAvatar,
			Language:    language,
			Stars:       stars,
			Category:    "",
			CreatedAt:   time.Now(),
		}

		// If client provided a created_at timestamp, try to parse it
		if req.CreatedAt != "" {
			if t, err := time.Parse(time.RFC3339, req.CreatedAt); err == nil {
				proj.CreatedAt = t
			}
		}

		if err := s.repo.InsertProject(r.Context(), proj); err != nil {
			log.Printf("AddProjectHandler insert failed: %v", err)
			http.Error(w, fmt.Sprintf("failed to add project: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(proj); err != nil {
			http.Error(w, fmt.Sprintf("failed to encode response: %v", err), http.StatusInternalServerError)
			return
		}
		return
	}

	// Otherwise, fetch and save project from GitHub as before
	project, err := s.GetProject(r.Context(), owner, repoName)
	if err != nil {
		log.Printf("AddProjectHandler failed to add project: %v", err)

		// If GitHub responded with a 403 or mentioned rate limiting, return 429 to the client
		errMsg := err.Error()
		if strings.Contains(errMsg, "status 403") || strings.Contains(strings.ToLower(errMsg), "rate limit") {
			http.Error(w, fmt.Sprintf("failed to add project: %v", err), http.StatusTooManyRequests)
			return
		}

		http.Error(w, fmt.Sprintf("failed to add project: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(project); err != nil {
		http.Error(w, fmt.Sprintf("failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}
