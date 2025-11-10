package taskhttprequests

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/task"
	taskservice "github.com/gabrielmrtt/taski/internal/task/services"
)

type CreateSubTaskRequest struct {
	Name string `json:"name"`
}

type CreateTaskRequest struct {
	ProjectId        string                  `json:"projectId"`
	StatusId         string                  `json:"statusId"`
	CategoryId       *string                 `json:"categoryId"`
	ParentTaskId     *string                 `json:"parentTaskId"`
	Name             string                  `json:"name"`
	Description      string                  `json:"description"`
	EstimatedMinutes *int16                  `json:"estimatedMinutes"`
	PriorityLevel    int8                    `json:"priorityLevel"`
	DueDate          *string                 `json:"dueDate"`
	SubTasks         []*CreateSubTaskRequest `json:"subTasks"`
	Users            []*string               `json:"users"`
	ChildrenTasks    []*string               `json:"childrenTasks"`
}

func (r *CreateTaskRequest) ToInput() taskservice.CreateTaskInput {
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

	var users []*core.Identity = make([]*core.Identity, len(r.Users))
	for i, user := range r.Users {
		identity := core.NewIdentity(*user)
		users[i] = &identity
	}

	var childrenTasks []*core.Identity = make([]*core.Identity, len(r.ChildrenTasks))
	for i, childTask := range r.ChildrenTasks {
		identity := core.NewIdentity(*childTask)
		childrenTasks[i] = &identity
	}

	var subTasks []*taskservice.CreateSubTaskInput = make([]*taskservice.CreateSubTaskInput, len(r.SubTasks))
	for i, subTask := range r.SubTasks {
		subTasks[i] = &taskservice.CreateSubTaskInput{
			Name: subTask.Name,
		}
	}

	var dueDate *core.DateTime = nil
	if r.DueDate != nil {
		d, err := core.NewDateTimeFromRFC3339(*r.DueDate)
		if err != nil {
			return taskservice.CreateTaskInput{}
		}
		dueDate = &d
	}

	return taskservice.CreateTaskInput{
		ProjectIdentity:    core.NewIdentity(r.ProjectId),
		StatusIdentity:     core.NewIdentity(r.StatusId),
		CategoryIdentity:   categoryIdentity,
		ParentTaskIdentity: parentTaskIdentity,
		Name:               r.Name,
		Description:        r.Description,
		EstimatedMinutes:   r.EstimatedMinutes,
		PriorityLevel:      task.TaskPriorityLevels(r.PriorityLevel),
		DueDate:            dueDate,
		SubTasks:           subTasks,
		Users:              users,
		ChildrenTasks:      childrenTasks,
	}
}
