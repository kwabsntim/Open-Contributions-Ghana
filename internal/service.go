package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Service struct {
	repo RepoInterface
}

func NewService(repo RepoInterface) ProjectService {
	return &Service{
		repo: repo,
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

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repository: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("repo not found: status %d", resp.StatusCode)
	}

	var repo Repo
	if err := json.NewDecoder(resp.Body).Decode(&repo); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

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

	// Save to database
	if err := s.repo.InsertProject(ctx, project); err != nil {
		return nil, fmt.Errorf("failed to save project: %w", err)
	}

	return project, nil
}
