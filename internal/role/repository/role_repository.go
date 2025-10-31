package rolerepo

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/role"
)

type RoleFilters struct {
	OrganizationIdentity *core.Identity
	Name                 *core.ComparableFilter[string]
	Description          *core.ComparableFilter[string]
	CreatedAt            *core.ComparableFilter[int64]
	UpdatedAt            *core.ComparableFilter[int64]
	DeletedAt            *core.ComparableFilter[int64]
}

type GetRoleByIdentityParams struct {
	RoleIdentity   core.Identity
	RelationsInput core.RelationsInput
}

type GetRoleByIdentityAndOrganizationIdentityParams struct {
	RoleIdentity         core.Identity
	OrganizationIdentity core.Identity
	RelationsInput       core.RelationsInput
}

type GetDefaultRoleParams struct {
	Slug role.DefaultRoleSlugs
}

type PaginateRolesParams struct {
	Filters        RoleFilters
	SortInput      core.SortInput
	Pagination     core.PaginationInput
	RelationsInput core.RelationsInput
	ShowDeleted    bool
}

type StoreRoleParams struct {
	Role *role.Role
}

type UpdateRoleParams struct {
	Role *role.Role
}

type DeleteRoleParams struct {
	RoleIdentity core.Identity
}

type ChangeRoleUsersToDefaultParams struct {
	RoleIdentity    core.Identity
	DefaultRoleSlug role.DefaultRoleSlugs
}

type RoleRepository interface {
	SetTransaction(tx core.Transaction) error

	GetRoleByIdentity(params GetRoleByIdentityParams) (*role.Role, error)
	GetRoleByIdentityAndOrganizationIdentity(params GetRoleByIdentityAndOrganizationIdentityParams) (*role.Role, error)
	GetSystemDefaultRole(params GetDefaultRoleParams) (*role.Role, error)
	PaginateRolesBy(params PaginateRolesParams) (*core.PaginationOutput[role.Role], error)
	ChangeRoleUsersToDefault(params ChangeRoleUsersToDefaultParams) error

	StoreRole(params StoreRoleParams) (*role.Role, error)
	UpdateRole(params UpdateRoleParams) error
	DeleteRole(params DeleteRoleParams) error
}
