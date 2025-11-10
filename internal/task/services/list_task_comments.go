package taskservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/task"
	taskrepo "github.com/gabrielmrtt/taski/internal/task/repository"
)

type ListTaskCommentsService struct {
	TaskCommentRepository taskrepo.TaskCommentRepository
	TaskRepository        taskrepo.TaskRepository
	TransactionRepository core.TransactionRepository
}

func NewListTaskCommentsService(
	taskCommentRepository taskrepo.TaskCommentRepository,
	taskRepository taskrepo.TaskRepository,
) *ListTaskCommentsService {
	return &ListTaskCommentsService{
		TaskCommentRepository: taskCommentRepository,
		TaskRepository:        taskRepository,
	}
}

type ListTaskCommentsInput struct {
	TaskIdentity   core.Identity
	Filters        taskrepo.TaskCommentFilters
	Pagination     core.PaginationInput
	SortInput      core.SortInput
	RelationsInput core.RelationsInput
}

func (i ListTaskCommentsInput) Validate() error { return nil }

func (s *ListTaskCommentsService) Execute(input ListTaskCommentsInput) (*core.PaginationOutput[task.TaskCommentDto], error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return nil, err
	}

	s.TaskCommentRepository.SetTransaction(tx)
	s.TaskRepository.SetTransaction(tx)

	tsk, err := s.TaskRepository.GetTaskByIdentity(taskrepo.GetTaskByIdentityParams{
		TaskIdentity: input.TaskIdentity,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if tsk == nil {
		tx.Rollback()
		return nil, core.NewNotFoundError("task not found")
	}

	comments, err := s.TaskCommentRepository.PaginateTaskCommentsBy(taskrepo.PaginateTaskCommentsParams{
		Filters:        input.Filters,
		Pagination:     input.Pagination,
		SortInput:      input.SortInput,
		RelationsInput: input.RelationsInput,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var commentsDto []task.TaskCommentDto = make([]task.TaskCommentDto, len(comments.Data))
	for i, comment := range comments.Data {
		commentsDto[i] = *task.TaskCommentToDto(&comment)
	}

	return &core.PaginationOutput[task.TaskCommentDto]{
		Data:    commentsDto,
		Page:    comments.Page,
		HasMore: comments.HasMore,
		Total:   comments.Total,
	}, nil
}
