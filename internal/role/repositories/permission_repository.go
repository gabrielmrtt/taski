package role_repositories

import (
	"github.com/gabrielmrtt/taski/internal/core"
	role_core "github.com/gabrielmrtt/taski/internal/role"
)

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

type ListPermissionsParams struct {
	Filters PermissionFilters
	Include map[string]any
}

type PaginatePermissionsParams struct {
	Filters    PermissionFilters
	SortInput  *core.SortInput
	Pagination *core.PaginationInput
}

type StorePermissionParams struct {
	Permission *role_core.Permission
}

type UpdatePermissionParams struct {
	Permission *role_core.Permission
}

type DeletePermissionParams struct {
	PermissionSlug string
}

type PermissionRepository interface {
	SetTransaction(tx core.Transaction) error

	GetPermissionBySlug(params GetPermissionBySlugParams) (*role_core.Permission, error)
	ListPermissionsBy(params ListPermissionsParams) (*[]role_core.Permission, error)
	PaginatePermissionsBy(params PaginatePermissionsParams) (*core.PaginationOutput[role_core.Permission], error)

	StorePermission(params StorePermissionParams) (*role_core.Permission, error)
	UpdatePermission(params UpdatePermissionParams) error
	DeletePermission(params DeletePermissionParams) error
}
