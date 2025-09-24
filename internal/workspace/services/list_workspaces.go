package workspace_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	workspace_core "github.com/gabrielmrtt/taski/internal/workspace"
	workspace_repositories "github.com/gabrielmrtt/taski/internal/workspace/repositories"
)

type ListWorkspacesService struct {
	WorkspaceRepository workspace_repositories.WorkspaceRepository
}

func NewListWorkspacesService(
	workspaceRepository workspace_repositories.WorkspaceRepository,
) *ListWorkspacesService {
	return &ListWorkspacesService{
		WorkspaceRepository: workspaceRepository,
	}
}

type ListWorkspacesInput struct {
	OrganizationIdentity core.Identity
	Filters              workspace_repositories.WorkspaceFilters
	SortInput            *core.SortInput
	Pagination           *core.PaginationInput
	RelationsInput       core.RelationsInput
	ShowDeleted          bool
}

func (i ListWorkspacesInput) Validate() error {
	return nil
}

func (s *ListWorkspacesService) Execute(input ListWorkspacesInput) (*core.PaginationOutput[workspace_core.WorkspaceDto], error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	input.Filters.OrganizationIdentity = input.OrganizationIdentity

	workspaces, err := s.WorkspaceRepository.PaginateWorkspacesBy(workspace_repositories.PaginateWorkspacesParams{
		Filters:     input.Filters,
		SortInput:   input.SortInput,
		Pagination:  input.Pagination,
		ShowDeleted: input.ShowDeleted,
	})
	if err != nil {
		return nil, err
	}

	var workspacesDto []workspace_core.WorkspaceDto = make([]workspace_core.WorkspaceDto, 0)
	for _, workspace := range workspaces.Data {
		workspacesDto = append(workspacesDto, *workspace_core.WorkspaceToDto(&workspace))
	}

	return &core.PaginationOutput[workspace_core.WorkspaceDto]{
		Data:    workspacesDto,
		Page:    workspaces.Page,
		HasMore: workspaces.HasMore,
		Total:   workspaces.Total,
	}, nil
}
