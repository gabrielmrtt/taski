package project_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	project_core "github.com/gabrielmrtt/taski/internal/project"
	project_repositories "github.com/gabrielmrtt/taski/internal/project/repositories"
)

type GetProjectService struct {
	ProjectRepository project_repositories.ProjectRepository
}

func NewGetProjectService(projectRepository project_repositories.ProjectRepository) *GetProjectService {
	return &GetProjectService{
		ProjectRepository: projectRepository,
	}
}

type GetProjectInput struct {
	OrganizationIdentity core.Identity
	WorkspaceIdentity    core.Identity
	ProjectIdentity      core.Identity
}

func (i GetProjectInput) Validate() error {
	return nil
}

func (s *GetProjectService) Execute(input GetProjectInput) (*project_core.ProjectDto, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	project, err := s.ProjectRepository.GetProjectByIdentity(project_repositories.GetProjectByIdentityParams{
		ProjectIdentity:      input.ProjectIdentity,
		WorkspaceIdentity:    &input.WorkspaceIdentity,
		OrganizationIdentity: &input.OrganizationIdentity,
	})
	if err != nil {
		return nil, err
	}

	if project == nil {
		return nil, core.NewNotFoundError("project not found")
	}

	return project_core.ProjectToDto(project), nil
}
