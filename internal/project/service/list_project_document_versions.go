package projectservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/project"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
)

type ListProjectDocumentVersionsService struct {
	ProjectRepository         projectrepo.ProjectRepository
	ProjectDocumentRepository projectrepo.ProjectDocumentRepository
}

func NewListProjectDocumentVersionsService(
	projectRepository projectrepo.ProjectRepository,
	projectDocumentRepository projectrepo.ProjectDocumentRepository,
) *ListProjectDocumentVersionsService {
	return &ListProjectDocumentVersionsService{
		ProjectRepository:         projectRepository,
		ProjectDocumentRepository: projectDocumentRepository,
	}
}

type ListProjectDocumentVersionsInput struct {
	Filters    projectrepo.ProjectDocumentVersionFilters
	SortInput  core.SortInput
	Pagination core.PaginationInput
}

func (i ListProjectDocumentVersionsInput) Validate() error {
	return nil
}

func (s *ListProjectDocumentVersionsService) Execute(input ListProjectDocumentVersionsInput) (*core.PaginationOutput[project.ProjectDocumentVersionDto], error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	projectDocumentVersions, err := s.ProjectDocumentRepository.PaginateProjectDocumentVersionsBy(projectrepo.PaginateProjectDocumentVersionsByParams{
		Filters:    input.Filters,
		SortInput:  input.SortInput,
		Pagination: input.Pagination,
	})
	if err != nil {
		return nil, err
	}

	var projectDocumentVersionsDto []project.ProjectDocumentVersionDto = make([]project.ProjectDocumentVersionDto, 0)
	for _, projectDocumentVersion := range projectDocumentVersions.Data {
		projectDocumentVersionsDto = append(projectDocumentVersionsDto, *project.ProjectDocumentVersionToDto(&projectDocumentVersion))
	}

	return &core.PaginationOutput[project.ProjectDocumentVersionDto]{
		Data:    projectDocumentVersionsDto,
		Page:    projectDocumentVersions.Page,
		HasMore: projectDocumentVersions.HasMore,
		Total:   projectDocumentVersions.Total,
	}, nil
}
