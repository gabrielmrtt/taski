package task

import (
	"slices"

	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/project"
	"github.com/gabrielmrtt/taski/internal/user"
)

type TaskUser struct {
	User *user.User
}

type SubTask struct {
	Identity    core.Identity
	Name        string
	CompletedAt *core.DateTime
}

type NewSubTaskInput struct {
	Name string
}

func NewSubTask(input NewSubTaskInput) (*SubTask, error) {
	if _, err := core.NewName(input.Name); err != nil {
		return nil, err
	}

	return &SubTask{
		Identity:    core.NewIdentity(SubTaskIdentityPrefix),
		Name:        input.Name,
		CompletedAt: nil,
	}, nil
}

func (s *SubTask) ChangeName(name string) error {
	if _, err := core.NewName(name); err != nil {
		return err
	}

	s.Name = name
	return nil
}

func (s *SubTask) Complete() {
	now := core.NewDateTime()
	s.CompletedAt = &now
}

func (s *SubTask) Uncomplete() {
	s.CompletedAt = nil
}

type Task struct {
	Identity                core.Identity
	ProjectIdentity         core.Identity
	Status                  *project.ProjectTaskStatus
	Category                *project.ProjectTaskCategory
	ParentTaskIdentity      *core.Identity
	Type                    TaskType
	Name                    string
	Description             string
	EstimatedMinutes        *int16
	PriorityLevel           TaskPriorityLevels
	DueDate                 *core.DateTime
	CompletedAt             *core.DateTime
	SubTasks                []*SubTask
	ChildrenTasks           []*Task
	Users                   []*TaskUser
	UserCompletedByIdentity *core.Identity
	UserCreatorIdentity     *core.Identity
	UserEditorIdentity      *core.Identity
	Timestamps              core.Timestamps
	DeletedAt               *core.DateTime
}

type NewTaskInput struct {
	ProjectIdentity     core.Identity
	Status              *project.ProjectTaskStatus
	Category            *project.ProjectTaskCategory
	ParentTaskIdentity  *core.Identity
	Name                string
	Description         string
	EstimatedMinutes    *int16
	PriorityLevel       TaskPriorityLevels
	DueDate             *core.DateTime
	SubTasks            []*SubTask
	Users               []*TaskUser
	ChildrenTasks       []*Task
	UserCreatorIdentity *core.Identity
}

func NewTask(input NewTaskInput) (*Task, error) {
	now := core.NewDateTime()

	if _, err := core.NewName(input.Name); err != nil {
		return nil, err
	}

	if _, err := core.NewDescription(input.Description); err != nil {
		return nil, err
	}

	if input.EstimatedMinutes != nil {
		if *input.EstimatedMinutes < 0 {
			return nil, core.NewInternalError("estimated minutes cannot be negative")
		}
	}

	var taskType TaskType = TaskTypeNormal

	if input.ChildrenTasks != nil {
		if len(input.ChildrenTasks) == 0 {
			return nil, core.NewInternalError("children tasks cannot be empty")
		}

		taskType = TaskTypeGroup
	}

	return &Task{
		Identity:            core.NewIdentity(TaskIdentityPrefix),
		ProjectIdentity:     input.ProjectIdentity,
		ParentTaskIdentity:  input.ParentTaskIdentity,
		Status:              input.Status,
		Category:            input.Category,
		Type:                taskType,
		Name:                input.Name,
		SubTasks:            input.SubTasks,
		ChildrenTasks:       input.ChildrenTasks,
		Description:         input.Description,
		EstimatedMinutes:    input.EstimatedMinutes,
		PriorityLevel:       input.PriorityLevel,
		DueDate:             input.DueDate,
		CompletedAt:         nil,
		UserCreatorIdentity: input.UserCreatorIdentity,
		UserEditorIdentity:  nil,
		Timestamps: core.Timestamps{
			CreatedAt: &now,
			UpdatedAt: nil,
		},
		DeletedAt: nil,
	}, nil
}

func (t *Task) ChangeName(name string, userEditorIdentity *core.Identity) error {
	if _, err := core.NewName(name); err != nil {
		return err
	}

	t.Name = name
	t.UserEditorIdentity = userEditorIdentity
	now := core.NewDateTime()
	t.Timestamps.UpdatedAt = &now
	return nil
}

func (t *Task) ChangeDescription(description string, userEditorIdentity *core.Identity) error {
	if _, err := core.NewDescription(description); err != nil {
		return err
	}

	t.Description = description
	t.UserEditorIdentity = userEditorIdentity
	now := core.NewDateTime()
	t.Timestamps.UpdatedAt = &now
	return nil
}

func (t *Task) ChangeEstimatedMinutes(estimatedMinutes int16, userEditorIdentity *core.Identity) error {
	if estimatedMinutes < 0 {
		return core.NewInternalError("estimated minutes cannot be negative")
	}

	t.EstimatedMinutes = &estimatedMinutes
	t.UserEditorIdentity = userEditorIdentity
	now := core.NewDateTime()
	t.Timestamps.UpdatedAt = &now
	return nil
}

func (t *Task) ChangePriorityLevel(priorityLevel TaskPriorityLevels, userEditorIdentity *core.Identity) error {
	t.PriorityLevel = priorityLevel
	t.UserEditorIdentity = userEditorIdentity
	now := core.NewDateTime()
	t.Timestamps.UpdatedAt = &now
	return nil
}

func (t *Task) ChangeDueDate(dueDate core.DateTime, userEditorIdentity *core.Identity) error {
	now := core.NewDateTime()

	if dueDate.IsBefore(now) {
		return core.NewInternalError("due date cannot be in the past")
	}

	t.DueDate = &dueDate
	t.UserEditorIdentity = userEditorIdentity
	t.Timestamps.UpdatedAt = &now
	return nil
}

