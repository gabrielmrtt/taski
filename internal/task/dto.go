package task

import "github.com/gabrielmrtt/taski/pkg/datetimeutils"

type SubTaskDto struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	CompletedAt string `json:"completedAt"`
}

func SubTaskToDto(subTask *SubTask) *SubTaskDto {
	var completedAt *string = nil
	if subTask.CompletedAt != nil {
		completedAtString := datetimeutils.EpochToRFC3339(*subTask.CompletedAt)
		completedAt = &completedAtString
	}

	return &SubTaskDto{
		Id:          subTask.Identity.Public,
		Name:        subTask.Name,
		CompletedAt: *completedAt,
	}
}

type TaskUserDto struct {
	UserId string `json:"userId"`
}

func TaskUserToDto(taskUser *TaskUser) *TaskUserDto {
	return &TaskUserDto{
		UserId: taskUser.User.Identity.Public,
	}
}

type TaskDto struct {
	Id               string         `json:"id"`
	Name             string         `json:"name"`
	Description      string         `json:"description"`
	EstimatedMinutes int16          `json:"estimatedMinutes"`
	PriorityLevel    int8           `json:"priorityLevel"`
	DueDate          string         `json:"dueDate"`
	CompletedAt      string         `json:"completedAt"`
	SubTasks         []*SubTaskDto  `json:"subTasks"`
	ChildrenTasks    []*TaskDto     `json:"childrenTasks"`
	Users            []*TaskUserDto `json:"users"`
	UserCreatorId    string         `json:"userCreatorId"`
	UserEditorId     *string        `json:"userEditorId"`
	CreatedAt        string         `json:"createdAt"`
	UpdatedAt        *string        `json:"updatedAt"`
}

func TaskToDto(task *Task) *TaskDto {
	var dueDate *string = nil
	if task.DueDate != nil {
		dueDateString := datetimeutils.EpochToRFC3339(*task.DueDate)
		dueDate = &dueDateString
	}

	var completedAt *string = nil
	if task.CompletedAt != nil {
		completedAtString := datetimeutils.EpochToRFC3339(*task.CompletedAt)
		completedAt = &completedAtString
	}

	var updatedAt *string = nil
	if task.Timestamps.UpdatedAt != nil {
		updatedAtString := datetimeutils.EpochToRFC3339(*task.Timestamps.UpdatedAt)
		updatedAt = &updatedAtString
	}

	var usersDto []*TaskUserDto = make([]*TaskUserDto, len(task.Users))
	for i, user := range task.Users {
		usersDto[i] = TaskUserToDto(user)
	}

	var subTasksDto []*SubTaskDto = make([]*SubTaskDto, len(task.SubTasks))
	for i, subTask := range task.SubTasks {
		subTasksDto[i] = SubTaskToDto(subTask)
	}

	var childrenTasksDto []*TaskDto = make([]*TaskDto, len(task.ChildrenTasks))
	for i, childTask := range task.ChildrenTasks {
		childrenTasksDto[i] = TaskToDto(childTask)
	}

	var userCreatorId *string = nil
	if task.UserCreatorIdentity != nil {
		userCreatorId = &task.UserCreatorIdentity.Public
	}

	var userEditorId *string = nil
	if task.UserEditorIdentity != nil {
		userEditorId = &task.UserEditorIdentity.Public
	}

	return &TaskDto{
		Id:               task.Identity.Public,
		Name:             task.Name,
		Description:      task.Description,
		EstimatedMinutes: *task.EstimatedMinutes,
		PriorityLevel:    int8(task.PriorityLevel),
		DueDate:          *dueDate,
		CompletedAt:      *completedAt,
		SubTasks:         subTasksDto,
		ChildrenTasks:    childrenTasksDto,
		Users:            usersDto,
		UserCreatorId:    *userCreatorId,
		UserEditorId:     userEditorId,
		CreatedAt:        datetimeutils.EpochToRFC3339(*task.Timestamps.CreatedAt),
		UpdatedAt:        updatedAt,
	}
}
