package project

import (
	"github.com/gabrielmrtt/taski/internal/user"
	"github.com/gabrielmrtt/taski/pkg/datetimeutils"
)

type ProjectDto struct {
	Id            string        `json:"id"`
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	Status        string        `json:"status"`
	Color         string        `json:"color"`
	PriorityLevel int8          `json:"priorityLevel"`
	StartAt       *string       `json:"startAt"`
	EndAt         *string       `json:"endAt"`
	WorkspaceId   string        `json:"workspaceId"`
	UserCreatorId string        `json:"userCreatorId"`
	UserEditorId  *string       `json:"userEditorId"`
	Creator       *user.UserDto `json:"creator,omitempty"`
	Editor        *user.UserDto `json:"editor,omitempty"`
	CreatedAt     string        `json:"createdAt"`
	UpdatedAt     *string       `json:"updatedAt"`
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

	var creator *user.UserDto = nil
	if project.Creator != nil {
		creator = user.UserToDto(project.Creator)
	}

	var editor *user.UserDto = nil
	if project.Editor != nil {
		editor = user.UserToDto(project.Editor)
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
		Creator:       creator,
		Editor:        editor,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

type ProjectTaskStatusDto struct {
	Id                       string `json:"id"`
	Name                     string `json:"name"`
	Color                    string `json:"color"`
	Order                    *int8  `json:"order"`
	ShouldSetTaskToCompleted bool   `json:"shouldSetTaskToCompleted"`
	IsDefault                bool   `json:"isDefault"`
}

func ProjectTaskStatusToDto(projectTaskStatus *ProjectTaskStatus) *ProjectTaskStatusDto {
	return &ProjectTaskStatusDto{
		Id:                       projectTaskStatus.Identity.Public,
		Name:                     projectTaskStatus.Name,
		Color:                    projectTaskStatus.Color,
		Order:                    projectTaskStatus.Order,
		ShouldSetTaskToCompleted: projectTaskStatus.ShouldSetTaskToCompleted,
		IsDefault:                projectTaskStatus.IsDefault,
	}
}

type ProjectTaskCategoryDto struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

func ProjectTaskCategoryToDto(projectTaskCategory *ProjectTaskCategory) *ProjectTaskCategoryDto {
	return &ProjectTaskCategoryDto{
		Id:    projectTaskCategory.Identity.Public,
		Name:  projectTaskCategory.Name,
		Color: projectTaskCategory.Color,
	}
}

type ProjectDocumentVersionDto struct {
	Id                              string                   `json:"id"`
	ProjectDocumentVersionManagerId string                   `json:"projectDocumentVersionManagerId"`
	Version                         string                   `json:"version"`
	Latest                          bool                     `json:"latest"`
	Title                           string                   `json:"title"`
	Content                         string                   `json:"content"`
	Files                           []ProjectDocumentFileDto `json:"files"`
	UserCreatorId                   string                   `json:"userCreatorId"`
	UserEditorId                    *string                  `json:"userEditorId"`
	Creator                         *user.UserDto            `json:"creator,omitempty"`
	Editor                          *user.UserDto            `json:"editor,omitempty"`
	CreatedAt                       string                   `json:"createdAt"`
	UpdatedAt                       *string                  `json:"updatedAt"`
}

type ProjectDocumentFileDto struct {
	FileId string `json:"fileId"`
}

func ProjectDocumentVersionToDto(projectDocumentVersion *ProjectDocumentVersion) *ProjectDocumentVersionDto {
	var files []ProjectDocumentFileDto = make([]ProjectDocumentFileDto, len(projectDocumentVersion.Document.Files))

	for i, file := range projectDocumentVersion.Document.Files {
		files[i] = ProjectDocumentFileDto{
			FileId: file.FileIdentity.Public,
		}
	}

	var userEditorId *string = nil
	if projectDocumentVersion.UserEditorIdentity != nil {
		userEditorId = &projectDocumentVersion.UserEditorIdentity.Public
	}

	var updatedAt *string = nil
	if projectDocumentVersion.Timestamps.UpdatedAt != nil {
		updatedAtString := datetimeutils.EpochToRFC3339(*projectDocumentVersion.Timestamps.UpdatedAt)
		updatedAt = &updatedAtString
	}

	var creator *user.UserDto = nil
	if projectDocumentVersion.Creator != nil {
		creator = user.UserToDto(projectDocumentVersion.Creator)
	}

	var editor *user.UserDto = nil
	if projectDocumentVersion.Editor != nil {
		editor = user.UserToDto(projectDocumentVersion.Editor)
	}

	return &ProjectDocumentVersionDto{
		Id:                              projectDocumentVersion.Identity.Public,
		ProjectDocumentVersionManagerId: projectDocumentVersion.ProjectDocumentVersionManagerIdentity.Public,
		Version:                         projectDocumentVersion.Version,
		Latest:                          projectDocumentVersion.Latest,
		Title:                           projectDocumentVersion.Document.Title,
		Content:                         projectDocumentVersion.Document.Content,
		Files:                           files,
		UserCreatorId:                   projectDocumentVersion.UserCreatorIdentity.Public,
		UserEditorId:                    userEditorId,
		Creator:                         creator,
		Editor:                          editor,
		CreatedAt:                       datetimeutils.EpochToRFC3339(*projectDocumentVersion.Timestamps.CreatedAt),
		UpdatedAt:                       updatedAt,
	}
}
