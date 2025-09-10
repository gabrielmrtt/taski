package organization_core

import "github.com/gabrielmrtt/taski/internal/core"

type OrganizationFilters struct {
	Name               *core.ComparableFilter[string]
	Status             *core.ComparableFilter[OrganizationStatuses]
	CreatedAt          *core.ComparableFilter[int64]
	UpdatedAt          *core.ComparableFilter[int64]
	DeletedAt          *core.ComparableFilter[int64]
	LoggedUserIdentity *core.Identity
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

type StoreOrganizationParams struct {
	Organization *Organization
}

type UpdateOrganizationParams struct {
	Organization *Organization
}

type DeleteOrganizationParams struct {
	OrganizationIdentity core.Identity
}

type OrganizationRepository interface {
	SetTransaction(tx core.Transaction) error

	GetOrganizationByIdentity(params GetOrganizationByIdentityParams) (*Organization, error)
	PaginateOrganizationsBy(params PaginateOrganizationsParams) (*core.PaginationOutput[Organization], error)

	StoreOrganization(params StoreOrganizationParams) (*Organization, error)
	UpdateOrganization(params UpdateOrganizationParams) error
	DeleteOrganization(params DeleteOrganizationParams) error
}

type OrganizationUserFilters struct {
	OrganizationIdentity core.Identity
	Name                 *core.ComparableFilter[string]
	Email                *core.ComparableFilter[string]
	DisplayName          *core.ComparableFilter[string]
	RolePublicId         *core.ComparableFilter[string]
	Status               *core.ComparableFilter[OrganizationUserStatuses]
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
	OrganizationUser *OrganizationUser
}

type UpdateOrganizationUserParams struct {
	OrganizationUser *OrganizationUser
}

type DeleteOrganizationUserParams struct {
	OrganizationIdentity core.Identity
	UserIdentity         core.Identity
}

type OrganizationUserRepository interface {
	SetTransaction(tx core.Transaction) error

	GetOrganizationUserByIdentity(params GetOrganizationUserByIdentityParams) (*OrganizationUser, error)
	PaginateOrganizationUsersBy(params PaginateOrganizationUsersParams) (*core.PaginationOutput[OrganizationUser], error)

	CreateOrganizationUser(params CreateOrganizationUserParams) (*OrganizationUser, error)
	UpdateOrganizationUser(params UpdateOrganizationUserParams) error
	DeleteOrganizationUser(params DeleteOrganizationUserParams) error
}
