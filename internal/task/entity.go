package task

import "github.com/gabrielmrtt/taski/internal/core"

type Task struct {
	Identity                core.Identity
	ProjectIdentity         core.Identity
	ProjectStatusIdentity   core.Identity
	ProjectCategoryIdentity core.Identity
	Name                    string
	Description             string
	EstimatedMinutes        *int16
	PriorityLevel           TaskPriorityLevels
	DueDate                 *int64
	CompletedAt             *int64
	UserCreatorIdentity     *core.Identity
	UserEditorIdentity      *core.Identity
	Timestamps              core.Timestamps
	DeletedAt               *int64
}
