package workspace_core

const WorkspaceIdentityPrefix = "wks"

type WorkspaceStatuses string

const (
	WorkspaceStatusActive   WorkspaceStatuses = "active"
	WorkspaceStatusInactive WorkspaceStatuses = "inactive"
	WorkspaceStatusArchived WorkspaceStatuses = "archived"
)

type WorkspaceUserStatuses string

const (
	WorkspaceUserStatusActive   WorkspaceUserStatuses = "active"
	WorkspaceUserStatusInactive WorkspaceUserStatuses = "inactive"
	WorkspaceUserStatusInvited  WorkspaceUserStatuses = "invited"
	WorkspaceUserStatusRefused  WorkspaceUserStatuses = "refused"
)
