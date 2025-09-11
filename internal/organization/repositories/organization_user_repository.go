package organization_repositories

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organization_core "github.com/gabrielmrtt/taski/internal/organization"
)

type OrganizationUserFilters struct {
	OrganizationIdentity core.Identity
	Name                 *core.ComparableFilter[string]
	Email                *core.ComparableFilter[string]
	DisplayName          *core.ComparableFilter[string]
	RolePublicId         *core.ComparableFilter[string]
	Status               *core.ComparableFilter[organization_core.OrganizationUserStatuses]
}

type GetOrganizationUserByIdentityParams struct {
	OrganizationIdentity core.Identity
	UserIdentity         core.Identity
}

type PaginateOrganizationUsersParams struct {
	Filters    OrganizationUserFilters
	SortInput  *core.SortInput
	Pagination *core.PaginationInput
}

type CreateOrganizationUserParams struct {
	OrganizationUser *organization_core.OrganizationUser
}

type UpdateOrganizationUserParams struct {
	OrganizationUser *organization_core.OrganizationUser
}

type DeleteOrganizationUserParams struct {
	OrganizationIdentity core.Identity
	UserIdentity         core.Identity
}

type OrganizationUserRepository interface {
	SetTransaction(tx core.Transaction) error

	GetOrganizationUserByIdentity(params GetOrganizationUserByIdentityParams) (*organization_core.OrganizationUser, error)
	PaginateOrganizationUsersBy(params PaginateOrganizationUsersParams) (*core.PaginationOutput[organization_core.OrganizationUser], error)

	CreateOrganizationUser(params CreateOrganizationUserParams) (*organization_core.OrganizationUser, error)
	UpdateOrganizationUser(params UpdateOrganizationUserParams) error
	DeleteOrganizationUser(params DeleteOrganizationUserParams) error
}
