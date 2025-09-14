package role_database_postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gabrielmrtt/taski/internal/core"
	core_database_postgres "github.com/gabrielmrtt/taski/internal/core/database/postgres"
	organization_core "github.com/gabrielmrtt/taski/internal/organization"
	role_core "github.com/gabrielmrtt/taski/internal/role"
	role_repositories "github.com/gabrielmrtt/taski/internal/role/repositories"
	user_core "github.com/gabrielmrtt/taski/internal/user"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type RoleTable struct {
	bun.BaseModel `bun:"table:roles,alias:roles"`

	InternalId             string  `bun:"internal_id,pk,notnull,type:uuid"`
	PublicId               string  `bun:"public_id,notnull,type:varchar(510)"`
	Name                   string  `bun:"name,notnull,type:varchar(255)"`
	Slug                   string  `bun:"slug,notnull,type:varchar(255)"`
	Description            string  `bun:"description,type:varchar(510)"`
	OrganizationInternalId *string `bun:"organization_internal_id,type:uuid"`
	UserCreatorInternalId  *string `bun:"user_creator_internal_id,type:uuid"`
	UserEditorInternalId   *string `bun:"user_editor_internal_id,type:uuid"`
	IsSystemDefault        bool    `bun:"is_system_default,notnull,type:boolean"`
	CreatedAt              int64   `bun:"created_at,notnull,type:bigint"`
	UpdatedAt              *int64  `bun:"updated_at,type:bigint"`
	DeletedAt              *int64  `bun:"deleted_at,type:bigint"`

	RolePermissions []*RolePermissionTable `bun:"rel:has-many,join:internal_id=role_internal_id"`
}

