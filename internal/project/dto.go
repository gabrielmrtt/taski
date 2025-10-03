package project

import "github.com/gabrielmrtt/taski/pkg/datetimeutils"

type ProjectDto struct {
	Id            string  `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Status        string  `json:"status"`
	Color         string  `json:"color"`
	PriorityLevel int8    `json:"priorityLevel"`
	StartAt       *string `json:"startAt"`
	EndAt         *string `json:"endAt"`
	WorkspaceId   string  `json:"workspaceId"`
	UserCreatorId string  `json:"userCreatorId"`
	UserEditorId  *string `json:"userEditorId"`
	CreatedAt     string  `json:"createdAt"`
	UpdatedAt     *string `json:"updatedAt"`
}

func ProjectToDto(project *Project) *ProjectDto {
	var startAt *string = nil
	if project.StartAt != nil {
		startAtString := datetimeutils.EpochToRFC3339(*project.StartAt)
		startAt = &startAtString
	}

	var endAt *string = nil
	if project.EndAt != nil {
		endAtString := datetimeutils.EpochToRFC3339(*project.EndAt)
		endAt = &endAtString
	}

	var userCreatorId *string = nil
	if project.UserCreatorIdentity != nil {
		userCreatorId = &project.UserCreatorIdentity.Public
	}

	var userEditorId *string = nil
	if project.UserEditorIdentity != nil {
		userEditorId = &project.UserEditorIdentity.Public
	}

	createdAt := datetimeutils.EpochToRFC3339(*project.Timestamps.CreatedAt)
	var updatedAt *string = nil
	if project.Timestamps.UpdatedAt != nil {
		updatedAtString := datetimeutils.EpochToRFC3339(*project.Timestamps.UpdatedAt)
		updatedAt = &updatedAtString
	}

	return &ProjectDto{
		Id:            project.Identity.Public,
		Name:          project.Name,
		Description:   project.Description,
		Status:        string(project.Status),
		Color:         project.Color,
		PriorityLevel: int8(project.PriorityLevel),
		StartAt:       startAt,
		EndAt:         endAt,
		WorkspaceId:   project.WorkspaceIdentity.Public,
		UserCreatorId: *userCreatorId,
		UserEditorId:  userEditorId,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}
