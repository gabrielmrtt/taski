package rolerepo

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/role"
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
}

type PaginatePermissionsParams struct {
	Filters    PermissionFilters
	SortInput  core.SortInput
	Pagination core.PaginationInput
}

type StorePermissionParams struct {
	Permission *role.Permission
}

type UpdatePermissionParams struct {
	Permission *role.Permission
}

type DeletePermissionParams struct {
	PermissionSlug string
}

type PermissionRepository interface {
	SetTransaction(tx core.Transaction) error

	GetPermissionBySlug(params GetPermissionBySlugParams) (*role.Permission, error)
	ListPermissionsBy(params ListPermissionsParams) (*[]role.Permission, error)
	PaginatePermissionsBy(params PaginatePermissionsParams) (*core.PaginationOutput[role.Permission], error)

	StorePermission(params StorePermissionParams) (*role.Permission, error)
	UpdatePermission(params UpdatePermissionParams) error
	DeletePermission(params DeletePermissionParams) error
}
