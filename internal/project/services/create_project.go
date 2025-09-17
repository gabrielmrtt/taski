package project_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	project_repositories "github.com/gabrielmrtt/taski/internal/project/repositories"
)

type CreateProjectService struct {
	ProjectRepository     project_repositories.ProjectRepository
	TransactionRepository core.TransactionRepository
}

func NewCreateProjectService(
	projectRepository project_repositories.ProjectRepository,
	transactionRepository core.TransactionRepository,
) *CreateProjectService {
	return &CreateProjectService{
		ProjectRepository:     projectRepository,
		TransactionRepository: transactionRepository,
	}
}

type CreateProjectInput struct {
	Name              string
	Description       string
	Color             string
	WorkspaceIdentity core.Identity
}
