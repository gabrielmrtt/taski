package workspaceservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/workspace"
	workspacerepo "github.com/gabrielmrtt/taski/internal/workspace/repository"
)

type ListWorkspacesService struct {
	WorkspaceRepository workspacerepo.WorkspaceRepository
}

func NewListWorkspacesService(
	workspaceRepository workspacerepo.WorkspaceRepository,
) *ListWorkspacesService {
	return &ListWorkspacesService{
		WorkspaceRepository: workspaceRepository,
	}
}

type ListWorkspacesInput struct {
	Filters        workspacerepo.WorkspaceFilters
	SortInput      core.SortInput
	Pagination     core.PaginationInput
	RelationsInput core.RelationsInput
}

func (i ListWorkspacesInput) Validate() error {
	return nil
}

func (s *ListWorkspacesService) Execute(input ListWorkspacesInput) (*core.PaginationOutput[workspace.WorkspaceDto], error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	workspaces, err := s.WorkspaceRepository.PaginateWorkspacesBy(workspacerepo.PaginateWorkspacesParams{
		Filters:        input.Filters,
		SortInput:      input.SortInput,
		Pagination:     input.Pagination,
		ShowDeleted:    false,
		RelationsInput: input.RelationsInput,
	})
	if err != nil {
		return nil, err
	}

	var workspacesDto []workspace.WorkspaceDto = make([]workspace.WorkspaceDto, 0)
	for _, wrk := range workspaces.Data {
		workspacesDto = append(workspacesDto, *workspace.WorkspaceToDto(&wrk))
	}

	return &core.PaginationOutput[workspace.WorkspaceDto]{
		Data:    workspacesDto,
		Page:    workspaces.Page,
		HasMore: workspaces.HasMore,
		Total:   workspaces.Total,
	}, nil
}
