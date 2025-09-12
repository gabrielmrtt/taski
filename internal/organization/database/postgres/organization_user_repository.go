package organization_database_postgres

import (
	"context"
	"database/sql"

	"github.com/gabrielmrtt/taski/internal/core"
	core_database_postgres "github.com/gabrielmrtt/taski/internal/core/database/postgres"
	organization_core "github.com/gabrielmrtt/taski/internal/organization"
	organization_repositories "github.com/gabrielmrtt/taski/internal/organization/repositories"
	role_database_postgres "github.com/gabrielmrtt/taski/internal/role/database/postgres"
	user_database_postgres "github.com/gabrielmrtt/taski/internal/user/database/postgres"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type OrganizationUserTable struct {
	bun.BaseModel `bun:"table:organization_user,alias:organization_user"`

	OrganizationInternalId string `bun:"organization_internal_id,pk,notnull,type:uuid"`
	UserInternalId         string `bun:"user_internal_id,pk,notnull,type:uuid"`
	RoleInternalId         string `bun:"role_internal_id,notnull,type:uuid"`
	Status                 string `bun:"status,notnull,type:varchar(100)"`

	User *user_database_postgres.UserTable `bun:"rel:has-one,join:user_internal_id=internal_id"`
	Role *role_database_postgres.RoleTable `bun:"rel:has-one,join:role_internal_id=internal_id"`
}

func (o *OrganizationUserTable) ToEntity() *organization_core.OrganizationUser {
	return &organization_core.OrganizationUser{
		OrganizationIdentity: core.NewIdentityFromInternal(uuid.MustParse(o.OrganizationInternalId), organization_core.OrganizationIdentityPrefix),
		User:                 o.User.ToEntity(),
		Role:                 o.Role.ToEntity(),
		Status:               organization_core.OrganizationUserStatuses(o.Status),
	}
}

type OrganizationUserPostgresRepository struct {
	db *bun.DB
	tx *core_database_postgres.TransactionPostgres
}

func NewOrganizationUserPostgresRepository() *OrganizationUserPostgresRepository {
	return &OrganizationUserPostgresRepository{db: core_database_postgres.DB, tx: nil}
}

func (r *OrganizationUserPostgresRepository) SetTransaction(tx core.Transaction) error {
	r.tx = tx.(*core_database_postgres.TransactionPostgres)
	return nil
}

func (r *OrganizationUserPostgresRepository) applyFilters(selectQuery *bun.SelectQuery, filters organization_repositories.OrganizationUserFilters) *bun.SelectQuery {
	selectQuery = selectQuery.Where("organization_user.organization_internal_id = ?", filters.OrganizationIdentity.Internal.String())

	if filters.Name != nil {
		selectQuery = core_database_postgres.ApplyComparableFilter(selectQuery, "user_credentials.name", filters.Name)
	}

	if filters.Email != nil {
		selectQuery = core_database_postgres.ApplyComparableFilter(selectQuery, "user_credentials.email", filters.Email)
	}

	if filters.DisplayName != nil {
		selectQuery = core_database_postgres.ApplyComparableFilter(selectQuery, "user_data.display_name", filters.DisplayName)
	}

	if filters.RolePublicId != nil {
		selectQuery = core_database_postgres.ApplyComparableFilter(selectQuery, "role.public_id", filters.RolePublicId)
	}

	if filters.Status != nil {
		selectQuery = core_database_postgres.ApplyComparableFilter(selectQuery, "organization_user.status", filters.Status)
	}

	return selectQuery
}

func (r *OrganizationUserPostgresRepository) GetOrganizationUserByIdentity(params organization_repositories.GetOrganizationUserByIdentityParams) (*organization_core.OrganizationUser, error) {
	var organizationUser *OrganizationUserTable = new(OrganizationUserTable)
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(organizationUser)
	selectQuery = selectQuery.Relation("Role.RolePermissions.Permission").Relation("User.Credentials").Relation("User.Data")
	selectQuery = core_database_postgres.ApplyRelations(selectQuery, params.RelationsInput)
	selectQuery = selectQuery.Where("organization_user.organization_internal_id = ? and organization_user.user_internal_id = ?", params.OrganizationIdentity.Internal.String(), params.UserIdentity.Internal.String())

	err := selectQuery.Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	if organizationUser.OrganizationInternalId == "" {
		return nil, nil
	}

	return organizationUser.ToEntity(), nil
}

func (r *OrganizationUserPostgresRepository) PaginateOrganizationUsersBy(params organization_repositories.PaginateOrganizationUsersParams) (*core.PaginationOutput[organization_core.OrganizationUser], error) {
	var organizationUsers []OrganizationUserTable = make([]OrganizationUserTable, 0)
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

	selectQuery = selectQuery.Model(&organizationUsers)
	selectQuery = selectQuery.Relation("Role.RolePermissions.Permission").Relation("User.Credentials").Relation("User.Data")
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
			return &core.PaginationOutput[organization_core.OrganizationUser]{
				Data:    []organization_core.OrganizationUser{},
				Page:    page,
				HasMore: false,
				Total:   0,
			}, nil
		}

		return nil, err
	}

	var organizationUserEntities []organization_core.OrganizationUser = make([]organization_core.OrganizationUser, 0)
	for _, organizationUser := range organizationUsers {
		organizationUserEntities = append(organizationUserEntities, *organizationUser.ToEntity())
	}

	return &core.PaginationOutput[organization_core.OrganizationUser]{
		Data:    organizationUserEntities,
		Page:    page,
		HasMore: core.HasMorePages(page, countBeforePagination, perPage),
		Total:   countBeforePagination,
	}, nil
}

func (r *OrganizationUserPostgresRepository) StoreOrganizationUser(params organization_repositories.StoreOrganizationUserParams) (*organization_core.OrganizationUser, error) {
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

	organizationUserTable := &OrganizationUserTable{
		OrganizationInternalId: params.OrganizationUser.OrganizationIdentity.Internal.String(),
		UserInternalId:         params.OrganizationUser.User.Identity.Internal.String(),
		RoleInternalId:         params.OrganizationUser.Role.Identity.Internal.String(),
		Status:                 string(params.OrganizationUser.Status),
	}

	_, err := tx.NewInsert().Model(organizationUserTable).Exec(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	if shouldCommit {
		err = tx.Commit()
	}

	return params.OrganizationUser, nil
}

func (r *OrganizationUserPostgresRepository) UpdateOrganizationUser(params organization_repositories.UpdateOrganizationUserParams) error {
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

	organizationUserTable := &OrganizationUserTable{
		OrganizationInternalId: params.OrganizationUser.OrganizationIdentity.Internal.String(),
		UserInternalId:         params.OrganizationUser.User.Identity.Internal.String(),
		RoleInternalId:         params.OrganizationUser.Role.Identity.Internal.String(),
		Status:                 string(params.OrganizationUser.Status),
	}

	_, err := tx.NewUpdate().Model(organizationUserTable).Where("internal_id = ?", params.OrganizationUser.User.Identity.Internal.String()).Exec(context.Background())
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

func (r *OrganizationUserPostgresRepository) DeleteOrganizationUser(params organization_repositories.DeleteOrganizationUserParams) error {
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

	_, err := tx.NewDelete().Model(&OrganizationUserTable{}).Where("organization_user.organization_internal_id = ? and organization_user.user_internal_id = ?", params.OrganizationIdentity.Internal.String(), params.UserIdentity.Internal.String()).Exec(context.Background())
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
