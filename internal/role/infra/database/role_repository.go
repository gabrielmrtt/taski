package roledatabase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gabrielmrtt/taski/internal/core"
	coredatabase "github.com/gabrielmrtt/taski/internal/core/database"
	"github.com/gabrielmrtt/taski/internal/organization"
	"github.com/gabrielmrtt/taski/internal/role"
	rolerepo "github.com/gabrielmrtt/taski/internal/role/repository"
	"github.com/gabrielmrtt/taski/internal/user"
	userdatabase "github.com/gabrielmrtt/taski/internal/user/infra/database"
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

	RolePermissions []*RolePermissionTable  `bun:"rel:has-many,join:internal_id=role_internal_id"`
	Creator         *userdatabase.UserTable `bun:"rel:has-one,join:user_creator_internal_id=internal_id"`
	Editor          *userdatabase.UserTable `bun:"rel:has-one,join:user_editor_internal_id=internal_id"`
}

func (r *RoleTable) ToEntity() *role.Role {

	var userCreatorIdentity *core.Identity
	var userEditorIdentity *core.Identity

	if r.UserCreatorInternalId != nil {
		identity := core.NewIdentityFromInternal(uuid.MustParse(*r.UserCreatorInternalId), user.UserIdentityPrefix)
		userCreatorIdentity = &identity
	}

	if r.UserEditorInternalId != nil {
		identity := core.NewIdentityFromInternal(uuid.MustParse(*r.UserEditorInternalId), user.UserIdentityPrefix)
		userEditorIdentity = &identity
	}

	var organizationIdentity *core.Identity
	if r.OrganizationInternalId != nil {
		identity := core.NewIdentityFromInternal(uuid.MustParse(*r.OrganizationInternalId), organization.OrganizationIdentityPrefix)
		organizationIdentity = &identity
	}

	var creator *user.User = nil
	if r.Creator != nil {
		creator = r.Creator.ToEntity()
	}

	var editor *user.User = nil
	if r.Editor != nil {
		editor = r.Editor.ToEntity()
	}

	var permissions []role.Permission = make([]role.Permission, 0)
	for _, rolePermission := range r.RolePermissions {
		if rolePermission.Permission != nil {
			permissions = append(permissions, *rolePermission.ToEntity())
		}
	}

	var createdAt *core.DateTime = nil
	if r.CreatedAt != 0 {
		createdAt = &core.DateTime{Value: r.CreatedAt}
	}

	var updatedAt *core.DateTime = nil
	if r.UpdatedAt != nil {
		updatedAt = &core.DateTime{Value: *r.UpdatedAt}
	}

	var deletedAt *core.DateTime = nil
	if r.DeletedAt != nil {
		deletedAt = &core.DateTime{Value: *r.DeletedAt}
	}

	return &role.Role{
		Identity:             core.NewIdentityFromInternal(uuid.MustParse(r.InternalId), role.RoleIdentityPrefix),
		Name:                 r.Name,
		Slug:                 r.Slug,
		Description:          r.Description,
		Permissions:          permissions,
		OrganizationIdentity: organizationIdentity,
		UserCreatorIdentity:  userCreatorIdentity,
		UserEditorIdentity:   userEditorIdentity,
		IsSystemDefault:      r.IsSystemDefault,
		Creator:              creator,
		Editor:               editor,
		Timestamps: core.Timestamps{
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		},
		DeletedAt: deletedAt,
	}
}

type RolePermissionTable struct {
	bun.BaseModel `bun:"table:role_permission,alias:role_permission"`

	RoleInternalId       string `bun:"role_internal_id,pk,notnull,type:uuid"`
	PermissionInternalId string `bun:"permission_internal_id,pk,notnull,type:uuid"`

	Permission *PermissionTable `bun:"rel:has-one,join:permission_internal_id=internal_id"`
}

func (r *RolePermissionTable) ToEntity() *role.Permission {
	return &role.Permission{
		Identity:    core.NewIdentityWithoutPublicFromInternal(uuid.MustParse(r.PermissionInternalId)),
		Name:        r.Permission.Name,
		Description: r.Permission.Description,
		Slug:        role.PermissionSlugs(r.Permission.Slug),
	}
}

type RoleBunRepository struct {
	db *bun.DB
	tx *coredatabase.TransactionBun
}

func NewRoleBunRepository(connection *bun.DB) *RoleBunRepository {
	return &RoleBunRepository{db: connection, tx: nil}
}

func (r *RoleBunRepository) SetTransaction(tx core.Transaction) error {
	if r.tx != nil && !r.tx.IsClosed() {
		return nil
	}

	r.tx = tx.(*coredatabase.TransactionBun)
	return nil
}

func (r *RoleBunRepository) applyFilters(selectQuery *bun.SelectQuery, filters rolerepo.RoleFilters) *bun.SelectQuery {
	selectQuery = selectQuery.Where("roles.organization_internal_id = ?", filters.OrganizationIdentity.Internal.String())

	if filters.Name != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "roles.name", filters.Name)
	}

	if filters.Description != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "roles.description", filters.Description)
	}

	if filters.CreatedAt != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "roles.created_at", filters.CreatedAt)
	}

	if filters.UpdatedAt != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "roles.updated_at", filters.UpdatedAt)
	}

	if filters.DeletedAt != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "roles.deleted_at", filters.DeletedAt)
	}

	return selectQuery
}

