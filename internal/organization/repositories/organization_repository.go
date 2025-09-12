package organization_repositories

import (
	"github.com/gabrielmrtt/taski/internal/core"
	organization_core "github.com/gabrielmrtt/taski/internal/organization"
)

type OrganizationFilters struct {
	Name               *core.ComparableFilter[string]
	Status             *core.ComparableFilter[organization_core.OrganizationStatuses]
	CreatedAt          *core.ComparableFilter[int64]
	UpdatedAt          *core.ComparableFilter[int64]
	DeletedAt          *core.ComparableFilter[int64]
	LoggedUserIdentity *core.Identity
}

type GetOrganizationByIdentityParams struct {
	OrganizationIdentity core.Identity
	RelationsInput       core.RelationsInput
}

type PaginateOrganizationsParams struct {
	Filters        OrganizationFilters
	ShowDeleted    bool
	SortInput      *core.SortInput
	Pagination     *core.PaginationInput
	RelationsInput core.RelationsInput
}

type StoreOrganizationParams struct {
	Organization *organization_core.Organization
}

type UpdateOrganizationParams struct {
	Organization *organization_core.Organization
}

type DeleteOrganizationParams struct {
	OrganizationIdentity core.Identity
}

type OrganizationRepository interface {
	SetTransaction(tx core.Transaction) error

	GetOrganizationByIdentity(params GetOrganizationByIdentityParams) (*organization_core.Organization, error)
	PaginateOrganizationsBy(params PaginateOrganizationsParams) (*core.PaginationOutput[organization_core.Organization], error)

	StoreOrganization(params StoreOrganizationParams) (*organization_core.Organization, error)
	UpdateOrganization(params UpdateOrganizationParams) error
	DeleteOrganization(params DeleteOrganizationParams) error
}
