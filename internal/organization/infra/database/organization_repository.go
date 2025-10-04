package organizationdatabase

import (
	"context"
	"database/sql"

	"github.com/gabrielmrtt/taski/internal/core"
	coredatabase "github.com/gabrielmrtt/taski/internal/core/database"
	"github.com/gabrielmrtt/taski/internal/organization"
	organizationrepo "github.com/gabrielmrtt/taski/internal/organization/repository"
	"github.com/gabrielmrtt/taski/internal/user"
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

func (o *OrganizationTable) ToEntity() *organization.Organization {
	var userCreatorIdentity *core.Identity
	var userEditorIdentity *core.Identity

	if o.UserCreatorInternalId != nil {
		identity := core.NewIdentityFromInternal(uuid.MustParse(*o.UserCreatorInternalId), user.UserIdentityPrefix)
		userCreatorIdentity = &identity
	}

	if o.UserEditorInternalId != nil {
		identity := core.NewIdentityFromInternal(uuid.MustParse(*o.UserEditorInternalId), user.UserIdentityPrefix)
		userEditorIdentity = &identity
	}

	return &organization.Organization{
		Identity:            core.NewIdentityFromInternal(uuid.MustParse(o.InternalId), organization.OrganizationIdentityPrefix),
		Name:                o.Name,
		Status:              organization.OrganizationStatuses(o.Status),
		UserCreatorIdentity: userCreatorIdentity,
		UserEditorIdentity:  userEditorIdentity,
		Timestamps: core.Timestamps{
			CreatedAt: &o.CreatedAt,
			UpdatedAt: o.UpdatedAt,
		},
		DeletedAt: o.DeletedAt,
	}
}

type OrganizationBunRepository struct {
	db *bun.DB
	tx *coredatabase.TransactionBun
}

func NewOrganizationBunRepository(connection *bun.DB) *OrganizationBunRepository {
	return &OrganizationBunRepository{db: connection, tx: nil}
}

func (r *OrganizationBunRepository) SetTransaction(tx core.Transaction) error {
	r.tx = tx.(*coredatabase.TransactionBun)
	return nil
}

func (r *OrganizationBunRepository) applyFilters(selectQuery *bun.SelectQuery, filters organizationrepo.OrganizationFilters) *bun.SelectQuery {
	if filters.Name != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "name", filters.Name)
	}

	if filters.Status != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "status", filters.Status)
	}

	if filters.CreatedAt != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "created_at", filters.CreatedAt)
	}

	if filters.UpdatedAt != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "updated_at", filters.UpdatedAt)
	}

	if filters.DeletedAt != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "deleted_at", filters.DeletedAt)
	}

	if filters.LoggedUserIdentity != nil {
		selectQuery = selectQuery.WhereGroup(" OR ", func(query *bun.SelectQuery) *bun.SelectQuery {
			query = query.Where("user_creator_internal_id = ?", filters.LoggedUserIdentity.Internal.String())
			query = query.WhereOr("internal_id IN (SELECT organization_internal_id FROM organization_user WHERE user_internal_id = ? AND organization_user.status = ?)", filters.LoggedUserIdentity.Internal.String(), organization.OrganizationUserStatusActive)
			return query
		})
	}

	return selectQuery
}

func (r *OrganizationBunRepository) GetOrganizationByIdentity(params organizationrepo.GetOrganizationByIdentityParams) (*organization.Organization, error) {
	var organization *OrganizationTable = new(OrganizationTable)
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(organization)
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

func (r *OrganizationBunRepository) PaginateOrganizationsBy(params organizationrepo.PaginateOrganizationsParams) (*core.PaginationOutput[organization.Organization], error) {
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
	selectQuery = r.applyFilters(selectQuery, params.Filters)

	if !params.ShowDeleted {
		selectQuery = selectQuery.Where("deleted_at IS NULL")
	}

	selectQuery = coredatabase.ApplySort(selectQuery, *params.SortInput)
	countBeforePagination, err := selectQuery.Count(context.Background())
	if err != nil {
		return nil, err
	}

	selectQuery = coredatabase.ApplyPagination(selectQuery, params.Pagination)
	err = selectQuery.Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return &core.PaginationOutput[organization.Organization]{
				Data:    []organization.Organization{},
				Page:    page,
				HasMore: false,
				Total:   0,
			}, nil
		}

		return nil, err
	}

	var organizationEntities []organization.Organization = make([]organization.Organization, 0)
	for _, organization := range organizations {
		organizationEntities = append(organizationEntities, *organization.ToEntity())
	}

	return &core.PaginationOutput[organization.Organization]{
		Data:    organizationEntities,
		Page:    page,
		HasMore: core.HasMorePages(page, countBeforePagination, perPage),
		Total:   countBeforePagination,
	}, nil
}

func (r *OrganizationBunRepository) PaginateInvitedOrganizationsBy(params organizationrepo.PaginateInvitedOrganizationsParams) (*core.PaginationOutput[organization.Organization], error) {
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
	selectQuery = selectQuery.Where("organization.internal_id IN (SELECT organization_user.organization_internal_id FROM organization_user WHERE organization_user.user_internal_id = ? AND organization_user.status = ?)", params.LoggedUserIdentity.Internal.String(), organization.OrganizationUserStatusInvited)
	countBeforePagination, err := selectQuery.Count(context.Background())
	if err != nil {
		return nil, err
	}

	selectQuery = coredatabase.ApplyPagination(selectQuery, params.Pagination)
	err = selectQuery.Scan(context.Background())
	if err != nil {
		return nil, err
	}

	var organizationEntities []organization.Organization = make([]organization.Organization, 0)
	for _, organization := range organizations {
		organizationEntities = append(organizationEntities, *organization.ToEntity())
	}

	return &core.PaginationOutput[organization.Organization]{
		Data:    organizationEntities,
		Page:    page,
		HasMore: core.HasMorePages(page, countBeforePagination, perPage),
		Total:   countBeforePagination,
	}, nil
}

func (r *OrganizationBunRepository) StoreOrganization(params organizationrepo.StoreOrganizationParams) (*organization.Organization, error) {
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
		if err != nil {
			return nil, err
		}
	}

	return params.Organization, nil
}

func (r *OrganizationBunRepository) UpdateOrganization(params organizationrepo.UpdateOrganizationParams) error {
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
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *OrganizationBunRepository) DeleteOrganization(params organizationrepo.DeleteOrganizationParams) error {
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
