package projectservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/project"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
)

type ListProjectDocumentsService struct {
	ProjectRepository         projectrepo.ProjectRepository
	ProjectDocumentRepository projectrepo.ProjectDocumentRepository
}

func NewListProjectDocumentsService(
	projectRepository projectrepo.ProjectRepository,
	projectDocumentRepository projectrepo.ProjectDocumentRepository,
) *ListProjectDocumentsService {
	return &ListProjectDocumentsService{
		ProjectRepository:         projectRepository,
		ProjectDocumentRepository: projectDocumentRepository,
	}
}

type ListProjectDocumentsInput struct {
	Filters    projectrepo.ProjectDocumentVersionManagerFilters
	SortInput  core.SortInput
	Pagination core.PaginationInput
}

func (i ListProjectDocumentsInput) Validate() error {
	return nil
}

func (s *ListProjectDocumentsService) Execute(input ListProjectDocumentsInput) (*core.PaginationOutput[project.ProjectDocumentVersionDto], error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	projectDocumentVersionManagers, err := s.ProjectDocumentRepository.PaginateProjectDocumentVersionManagersBy(projectrepo.PaginateProjectDocumentVersionManagersByParams{
		Filters:    input.Filters,
		SortInput:  input.SortInput,
		Pagination: input.Pagination,
	})
	if err != nil {
		return nil, err
	}

	var projectDocumentVersionsDto []project.ProjectDocumentVersionDto = make([]project.ProjectDocumentVersionDto, 0)
	for _, projectDocumentVersionManager := range projectDocumentVersionManagers.Data {
		projectDocumentVersionsDto = append(projectDocumentVersionsDto, *project.ProjectDocumentVersionToDto(projectDocumentVersionManager.LatestVersion))
	}

	return &core.PaginationOutput[project.ProjectDocumentVersionDto]{
		Data:    projectDocumentVersionsDto,
		Page:    projectDocumentVersionManagers.Page,
		HasMore: projectDocumentVersionManagers.HasMore,
		Total:   projectDocumentVersionManagers.Total,
	}, nil
}
