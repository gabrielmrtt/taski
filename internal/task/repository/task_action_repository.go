package taskrepo

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/task"
)

type TaskActionFilters struct {
	TaskIdentity *core.Identity
	CreatedAt    *core.ComparableFilter[int64]
	Type         *core.ComparableFilter[task.TaskActionType]
}

type PaginateTaskActionsParams struct {
	Filters        TaskActionFilters
	Pagination     core.PaginationInput
	SortInput      core.SortInput
	RelationsInput core.RelationsInput
}

type StoreTaskActionParams struct {
	TaskAction *task.TaskAction
}

type DeleteTaskActionParams struct {
	TaskActionIdentity core.Identity
}

type TaskActionRepository interface {
	SetTransaction(tx core.Transaction) error

	PaginateTaskActionsBy(params PaginateTaskActionsParams) (*core.PaginationOutput[task.TaskAction], error)
	StoreTaskAction(params StoreTaskActionParams) (*task.TaskAction, error)
	DeleteTaskAction(params DeleteTaskActionParams) error
}
