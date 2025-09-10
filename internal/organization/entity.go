package organization_core

import (
	"github.com/gabrielmrtt/taski/internal/core"
	role_core "github.com/gabrielmrtt/taski/internal/role"
	user_core "github.com/gabrielmrtt/taski/internal/user"
	"github.com/gabrielmrtt/taski/pkg/datetimeutils"
)

type OrganizationStatuses string

const (
	OrganizationStatusActive   OrganizationStatuses = "active"
	OrganizationStatusInactive OrganizationStatuses = "inactive"
)

const OrganizationIdentityPrefix = "org"

type Organization struct {
	Identity            core.Identity
	Name                string
	Status              OrganizationStatuses
	UserCreatorIdentity *core.Identity
	UserEditorIdentity  *core.Identity
	core.Timestamps
	DeletedAt *int64
}

type NewOrganizationInput struct {
	Name                string
	UserCreatorIdentity *core.Identity
}

func NewOrganization(input NewOrganizationInput) (*Organization, error) {
	now := datetimeutils.EpochNow()

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
	now := datetimeutils.EpochNow()
	o.Timestamps.UpdatedAt = &now
	return nil
}

func (o *Organization) Delete() {
	now := datetimeutils.EpochNow()
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

type OrganizationUserStatuses string

const (
	OrganizationUserStatusActive   OrganizationUserStatuses = "active"
	OrganizationUserStatusInactive OrganizationUserStatuses = "inactive"
	OrganizationUserStatusInvited  OrganizationUserStatuses = "invited"
	OrganizationUserStatusRefused  OrganizationUserStatuses = "refused"
)

type OrganizationUser struct {
	OrganizationIdentity core.Identity
	User                 *user_core.User
	Role                 *role_core.Role
	Status               OrganizationUserStatuses
}

type NewOrganizationUserInput struct {
	OrganizationIdentity core.Identity
	User                 *user_core.User
	Role                 *role_core.Role
	Status               OrganizationUserStatuses
}

func NewOrganizationUser(input NewOrganizationUserInput) (*OrganizationUser, error) {
	return &OrganizationUser{
		OrganizationIdentity: input.OrganizationIdentity,
		User:                 input.User,
		Role:                 input.Role,
		Status:               input.Status,
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
}

func (o *OrganizationUser) RefuseInvitation() {
	o.Status = OrganizationUserStatusRefused
}
