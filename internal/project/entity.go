package project_core

import (
	"github.com/gabrielmrtt/taski/internal/core"
	user_core "github.com/gabrielmrtt/taski/internal/user"
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
	User            user_core.User
	Status          ProjectUserStatuses
}

type NewProjectUserInput struct {
	ProjectIdentity core.Identity
	User            user_core.User
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
