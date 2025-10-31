package projectservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	project "github.com/gabrielmrtt/taski/internal/project"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
)

type ListProjectsService struct {
	ProjectRepository projectrepo.ProjectRepository
}

func NewListProjectsService(
	projectRepository projectrepo.ProjectRepository,
) *ListProjectsService {
	return &ListProjectsService{
		ProjectRepository: projectRepository,
	}
}

type ListProjectsInput struct {
	Filters        projectrepo.ProjectFilters
	SortInput      core.SortInput
	Pagination     core.PaginationInput
	RelationsInput core.RelationsInput
}

func (i ListProjectsInput) Validate() error {
	return nil
}

func (s *ListProjectsService) Execute(input ListProjectsInput) (*core.PaginationOutput[project.ProjectDto], error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	prjs, err := s.ProjectRepository.PaginateProjectsBy(projectrepo.PaginateProjectsParams{
		Filters:        input.Filters,
		SortInput:      input.SortInput,
		Pagination:     input.Pagination,
		ShowDeleted:    false,
		RelationsInput: input.RelationsInput,
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
