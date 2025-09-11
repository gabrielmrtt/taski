package role_repositories

import (
	"github.com/gabrielmrtt/taski/internal/core"
	role_core "github.com/gabrielmrtt/taski/internal/role"
)

type RoleFilters struct {
	OrganizationIdentity core.Identity
	Name                 *core.ComparableFilter[string]
	Description          *core.ComparableFilter[string]
	CreatedAt            *core.ComparableFilter[int64]
	UpdatedAt            *core.ComparableFilter[int64]
	DeletedAt            *core.ComparableFilter[int64]
}

type GetRoleByIdentityParams struct {
	RoleIdentity core.Identity
}

type GetRoleByIdentityAndOrganizationIdentityParams struct {
	RoleIdentity         core.Identity
	OrganizationIdentity core.Identity
}

type GetDefaultRoleParams struct {
	Slug string
}

type PaginateRolesParams struct {
	Filters    RoleFilters
	SortInput  *core.SortInput
	Pagination *core.PaginationInput
}

type StoreRoleParams struct {
	Role *role_core.Role
}

type UpdateRoleParams struct {
	Role *role_core.Role
}

type DeleteRoleParams struct {
	RoleIdentity core.Identity
}

type ChangeRoleUsersToDefaultParams struct {
	RoleIdentity    core.Identity
	DefaultRoleSlug string
}

type RoleRepository interface {
	SetTransaction(tx core.Transaction) error

	GetRoleByIdentity(params GetRoleByIdentityParams) (*role_core.Role, error)
	GetRoleByIdentityAndOrganizationIdentity(params GetRoleByIdentityAndOrganizationIdentityParams) (*role_core.Role, error)
	GetSystemDefaultRole(params GetDefaultRoleParams) (*role_core.Role, error)
	PaginateRolesBy(params PaginateRolesParams) (*core.PaginationOutput[role_core.Role], error)
	ChangeRoleUsersToDefault(params ChangeRoleUsersToDefaultParams) error

	CheckIfOrganizatonHasUser(organizationIdentity core.Identity, userIdentity core.Identity) (bool, error)

	StoreRole(params StoreRoleParams) (*role_core.Role, error)
	UpdateRole(params UpdateRoleParams) error
	DeleteRole(params DeleteRoleParams) error
}
