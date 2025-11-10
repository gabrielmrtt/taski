package workspace

import (
	"github.com/gabrielmrtt/taski/internal/organization"
	"github.com/gabrielmrtt/taski/internal/user"
)

type WorkspaceDto struct {
	Id             string                        `json:"id"`
	Name           string                        `json:"name"`
	Description    string                        `json:"description"`
	Color          string                        `json:"color"`
	Status         string                        `json:"status"`
	OrganizationId string                        `json:"organizationId"`
	UserCreatorId  string                        `json:"userCreatorId"`
	UserEditorId   *string                       `json:"userEditorId"`
	Creator        *user.UserDto                 `json:"creator,omitempty"`
	Editor         *user.UserDto                 `json:"editor,omitempty"`
	Organization   *organization.OrganizationDto `json:"organization,omitempty"`
	CreatedAt      string                        `json:"createdAt"`
	UpdatedAt      *string                       `json:"updatedAt"`
}

func WorkspaceToDto(workspace *Workspace) *WorkspaceDto {
	createdAt := workspace.Timestamps.CreatedAt.ToRFC3339()

	var updatedAt *string = nil
	if workspace.Timestamps.UpdatedAt != nil {
		updatedAtString := workspace.Timestamps.UpdatedAt.ToRFC3339()
		updatedAt = &updatedAtString
	}

	var userCreatorId *string = nil
	if workspace.UserCreatorIdentity != nil {
		userCreatorId = &workspace.UserCreatorIdentity.Public
	}

	var userEditorId *string = nil
	if workspace.UserEditorIdentity != nil {
		userEditorId = &workspace.UserEditorIdentity.Public
	}

	var creator *user.UserDto = nil
	if workspace.Creator != nil {
		creator = user.UserToDto(workspace.Creator)
	}

	var editor *user.UserDto = nil
	if workspace.Editor != nil {
		editor = user.UserToDto(workspace.Editor)
	}

	var org *organization.OrganizationDto = nil
	if workspace.Organization != nil {
		org = organization.OrganizationToDto(workspace.Organization)
	}

	return &WorkspaceDto{
		Id:             workspace.Identity.Public,
		Name:           workspace.Name,
		Description:    workspace.Description,
		Color:          workspace.Color,
		Status:         string(workspace.Status),
		OrganizationId: workspace.OrganizationIdentity.Public,
		UserCreatorId:  *userCreatorId,
		UserEditorId:   userEditorId,
		Creator:        creator,
		Editor:         editor,
		Organization:   org,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}
