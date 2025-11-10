package organization

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/role"
	"github.com/gabrielmrtt/taski/internal/user"
)

type Organization struct {
	Identity            core.Identity
	Name                string
	Status              OrganizationStatuses
	UserCreatorIdentity *core.Identity
	UserEditorIdentity  *core.Identity
	Creator             *user.User
	Editor              *user.User
	DeletedAt           *core.DateTime
	core.Timestamps
}

type NewOrganizationInput struct {
	Name                string
	UserCreatorIdentity *core.Identity
}

func NewOrganization(input NewOrganizationInput) (*Organization, error) {
	now := core.NewDateTime()

	return &Organization{
		Identity:            core.NewIdentity(OrganizationIdentityPrefix),
		Name:                input.Name,
		Status:              OrganizationStatusActive,
		UserCreatorIdentity: input.UserCreatorIdentity,
		UserEditorIdentity:  nil,
		Timestamps: core.Timestamps{
			CreatedAt: &now,
			UpdatedAt: nil,
		},
		DeletedAt: nil,
	}, nil
}

func (o *Organization) ChangeName(name string, userEditorIdentity *core.Identity) error {
	nameValueObject, err := core.NewName(name)
	if err != nil {
		return err
	}

	o.Name = nameValueObject.Value
	o.UserEditorIdentity = userEditorIdentity
	now := core.NewDateTime()
	o.Timestamps.UpdatedAt = &now
	return nil
}

func (o *Organization) Delete() {
	now := core.NewDateTime()
	o.DeletedAt = &now
}

func (o *Organization) IsDeleted() bool {
	return o.DeletedAt != nil
}

func (o *Organization) IsActive() bool {
	return o.Status == OrganizationStatusActive
}

func (o *Organization) IsInactive() bool {
	return o.Status == OrganizationStatusInactive
}

type OrganizationUser struct {
	OrganizationIdentity core.Identity
	User                 user.User
	Role                 role.Role
	Status               OrganizationUserStatuses
	LastAccessAt         *core.DateTime
}

type NewOrganizationUserInput struct {
	OrganizationIdentity core.Identity
	User                 user.User
	Role                 role.Role
	Status               OrganizationUserStatuses
}

func NewOrganizationUser(input NewOrganizationUserInput) (*OrganizationUser, error) {
	now := core.NewDateTime()

	return &OrganizationUser{
		OrganizationIdentity: input.OrganizationIdentity,
		User:                 input.User,
		Role:                 input.Role,
		Status:               input.Status,
		LastAccessAt:         &now,
	}, nil
}

func (o *OrganizationUser) Activate() {
	o.Status = OrganizationUserStatusActive
}

func (o *OrganizationUser) Deactivate() {
	o.Status = OrganizationUserStatusInactive
}

func (o *OrganizationUser) IsActive() bool {
	return o.Status == OrganizationUserStatusActive
}

func (o *OrganizationUser) IsInactive() bool {
	return o.Status == OrganizationUserStatusInactive
}

func (o *OrganizationUser) IsInvited() bool {
	return o.Status == OrganizationUserStatusInvited
}

func (o *OrganizationUser) Invite() {
	o.Status = OrganizationUserStatusInvited
}

func (o *OrganizationUser) AcceptInvitation() {
	o.Status = OrganizationUserStatusActive
	now := core.NewDateTime()
	o.LastAccessAt = &now
}

func (o *OrganizationUser) RefuseInvitation() {
	o.Status = OrganizationUserStatusRefused
}

func (o *OrganizationUser) ChangeRole(role role.Role) {
	o.Role = role
}

func (o *OrganizationUser) CanExecuteAction(permissionSlug role.PermissionSlugs) bool {
	return o.Role.HasPermission(permissionSlug)
}

func (o *OrganizationUser) Access() {
	now := core.NewDateTime()
	o.LastAccessAt = &now
}