func (t *Task) ChangeStatus(status *project.ProjectTaskStatus, userEditorIdentity *core.Identity) error {
	t.Status = status
	t.UserEditorIdentity = userEditorIdentity
	now := core.NewDateTime()
	t.Timestamps.UpdatedAt = &now
	return nil
}

func (t *Task) ChangeCategory(category *project.ProjectTaskCategory, userEditorIdentity *core.Identity) error {
	t.Category = category
	t.UserEditorIdentity = userEditorIdentity
	now := core.NewDateTime()
	t.Timestamps.UpdatedAt = &now
	return nil
}

func (t *Task) Complete() {
	now := core.NewDateTime()
	t.CompletedAt = &now
}

func (t *Task) Delete() {
	now := core.NewDateTime()
	t.DeletedAt = &now
}

func (t *Task) IsCompleted() bool {
	return t.CompletedAt != nil
}

func (t *Task) IsDeleted() bool {
	return t.DeletedAt != nil
}

func (t *Task) IsOverdue() bool {
	now := core.NewDateTime()
	return t.DueDate != nil && t.DueDate.IsBefore(now) && !t.IsCompleted()
}

func (t *Task) AddUser(user *TaskUser) {
	t.Users = append(t.Users, user)
	now := core.NewDateTime()
	t.Timestamps.UpdatedAt = &now
}

func (t *Task) RemoveUser(user *TaskUser) {
	t.Users = slices.DeleteFunc(t.Users, func(u *TaskUser) bool {
		return u.User.Identity == user.User.Identity
	})
	now := core.NewDateTime()
	t.Timestamps.UpdatedAt = &now
}

func (t *Task) ClearUsers() {
	t.Users = []*TaskUser{}
	now := core.NewDateTime()
	t.Timestamps.UpdatedAt = &now
}

func (t *Task) AddSubTask(subTask *SubTask) {
	t.SubTasks = append(t.SubTasks, subTask)
	now := core.NewDateTime()
	t.Timestamps.UpdatedAt = &now
}

func (t *Task) RemoveSubTask(subTask *SubTask) {
	t.SubTasks = slices.DeleteFunc(t.SubTasks, func(s *SubTask) bool {
		return s.Identity == subTask.Identity
	})
	now := core.NewDateTime()
	t.Timestamps.UpdatedAt = &now
}

func (t *Task) AddChildTask(childTask *Task) {
	t.ChildrenTasks = append(t.ChildrenTasks, childTask)
	childTask.ParentTaskIdentity = &t.Identity
	now := core.NewDateTime()
	t.Timestamps.UpdatedAt = &now
	childTask.Timestamps.UpdatedAt = &now

	t.Type = TaskTypeGroup
}

func (t *Task) GetSubTaskByIdentity(subTaskIdentity core.Identity) *SubTask {
	for _, subTask := range t.SubTasks {
		if subTask.Identity.Equals(subTaskIdentity) {
			return subTask
		}
	}

	return nil
}

func (t *Task) RemoveChildTask(childTask *Task) {
	t.ChildrenTasks = slices.DeleteFunc(t.ChildrenTasks, func(c *Task) bool {
		return t.Identity.Equals(childTask.Identity)
	})

	now := core.NewDateTime()
	t.Timestamps.UpdatedAt = &now
	childTask.Timestamps.UpdatedAt = &now

	if len(t.ChildrenTasks) == 0 {
		t.Type = TaskTypeNormal
	}

	childTask.ParentTaskIdentity = nil
}

type TaskCommentFile struct {
	Identity     core.Identity
	FileIdentity core.Identity
}

type TaskComment struct {
	Identity     core.Identity
	TaskIdentity core.Identity
	Content      string
	Files        []TaskCommentFile
	Author       *user.User
	Timestamps   core.Timestamps
}

type NewTaskCommentInput struct {
	TaskIdentity core.Identity
	Content      string
	Files        []TaskCommentFile
	Author       *user.User
}

func NewTaskComment(input NewTaskCommentInput) (*TaskComment, error) {
	if _, err := NewTaskCommentContent(input.Content); err != nil {
		return nil, err
	}

	now := core.NewDateTime()

	return &TaskComment{
		Identity:     core.NewIdentity(TaskCommentIdentityPrefix),
		TaskIdentity: input.TaskIdentity,
		Content:      input.Content,
		Files:        input.Files,
		Author:       input.Author,
		Timestamps: core.Timestamps{
			CreatedAt: &now,
			UpdatedAt: nil,
		},
	}, nil
}

func (t *TaskComment) ChangeContent(content string) error {
	if _, err := NewTaskCommentContent(content); err != nil {
		return err
	}

	t.Content = content
	now := core.NewDateTime()
	t.Timestamps.UpdatedAt = &now
	return nil
}

func (t *TaskComment) AddFile(file TaskCommentFile) {
	t.Files = append(t.Files, file)
	now := core.NewDateTime()
	t.Timestamps.UpdatedAt = &now
}

func (t *TaskComment) RemoveFile(file TaskCommentFile) {
	t.Files = slices.DeleteFunc(t.Files, func(f TaskCommentFile) bool {
		return f.Identity == file.Identity
	})
	now := core.NewDateTime()
	t.Timestamps.UpdatedAt = &now
}

func (t *TaskComment) ClearAllFiles() {
	t.Files = []TaskCommentFile{}
	now := core.NewDateTime()
	t.Timestamps.UpdatedAt = &now
}
