package internal

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository struct {
	Db *sql.DB
}

func NewRepository(db *sql.DB) RepoInterface {
	return &Repository{
		Db: db,
	}
}

func (r Repository) InsertProject(ctx context.Context, project *Project) error {
	query := `
        INSERT INTO projects (name, description, github_url, owner_name, owner_avatar, language, stars, category, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id
    `

	err := r.Db.QueryRowContext(
		ctx,
		query,
		project.Name,
		project.Description,
		project.GithubURL,
		project.OwnerName,
		project.OwnerAvatar,
		project.Language,
		project.Stars,
		project.Category,
		project.CreatedAt,
	).Scan(&project.ID)

	if err != nil {
		return fmt.Errorf("failed to insert project: %w", err)
	}

	return nil
}

func (r Repository) GetAllProjects(ctx context.Context) ([]*Project, error) {
	query := `
		SELECT id, name, description, github_url, owner_name, owner_avatar, language, stars, category, created_at
		FROM projects
		ORDER BY stars DESC
	`

	rows, err := r.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query projects: %w", err)
	}
	defer rows.Close()

	var projects []*Project
	for rows.Next() {
		var project Project
		err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.Description,
			&project.GithubURL,
			&project.OwnerName,
			&project.OwnerAvatar,
			&project.Language,
			&project.Stars,
			&project.Category,
			&project.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}
		projects = append(projects, &project)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating projects: %w", err)
	}

	return projects, nil
}
