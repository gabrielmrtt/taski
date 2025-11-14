package task

var TaskIdentityPrefix = "tsk"

var SubTaskIdentityPrefix = "sts"

var TaskCommentIdentityPrefix = "tsc"

var TaskActionIdentityPrefix = "tsa"

type TaskPriorityLevels int8

const (
	TaskPriorityLevelNone     TaskPriorityLevels = 0
	TaskPriorityLevelLow      TaskPriorityLevels = 1
	TaskPriorityLevelMedium   TaskPriorityLevels = 2
	TaskPriorityLevelHigh     TaskPriorityLevels = 3
	TaskPriorityLevelCritical TaskPriorityLevels = 4
	TaskPriorityLevelUrgent   TaskPriorityLevels = 5
)

type TaskType string

const (
	TaskTypeNormal TaskType = "normal"
	TaskTypeGroup  TaskType = "group"
)

type TaskActionType string

const (
	TaskActionTypeChangeStatus    TaskActionType = "task_status_changed"
	TaskActionTypeCreate          TaskActionType = "task_created"
	TaskActionTypeUpdate          TaskActionType = "task_updated"
	TaskActionTypeDelete          TaskActionType = "task_deleted"
	TaskActionTypeAddSubTask      TaskActionType = "sub_task_created"
	TaskActionTypeRemoveSubTask   TaskActionType = "sub_task_removed"
	TaskActionTypeAddComment      TaskActionType = "comment_created"
	TaskActionTypeUpdateComment   TaskActionType = "comment_updated"
	TaskActionTypeDeleteComment   TaskActionType = "comment_deleted"
	TaskActionTypeComplete        TaskActionType = "task_completed"
	TaskActionTypeSubTaskComplete TaskActionType = "sub_task_completed"
)
