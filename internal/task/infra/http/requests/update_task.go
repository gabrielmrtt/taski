package taskhttprequests

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/task"
	taskservice "github.com/gabrielmrtt/taski/internal/task/services"
)

type UpdateTaskRequest struct {
	StatusId         *string   `json:"statusId"`
	CategoryId       *string   `json:"categoryId"`
	ParentTaskId     *string   `json:"parentTaskId"`
	Name             *string   `json:"name"`
	Description      *string   `json:"description"`
	EstimatedMinutes *int16    `json:"estimatedMinutes"`
	PriorityLevel    *int8     `json:"priorityLevel"`
	DueDate          *string   `json:"dueDate"`
	Users            *[]string `json:"users"`
	ChildrenTasks    *[]string `json:"childrenTasks"`
}

func (r *UpdateTaskRequest) ToInput() taskservice.UpdateTaskInput {
	var statusIdentity *core.Identity = nil
	if r.StatusId != nil {
		identity := core.NewIdentity(*r.StatusId)
		statusIdentity = &identity
	}

	var categoryIdentity *core.Identity = nil
	if r.CategoryId != nil {
		identity := core.NewIdentity(*r.CategoryId)
		categoryIdentity = &identity
	}

	var parentTaskIdentity *core.Identity = nil
	if r.ParentTaskId != nil {
		identity := core.NewIdentity(*r.ParentTaskId)
		parentTaskIdentity = &identity
	}

	var users []*core.Identity = make([]*core.Identity, len(*r.Users))
	for i, user := range *r.Users {
		identity := core.NewIdentity(user)
		users[i] = &identity
	}

	var childrenTasks []*core.Identity = make([]*core.Identity, len(*r.ChildrenTasks))
	for i, childTask := range *r.ChildrenTasks {
		identity := core.NewIdentity(childTask)
		childrenTasks[i] = &identity
	}

	var priorityLevel *task.TaskPriorityLevels = nil
	if r.PriorityLevel != nil {
		p := task.TaskPriorityLevels(*r.PriorityLevel)
		priorityLevel = &p
	}

	var dueDate *core.DateTime = nil
	if r.DueDate != nil {
		d, err := core.NewDateTimeFromRFC3339(*r.DueDate)
		if err != nil {
			return taskservice.UpdateTaskInput{}
		}
		dueDate = &d
	}

	return taskservice.UpdateTaskInput{
		StatusIdentity:     statusIdentity,
		CategoryIdentity:   categoryIdentity,
		ParentTaskIdentity: parentTaskIdentity,
		Name:               r.Name,
		Description:        r.Description,
		EstimatedMinutes:   r.EstimatedMinutes,
		PriorityLevel:      priorityLevel,
		DueDate:            dueDate,
		Users:              users,
		ChildrenTasks:      childrenTasks,
	}
}
