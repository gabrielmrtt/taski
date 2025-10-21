package taskservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
	"github.com/gabrielmrtt/taski/internal/task"
	taskrepo "github.com/gabrielmrtt/taski/internal/task/repository"
)

type ListTasksService struct {
	TaskRepository        taskrepo.TaskRepository
	ProjectRepository     projectrepo.ProjectRepository
	TransactionRepository core.TransactionRepository
}

func NewListTasksService(
	taskRepository taskrepo.TaskRepository,
	projectRepository projectrepo.ProjectRepository,
	transactionRepository core.TransactionRepository,
) *ListTasksService {
	return &ListTasksService{
		TaskRepository:        taskRepository,
		ProjectRepository:     projectRepository,
		TransactionRepository: transactionRepository,
	}
}

type ListTasksInput struct {
	ProjectIdentity core.Identity
	Filters         taskrepo.TaskFilters
	Pagination      core.PaginationInput
	SortInput       core.SortInput
}

func (i ListTasksInput) Validate() error { return nil }

func (s *ListTasksService) Execute(input ListTasksInput) (*core.PaginationOutput[task.TaskDto], error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return nil, err
	}

	s.TaskRepository.SetTransaction(tx)
	s.ProjectRepository.SetTransaction(tx)

	prj, err := s.ProjectRepository.GetProjectByIdentity(projectrepo.GetProjectByIdentityParams{
		ProjectIdentity: input.ProjectIdentity,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if prj == nil {
		tx.Rollback()
		return nil, core.NewNotFoundError("project not found")
	}

	tasks, err := s.TaskRepository.PaginateTasksBy(taskrepo.PaginateTasksParams{
		Filters:    input.Filters,
		Pagination: input.Pagination,
		SortInput:  input.SortInput,
	})
	if err != nil {
		tx.Rollback()
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
