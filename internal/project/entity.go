package project

import (
	"slices"

	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/user"
	"github.com/gabrielmrtt/taski/pkg/datetimeutils"
)

type Project struct {
	Identity            core.Identity
	WorkspaceIdentity   core.Identity
	Name                string
	Description         string
	Status              ProjectStatuses
	Color               string
	PriorityLevel       ProjectPriorityLevels
	UserCreatorIdentity *core.Identity
	UserEditorIdentity  *core.Identity
	StartAt             *int64
	EndAt               *int64
	Timestamps          core.Timestamps
	DeletedAt           *int64
}

type NewProjectInput struct {
	Name                string
	Description         string
	Color               string
	WorkspaceIdentity   core.Identity
	PriorityLevel       ProjectPriorityLevels
	StartAt             *int64
	EndAt               *int64
	UserCreatorIdentity *core.Identity
}

func NewProject(input NewProjectInput) (*Project, error) {
	now := datetimeutils.EpochNow()

	if _, err := core.NewName(input.Name); err != nil {
		return nil, err
	}

	if _, err := core.NewDescription(input.Description); err != nil {
		return nil, err
	}

	if _, err := core.NewColor(input.Color); err != nil {
		return nil, err
	}

	if input.StartAt != nil {
		if *input.StartAt < now {
			return nil, core.NewInternalError("start at cannot be in the past")
		}
	} else {
		input.StartAt = &now
	}

	if input.EndAt != nil {
		if *input.EndAt < now {
			return nil, core.NewInternalError("end at cannot be in the past")
		}
	} else {
		input.EndAt = &now
	}

	return &Project{
		Identity:            core.NewIdentity(ProjectIdentityPrefix),
		WorkspaceIdentity:   input.WorkspaceIdentity,
		Name:                input.Name,
		Description:         input.Description,
		Status:              ProjectStatusOngoing,
		Color:               input.Color,
		PriorityLevel:       input.PriorityLevel,
		UserCreatorIdentity: input.UserCreatorIdentity,
		UserEditorIdentity:  nil,
		StartAt:             input.StartAt,
		EndAt:               input.EndAt,
		Timestamps: core.Timestamps{
			CreatedAt: &now,
			UpdatedAt: nil,
		},
		DeletedAt: nil,
	}, nil
}

func (p *Project) ChangeName(name string, userEditorIdentity *core.Identity) error {
	if _, err := core.NewName(name); err != nil {
		return err
	}

	p.Name = name
	p.UserEditorIdentity = userEditorIdentity
	now := datetimeutils.EpochNow()
	p.Timestamps.UpdatedAt = &now
	return nil
}

func (p *Project) ChangeDescription(description string, userEditorIdentity *core.Identity) error {
	if _, err := core.NewDescription(description); err != nil {
		return err
	}

	p.Description = description
	p.UserEditorIdentity = userEditorIdentity
	now := datetimeutils.EpochNow()
	p.Timestamps.UpdatedAt = &now
	return nil
}

func (p *Project) ChangeColor(color string, userEditorIdentity *core.Identity) error {
	if _, err := core.NewColor(color); err != nil {
		return err
	}

	p.Color = color
	p.UserEditorIdentity = userEditorIdentity
	now := datetimeutils.EpochNow()
	p.Timestamps.UpdatedAt = &now
	return nil
}

func (p *Project) ChangeStatus(status ProjectStatuses, userEditorIdentity *core.Identity) error {
	if p.Status == status {
		return nil
	}

	p.Status = status
	p.UserEditorIdentity = userEditorIdentity
	now := datetimeutils.EpochNow()
	p.Timestamps.UpdatedAt = &now
	return nil
}

func (p *Project) ChangeStartAt(startAt int64, userEditorIdentity *core.Identity) error {
	now := datetimeutils.EpochNow()

	if startAt < now {
		return core.NewInternalError("start at cannot be in the past")
	}

	if startAt > *p.EndAt {
		return core.NewInternalError("start at cannot be after end at")
	}

	p.StartAt = &startAt
	p.UserEditorIdentity = userEditorIdentity
	p.Timestamps.UpdatedAt = &now
	return nil
}

