package task

import (
	"github.com/gabrielmrtt/taski/internal/user"
)

type SubTaskDto struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	CompletedAt string `json:"completedAt"`
}

func SubTaskToDto(subTask *SubTask) *SubTaskDto {
	var completedAt *string = nil
	if subTask.CompletedAt != nil {
		completedAtString := subTask.CompletedAt.ToRFC3339()
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
		dueDateString := task.DueDate.ToRFC3339()
		dueDate = &dueDateString
	}

	var completedAt *string = nil
	if task.CompletedAt != nil {
		completedAtString := task.CompletedAt.ToRFC3339()
		completedAt = &completedAtString
	}

	var updatedAt *string = nil
	if task.Timestamps.UpdatedAt != nil {
		updatedAtString := task.Timestamps.UpdatedAt.ToRFC3339()
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
		CreatedAt:        task.Timestamps.CreatedAt.ToRFC3339(),
		UpdatedAt:        updatedAt,
	}
}

type TaskCommentFileDto struct {
	FileId string `json:"fileId"`
}

type TaskCommentDto struct {
	Id        string               `json:"id"`
	Content   string               `json:"content"`
	Files     []TaskCommentFileDto `json:"files"`
	Author    *user.UserDto        `json:"author"`
	CreatedAt string               `json:"createdAt"`
	UpdatedAt *string              `json:"updatedAt"`
}

func TaskCommentToDto(taskComment *TaskComment) *TaskCommentDto {
	var updatedAt *string = nil
	if taskComment.Timestamps.UpdatedAt != nil {
		updatedAtString := taskComment.Timestamps.UpdatedAt.ToRFC3339()
		updatedAt = &updatedAtString
	}

	var taskCommentFilesDto []TaskCommentFileDto = make([]TaskCommentFileDto, len(taskComment.Files))
	for i, file := range taskComment.Files {
		taskCommentFilesDto[i] = TaskCommentFileDto{
			FileId: file.FileIdentity.Public,
		}
	}

	var author *user.UserDto = nil
	if taskComment.Author != nil {
		author = user.UserToDto(taskComment.Author)
	}

	return &TaskCommentDto{
		Id:        taskComment.Identity.Public,
		Content:   taskComment.Content,
		Files:     taskCommentFilesDto,
		Author:    author,
		CreatedAt: taskComment.Timestamps.CreatedAt.ToRFC3339(),
		UpdatedAt: updatedAt,
	}
}
