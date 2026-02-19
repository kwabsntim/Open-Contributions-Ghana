package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Service struct {
	repo        RepoInterface
	githubToken string
}

func NewService(repo RepoInterface) ProjectService {
	// Load GitHub token from environment if available
	config := LoadConfig()
	return &Service{
		repo:        repo,
		githubToken: config.GitHubToken,
	}
}
func (s *Service) GetProject(ctx context.Context, owner, reponame string) (*Project, error) {
	// Fetch from GitHub API
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, reponame)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	// Add GitHub token if available for higher rate limits
	if s.githubToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.githubToken))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repository: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Read response body to surface GitHub error messages (e.g., rate limit)
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("repo not found: status %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var repo Repo
	if err := json.NewDecoder(resp.Body).Decode(&repo); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Debug: Log the repo data
	fmt.Printf("Fetched repo: Name=%s, Owner=%s, Stars=%d\n", repo.Name, repo.Owner.Login, repo.StargazersCount)

	// Map GitHub Repo to Project model
	project := &Project{
		Name:        repo.Name,
		Description: repo.Description,
		GithubURL:   repo.HTMLURL,
		OwnerName:   repo.Owner.Login,
		OwnerAvatar: repo.Owner.AvatarURL,
		Language:    repo.Language,
		Stars:       repo.StargazersCount,
		Category:    "", // TODO: Add logic to determine category
		CreatedAt:   time.Now(),
	}

	// Debug: Log the project before insert
	fmt.Printf("Inserting project: Name=%s, Owner=%s, URL=%s\n", project.Name, project.OwnerName, project.GithubURL)

	// Save to database
	if err := s.repo.InsertProject(ctx, project); err != nil {
		return nil, fmt.Errorf("failed to save project: %w", err)
	}

	return project, nil
}
func (s *Service) GetAllProjects(ctx context.Context) ([]*Project, error) {
	return s.repo.GetAllProjects(ctx)
}