func (r *RoleBunRepository) GetRoleByIdentity(params rolerepo.GetRoleByIdentityParams) (*role.Role, error) {
	var role *RoleTable = new(RoleTable)
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(role)
	selectQuery = selectQuery.Relation("RolePermissions.Permission")
	selectQuery = coredatabase.ApplyRelations(selectQuery, params.RelationsInput)
	selectQuery = selectQuery.Where("roles.internal_id = ?", params.RoleIdentity.Internal.String())

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

func (r *RoleBunRepository) GetRoleByIdentityAndOrganizationIdentity(params rolerepo.GetRoleByIdentityAndOrganizationIdentityParams) (*role.Role, error) {
	var role *RoleTable = new(RoleTable)
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(role)
	selectQuery = selectQuery.Relation("RolePermissions.Permission")
	selectQuery = coredatabase.ApplyRelations(selectQuery, params.RelationsInput)
	selectQuery = selectQuery.Where("roles.internal_id = ? AND roles.organization_internal_id = ?", params.RoleIdentity.Internal.String(), params.OrganizationIdentity.Internal.String())
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

func (r *RoleBunRepository) GetSystemDefaultRole(params rolerepo.GetDefaultRoleParams) (*role.Role, error) {
	var role *RoleTable = new(RoleTable)
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(role)
	selectQuery = selectQuery.Relation("RolePermissions.Permission")
	selectQuery = selectQuery.Where("roles.slug = ? AND roles.is_system_default = TRUE", string(params.Slug))
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

func (r *RoleBunRepository) PaginateRolesBy(params rolerepo.PaginateRolesParams) (*core.PaginationOutput[role.Role], error) {
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
	selectQuery = coredatabase.ApplyRelations(selectQuery, params.RelationsInput)
	selectQuery = r.applyFilters(selectQuery, params.Filters)

	if !params.ShowDeleted {
		selectQuery = selectQuery.Where("roles.deleted_at IS NULL")
	}

	countBeforePagination, err := selectQuery.Count(context.Background())
	if err != nil {
		return nil, err
	}

	selectQuery = coredatabase.ApplySort(selectQuery, params.SortInput)
	selectQuery = coredatabase.ApplyPagination(selectQuery, params.Pagination)
	err = selectQuery.Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return &core.PaginationOutput[role.Role]{
				Data:    []role.Role{},
				Page:    page,
				HasMore: false,
				Total:   0,
			}, nil
		}

		return nil, err
	}

	var roleEntities []role.Role = make([]role.Role, 0)
	for _, role := range roles {
		roleEntities = append(roleEntities, *role.ToEntity())
	}

	return &core.PaginationOutput[role.Role]{
		Data:    roleEntities,
		Page:    page,
		HasMore: core.HasMorePages(page, countBeforePagination, perPage),
		Total:   countBeforePagination,
	}, nil
}

func (r *RoleBunRepository) StoreRole(params rolerepo.StoreRoleParams) (*role.Role, error) {
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

	var createdAt *int64 = nil
	if params.Role.Timestamps.CreatedAt != nil {
		createdAt = &params.Role.Timestamps.CreatedAt.Value
	}

	var updatedAt *int64 = nil
	if params.Role.Timestamps.UpdatedAt != nil {
		updatedAt = &params.Role.Timestamps.UpdatedAt.Value
	}

	var deletedAt *int64 = nil
	if params.Role.DeletedAt != nil {
		deletedAt = &params.Role.DeletedAt.Value
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
		CreatedAt:              *createdAt,
		UpdatedAt:              updatedAt,
		DeletedAt:              deletedAt,
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

func (r *RoleBunRepository) UpdateRole(params rolerepo.UpdateRoleParams) error {
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

	var createdAt *int64 = nil
	if params.Role.Timestamps.CreatedAt != nil {
		createdAt = &params.Role.Timestamps.CreatedAt.Value
	}

	var updatedAt *int64 = nil
	if params.Role.Timestamps.UpdatedAt != nil {
		updatedAt = &params.Role.Timestamps.UpdatedAt.Value
	}

	var deletedAt *int64 = nil
	if params.Role.DeletedAt != nil {
		deletedAt = &params.Role.DeletedAt.Value
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
		CreatedAt:              *createdAt,
		UpdatedAt:              updatedAt,
		DeletedAt:              deletedAt,
	}

	_, err := tx.NewUpdate().Model(roleTable).Where("roles.internal_id = ?", params.Role.Identity.Internal.String()).Exec(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
	}

	if params.Role.Permissions != nil {
		_, err := tx.NewDelete().Model(&RolePermissionTable{}).Where("role_permission.role_internal_id = ?", params.Role.Identity.Internal.String()).Exec(context.Background())
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

func (r *RoleBunRepository) DeleteRole(params rolerepo.DeleteRoleParams) error {
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

	_, err := tx.NewDelete().Model(&RoleTable{}).Where("roles.internal_id = ?", params.RoleIdentity.Internal.String()).Exec(context.Background())
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

func (r *RoleBunRepository) ChangeRoleUsersToDefault(params rolerepo.ChangeRoleUsersToDefaultParams) error {
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

	defaultRole, err := r.GetSystemDefaultRole(rolerepo.GetDefaultRoleParams{
		Slug: params.DefaultRoleSlug,
	})
	if err != nil {
		return err
	}

	if defaultRole == nil {
		return errors.New("default role not found")
	}

	_, err = tx.NewRaw("UPDATE organization_user SET organization_user.role_internal_id = ? WHERE organization_user.role_internal_id = ?", defaultRole.Identity.Internal.String(), params.RoleIdentity.Internal.String()).Exec(context.Background())
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