func (p *Project) ChangeEndAt(endAt int64, userEditorIdentity *core.Identity) error {
	now := datetimeutils.EpochNow()

	if endAt < now {
		return core.NewInternalError("end at cannot be in the past")
	}

	if endAt < *p.StartAt {
		return core.NewInternalError("end at cannot be before start at")
	}

	p.EndAt = &endAt
	p.UserEditorIdentity = userEditorIdentity
	p.Timestamps.UpdatedAt = &now
	return nil
}

func (p *Project) ChangePriorityLevel(priorityLevel ProjectPriorityLevels, userEditorIdentity *core.Identity) error {
	if p.PriorityLevel == priorityLevel {
		return nil
	}

	p.PriorityLevel = priorityLevel
	p.UserEditorIdentity = userEditorIdentity
	now := datetimeutils.EpochNow()
	p.Timestamps.UpdatedAt = &now
	return nil
}

func (p *Project) Delete() {
	now := datetimeutils.EpochNow()
	p.DeletedAt = &now
}

func (p *Project) IsPaused() bool {
	return p.Status == ProjectStatusPaused
}

func (p *Project) IsOngoing() bool {
	return p.Status == ProjectStatusOngoing
}

func (p *Project) IsCompleted() bool {
	return p.Status == ProjectStatusCompleted
}

func (p *Project) IsCancelled() bool {
	return p.Status == ProjectStatusCancelled
}

func (p *Project) IsArchived() bool {
	return p.Status == ProjectStatusArchived
}

func (p *Project) IsDeleted() bool {
	return p.DeletedAt != nil
}

func (p *Project) HasEnded() bool {
	now := datetimeutils.EpochNow()

	return p.EndAt != nil && *p.EndAt < now
}

func (p *Project) HasStarted() bool {
	now := datetimeutils.EpochNow()

	return p.StartAt != nil && *p.StartAt <= now
}

type ProjectUser struct {
	ProjectIdentity core.Identity
	User            user.User
	Status          ProjectUserStatuses
}

type NewProjectUserInput struct {
	ProjectIdentity core.Identity
	User            user.User
	Status          ProjectUserStatuses
}

func NewProjectUser(input NewProjectUserInput) (*ProjectUser, error) {
	return &ProjectUser{
		ProjectIdentity: input.ProjectIdentity,
		User:            input.User,
		Status:          input.Status,
	}, nil
}

func (p *ProjectUser) Activate() {
	p.Status = ProjectUserStatusActive
}

func (p *ProjectUser) Deactivate() {
	p.Status = ProjectUserStatusInactive
}

func (p *ProjectUser) Invite() {
	p.Status = ProjectUserStatusInvited
}

func (p *ProjectUser) IsActive() bool {
	return p.Status == ProjectUserStatusActive
}

func (p *ProjectUser) IsInactive() bool {
	return p.Status == ProjectUserStatusInactive
}

func (p *ProjectUser) IsInvited() bool {
	return p.Status == ProjectUserStatusInvited
}

func (p *ProjectUser) AcceptInvitation() {
	p.Status = ProjectUserStatusActive
}

func (p *ProjectUser) RefuseInvitation() {
	p.Status = ProjectUserStatusRefused
}

type ProjectTaskStatus struct {
	Identity                 core.Identity
	ProjectIdentity          core.Identity
	Name                     string
	Color                    string
	Order                    *int8
	ShouldSetTaskToCompleted bool
	IsDefault                bool
	DeletedAt                *int64
}

type NewProjectTaskStatusInput struct {
	Name                     string
	Color                    string
	Order                    *int8
	ShouldSetTaskToCompleted bool
	IsDefault                bool
	ProjectIdentity          core.Identity
}

