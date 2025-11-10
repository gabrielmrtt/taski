package taskservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/task"
	taskrepo "github.com/gabrielmrtt/taski/internal/task/repository"
)

type GetTaskService struct {
	TaskRepository taskrepo.TaskRepository
}

func NewGetTaskService(
	taskRepository taskrepo.TaskRepository,
) *GetTaskService {
	return &GetTaskService{
		TaskRepository: taskRepository,
	}
}

type GetTaskInput struct {
	TaskIdentity         core.Identity
	OrganizationIdentity *core.Identity
	RelationsInput       core.RelationsInput
}

func (i GetTaskInput) Validate() error { return nil }

func (s *GetTaskService) Execute(input GetTaskInput) (*task.TaskDto, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	tsk, err := s.TaskRepository.GetTaskByIdentity(taskrepo.GetTaskByIdentityParams{
		TaskIdentity:         input.TaskIdentity,
		OrganizationIdentity: input.OrganizationIdentity,
		RelationsInput:       input.RelationsInput,
	})
	if err != nil {
		return nil, err
	}

	if tsk == nil {
		return nil, core.NewNotFoundError("task not found")
	}

	return task.TaskToDto(tsk), nil
}
