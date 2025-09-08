package role_core

import (
	"github.com/gabrielmrtt/taski/internal/core"
)

type GetRoleByIdentityParams struct {
	Identity core.Identity
	Include  map[string]any
}

type GetRoleByIdentityAndOrganizationIdentityParams struct {
	Identity             core.Identity
	OrganizationIdentity core.Identity
	Include              map[string]any
}

type GetDefaultRoleParams struct {
	Slug string
}

type RoleFilters struct {
	OrganizationIdentity core.Identity
	Name                 *core.ComparableFilter[string]
	Description          *core.ComparableFilter[string]
	CreatedAt            *core.ComparableFilter[int64]
	UpdatedAt            *core.ComparableFilter[int64]
	DeletedAt            *core.ComparableFilter[int64]
}

type ListRolesParams struct {
	Filters RoleFilters
	Include map[string]any
}

type PaginateRolesParams struct {
	Filters    RoleFilters
	Include    map[string]any
	SortInput  *core.SortInput
	Pagination *core.PaginationInput
}

type RoleRepository interface {
	SetTransaction(tx core.Transaction) error

	GetRoleByIdentity(params GetRoleByIdentityParams) (*Role, error)
	GetRoleByIdentityAndOrganizationIdentity(params GetRoleByIdentityAndOrganizationIdentityParams) (*Role, error)
	GetSystemDefaultRole(params GetDefaultRoleParams) (*Role, error)
	ListRolesBy(params ListRolesParams) (*[]Role, error)
	PaginateRolesBy(params PaginateRolesParams) (*core.PaginationOutput[Role], error)
	ChangeRoleUsersToDefault(roleIdentity core.Identity, defaultRoleSlug string) error

	CheckIfOrganizatonHasUser(organizationIdentity core.Identity, userIdentity core.Identity) (bool, error)

	StoreRole(role *Role) (*Role, error)
	UpdateRole(role *Role) error
	DeleteRole(roleIdentity core.Identity) error
}

type GetPermissionByIdentityParams struct {
	Identity core.Identity
	Include  map[string]any
}

type GetPermissionBySlugParams struct {
	Slug    string
	Include map[string]any
}

type PermissionFilters struct {
	Name *core.ComparableFilter[string]
	Slug *core.ComparableFilter[string]
}

type ListPermissionsParams struct {
	Filters PermissionFilters
	Include map[string]any
}

type PaginatePermissionsParams struct {
	Filters    PermissionFilters
	Include    map[string]any
	Pagination *core.PaginationInput
}

type PermissionRepository interface {
	SetTransaction(tx core.Transaction) error

	GetPermissionBySlug(params GetPermissionBySlugParams) (*Permission, error)
	ListPermissionsBy(params ListPermissionsParams) (*[]Permission, error)
	PaginatePermissionsBy(params PaginatePermissionsParams) (*core.PaginationOutput[Permission], error)

	StorePermission(permission *Permission) (*Permission, error)
	UpdatePermission(permission *Permission) error
	DeletePermission(permissionSlug string) error
}
