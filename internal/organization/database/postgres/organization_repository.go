package organization_database_postgres

import (
	"context"
	"database/sql"

	"github.com/gabrielmrtt/taski/internal/core"
	core_database_postgres "github.com/gabrielmrtt/taski/internal/core/database/postgres"
	organization_core "github.com/gabrielmrtt/taski/internal/organization"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type OrganizationTable struct {
	bun.BaseModel `bun:"table:organization,alias:organization"`

	InternalId            string  `bun:"internal_id,pk,notnull,type:uuid"`
	PublicId              string  `bun:"public_id,notnull,type:varchar(510)"`
	Name                  string  `bun:"name,notnull,type:varchar(255)"`
	Status                string  `bun:"status,notnull,type:varchar(100)"`
	UserCreatorInternalId *string `bun:"user_creator_internal_id,notnull,type:uuid"`
	UserEditorInternalId  *string `bun:"user_editor_internal_id,type:uuid"`
	CreatedAt             int64   `bun:"created_at,notnull,type:bigint"`
	UpdatedAt             *int64  `bun:"updated_at,type:bigint"`
	DeletedAt             *int64  `bun:"deleted_at,type:bigint"`
}

func (o *OrganizationTable) ToEntity() *organization_core.Organization {
	var userCreatorIdentity *core.Identity
	var userEditorIdentity *core.Identity

	if o.UserCreatorInternalId != nil {
		identity := core.NewIdentityFromInternal(uuid.MustParse(*o.UserCreatorInternalId), "usr")
		userCreatorIdentity = &identity
	}

	if o.UserEditorInternalId != nil {
		identity := core.NewIdentityFromInternal(uuid.MustParse(*o.UserEditorInternalId), "usr")
		userEditorIdentity = &identity
	}

	return &organization_core.Organization{
		Identity:            core.NewIdentityFromInternal(uuid.MustParse(o.InternalId), "org"),
		Name:                o.Name,
		Status:              organization_core.OrganizationStatuses(o.Status),
		UserCreatorIdentity: userCreatorIdentity,
		UserEditorIdentity:  userEditorIdentity,
		Timestamps: core.Timestamps{
			CreatedAt: &o.CreatedAt,
			UpdatedAt: o.UpdatedAt,
		},
		DeletedAt: o.DeletedAt,
	}
}

type OrganizationPostgresRepository struct {
	db *bun.DB
	tx *core_database_postgres.TransactionPostgres
}

func NewOrganizationPostgresRepository() *OrganizationPostgresRepository {
	return &OrganizationPostgresRepository{db: core_database_postgres.DB, tx: nil}
}

func (r *OrganizationPostgresRepository) SetTransaction(tx core.Transaction) error {
	r.tx = tx.(*core_database_postgres.TransactionPostgres)
	return nil
}

func (r *OrganizationPostgresRepository) applyFilters(selectQuery *bun.SelectQuery, filters organization_core.OrganizationFilters) *bun.SelectQuery {
	if filters.Name != nil {
		selectQuery = core_database_postgres.ApplyComparableFilter(selectQuery, "name", filters.Name)
	}

	if filters.Status != nil {
		selectQuery = core_database_postgres.ApplyComparableFilter(selectQuery, "status", filters.Status)
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

	if filters.LoggedUserIdentity != nil {
		selectQuery = selectQuery.WhereGroup(" OR ", func(query *bun.SelectQuery) *bun.SelectQuery {
			query = query.Where("user_creator_internal_id = ?", filters.LoggedUserIdentity.Internal.String())
			query = query.WhereOr("internal_id IN (SELECT organization_internal_id FROM organization_user WHERE user_internal_id = ?)", filters.LoggedUserIdentity.Internal.String())
			return query
		})
	}

	return selectQuery
}

func (r *OrganizationPostgresRepository) GetOrganizationByIdentity(params organization_core.GetOrganizationByIdentityParams) (*organization_core.Organization, error) {
	var organization OrganizationTable
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(&organization).Where("internal_id = ?", params.Identity.Internal.String())

	err := selectQuery.Scan(context.Background())

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return organization.ToEntity(), nil
}

func (r *OrganizationPostgresRepository) ListOrganizationsBy(params organization_core.ListOrganizationsParams) (*[]organization_core.Organization, error) {
	var organizations []OrganizationTable
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(&organizations)
	selectQuery = r.applyFilters(selectQuery, params.Filters)

	err := selectQuery.Scan(context.Background())

	if err != nil {
		if err == sql.ErrNoRows {
			return &[]organization_core.Organization{}, nil
		}
	}

	var organizationEntities []organization_core.Organization

	for _, organization := range organizations {
		organizationEntities = append(organizationEntities, *organization.ToEntity())
	}

	return &organizationEntities, nil
}

func (r *OrganizationPostgresRepository) PaginateOrganizationsBy(params organization_core.PaginateOrganizationsParams) (*core.PaginationOutput[organization_core.Organization], error) {
	var organizations []OrganizationTable
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

	selectQuery = selectQuery.Model(&organizations)
	selectQuery = r.applyFilters(selectQuery, params.Filters)

	if !params.ShowDeleted {
		selectQuery = selectQuery.Where("deleted_at IS NULL")
	}

	countBeforePagination, err := selectQuery.Count(context.Background())

	if err != nil {
		return nil, err
	}

	selectQuery = core_database_postgres.ApplyPagination(selectQuery, params.Pagination)

	err = selectQuery.Scan(context.Background(), &organizations)

	if err != nil {
		return nil, err
	}

	var organizationEntities []organization_core.Organization = []organization_core.Organization{}

	for _, organization := range organizations {
		organizationEntities = append(organizationEntities, *organization.ToEntity())
	}

	return &core.PaginationOutput[organization_core.Organization]{
		Data:    organizationEntities,
		Page:    page,
		HasMore: core.HasMorePages(page, countBeforePagination, perPage),
		Total:   countBeforePagination,
	}, nil
}

func (r *OrganizationPostgresRepository) StoreOrganization(organization *organization_core.Organization) (*organization_core.Organization, error) {
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

	var userCreatorInternalId *string
	if organization.UserCreatorIdentity != nil {
		identity := organization.UserCreatorIdentity.Internal.String()
		userCreatorInternalId = &identity
	}

	var userEditorInternalId *string
	if organization.UserEditorIdentity != nil {
		identity := organization.UserEditorIdentity.Internal.String()
		userEditorInternalId = &identity
	}

	organizationTable := &OrganizationTable{
		InternalId:            organization.Identity.Internal.String(),
		PublicId:              organization.Identity.Public,
		Name:                  organization.Name,
		Status:                string(organization.Status),
		UserCreatorInternalId: userCreatorInternalId,
		UserEditorInternalId:  userEditorInternalId,
		CreatedAt:             *organization.Timestamps.CreatedAt,
		UpdatedAt:             organization.Timestamps.UpdatedAt,
		DeletedAt:             organization.DeletedAt,
	}

	_, err := tx.NewInsert().Model(organizationTable).Exec(context.Background())

	if err != nil {
		return nil, err
	}

	if shouldCommit {
		err = tx.Commit()
	}

	return organizationTable.ToEntity(), nil
}

func (r *OrganizationPostgresRepository) UpdateOrganization(organization *organization_core.Organization) error {
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

	var userCreatorInternalId *string
	if organization.UserCreatorIdentity != nil {
		identity := organization.UserCreatorIdentity.Internal.String()
		userCreatorInternalId = &identity
	}

	var userEditorInternalId *string
	if organization.UserEditorIdentity != nil {
		identity := organization.UserEditorIdentity.Internal.String()
		userEditorInternalId = &identity
	}

	organizationTable := &OrganizationTable{
		InternalId:            organization.Identity.Internal.String(),
		PublicId:              organization.Identity.Public,
		Name:                  organization.Name,
		Status:                string(organization.Status),
		UserCreatorInternalId: userCreatorInternalId,
		UserEditorInternalId:  userEditorInternalId,
		CreatedAt:             *organization.Timestamps.CreatedAt,
		UpdatedAt:             organization.Timestamps.UpdatedAt,
		DeletedAt:             organization.DeletedAt,
	}

	_, err := tx.NewUpdate().Model(organizationTable).Where("internal_id = ?", organization.Identity.Internal.String()).Exec(context.Background())

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
	}

	if shouldCommit {
		err = tx.Commit()
	}

	return nil
}

func (r *OrganizationPostgresRepository) DeleteOrganization(organizationIdentity core.Identity) error {
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

	_, err := tx.NewDelete().Model(&OrganizationTable{}).Where("internal_id = ?", organizationIdentity.Internal.String()).Exec(context.Background())

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
	}

	if shouldCommit {
		err = tx.Commit()
	}

	return nil
}
