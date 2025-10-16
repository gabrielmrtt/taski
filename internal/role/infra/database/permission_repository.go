package roledatabase

import (
	"context"
	"database/sql"

	"github.com/gabrielmrtt/taski/internal/core"
	coredatabase "github.com/gabrielmrtt/taski/internal/core/database"
	"github.com/gabrielmrtt/taski/internal/role"
	rolerepo "github.com/gabrielmrtt/taski/internal/role/repository"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type PermissionTable struct {
	bun.BaseModel `bun:"table:permissions,alias:permissions"`

	InternalId  string `bun:"internal_id,pk,notnull,type:uuid"`
	Slug        string `bun:"slug,notnull,type:varchar(510)"`
	Name        string `bun:"name,notnull,type:varchar(255)"`
	Description string `bun:"description,type:varchar(510)"`
}

func (p *PermissionTable) ToEntity() *role.Permission {
	return &role.Permission{
		Identity:    core.NewIdentityWithoutPublicFromInternal(uuid.MustParse(p.InternalId)),
		Name:        p.Name,
		Description: p.Description,
		Slug:        role.PermissionSlugs(p.Slug),
	}
}

type PermissionBunRepository struct {
	db *bun.DB
	tx *coredatabase.TransactionBun
}

func NewPermissionBunRepository(connection *bun.DB) *PermissionBunRepository {
	return &PermissionBunRepository{db: connection, tx: nil}
}

func (r *PermissionBunRepository) SetTransaction(tx core.Transaction) error {
	if r.tx != nil && !r.tx.IsClosed() {
		return nil
	}

	r.tx = tx.(*coredatabase.TransactionBun)
	return nil
}

func (r *PermissionBunRepository) applyFilters(selectQuery *bun.SelectQuery, filters rolerepo.PermissionFilters) *bun.SelectQuery {
	if filters.Name != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "name", filters.Name)
	}

	if filters.Slug != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "slug", filters.Slug)
	}

	return selectQuery
}

func (r *PermissionBunRepository) GetPermissionBySlug(params rolerepo.GetPermissionBySlugParams) (*role.Permission, error) {
	var permission *PermissionTable = new(PermissionTable)
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(permission).Where("slug = ?", params.Slug)

	err := selectQuery.Scan(context.Background())

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if permission.InternalId == "" {
		return nil, nil
	}

	return permission.ToEntity(), nil
}

func (r *PermissionBunRepository) ListPermissionsBy(params rolerepo.ListPermissionsParams) (*[]role.Permission, error) {
	var permissions []PermissionTable
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(&permissions)
	selectQuery = r.applyFilters(selectQuery, params.Filters)
	err := selectQuery.Scan(context.Background())
	if err != nil {
		return nil, err
	}

	var permissionEntities []role.Permission
	for _, permission := range permissions {
		permissionEntities = append(permissionEntities, *permission.ToEntity())
	}

	return &permissionEntities, nil
}

func (r *PermissionBunRepository) PaginatePermissionsBy(params rolerepo.PaginatePermissionsParams) (*core.PaginationOutput[role.Permission], error) {
	var permissions []PermissionTable = make([]PermissionTable, 0)
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

	selectQuery = coredatabase.ApplySort(selectQuery, params.SortInput)
	selectQuery = coredatabase.ApplyPagination(selectQuery, params.Pagination)
	err = selectQuery.Scan(context.Background())
	if err != nil {
		return nil, err
	}

	var permissionEntities []role.Permission = make([]role.Permission, 0)
	for _, permission := range permissions {
		permissionEntities = append(permissionEntities, *permission.ToEntity())
	}

	return &core.PaginationOutput[role.Permission]{
		Data:    permissionEntities,
		Page:    page,
		HasMore: core.HasMorePages(page, countBeforePagination, perPage),
		Total:   countBeforePagination,
	}, nil
}

func (r *PermissionBunRepository) StorePermission(params rolerepo.StorePermissionParams) (*role.Permission, error) {
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
		Slug:        string(params.Permission.Slug),
		Name:        params.Permission.Name,
		Description: params.Permission.Description,
	}

	_, err := tx.NewInsert().Model(permissionTable).Exec(context.Background())
	if err != nil {
		return nil, err
	}

	if shouldCommit {
		err = tx.Commit()
		if err != nil {
			return nil, err
		}
	}

	return params.Permission, nil
}

func (r *PermissionBunRepository) UpdatePermission(params rolerepo.UpdatePermissionParams) error {
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
		Slug:        string(params.Permission.Slug),
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

func (r *PermissionBunRepository) DeletePermission(params rolerepo.DeletePermissionParams) error {
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
