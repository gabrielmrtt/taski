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

type PermissionTable struct {
	bun.BaseModel `bun:"table:permissions,alias:permissions"`

	InternalId  string `bun:"internal_id,pk,notnull,type:uuid"`
	Slug        string `bun:"slug,pk,notnull,type:varchar(510)"`
	Name        string `bun:"name,notnull,type:varchar(255)"`
	Description string `bun:"description,type:varchar(510)"`
}

func (p *PermissionTable) ToEntity() *role_core.Permission {
	return &role_core.Permission{
		Identity:    core.NewIdentityFromInternal(uuid.MustParse(p.InternalId), "permission"),
		Name:        p.Name,
		Description: p.Description,
		Slug:        p.Slug,
	}
}

type PermissionPostgresRepository struct {
	db *bun.DB
	tx *core_database_postgres.TransactionPostgres
}

func NewPermissionPostgresRepository() *PermissionPostgresRepository {
	return &PermissionPostgresRepository{db: core_database_postgres.DB, tx: nil}
}

func (r *PermissionPostgresRepository) SetTransaction(tx core.Transaction) error {
	r.tx = tx.(*core_database_postgres.TransactionPostgres)
	return nil
}

func (r *PermissionPostgresRepository) applyFilters(selectQuery *bun.SelectQuery, filters role_core.PermissionFilters) *bun.SelectQuery {
	if filters.Name != nil {
		selectQuery = core_database_postgres.ApplyComparableFilter(selectQuery, "name", filters.Name)
	}

	if filters.Slug != nil {
		selectQuery = core_database_postgres.ApplyComparableFilter(selectQuery, "slug", filters.Slug)
	}

	return selectQuery
}

func (r *PermissionPostgresRepository) GetPermissionBySlug(params role_core.GetPermissionBySlugParams) (*role_core.Permission, error) {
	var permission PermissionTable
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(&permission).Where("slug = ?", params.Slug)

	err := selectQuery.Scan(context.Background())

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
	}

	return permission.ToEntity(), nil
}

func (r *PermissionPostgresRepository) PaginatePermissionsBy(params role_core.PaginatePermissionsParams) (*core.PaginationOutput[role_core.Permission], error) {
	var permissions []PermissionTable
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

	selectQuery = selectQuery.Model(&permissions)
	selectQuery = r.applyFilters(selectQuery, params.Filters)
	countBeforePagination, err := selectQuery.Count(context.Background())
	if err != nil {
		return nil, err
	}

	selectQuery = core_database_postgres.ApplyPagination(selectQuery, params.Pagination)
	err = selectQuery.Scan(context.Background(), &permissions)
	if err != nil {
		return nil, err
	}

	var permissionEntities []role_core.Permission
	for _, permission := range permissions {
		permissionEntities = append(permissionEntities, *permission.ToEntity())
	}

	return &core.PaginationOutput[role_core.Permission]{
		Data:    permissionEntities,
		Page:    page,
		HasMore: core.HasMorePages(page, countBeforePagination, perPage),
		Total:   countBeforePagination,
	}, nil
}

func (r *PermissionPostgresRepository) StorePermission(params role_core.StorePermissionParams) (*role_core.Permission, error) {
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

	permissionTable := &PermissionTable{
		InternalId:  params.Permission.Identity.Internal.String(),
		Slug:        params.Permission.Slug,
		Name:        params.Permission.Name,
		Description: params.Permission.Description,
	}

	_, err := tx.NewInsert().Model(permissionTable).Exec(context.Background())
	if err != nil {
		return nil, err
	}

	if shouldCommit {
		err = tx.Commit()
	}

	return permissionTable.ToEntity(), nil
}

func (r *PermissionPostgresRepository) UpdatePermission(params role_core.UpdatePermissionParams) error {
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

	permissionTable := &PermissionTable{
		InternalId:  params.Permission.Identity.Internal.String(),
		Slug:        params.Permission.Slug,
		Name:        params.Permission.Name,
		Description: params.Permission.Description,
	}

	_, err := tx.NewUpdate().Model(permissionTable).Where("slug = ?", params.Permission.Slug).Exec(context.Background())
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

func (r *PermissionPostgresRepository) DeletePermission(params role_core.DeletePermissionParams) error {
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

	_, err := tx.NewDelete().Model(&PermissionTable{}).Where("slug = ?", params.PermissionSlug).Exec(context.Background())
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