func NewProjectTaskStatus(input NewProjectTaskStatusInput) (*ProjectTaskStatus, error) {
	if _, err := core.NewName(input.Name); err != nil {
		return nil, err
	}

	if _, err := core.NewColor(input.Color); err != nil {
		return nil, err
	}

	if input.IsDefault && input.ShouldSetTaskToCompleted {
		return nil, core.NewConflictError("project status should not be set to completed and default at the same time")
	}

	return &ProjectTaskStatus{
		Identity:                 core.NewIdentity(ProjectTaskStatusIdentityPrefix),
		ProjectIdentity:          input.ProjectIdentity,
		Name:                     input.Name,
		Color:                    input.Color,
		Order:                    input.Order,
		ShouldSetTaskToCompleted: input.ShouldSetTaskToCompleted,
		IsDefault:                input.IsDefault,
		DeletedAt:                nil,
	}, nil
}

func (s *ProjectTaskStatus) ChangeName(name string) error {
	if _, err := core.NewName(name); err != nil {
		return err
	}

	s.Name = name

	return nil
}

func (s *ProjectTaskStatus) ChangeColor(color string) error {
	if _, err := core.NewColor(color); err != nil {
		return err
	}

	s.Color = color

	return nil
}

func (s *ProjectTaskStatus) ChangeOrder(order int8) error {
	s.Order = &order
	return nil
}

func (s *ProjectTaskStatus) SetTaskToCompleted(v bool) error {
	s.ShouldSetTaskToCompleted = v
	return nil
}

func (s *ProjectTaskStatus) SetIsDefault(v bool) error {
	if s.ShouldSetTaskToCompleted && v {
		return core.NewConflictError("project status should not be set to completed and default at the same time")
	}

	s.IsDefault = v
	return nil
}

func (s *ProjectTaskStatus) Delete() {
	now := datetimeutils.EpochNow()
	s.DeletedAt = &now
}

func (s *ProjectTaskStatus) IsDeleted() bool {
	return s.DeletedAt != nil
}

type ProjectTaskCategory struct {
	Identity        core.Identity
	ProjectIdentity core.Identity
	Name            string
	Color           string
	DeletedAt       *int64
}

type NewProjectTaskCategoryInput struct {
	Name            string
	Color           string
	ProjectIdentity core.Identity
}

func NewProjectTaskCategory(input NewProjectTaskCategoryInput) (*ProjectTaskCategory, error) {
	if _, err := core.NewName(input.Name); err != nil {
		return nil, err
	}

	if _, err := core.NewColor(input.Color); err != nil {
		return nil, err
	}

	return &ProjectTaskCategory{
		Identity:        core.NewIdentity(ProjectTaskCategoryIdentityPrefix),
		ProjectIdentity: input.ProjectIdentity,
		Name:            input.Name,
		Color:           input.Color,
		DeletedAt:       nil,
	}, nil
}

func (c *ProjectTaskCategory) ChangeName(name string) error {
	if _, err := core.NewName(name); err != nil {
		return err
	}

	c.Name = name
	return nil
}

func (c *ProjectTaskCategory) ChangeColor(color string) error {
	if _, err := core.NewColor(color); err != nil {
		return err
	}

	c.Color = color
	return nil
}

func (c *ProjectTaskCategory) Delete() {
	now := datetimeutils.EpochNow()
	c.DeletedAt = &now
}

func (c *ProjectTaskCategory) IsDeleted() bool {
	return c.DeletedAt != nil
}

type ProjectDocumentVersionManager struct {
	Identity        core.Identity
	ProjectIdentity core.Identity
	LatestVersion   *ProjectDocumentVersion
}

type ProjectDocumentVersion struct {
	Identity                              core.Identity
	ProjectDocumentVersionManagerIdentity core.Identity
	Version                               string
	Document                              ProjectDocument
	UserCreatorIdentity                   *core.Identity
	UserEditorIdentity                    *core.Identity
	Latest                                bool
	Timestamps                            core.Timestamps
}

type ProjectDocument struct {
	Identity core.Identity
	Title    string
	Content  string
	Files    []ProjectDocumentFile
}

