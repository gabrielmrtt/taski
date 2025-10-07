package projectservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/project"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
)

type ListProjectTaskStatusesService struct {
	ProjectTaskStatusRepository projectrepo.ProjectTaskStatusRepository
}

func NewListProjectTaskStatusesService(
	projectTaskStatusRepository projectrepo.ProjectTaskStatusRepository,
) *ListProjectTaskStatusesService {
	return &ListProjectTaskStatusesService{
		ProjectTaskStatusRepository: projectTaskStatusRepository,
	}
}

type ListProjectTaskStatusesInput struct {
	Filters    projectrepo.ProjectTaskStatusFilters
	SortInput  core.SortInput
	Pagination core.PaginationInput
}

func (i ListProjectTaskStatusesInput) Validate() error {
	return nil
}

func (s *ListProjectTaskStatusesService) Execute(input ListProjectTaskStatusesInput) (*core.PaginationOutput[project.ProjectTaskStatusDto], error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	projectTaskStatuses, err := s.ProjectTaskStatusRepository.PaginateProjectTaskStatusesBy(projectrepo.PaginateProjectTaskStatusesParams{
		Filters:     input.Filters,
		SortInput:   input.SortInput,
		Pagination:  input.Pagination,
		ShowDeleted: false,
	})
	if err != nil {
		return nil, err
	}

	var projectTaskStatusesDto []project.ProjectTaskStatusDto = make([]project.ProjectTaskStatusDto, 0)
	for _, projectTaskStatus := range projectTaskStatuses.Data {
		projectTaskStatusesDto = append(projectTaskStatusesDto, *project.ProjectTaskStatusToDto(&projectTaskStatus))
	}

	return &core.PaginationOutput[project.ProjectTaskStatusDto]{
		Data:    projectTaskStatusesDto,
		Page:    projectTaskStatuses.Page,
		HasMore: projectTaskStatuses.HasMore,
		Total:   projectTaskStatuses.Total,
	}, nil
}
