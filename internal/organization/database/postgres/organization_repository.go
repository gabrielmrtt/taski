package organization_database_postgres

import (
	"context"
	"database/sql"

	"github.com/gabrielmrtt/taski/internal/core"
	core_database_postgres "github.com/gabrielmrtt/taski/internal/core/database/postgres"
	organization_core "github.com/gabrielmrtt/taski/internal/organization"
	role_database_postgres "github.com/gabrielmrtt/taski/internal/role/database/postgres"
	user_database_postgres "github.com/gabrielmrtt/taski/internal/user/database/postgres"
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
		OrganizationIdentity: core.NewIdentityFromInternal(uuid.MustParse(o.OrganizationInternalId), "org"),
		User:                 o.User.ToEntity(),
		Role:                 o.Role.ToEntity(),
		Status:               organization_core.OrganizationUserStatuses(o.Status),
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

func (r *OrganizationPostgresRepository) applyOrganizationFilters(selectQuery *bun.SelectQuery, filters organization_core.OrganizationFilters) *bun.SelectQuery {
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

func (r *OrganizationPostgresRepository) applyOrganizationUserFilters(selectQuery *bun.SelectQuery, filters organization_core.OrganizationUserFilters) *bun.SelectQuery {
	if filters.Name != nil {
		selectQuery = core_database_postgres.ApplyComparableFilter(selectQuery, "user_credentials.name", filters.Name)
	}

	if filters.Email != nil {
		selectQuery = core_database_postgres.ApplyComparableFilter(selectQuery, "user_credentials.email", filters.Email)
	}

	if filters.DisplayName != nil {
		selectQuery = core_database_postgres.ApplyComparableFilter(selectQuery, "user_data.display_name", filters.DisplayName)
	}

	if filters.RoleInternalId != nil {
		selectQuery = core_database_postgres.ApplyComparableFilter(selectQuery, "role_internal_id", filters.RoleInternalId)
	}

	if filters.Status != nil {
		selectQuery = core_database_postgres.ApplyComparableFilter(selectQuery, "status", filters.Status)
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
	selectQuery = r.applyOrganizationFilters(selectQuery, params.Filters)

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
	selectQuery = r.applyOrganizationFilters(selectQuery, params.Filters)

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

func (r *OrganizationPostgresRepository) CreateOrganizationUser(organizationUser *organization_core.OrganizationUser) (*organization_core.OrganizationUser, error) {
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
		OrganizationInternalId: organizationUser.OrganizationIdentity.Internal.String(),
		UserInternalId:         organizationUser.User.Identity.Internal.String(),
		RoleInternalId:         organizationUser.Role.Identity.Internal.String(),
		Status:                 string(organizationUser.Status),
	}

	_, err := tx.NewInsert().Model(organizationUserTable).Exec(context.Background())

	if err != nil {
		return nil, err
	}

	if shouldCommit {
		err = tx.Commit()
	}

	return organizationUser, nil
}

func (r *OrganizationPostgresRepository) UpdateOrganizationUser(organizationUser *organization_core.OrganizationUser) error {
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
		OrganizationInternalId: organizationUser.OrganizationIdentity.Internal.String(),
		UserInternalId:         organizationUser.User.Identity.Internal.String(),
		RoleInternalId:         organizationUser.Role.Identity.Internal.String(),
		Status:                 string(organizationUser.Status),
	}

	_, err := tx.NewUpdate().Model(organizationUserTable).Where("internal_id = ?", organizationUser.User.Identity.Internal.String()).Exec(context.Background())

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

func (r *OrganizationPostgresRepository) DeleteOrganizationUser(organizationIdentity core.Identity, userIdentity core.Identity) error {
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

	_, err := tx.NewDelete().Model(&OrganizationUserTable{}).Where("organization_internal_id = ? and user_internal_id = ?", organizationIdentity.Internal.String(), userIdentity.Internal.String()).Exec(context.Background())

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

func (r *OrganizationPostgresRepository) CheckIfOrganizationHasUser(organizationIdentity core.Identity, userIdentity core.Identity) (bool, error) {
	var count int

	err := r.db.NewRaw(
		"SELECT COUNT(*) FROM organization_user WHERE organization_internal_id = ? and user_internal_id = ? AND status = ?",
		organizationIdentity.Internal.String(),
		userIdentity.Internal.String(),
		organization_core.OrganizationUserStatusActive,
	).Scan(context.Background(), &count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *OrganizationPostgresRepository) GetOrganizationUserByIdentity(organizationIdentity core.Identity, userIdentity core.Identity) (*organization_core.OrganizationUser, error) {
	var organizationUser OrganizationUserTable
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(&organizationUser).Relation("User").Relation("Role").Relation("Organization").Relation("User.UserCredentials").Relation("User.UserData").Where("organization_internal_id = ? and user_internal_id = ?", organizationIdentity.Internal.String(), userIdentity.Internal.String())

	err := selectQuery.Scan(context.Background())

	if err != nil {
		return nil, err
	}

	return organizationUser.ToEntity(), nil
}

func (r *OrganizationPostgresRepository) ListOrganizationUsersBy(params organization_core.ListOrganizationUsersParams) (*[]organization_core.OrganizationUser, error) {
	var organizationUsers []OrganizationUserTable
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(&organizationUsers).Relation("User").Relation("Role").Relation("Organization").Relation("User.UserCredentials").Relation("User.UserData")
	selectQuery = r.applyOrganizationUserFilters(selectQuery, params.Filters)

	err := selectQuery.Scan(context.Background())

	if err != nil {
		if err == sql.ErrNoRows {
			return &[]organization_core.OrganizationUser{}, nil
		}
	}

	var organizationUserEntities []organization_core.OrganizationUser

	for _, organizationUser := range organizationUsers {
		organizationUserEntities = append(organizationUserEntities, *organizationUser.ToEntity())
	}

	return &organizationUserEntities, nil
}

func (r *OrganizationPostgresRepository) PaginateOrganizationUsersBy(params organization_core.PaginateOrganizationUsersParams) (*core.PaginationOutput[organization_core.OrganizationUser], error) {
	var organizationUsers []OrganizationUserTable
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

	selectQuery = selectQuery.Model(&organizationUsers).Relation("User").Relation("Role").Relation("Organization").Relation("User.UserCredentials").Relation("User.UserData")
	selectQuery = r.applyOrganizationUserFilters(selectQuery, params.Filters)

	countBeforePagination, err := selectQuery.Count(context.Background())

	if err != nil {
		return nil, err
	}

	selectQuery = core_database_postgres.ApplyPagination(selectQuery, params.Pagination)

	err = selectQuery.Scan(context.Background(), &organizationUsers)

	if err != nil {
		return nil, err
	}

	var organizationUserEntities []organization_core.OrganizationUser

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
