package projectservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	project "github.com/gabrielmrtt/taski/internal/project"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
	workspacerepo "github.com/gabrielmrtt/taski/internal/workspace/repository"
)

type ListProjectsService struct {
	ProjectRepository   projectrepo.ProjectRepository
	WorkspaceRepository workspacerepo.WorkspaceRepository
}

func NewListProjectsService(
	projectRepository projectrepo.ProjectRepository,
	workspaceRepository workspacerepo.WorkspaceRepository,
) *ListProjectsService {
	return &ListProjectsService{
		ProjectRepository:   projectRepository,
		WorkspaceRepository: workspaceRepository,
	}
}

type ListProjectsInput struct {
	WorkspaceIdentity    core.Identity
	OrganizationIdentity core.Identity
	Filters              projectrepo.ProjectFilters
	SortInput            *core.SortInput
	Pagination           *core.PaginationInput
}

func (i ListProjectsInput) Validate() error {
	return nil
}

func (s *ListProjectsService) Execute(input ListProjectsInput) (*core.PaginationOutput[project.ProjectDto], error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	wrk, err := s.WorkspaceRepository.GetWorkspaceByIdentity(workspacerepo.GetWorkspaceByIdentityParams{
		WorkspaceIdentity:    input.WorkspaceIdentity,
		OrganizationIdentity: &input.OrganizationIdentity,
	})
	if err != nil {
		return nil, err
	}

	if wrk == nil {
		return nil, core.NewNotFoundError("workspace not found")
	}

	input.Filters.WorkspaceIdentity = input.WorkspaceIdentity

	prjs, err := s.ProjectRepository.PaginateProjectsBy(projectrepo.PaginateProjectsParams{
		Filters:    input.Filters,
		SortInput:  input.SortInput,
		Pagination: input.Pagination,
	})
	if err != nil {
		return nil, err
	}

	var projectsDto []project.ProjectDto = make([]project.ProjectDto, 0)
	for _, prj := range prjs.Data {
		projectsDto = append(projectsDto, *project.ProjectToDto(&prj))
	}

	return &core.PaginationOutput[project.ProjectDto]{
		Data:    projectsDto,
		Page:    prjs.Page,
		HasMore: prjs.HasMore,
		Total:   prjs.Total,
	}, nil
}
