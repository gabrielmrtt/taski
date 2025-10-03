package project

const ProjectIdentityPrefix = "prj"

type ProjectStatuses string

const (
	ProjectStatusPaused    ProjectStatuses = "paused"
	ProjectStatusOngoing   ProjectStatuses = "ongoing"
	ProjectStatusCompleted ProjectStatuses = "completed"
	ProjectStatusCancelled ProjectStatuses = "cancelled"
	ProjectStatusArchived  ProjectStatuses = "archived"
)

type ProjectPriorityLevels int8

const (
	ProjectPriorityLevelNone     ProjectPriorityLevels = 0
	ProjectPriorityLevelLow      ProjectPriorityLevels = 1
	ProjectPriorityLevelMedium   ProjectPriorityLevels = 2
	ProjectPriorityLevelHigh     ProjectPriorityLevels = 3
	ProjectPriorityLevelCritical ProjectPriorityLevels = 4
	ProjectPriorityLevelUrgent   ProjectPriorityLevels = 5
)

type ProjectUserStatuses string

const (
	ProjectUserStatusActive   ProjectUserStatuses = "active"
	ProjectUserStatusInactive ProjectUserStatuses = "inactive"
	ProjectUserStatusInvited  ProjectUserStatuses = "invited"
	ProjectUserStatusRefused  ProjectUserStatuses = "refused"
)
