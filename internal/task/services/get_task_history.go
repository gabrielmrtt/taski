package taskservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/task"
	taskrepo "github.com/gabrielmrtt/taski/internal/task/repository"
)

type GetTaskHistoryService struct {
	TaskActionRepository taskrepo.TaskActionRepository
	TaskRepository       taskrepo.TaskRepository
}

func NewGetTaskHistoryService(
	taskActionRepository taskrepo.TaskActionRepository,
	taskRepository taskrepo.TaskRepository,
) *GetTaskHistoryService {
	return &GetTaskHistoryService{
		TaskActionRepository: taskActionRepository,
		TaskRepository:       taskRepository,
	}
}

type GetTaskHistoryInput struct {
	OrganizationIdentity *core.Identity
	TaskIdentity         core.Identity
	Filters              taskrepo.TaskActionFilters
	SortInput            core.SortInput
	PaginationInput      core.PaginationInput
	RelationsInput       core.RelationsInput
}

func (i GetTaskHistoryInput) Validate() error {
	return nil
}

func (s *GetTaskHistoryService) Execute(input GetTaskHistoryInput) (*core.PaginationOutput[task.TaskActionDto], error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	taskActions, err := s.TaskActionRepository.PaginateTaskActionsBy(taskrepo.PaginateTaskActionsParams{
		Filters: taskrepo.TaskActionFilters{
			TaskIdentity: &input.TaskIdentity,
		},
		Pagination:     input.PaginationInput,
		SortInput:      input.SortInput,
		RelationsInput: input.RelationsInput,
	})
	if err != nil {
		return nil, err
	}

	var taskActionsDto []task.TaskActionDto = make([]task.TaskActionDto, len(taskActions.Data))
	for i, taskAction := range taskActions.Data {
		taskActionsDto[i] = *task.TaskActionToDto(&taskAction)
	}

	return &core.PaginationOutput[task.TaskActionDto]{
		Data:    taskActionsDto,
		Page:    taskActions.Page,
		HasMore: taskActions.HasMore,
		Total:   taskActions.Total,
	}, nil
}
