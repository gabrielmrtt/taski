package workspace

import "github.com/gabrielmrtt/taski/pkg/datetimeutils"

type WorkspaceDto struct {
	Id             string  `json:"id"`
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	Color          string  `json:"color"`
	Status         string  `json:"status"`
	OrganizationId string  `json:"organizationId"`
	UserCreatorId  string  `json:"userCreatorId"`
	UserEditorId   *string `json:"userEditorId"`
	CreatedAt      string  `json:"createdAt"`
	UpdatedAt      *string `json:"updatedAt"`
}

func WorkspaceToDto(workspace *Workspace) *WorkspaceDto {
	createdAt := datetimeutils.EpochToRFC3339(*workspace.Timestamps.CreatedAt)

	var updatedAt *string = nil
	if workspace.Timestamps.UpdatedAt != nil {
		updatedAtString := datetimeutils.EpochToRFC3339(*workspace.Timestamps.UpdatedAt)
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

	return &WorkspaceDto{
		Id:             workspace.Identity.Public,
		Name:           workspace.Name,
		Description:    workspace.Description,
		Color:          workspace.Color,
		Status:         string(workspace.Status),
		OrganizationId: workspace.OrganizationIdentity.Public,
		UserCreatorId:  *userCreatorId,
		UserEditorId:   userEditorId,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}
