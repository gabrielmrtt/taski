package organizationrepo

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/organization"
)

type OrganizationFilters struct {
	Name                      *core.ComparableFilter[string]
	Status                    *core.ComparableFilter[organization.OrganizationStatuses]
	CreatedAt                 *core.ComparableFilter[int64]
	UpdatedAt                 *core.ComparableFilter[int64]
	DeletedAt                 *core.ComparableFilter[int64]
	AuthenticatedUserIdentity *core.Identity
}

type GetOrganizationByIdentityParams struct {
	OrganizationIdentity core.Identity
}

type PaginateOrganizationsParams struct {
	Filters     OrganizationFilters
	ShowDeleted bool
	SortInput   *core.SortInput
	Pagination  *core.PaginationInput
}

type PaginateInvitedOrganizationsParams struct {
	AuthenticatedUserIdentity core.Identity
	SortInput                 *core.SortInput
	Pagination                *core.PaginationInput
}

type StoreOrganizationParams struct {
	Organization *organization.Organization
}

type UpdateOrganizationParams struct {
	Organization *organization.Organization
}

type DeleteOrganizationParams struct {
	OrganizationIdentity core.Identity
}

type OrganizationRepository interface {
	SetTransaction(tx core.Transaction) error

	GetOrganizationByIdentity(params GetOrganizationByIdentityParams) (*organization.Organization, error)
	PaginateOrganizationsBy(params PaginateOrganizationsParams) (*core.PaginationOutput[organization.Organization], error)
	PaginateInvitedOrganizationsBy(params PaginateInvitedOrganizationsParams) (*core.PaginationOutput[organization.Organization], error)

	StoreOrganization(params StoreOrganizationParams) (*organization.Organization, error)
	UpdateOrganization(params UpdateOrganizationParams) error
	DeleteOrganization(params DeleteOrganizationParams) error
}
