package taskrepo

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/task"
)

type TaskFilters struct {
	ProjectIdentity      *core.Identity
	TaskStatusIdentity   *core.Identity
	TaskCategoryIdentity *core.Identity
	ParentTaskIdentity   *core.Identity
	Name                 *core.ComparableFilter[string]
	CompletedAt          *core.ComparableFilter[int64]
	CreatedAt            *core.ComparableFilter[int64]
	UpdatedAt            *core.ComparableFilter[int64]
	DueDate              *core.ComparableFilter[int64]
	Type                 *core.ComparableFilter[task.TaskType]
	Priority             *core.ComparableFilter[task.TaskPriorityLevels]
}

type GetTaskByIdentityParams struct {
	TaskIdentity    core.Identity
	ProjectIdentity *core.Identity
}

type GetTasksByParentTaskIdentityParams struct {
	ParentTaskIdentity core.Identity
	ProjectIdentity    *core.Identity
}

type PaginateTasksParams struct {
	Filters    TaskFilters
	Pagination core.PaginationInput
	SortInput  core.SortInput
}

type StoreTaskParams struct {
	Task *task.Task
}

type UpdateTaskParams struct {
	Task *task.Task
}

type DeleteTaskParams struct {
	TaskIdentity core.Identity
}

type AddSubTaskParams struct {
	Task    *task.Task
	SubTask *task.SubTask
}

type UpdateSubTaskParams struct {
	Task    *task.Task
	SubTask *task.SubTask
}

type RemoveSubTaskParams struct {
	Task    *task.Task
	SubTask *task.SubTask
}

type TaskRepository interface {
	SetTransaction(tx core.Transaction) error

	GetTaskByIdentity(params GetTaskByIdentityParams) (*task.Task, error)
	GetTasksByParentTaskIdentity(params GetTasksByParentTaskIdentityParams) ([]*task.Task, error)
	PaginateTasksBy(params PaginateTasksParams) (*core.PaginationOutput[task.Task], error)

	AddSubTask(params AddSubTaskParams) error
	UpdateSubTask(params UpdateSubTaskParams) error
	RemoveSubTask(params RemoveSubTaskParams) error

	StoreTask(params StoreTaskParams) (*task.Task, error)
	UpdateTask(params UpdateTaskParams) error
	DeleteTask(params DeleteTaskParams) error
}
