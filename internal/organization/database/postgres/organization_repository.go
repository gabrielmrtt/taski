package organization_database_postgres

import (
	"context"
	"database/sql"

	"github.com/gabrielmrtt/taski/internal/core"
	core_database_postgres "github.com/gabrielmrtt/taski/internal/core/database/postgres"
	organization_core "github.com/gabrielmrtt/taski/internal/organization"
	organization_repositories "github.com/gabrielmrtt/taski/internal/organization/repositories"
	user_core "github.com/gabrielmrtt/taski/internal/user"
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
		identity := core.NewIdentityFromInternal(uuid.MustParse(*o.UserCreatorInternalId), user_core.UserIdentityPrefix)
		userCreatorIdentity = &identity
	}

	if o.UserEditorInternalId != nil {
		identity := core.NewIdentityFromInternal(uuid.MustParse(*o.UserEditorInternalId), user_core.UserIdentityPrefix)
		userEditorIdentity = &identity
	}

	return &organization_core.Organization{
		Identity:            core.NewIdentityFromInternal(uuid.MustParse(o.InternalId), organization_core.OrganizationIdentityPrefix),
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

func (r *OrganizationPostgresRepository) applyFilters(selectQuery *bun.SelectQuery, filters organization_repositories.OrganizationFilters) *bun.SelectQuery {
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
			query = query.WhereOr("internal_id IN (SELECT organization_internal_id FROM organization_user WHERE user_internal_id = ? AND organization_user.status = ?)", filters.LoggedUserIdentity.Internal.String(), organization_core.OrganizationUserStatusActive)
			return query
		})
	}

	return selectQuery
}

func (r *OrganizationPostgresRepository) GetOrganizationByIdentity(params organization_repositories.GetOrganizationByIdentityParams) (*organization_core.Organization, error) {
	var organization *OrganizationTable = new(OrganizationTable)
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(organization)
	selectQuery = core_database_postgres.ApplyRelations(selectQuery, params.RelationsInput)
	selectQuery = selectQuery.Where("internal_id = ?", params.OrganizationIdentity.Internal.String())
	err := selectQuery.Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	if organization.InternalId == "" {
		return nil, nil
	}

	return organization.ToEntity(), nil
}

func (r *OrganizationPostgresRepository) PaginateOrganizationsBy(params organization_repositories.PaginateOrganizationsParams) (*core.PaginationOutput[organization_core.Organization], error) {
	var organizations []OrganizationTable = make([]OrganizationTable, 0)
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
	selectQuery = core_database_postgres.ApplyRelations(selectQuery, params.RelationsInput)
	selectQuery = r.applyFilters(selectQuery, params.Filters)

	if !params.ShowDeleted {
		selectQuery = selectQuery.Where("deleted_at IS NULL")
	}

	selectQuery = core_database_postgres.ApplySort(selectQuery, *params.SortInput)
	countBeforePagination, err := selectQuery.Count(context.Background())
	if err != nil {
		return nil, err
	}

	selectQuery = core_database_postgres.ApplyPagination(selectQuery, params.Pagination)
	err = selectQuery.Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return &core.PaginationOutput[organization_core.Organization]{
				Data:    []organization_core.Organization{},
				Page:    page,
				HasMore: false,
				Total:   0,
			}, nil
		}

		return nil, err
	}

	var organizationEntities []organization_core.Organization = make([]organization_core.Organization, 0)
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

func (r *OrganizationPostgresRepository) PaginateInvitedOrganizationsBy(params organization_repositories.PaginateInvitedOrganizationsParams) (*core.PaginationOutput[organization_core.Organization], error) {
	var organizations []OrganizationTable = make([]OrganizationTable, 0)
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
	selectQuery = core_database_postgres.ApplyRelations(selectQuery, params.RelationsInput)
	selectQuery = selectQuery.Where("organization.internal_id IN (SELECT organization_user.organization_internal_id FROM organization_user WHERE organization_user.user_internal_id = ? AND organization_user.status = ?)", params.LoggedUserIdentity.Internal.String(), organization_core.OrganizationUserStatusInvited)
	countBeforePagination, err := selectQuery.Count(context.Background())
	if err != nil {
		return nil, err
	}

	selectQuery = core_database_postgres.ApplyPagination(selectQuery, params.Pagination)
	err = selectQuery.Scan(context.Background())
	if err != nil {
		return nil, err
	}

	var organizationEntities []organization_core.Organization = make([]organization_core.Organization, 0)
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

func (r *OrganizationPostgresRepository) StoreOrganization(params organization_repositories.StoreOrganizationParams) (*organization_core.Organization, error) {
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
	if params.Organization.UserCreatorIdentity != nil {
		identity := params.Organization.UserCreatorIdentity.Internal.String()
		userCreatorInternalId = &identity
	}

	var userEditorInternalId *string
	if params.Organization.UserEditorIdentity != nil {
		identity := params.Organization.UserEditorIdentity.Internal.String()
		userEditorInternalId = &identity
	}

	organizationTable := &OrganizationTable{
		InternalId:            params.Organization.Identity.Internal.String(),
		PublicId:              params.Organization.Identity.Public,
		Name:                  params.Organization.Name,
		Status:                string(params.Organization.Status),
		UserCreatorInternalId: userCreatorInternalId,
		UserEditorInternalId:  userEditorInternalId,
		CreatedAt:             *params.Organization.Timestamps.CreatedAt,
		UpdatedAt:             params.Organization.Timestamps.UpdatedAt,
		DeletedAt:             params.Organization.DeletedAt,
	}

	_, err := tx.NewInsert().Model(organizationTable).Exec(context.Background())
	if err != nil {
		return nil, err
	}

	if shouldCommit {
		err = tx.Commit()
	}

	return params.Organization, nil
}

func (r *OrganizationPostgresRepository) UpdateOrganization(params organization_repositories.UpdateOrganizationParams) error {
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
	if params.Organization.UserCreatorIdentity != nil {
		identity := params.Organization.UserCreatorIdentity.Internal.String()
		userCreatorInternalId = &identity
	}

	var userEditorInternalId *string
	if params.Organization.UserEditorIdentity != nil {
		identity := params.Organization.UserEditorIdentity.Internal.String()
		userEditorInternalId = &identity
	}

	organizationTable := &OrganizationTable{
		InternalId:            params.Organization.Identity.Internal.String(),
		PublicId:              params.Organization.Identity.Public,
		Name:                  params.Organization.Name,
		Status:                string(params.Organization.Status),
		UserCreatorInternalId: userCreatorInternalId,
		UserEditorInternalId:  userEditorInternalId,
		CreatedAt:             *params.Organization.Timestamps.CreatedAt,
		UpdatedAt:             params.Organization.Timestamps.UpdatedAt,
		DeletedAt:             params.Organization.DeletedAt,
	}

	_, err := tx.NewUpdate().Model(organizationTable).Where("internal_id = ?", params.Organization.Identity.Internal.String()).Exec(context.Background())
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

func (r *OrganizationPostgresRepository) DeleteOrganization(params organization_repositories.DeleteOrganizationParams) error {
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

	_, err := tx.NewDelete().Model(&OrganizationTable{}).Where("internal_id = ?", params.OrganizationIdentity.Internal.String()).Exec(context.Background())
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
