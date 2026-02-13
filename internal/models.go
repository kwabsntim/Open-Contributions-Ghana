package internal

import "time"

// Repo represents the GitHub API response
type Repo struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	HTMLURL         string `json:"html_url"`
	StargazersCount int    `json:"stargazers_count"`
	Language        string `json:"language"`
	Owner           struct {
		Login     string `json:"login"`
		AvatarURL string `json:"avatar_url"`
	} `json:"owner"`
	CreatedAt string `json:"created_at"`
}

type Project struct {
	ID          int       `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`                 // repo.Name
	Description string    `db:"description" json:"description"`   // repo.Description
	GithubURL   string    `db:"github_url" json:"github_url"`     // repo.HTMLURL
	OwnerName   string    `db:"owner_name" json:"owner_name"`     // repo.Owner.Login
	OwnerAvatar string    `db:"owner_avatar" json:"owner_avatar"` // repo.Owner.AvatarURL
	Language    string    `db:"language" json:"language"`         // repo.Language (e.g., Go, TypeScript)
	Stars       int       `db:"stars" json:"stars"`               // repo.StargazersCount
	Category    string    `db:"category" json:"category"`         // Local tag (e.g., "Web", "Fintech", "Health")
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}