type ProjectDocumentFile struct {
	Identity     core.Identity
	FileIdentity core.Identity
}

type NewProjectDocumentInput struct {
	ProjectIdentity                       core.Identity
	ProjectDocumentVersionManagerIdentity core.Identity
	Title                                 string
	Content                               string
	Version                               string
	Files                                 []ProjectDocumentFile
	UserCreatorIdentity                   *core.Identity
}

func NewProjectDocument(input NewProjectDocumentInput) (*ProjectDocumentVersion, error) {
	if _, err := NewProjectDocumentTitle(input.Title); err != nil {
		return nil, err
	}

	if _, err := NewProjectDocumentContent(input.Content); err != nil {
		return nil, err
	}

	now := datetimeutils.EpochNow()

	projectDocument := &ProjectDocument{
		Identity: core.NewIdentityWithoutPublic(),
		Title:    input.Title,
		Content:  input.Content,
		Files:    input.Files,
	}

	projectDocumentVersion := &ProjectDocumentVersion{
		Identity:                              core.NewIdentity(ProjectDocumentVersionIdentityPrefix),
		ProjectDocumentVersionManagerIdentity: input.ProjectDocumentVersionManagerIdentity,
		Document:                              *projectDocument,
		UserCreatorIdentity:                   input.UserCreatorIdentity,
		Latest:                                true,
		Version:                               input.Version,
		Timestamps: core.Timestamps{
			CreatedAt: &now,
			UpdatedAt: nil,
		},
	}

	return projectDocumentVersion, nil
}

func (v *ProjectDocumentVersion) ChangeTitle(title string, userEditorIdentity *core.Identity) error {
	if _, err := NewProjectDocumentTitle(title); err != nil {
		return err
	}

	v.Document.Title = title
	v.UserEditorIdentity = userEditorIdentity
	now := datetimeutils.EpochNow()
	v.Timestamps.UpdatedAt = &now
	return nil
}

func (v *ProjectDocumentVersion) ChangeContent(content string, userEditorIdentity *core.Identity) error {
	if _, err := NewProjectDocumentContent(content); err != nil {
		return err
	}

	v.Document.Content = content
	v.UserEditorIdentity = userEditorIdentity
	now := datetimeutils.EpochNow()
	v.Timestamps.UpdatedAt = &now
	return nil
}

func (v *ProjectDocumentVersion) ClearAllFiles() {
	v.Document.Files = []ProjectDocumentFile{}
	now := datetimeutils.EpochNow()
	v.Timestamps.UpdatedAt = &now
}

func (v *ProjectDocumentVersion) AddFile(file ProjectDocumentFile) {
	v.Document.Files = append(v.Document.Files, file)
	now := datetimeutils.EpochNow()
	v.Timestamps.UpdatedAt = &now
}

func (v *ProjectDocumentVersion) RemoveFile(file ProjectDocumentFile) {
	v.Document.Files = slices.DeleteFunc(v.Document.Files, func(f ProjectDocumentFile) bool {
		return f.Identity == file.Identity
	})
	now := datetimeutils.EpochNow()
	v.Timestamps.UpdatedAt = &now
}

func (v *ProjectDocumentVersion) IsLatest() bool {
	return v.Latest
}

func (v *ProjectDocumentVersion) NewVersion(version string) *ProjectDocumentVersion {
	now := datetimeutils.EpochNow()

	v.Latest = false
	return &ProjectDocumentVersion{
		Identity:                              core.NewIdentity(ProjectDocumentVersionIdentityPrefix),
		ProjectDocumentVersionManagerIdentity: v.ProjectDocumentVersionManagerIdentity,
		Version:                               version,
		Document:                              v.Document,
		UserCreatorIdentity:                   v.UserCreatorIdentity,
		UserEditorIdentity:                    v.UserEditorIdentity,
		Latest:                                true,
		Timestamps: core.Timestamps{
			CreatedAt: &now,
			UpdatedAt: nil,
		},
	}
}
