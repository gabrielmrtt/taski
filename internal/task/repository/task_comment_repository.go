package taskrepo

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/task"
)

type TaskCommentFilters struct {
	TaskIdentity   *core.Identity
	AuthorIdentity *core.Identity
	CreatedAt      *core.ComparableFilter[int64]
	UpdatedAt      *core.ComparableFilter[int64]
}

type GetTaskCommentByIdentityParams struct {
	TaskCommentIdentity core.Identity
	TaskIdentity        *core.Identity
	RelationsInput      core.RelationsInput
}

type PaginateTaskCommentsParams struct {
	Filters        TaskCommentFilters
	Pagination     core.PaginationInput
	SortInput      core.SortInput
	RelationsInput core.RelationsInput
}

type StoreTaskCommentParams struct {
	TaskComment *task.TaskComment
}

type UpdateTaskCommentParams struct {
	TaskComment *task.TaskComment
}

type DeleteTaskCommentParams struct {
	TaskCommentIdentity core.Identity
}

type TaskCommentRepository interface {
	SetTransaction(tx core.Transaction) error

	GetTaskCommentByIdentity(params GetTaskCommentByIdentityParams) (*task.TaskComment, error)
	PaginateTaskCommentsBy(params PaginateTaskCommentsParams) (*core.PaginationOutput[task.TaskComment], error)

	StoreTaskComment(params StoreTaskCommentParams) (*task.TaskComment, error)
	UpdateTaskComment(params UpdateTaskCommentParams) error
	DeleteTaskComment(params DeleteTaskCommentParams) error
}
