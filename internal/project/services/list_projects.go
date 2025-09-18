package project_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	project_core "github.com/gabrielmrtt/taski/internal/project"
	project_repositories "github.com/gabrielmrtt/taski/internal/project/repositories"
	workspace_repositories "github.com/gabrielmrtt/taski/internal/workspace/repositories"
)

type ListProjectsService struct {
	ProjectRepository   project_repositories.ProjectRepository
	WorkspaceRepository workspace_repositories.WorkspaceRepository
}

func NewListProjectsService(
	projectRepository project_repositories.ProjectRepository,
	workspaceRepository workspace_repositories.WorkspaceRepository,
) *ListProjectsService {
	return &ListProjectsService{
		ProjectRepository:   projectRepository,
		WorkspaceRepository: workspaceRepository,
	}
}

type ListProjectsInput struct {
	WorkspaceIdentity    core.Identity
	OrganizationIdentity core.Identity
	Filters              project_repositories.ProjectFilters
	SortInput            *core.SortInput
	Pagination           *core.PaginationInput
}

func (i ListProjectsInput) Validate() error {
	return nil
}

func (s *ListProjectsService) Execute(input ListProjectsInput) (*core.PaginationOutput[project_core.ProjectDto], error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	workspace, err := s.WorkspaceRepository.GetWorkspaceByIdentity(workspace_repositories.GetWorkspaceByIdentityParams{
		WorkspaceIdentity:    input.WorkspaceIdentity,
		OrganizationIdentity: &input.OrganizationIdentity,
	})
	if err != nil {
		return nil, err
	}

	if workspace == nil {
		return nil, core.NewNotFoundError("workspace not found")
	}

	input.Filters.WorkspaceIdentity = input.WorkspaceIdentity

	projects, err := s.ProjectRepository.PaginateProjectsBy(project_repositories.PaginateProjectsParams{
		Filters:    input.Filters,
		SortInput:  input.SortInput,
		Pagination: input.Pagination,
	})
	if err != nil {
		return nil, err
	}

	var projectsDto []project_core.ProjectDto = make([]project_core.ProjectDto, 0)
	for _, project := range projects.Data {
		projectsDto = append(projectsDto, *project_core.ProjectToDto(&project))
	}

	return &core.PaginationOutput[project_core.ProjectDto]{
		Data:    projectsDto,
		Page:    projects.Page,
		HasMore: projects.HasMore,
		Total:   projects.Total,
	}, nil
}
