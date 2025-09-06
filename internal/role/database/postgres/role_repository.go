package role_database_postgres

import (
	"context"
	"database/sql"

	"github.com/gabrielmrtt/taski/internal/core"
	core_database_postgres "github.com/gabrielmrtt/taski/internal/core/database/postgres"
	role_core "github.com/gabrielmrtt/taski/internal/role"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type RoleTable struct {
	bun.BaseModel `bun:"table:roles,alias:roles"`

	InternalId             string `bun:"internal_id,pk,notnull,type:uuid"`
	PublicId               string `bun:"public_id,notnull,type:varchar(510)"`
	Name                   string `bun:"name,notnull,type:varchar(255)"`
	Description            string `bun:"description,type:varchar(510)"`
	OrganizationInternalId string `bun:"organization_internal_id,notnull,type:uuid"`
	UserCreatorInternalId  string `bun:"user_creator_internal_id,type:uuid"`
	UserEditorInternalId   string `bun:"user_editor_internal_id,type:uuid"`
	IsSystemDefault        bool   `bun:"is_system_default,notnull,type:boolean"`
	CreatedAt              int64  `bun:"created_at,notnull,type:bigint"`
	UpdatedAt              *int64 `bun:"updated_at,type:bigint"`
	DeletedAt              *int64 `bun:"deleted_at,type:bigint"`

	RolePermission *RolePermissionTable `bun:"rel:has-many,join:internal_id=role_internal_id"`
}

func (r *RoleTable) ToEntity() *role_core.Role {

	var userCreatorIdentity *core.Identity
	var userEditorIdentity *core.Identity

	if r.UserCreatorInternalId != "" {
		identity := core.NewIdentityFromInternal(uuid.MustParse(r.UserCreatorInternalId), "user")
		userCreatorIdentity = &identity
	}

	if r.UserEditorInternalId != "" {
		identity := core.NewIdentityFromInternal(uuid.MustParse(r.UserEditorInternalId), "user")
		userEditorIdentity = &identity
	}

	return &role_core.Role{
		Identity:             core.NewIdentityFromInternal(uuid.MustParse(r.InternalId), "role"),
		Name:                 r.Name,
		Description:          r.Description,
		OrganizationIdentity: core.NewIdentityFromInternal(uuid.MustParse(r.OrganizationInternalId), "organization"),
		UserCreatorIdentity:  userCreatorIdentity,
		UserEditorIdentity:   userEditorIdentity,
		IsSystemDefault:      r.IsSystemDefault,
		Timestamps: core.Timestamps{
			CreatedAt: &r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		},
		DeletedAt: r.DeletedAt,
	}
}

type RolePermissionTable struct {
	bun.BaseModel `bun:"table:role_permission,alias:role_permission"`

	RoleInternalId string `bun:"role_internal_id,notnull,type:uuid"`
	PermissionSlug string `bun:"permission_slug,notnull,type:varchar(510)"`

	Permission *PermissionTable `bun:"rel:has-one,join:permission_slug=slug"`
}

func (r *RolePermissionTable) ToEntity() *role_core.Permission {
	return &role_core.Permission{
		Name:        r.Permission.Name,
		Description: r.Permission.Description,
		Slug:        r.Permission.Slug,
	}
}

type RolePostgresRepository struct {
	db *bun.DB
	tx *core_database_postgres.TransactionPostgres
}

func NewRolePostgresRepository() *RolePostgresRepository {
	return &RolePostgresRepository{db: core_database_postgres.DB, tx: nil}
}

func (r *RolePostgresRepository) SetTransaction(tx core.Transaction) error {
	r.tx = tx.(*core_database_postgres.TransactionPostgres)
	return nil
}

func (r *RolePostgresRepository) applyFilters(selectQuery *bun.SelectQuery, filters role_core.RoleFilters) *bun.SelectQuery {
	if filters.Name != nil {
		selectQuery = core_database_postgres.ApplyComparableFilter(selectQuery, "name", filters.Name)
	}

	if filters.Description != nil {
		selectQuery = core_database_postgres.ApplyComparableFilter(selectQuery, "description", filters.Description)
	}

	if filters.CreatedAt != nil {
		selectQuery = core_database_postgres.ApplyComparableFilter(selectQuery, "created_at", filters.CreatedAt)
	}

	if filters.UpdatedAt != nil {
		selectQuery = core_database_postgres.ApplyComparableFilter(selectQuery, "updated_at", filters.UpdatedAt)
	}

	if filters.DeletedAt != nil {
		selectQuery = core_database_postgres.ApplyComparableFilter(selectQuery, "deleted_at", filters.DeletedAt)
	}

	return selectQuery
}

func (r *RolePostgresRepository) GetRoleByIdentity(params role_core.GetRoleByIdentityParams) (*role_core.Role, error) {
	var role RoleTable
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(&role).Relation("RolePermission.Permission").Where("internal_id = ?", params.Identity.Internal.String())

	err := selectQuery.Scan(context.Background())

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return role.ToEntity(), nil
}

func (r *RolePostgresRepository) ListRolesBy(params role_core.ListRolesParams) (*[]role_core.Role, error) {
	var roles []RoleTable
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(&roles).Relation("RolePermission.Permission")
	selectQuery = r.applyFilters(selectQuery, params.Filters)

	err := selectQuery.Scan(context.Background())

	if err != nil {
		if err == sql.ErrNoRows {
			return &[]role_core.Role{}, nil
		}
	}

	var roleEntities []role_core.Role

	for _, role := range roles {
		roleEntities = append(roleEntities, *role.ToEntity())
	}

	return &roleEntities, nil
}

func (r *RolePostgresRepository) PaginateRolesBy(params role_core.PaginateRolesParams) (*core.PaginationOutput[role_core.Role], error) {
	var roles []RoleTable
	var selectQuery *bun.SelectQuery
	var perPage int = 10
	var page int = 1

	if params.Pagination.PerPage != nil {
		perPage = *params.Pagination.PerPage
	}

	if params.Pagination.Page != nil {
		page = *params.Pagination.Page
	}

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(&roles).Relation("RolePermission.Permission")
	selectQuery = r.applyFilters(selectQuery, params.Filters)

	countBeforePagination, err := selectQuery.Count(context.Background())

	if err != nil {
		return nil, err
	}

	selectQuery = core_database_postgres.ApplyPagination(selectQuery, params.Pagination)

	err = selectQuery.Scan(context.Background(), &roles)

	if err != nil {
		if err == sql.ErrNoRows {
			return &core.PaginationOutput[role_core.Role]{
				Data:    []role_core.Role{},
				Page:    page,
				HasMore: false,
				Total:   0,
			}, nil
		}

		return nil, err
	}

	var roleEntities []role_core.Role

	for _, role := range roles {
		roleEntities = append(roleEntities, *role.ToEntity())
	}

	return &core.PaginationOutput[role_core.Role]{
		Data:    roleEntities,
		Page:    page,
		HasMore: core.HasMorePages(page, countBeforePagination, perPage),
		Total:   countBeforePagination,
	}, nil
}

func (r *RolePostgresRepository) StoreRole(role *role_core.Role) (*role_core.Role, error) {
	var tx bun.Tx
	var shouldCommit bool = false

	if r.tx != nil && !r.tx.IsClosed() {
		tx = *r.tx.Tx
	} else {
		var err error
		tx, err = r.db.BeginTx(context.Background(), nil)
		shouldCommit = true

		if err != nil {
			return nil, err
		}
	}

	roleTable := &RoleTable{
		InternalId:             role.Identity.Internal.String(),
		PublicId:               role.Identity.Public,
		Name:                   role.Name,
		Description:            role.Description,
		OrganizationInternalId: role.OrganizationIdentity.Internal.String(),
		UserCreatorInternalId:  role.UserCreatorIdentity.Internal.String(),
		UserEditorInternalId:   role.UserEditorIdentity.Internal.String(),
		IsSystemDefault:        role.IsSystemDefault,
		CreatedAt:              *role.Timestamps.CreatedAt,
		UpdatedAt:              role.Timestamps.UpdatedAt,
		DeletedAt:              role.DeletedAt,
	}

	_, err := tx.NewInsert().Model(roleTable).Exec(context.Background())

	if err != nil {
		return nil, err
	}

	for _, permission := range role.Permissions {
		permissionTable := &RolePermissionTable{
			RoleInternalId: roleTable.InternalId,
			PermissionSlug: permission.Slug,
		}

		_, err = tx.NewInsert().Model(permissionTable).Exec(context.Background())

		if err != nil {
			return nil, err
		}
	}

	if shouldCommit {
		err = tx.Commit()
		if err != nil {
			return nil, err
		}
	}

	return roleTable.ToEntity(), nil
}

func (r *RolePostgresRepository) UpdateRole(role *role_core.Role) error {
	var tx bun.Tx
	var shouldCommit bool = false

	if r.tx != nil && !r.tx.IsClosed() {
		tx = *r.tx.Tx
	} else {
		var err error
		tx, err = r.db.BeginTx(context.Background(), nil)
		shouldCommit = true

		if err != nil {
			return err
		}
	}

	roleTable := &RoleTable{
		InternalId:             role.Identity.Internal.String(),
		PublicId:               role.Identity.Public,
		Name:                   role.Name,
		Description:            role.Description,
		OrganizationInternalId: role.OrganizationIdentity.Internal.String(),
		UserCreatorInternalId:  role.UserCreatorIdentity.Internal.String(),
		UserEditorInternalId:   role.UserEditorIdentity.Internal.String(),
		IsSystemDefault:        role.IsSystemDefault,
		CreatedAt:              *role.Timestamps.CreatedAt,
		UpdatedAt:              role.Timestamps.UpdatedAt,
		DeletedAt:              role.DeletedAt,
	}

	_, err := tx.NewUpdate().Model(roleTable).Where("internal_id = ?", role.Identity.Internal.String()).Exec(context.Background())

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
	}

	if role.Permissions != nil {
		_, err := tx.NewDelete().Model(&RolePermissionTable{}).Where("role_internal_id = ?", role.Identity.Internal.String()).Exec(context.Background())

		if err != nil {
			return err
		}

		for _, permission := range role.Permissions {
			permissionTable := &RolePermissionTable{
				RoleInternalId: roleTable.InternalId,
				PermissionSlug: permission.Slug,
			}

			_, err = tx.NewInsert().Model(permissionTable).Exec(context.Background())

			if err != nil {
				return err
			}
		}
	}

	if shouldCommit {
		err = tx.Commit()
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *RolePostgresRepository) DeleteRole(roleIdentity core.Identity) error {
	var tx bun.Tx
	var shouldCommit bool = false

	if r.tx != nil && !r.tx.IsClosed() {
		tx = *r.tx.Tx
	} else {
		var err error
		tx, err = r.db.BeginTx(context.Background(), nil)
		shouldCommit = true

		if err != nil {
			return err
		}
	}

	_, err := tx.NewDelete().Model(&RoleTable{}).Where("internal_id = ?", roleIdentity.Internal.String()).Exec(context.Background())

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}

		return err
	}

	if shouldCommit {
		err = tx.Commit()
		if err != nil {
			return err
		}
	}

	return nil
}
