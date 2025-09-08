package organization_core

import "github.com/gabrielmrtt/taski/internal/core"

type GetOrganizationByIdentityParams struct {
	Identity core.Identity
	Include  map[string]any
}

type OrganizationFilters struct {
	Name               *core.ComparableFilter[string]
	Status             *core.ComparableFilter[OrganizationStatuses]
	CreatedAt          *core.ComparableFilter[int64]
	UpdatedAt          *core.ComparableFilter[int64]
	DeletedAt          *core.ComparableFilter[int64]
	LoggedUserIdentity *core.Identity
}

type OrganizationUserFilters struct {
	Name           *core.ComparableFilter[string]
	Email          *core.ComparableFilter[string]
	DisplayName    *core.ComparableFilter[string]
	RoleInternalId *core.ComparableFilter[string]
	Status         *core.ComparableFilter[OrganizationUserStatuses]
}

type ListOrganizationsParams struct {
	Filters OrganizationFilters
	Include map[string]any
}

type ListOrganizationUsersParams struct {
	Filters OrganizationUserFilters
	Include map[string]any
}

type PaginateOrganizationsParams struct {
	Filters    OrganizationFilters
	Include    map[string]any
	SortInput  *core.SortInput
	Pagination *core.PaginationInput
}

type PaginateOrganizationUsersParams struct {
	Filters    OrganizationUserFilters
	Include    map[string]any
	SortInput  *core.SortInput
	Pagination *core.PaginationInput
}

type OrganizationRepository interface {
	SetTransaction(tx core.Transaction) error

	GetOrganizationByIdentity(params GetOrganizationByIdentityParams) (*Organization, error)
	ListOrganizationsBy(filters ListOrganizationsParams) (*[]Organization, error)
	PaginateOrganizationsBy(params PaginateOrganizationsParams) (*core.PaginationOutput[Organization], error)

	StoreOrganization(organization *Organization) (*Organization, error)
	UpdateOrganization(organization *Organization) error
	DeleteOrganization(identity core.Identity) error

	ListOrganizationUsersBy(filters ListOrganizationUsersParams) (*[]OrganizationUser, error)
	PaginateOrganizationUsersBy(params PaginateOrganizationUsersParams) (*core.PaginationOutput[OrganizationUser], error)

	CreateOrganizationUser(organizationUser *OrganizationUser) (*OrganizationUser, error)
	UpdateOrganizationUser(organizationUser *OrganizationUser) error
	DeleteOrganizationUser(identity core.Identity) error

	CheckIfOrganizationHasUser(organizationIdentity core.Identity, userIdentity core.Identity) (bool, error)
}
