package task

var TaskIdentityPrefix = "tsk"

type TaskPriorityLevels int8

const (
	TaskPriorityLevelNone     TaskPriorityLevels = 0
	TaskPriorityLevelLow      TaskPriorityLevels = 1
	TaskPriorityLevelMedium   TaskPriorityLevels = 2
	TaskPriorityLevelHigh     TaskPriorityLevels = 3
	TaskPriorityLevelCritical TaskPriorityLevels = 4
	TaskPriorityLevelUrgent   TaskPriorityLevels = 5
)
