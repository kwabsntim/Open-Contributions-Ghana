package internal

import (
	"context"
	"database/sql"
	"fmt"
	"log"
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
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
        RETURNING id
    `

	// First, try INSERT ... RETURNING id (works with SQLite newer versions and some SQL engines)
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

	if err == nil {
		return nil
	}

	// Log the initial error and attempt a fallback INSERT without RETURNING.
	log.Printf("InsertProject: RETURNING insert failed, attempting Exec fallback: %v", err)

	execQuery := `
		INSERT INTO projects (name, description, github_url, owner_name, owner_avatar, language, stars, category, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, execErr := r.Db.ExecContext(
		ctx,
		execQuery,
		project.Name,
		project.Description,
		project.GithubURL,
		project.OwnerName,
		project.OwnerAvatar,
		project.Language,
		project.Stars,
		project.Category,
		project.CreatedAt,
	)

	if execErr == nil {
		if id64, idErr := res.LastInsertId(); idErr == nil {
			project.ID = int(id64)
			return nil
		}
		// If LastInsertId not supported, return success without ID set.
		return nil
	}

	// Both attempts failed
	log.Printf("InsertProject Exec fallback error: %v", execErr)
	return fmt.Errorf("failed to insert project (returning: %v, exec: %v)", err, execErr)
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
