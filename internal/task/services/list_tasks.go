package taskservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/task"
	taskrepo "github.com/gabrielmrtt/taski/internal/task/repository"
)

type ListTasksService struct {
	TaskRepository taskrepo.TaskRepository
}

func NewListTasksService(
	taskRepository taskrepo.TaskRepository,
) *ListTasksService {
	return &ListTasksService{
		TaskRepository: taskRepository,
	}
}

type ListTasksInput struct {
	Filters        taskrepo.TaskFilters
	Pagination     core.PaginationInput
	SortInput      core.SortInput
	RelationsInput core.RelationsInput
}

func (i ListTasksInput) Validate() error { return nil }

func (s *ListTasksService) Execute(input ListTasksInput) (*core.PaginationOutput[task.TaskDto], error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	tasks, err := s.TaskRepository.PaginateTasksBy(taskrepo.PaginateTasksParams{
		Filters:        input.Filters,
		Pagination:     input.Pagination,
		SortInput:      input.SortInput,
		RelationsInput: input.RelationsInput,
	})
	if err != nil {
		return nil, err
	}

	var tasksDto []task.TaskDto = make([]task.TaskDto, len(tasks.Data))
	for i, tsk := range tasks.Data {
		tasksDto[i] = *task.TaskToDto(&tsk)
	}

	return &core.PaginationOutput[task.TaskDto]{
		Data:    tasksDto,
		Page:    tasks.Page,
		HasMore: tasks.HasMore,
		Total:   tasks.Total,
	}, nil
}
