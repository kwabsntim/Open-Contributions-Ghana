package internal

import "context"

type ProjectService interface {
	GetProject(ctx context.Context, owner, reponame string) (*Project, error)
}

type RepoInterface interface {
	InsertProject(ctx context.Context, project *Project) error
	GetAllProjects(ctx context.Context) ([]*Project, error)
}
