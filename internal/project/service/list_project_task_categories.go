package projectservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/project"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
)

type ListProjectTaskCategoriesService struct {
	ProjectTaskCategoryRepository projectrepo.ProjectTaskCategoryRepository
}

func NewListProjectTaskCategoriesService(
	projectTaskCategoryRepository projectrepo.ProjectTaskCategoryRepository,
) *ListProjectTaskCategoriesService {
	return &ListProjectTaskCategoriesService{
		ProjectTaskCategoryRepository: projectTaskCategoryRepository,
	}
}

type ListProjectTaskCategoriesInput struct {
	Filters        projectrepo.ProjectTaskCategoryFilters
	SortInput      core.SortInput
	Pagination     core.PaginationInput
	RelationsInput core.RelationsInput
}

func (i ListProjectTaskCategoriesInput) Validate() error {
	return nil
}

func (s *ListProjectTaskCategoriesService) Execute(input ListProjectTaskCategoriesInput) (*core.PaginationOutput[project.ProjectTaskCategoryDto], error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	projectTaskCategories, err := s.ProjectTaskCategoryRepository.PaginateProjectTaskCategoryBy(projectrepo.PaginateProjectTaskCategoryParams{
		Filters:        input.Filters,
		SortInput:      input.SortInput,
		Pagination:     input.Pagination,
		ShowDeleted:    false,
		RelationsInput: input.RelationsInput,
	})
	if err != nil {
		return nil, err
	}

	var projectTaskCategoriesDto []project.ProjectTaskCategoryDto = make([]project.ProjectTaskCategoryDto, 0)
	for _, projectTaskCategory := range projectTaskCategories.Data {
		projectTaskCategoriesDto = append(projectTaskCategoriesDto, *project.ProjectTaskCategoryToDto(&projectTaskCategory))
	}

	return &core.PaginationOutput[project.ProjectTaskCategoryDto]{
		Data:    projectTaskCategoriesDto,
		Page:    projectTaskCategories.Page,
		HasMore: projectTaskCategories.HasMore,
		Total:   projectTaskCategories.Total,
	}, nil
}
