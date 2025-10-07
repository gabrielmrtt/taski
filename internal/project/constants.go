package project

const ProjectIdentityPrefix = "prj"

const ProjectTaskStatusIdentityPrefix = "prs"

const ProjectTaskCategoryIdentityPrefix = "ptc"

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

var DefaultProjectTaskStatuses = []ProjectTaskStatus{
	{
		Name:                     "Pending",
		Color:                    "#ADADAD",
		ShouldSetTaskToCompleted: false,
		Order:                    &[]int8{1}[0],
		IsDefault:                true,
	},
	{
		Name:                     "Ongoing",
		Color:                    "#0F53BF",
		ShouldSetTaskToCompleted: false,
		Order:                    &[]int8{2}[0],
		IsDefault:                false,
	},
	{
		Name:                     "Paused",
		Color:                    "#FFA500",
		ShouldSetTaskToCompleted: false,
		Order:                    nil,
		IsDefault:                false,
	},
	{
		Name:                     "Completed",
		Color:                    "#219653",
		ShouldSetTaskToCompleted: true,
		Order:                    nil,
		IsDefault:                false,
	},
}
