package projectservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	project "github.com/gabrielmrtt/taski/internal/project"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
)

type GetProjectService struct {
	ProjectRepository projectrepo.ProjectRepository
}

func NewGetProjectService(projectRepository projectrepo.ProjectRepository) *GetProjectService {
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

func (s *GetProjectService) Execute(input GetProjectInput) (*project.ProjectDto, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	prj, err := s.ProjectRepository.GetProjectByIdentity(projectrepo.GetProjectByIdentityParams{
		ProjectIdentity:      input.ProjectIdentity,
		WorkspaceIdentity:    &input.WorkspaceIdentity,
		OrganizationIdentity: &input.OrganizationIdentity,
	})
	if err != nil {
		return nil, err
	}

	if prj == nil {
		return nil, core.NewNotFoundError("project not found")
	}

	return project.ProjectToDto(prj), nil
}