func (r *RoleTable) ToEntity() *role_core.Role {

	var userCreatorIdentity *core.Identity
	var userEditorIdentity *core.Identity

	if r.UserCreatorInternalId != nil {
		identity := core.NewIdentityFromInternal(uuid.MustParse(*r.UserCreatorInternalId), user_core.UserIdentityPrefix)
		userCreatorIdentity = &identity
	}

	if r.UserEditorInternalId != nil {
		identity := core.NewIdentityFromInternal(uuid.MustParse(*r.UserEditorInternalId), user_core.UserIdentityPrefix)
		userEditorIdentity = &identity
	}

	var organizationIdentity *core.Identity
	if r.OrganizationInternalId != nil {
		identity := core.NewIdentityFromInternal(uuid.MustParse(*r.OrganizationInternalId), organization_core.OrganizationIdentityPrefix)
		organizationIdentity = &identity
	}

	var permissions []role_core.Permission = make([]role_core.Permission, 0)
	for _, rolePermission := range r.RolePermissions {
		if rolePermission.Permission != nil {
			permissions = append(permissions, *rolePermission.ToEntity())
		}
	}

	return &role_core.Role{
		Identity:             core.NewIdentityFromInternal(uuid.MustParse(r.InternalId), role_core.RoleIdentityPrefix),
		Name:                 r.Name,
		Slug:                 r.Slug,
		Description:          r.Description,
		Permissions:          permissions,
		OrganizationIdentity: organizationIdentity,
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

	RoleInternalId       string `bun:"role_internal_id,pk,notnull,type:uuid"`
	PermissionInternalId string `bun:"permission_internal_id,pk,notnull,type:uuid"`

	Permission *PermissionTable `bun:"rel:has-one,join:permission_internal_id=internal_id"`
}

func (r *RolePermissionTable) ToEntity() *role_core.Permission {
	return &role_core.Permission{
		Identity:    core.NewIdentityWithoutPublicFromInternal(uuid.MustParse(r.PermissionInternalId)),
		Name:        r.Permission.Name,
		Description: r.Permission.Description,
		Slug:        role_core.PermissionSlugs(r.Permission.Slug),
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

func (r *RolePostgresRepository) applyFilters(selectQuery *bun.SelectQuery, filters role_repositories.RoleFilters) *bun.SelectQuery {
	selectQuery = selectQuery.Where("organization_internal_id = ?", filters.OrganizationIdentity.Internal.String())

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

func (r *RolePostgresRepository) GetRoleByIdentity(params role_repositories.GetRoleByIdentityParams) (*role_core.Role, error) {
	var role *RoleTable = new(RoleTable)
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(role)
	selectQuery = selectQuery.Relation("RolePermissions.Permission")
	selectQuery = core_database_postgres.ApplyRelations(selectQuery, params.RelationsInput)
	selectQuery = selectQuery.Where("internal_id = ?", params.RoleIdentity.Internal.String())

	err := selectQuery.Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	if role.InternalId == "" {
		return nil, nil
	}

	return role.ToEntity(), nil
}

func (r *RolePostgresRepository) GetRoleByIdentityAndOrganizationIdentity(params role_repositories.GetRoleByIdentityAndOrganizationIdentityParams) (*role_core.Role, error) {
	var role *RoleTable = new(RoleTable)
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(role)
	selectQuery = selectQuery.Relation("RolePermissions.Permission")
	selectQuery = core_database_postgres.ApplyRelations(selectQuery, params.RelationsInput)
	selectQuery = selectQuery.Where("internal_id = ? AND organization_internal_id = ?", params.RoleIdentity.Internal.String(), params.OrganizationIdentity.Internal.String())
	err := selectQuery.Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
	}

	if role.InternalId == "" {
		return nil, nil
	}

	return role.ToEntity(), nil
}

func (r *RolePostgresRepository) GetSystemDefaultRole(params role_repositories.GetDefaultRoleParams) (*role_core.Role, error) {
	var role *RoleTable = new(RoleTable)
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(role)
	selectQuery = selectQuery.Relation("RolePermissions.Permission")
	selectQuery = core_database_postgres.ApplyRelations(selectQuery, params.RelationsInput)
	selectQuery = selectQuery.Where("slug = ? AND is_system_default = TRUE", string(params.Slug))
	err := selectQuery.Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	if role.InternalId == "" {
		return nil, nil
	}

	return role.ToEntity(), nil
}

func (r *RolePostgresRepository) PaginateRolesBy(params role_repositories.PaginateRolesParams) (*core.PaginationOutput[role_core.Role], error) {
	var roles []RoleTable = make([]RoleTable, 0)
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

	selectQuery = selectQuery.Model(&roles)
	selectQuery = selectQuery.Relation("RolePermissions.Permission")
	selectQuery = core_database_postgres.ApplyRelations(selectQuery, params.RelationsInput)
	selectQuery = r.applyFilters(selectQuery, params.Filters)
	countBeforePagination, err := selectQuery.Count(context.Background())
	if err != nil {
		return nil, err
	}

	selectQuery = core_database_postgres.ApplyPagination(selectQuery, params.Pagination)
	err = selectQuery.Scan(context.Background())
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

	var roleEntities []role_core.Role = make([]role_core.Role, 0)
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

func (r *RolePostgresRepository) StoreRole(params role_repositories.StoreRoleParams) (*role_core.Role, error) {
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

	var organizationInternalId *string
	if params.Role.OrganizationIdentity != nil {
		identity := params.Role.OrganizationIdentity.Internal.String()
		organizationInternalId = &identity
	}

	var userCreatorInternalId *string
	if params.Role.UserCreatorIdentity != nil {
		identity := params.Role.UserCreatorIdentity.Internal.String()
		userCreatorInternalId = &identity
	}

	var userEditorInternalId *string
	if params.Role.UserEditorIdentity != nil {
		identity := params.Role.UserEditorIdentity.Internal.String()
		userEditorInternalId = &identity
	}

	roleTable := &RoleTable{
		InternalId:             params.Role.Identity.Internal.String(),
		PublicId:               params.Role.Identity.Public,
		Name:                   params.Role.Name,
		Slug:                   params.Role.Slug,
		Description:            params.Role.Description,
		OrganizationInternalId: organizationInternalId,
		UserCreatorInternalId:  userCreatorInternalId,
		UserEditorInternalId:   userEditorInternalId,
		IsSystemDefault:        params.Role.IsSystemDefault,
		CreatedAt:              *params.Role.Timestamps.CreatedAt,
		UpdatedAt:              params.Role.Timestamps.UpdatedAt,
		DeletedAt:              params.Role.DeletedAt,
	}

	_, err := tx.NewInsert().Model(roleTable).Exec(context.Background())
	if err != nil {
		return nil, err
	}

	for _, permission := range params.Role.Permissions {
		permissionTable := &RolePermissionTable{
			RoleInternalId:       roleTable.InternalId,
			PermissionInternalId: permission.Identity.Internal.String(),
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

	return params.Role, nil
}

func (r *RolePostgresRepository) UpdateRole(params role_repositories.UpdateRoleParams) error {
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

	var organizationInternalId *string
	if params.Role.OrganizationIdentity != nil {
		identity := params.Role.OrganizationIdentity.Internal.String()
		organizationInternalId = &identity
	}

	var userCreatorInternalId *string
	if params.Role.UserCreatorIdentity != nil {
		identity := params.Role.UserCreatorIdentity.Internal.String()
		userCreatorInternalId = &identity
	}

	var userEditorInternalId *string
	if params.Role.UserEditorIdentity != nil {
		identity := params.Role.UserEditorIdentity.Internal.String()
		userEditorInternalId = &identity
	}

	roleTable := &RoleTable{
		InternalId:             params.Role.Identity.Internal.String(),
		PublicId:               params.Role.Identity.Public,
		Name:                   params.Role.Name,
		Slug:                   params.Role.Slug,
		Description:            params.Role.Description,
		OrganizationInternalId: organizationInternalId,
		UserCreatorInternalId:  userCreatorInternalId,
		UserEditorInternalId:   userEditorInternalId,
		IsSystemDefault:        params.Role.IsSystemDefault,
		CreatedAt:              *params.Role.Timestamps.CreatedAt,
		UpdatedAt:              params.Role.Timestamps.UpdatedAt,
		DeletedAt:              params.Role.DeletedAt,
	}

	_, err := tx.NewUpdate().Model(roleTable).Where("internal_id = ?", params.Role.Identity.Internal.String()).Exec(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
	}

	if params.Role.Permissions != nil {
		_, err := tx.NewDelete().Model(&RolePermissionTable{}).Where("role_internal_id = ?", params.Role.Identity.Internal.String()).Exec(context.Background())
		if err != nil {
			return err
		}

		for _, permission := range params.Role.Permissions {
			permissionTable := &RolePermissionTable{
				RoleInternalId:       roleTable.InternalId,
				PermissionInternalId: permission.Identity.Internal.String(),
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

func (r *RolePostgresRepository) DeleteRole(params role_repositories.DeleteRoleParams) error {
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

	_, err := tx.NewDelete().Model(&RoleTable{}).Where("internal_id = ?", params.RoleIdentity.Internal.String()).Exec(context.Background())
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

func (r *RolePostgresRepository) ChangeRoleUsersToDefault(params role_repositories.ChangeRoleUsersToDefaultParams) error {
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

	defaultRole, err := r.GetSystemDefaultRole(role_repositories.GetDefaultRoleParams{
		Slug: params.DefaultRoleSlug,
	})
	if err != nil {
		return err
	}

	if defaultRole == nil {
		return errors.New("default role not found")
	}

	_, err = tx.NewRaw("UPDATE organization_user SET role_internal_id = ? WHERE role_internal_id = ?", defaultRole.Identity.Internal.String(), params.RoleIdentity.Internal.String()).Exec(context.Background())
	if err != nil {
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
