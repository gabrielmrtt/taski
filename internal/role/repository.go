package role_core

import (
	"github.com/gabrielmrtt/taski/internal/core"
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
	Role *Role
}

type UpdateRoleParams struct {
	Role *Role
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

	GetRoleByIdentity(params GetRoleByIdentityParams) (*Role, error)
	GetRoleByIdentityAndOrganizationIdentity(params GetRoleByIdentityAndOrganizationIdentityParams) (*Role, error)
	GetSystemDefaultRole(params GetDefaultRoleParams) (*Role, error)
	PaginateRolesBy(params PaginateRolesParams) (*core.PaginationOutput[Role], error)
	ChangeRoleUsersToDefault(params ChangeRoleUsersToDefaultParams) error

	CheckIfOrganizatonHasUser(organizationIdentity core.Identity, userIdentity core.Identity) (bool, error)

	StoreRole(params StoreRoleParams) (*Role, error)
	UpdateRole(params UpdateRoleParams) error
	DeleteRole(params DeleteRoleParams) error
}

type PermissionFilters struct {
	Name *core.ComparableFilter[string]
	Slug *core.ComparableFilter[string]
}

type GetPermissionByIdentityParams struct {
	Identity core.Identity
}

type GetPermissionBySlugParams struct {
	Slug string
}

type PaginatePermissionsParams struct {
	Filters    PermissionFilters
	SortInput  *core.SortInput
	Pagination *core.PaginationInput
}

type StorePermissionParams struct {
	Permission *Permission
}

type UpdatePermissionParams struct {
	Permission *Permission
}

type DeletePermissionParams struct {
	PermissionSlug string
}

type PermissionRepository interface {
	SetTransaction(tx core.Transaction) error

	GetPermissionBySlug(params GetPermissionBySlugParams) (*Permission, error)
	PaginatePermissionsBy(params PaginatePermissionsParams) (*core.PaginationOutput[Permission], error)

	StorePermission(params StorePermissionParams) (*Permission, error)
	UpdatePermission(params UpdatePermissionParams) error
	DeletePermission(params DeletePermissionParams) error
}
