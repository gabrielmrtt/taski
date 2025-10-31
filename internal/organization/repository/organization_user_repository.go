package organizationrepo

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/organization"
)

type OrganizationUserFilters struct {
	OrganizationIdentity core.Identity
	Name                 *core.ComparableFilter[string]
	Email                *core.ComparableFilter[string]
	DisplayName          *core.ComparableFilter[string]
	RolePublicId         *core.ComparableFilter[string]
	UserPublicId         *core.ComparableFilter[string]
	Status               *core.ComparableFilter[organization.OrganizationUserStatuses]
}

type GetOrganizationUserByIdentityParams struct {
	OrganizationIdentity core.Identity
	UserIdentity         core.Identity
	RelationsInput       core.RelationsInput
}

type PaginateOrganizationUsersParams struct {
	Filters        OrganizationUserFilters
	SortInput      core.SortInput
	Pagination     core.PaginationInput
	RelationsInput core.RelationsInput
}

type StoreOrganizationUserParams struct {
	OrganizationUser *organization.OrganizationUser
}

type UpdateOrganizationUserParams struct {
	OrganizationUser *organization.OrganizationUser
}

type DeleteOrganizationUserParams struct {
	OrganizationIdentity core.Identity
	UserIdentity         core.Identity
}

type GetLastAccessedOrganizationUserByUserIdentityParams struct {
	UserIdentity core.Identity
}

type OrganizationUserRepository interface {
	SetTransaction(tx core.Transaction) error

	GetLastAccessedOrganizationUserByUserIdentity(params GetLastAccessedOrganizationUserByUserIdentityParams) (*organization.OrganizationUser, error)
	GetOrganizationUserByIdentity(params GetOrganizationUserByIdentityParams) (*organization.OrganizationUser, error)
	PaginateOrganizationUsersBy(params PaginateOrganizationUsersParams) (*core.PaginationOutput[organization.OrganizationUser], error)

	StoreOrganizationUser(params StoreOrganizationUserParams) (*organization.OrganizationUser, error)
	UpdateOrganizationUser(params UpdateOrganizationUserParams) error
	DeleteOrganizationUser(params DeleteOrganizationUserParams) error
}
