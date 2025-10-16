package organizationdatabase

import (
	"context"
	"database/sql"

	"github.com/gabrielmrtt/taski/internal/core"
	coredatabase "github.com/gabrielmrtt/taski/internal/core/database"
	"github.com/gabrielmrtt/taski/internal/organization"
	organizationrepo "github.com/gabrielmrtt/taski/internal/organization/repository"
	"github.com/gabrielmrtt/taski/internal/role"
	roledatabase "github.com/gabrielmrtt/taski/internal/role/infra/database"
	"github.com/gabrielmrtt/taski/internal/user"
	userdatabase "github.com/gabrielmrtt/taski/internal/user/infra/database"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type OrganizationUserTable struct {
	bun.BaseModel `bun:"table:organization_user,alias:organization_user"`

	OrganizationInternalId string `bun:"organization_internal_id,pk,notnull,type:uuid"`
	UserInternalId         string `bun:"user_internal_id,pk,notnull,type:uuid"`
	RoleInternalId         string `bun:"role_internal_id,notnull,type:uuid"`
	Status                 string `bun:"status,notnull,type:varchar(100)"`
	LastAccessAt           *int64 `bun:"last_access_at,type:bigint"`

	User *userdatabase.UserTable `bun:"rel:has-one,join:user_internal_id=internal_id"`
	Role *roledatabase.RoleTable `bun:"rel:has-one,join:role_internal_id=internal_id"`
}

func (o *OrganizationUserTable) ToEntity() *organization.OrganizationUser {
	var user *user.User = nil
	var role *role.Role = nil

	if o.User != nil {
		user = o.User.ToEntity()
	}

	if o.Role != nil {
		role = o.Role.ToEntity()
	}

	return &organization.OrganizationUser{
		OrganizationIdentity: core.NewIdentityFromInternal(uuid.MustParse(o.OrganizationInternalId), organization.OrganizationIdentityPrefix),
		User:                 *user,
		Role:                 *role,
		Status:               organization.OrganizationUserStatuses(o.Status),
		LastAccessAt:         o.LastAccessAt,
	}
}

type OrganizationUserBunRepository struct {
	db *bun.DB
	tx *coredatabase.TransactionBun
}

func NewOrganizationUserBunRepository(connection *bun.DB) *OrganizationUserBunRepository {
	return &OrganizationUserBunRepository{db: connection, tx: nil}
}

func (r *OrganizationUserBunRepository) SetTransaction(tx core.Transaction) error {
	if r.tx != nil && !r.tx.IsClosed() {
		return nil
	}

	r.tx = tx.(*coredatabase.TransactionBun)
	return nil
}

func (r *OrganizationUserBunRepository) applyFilters(selectQuery *bun.SelectQuery, filters organizationrepo.OrganizationUserFilters) *bun.SelectQuery {
	selectQuery = selectQuery.Where("organization_user.organization_internal_id = ?", filters.OrganizationIdentity.Internal.String())

	if filters.Name != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "user_credentials.name", filters.Name)
	}

	if filters.Email != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "user_credentials.email", filters.Email)
	}

	if filters.DisplayName != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "user_data.display_name", filters.DisplayName)
	}

	if filters.RolePublicId != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "role.public_id", filters.RolePublicId)
	}

	if filters.Status != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "organization_user.status", filters.Status)
	}

	return selectQuery
}

func (r *OrganizationUserBunRepository) GetLastAccessedOrganizationUserByUserIdentity(params organizationrepo.GetLastAccessedOrganizationUserByUserIdentityParams) (*organization.OrganizationUser, error) {
	var organizationUser *OrganizationUserTable = new(OrganizationUserTable)

	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(organizationUser)
	selectQuery = selectQuery.Relation("Role.RolePermissions.Permission").Relation("User.Credentials").Relation("User.Data")
	selectQuery = selectQuery.Where("organization_user.user_internal_id = ?", params.UserIdentity.Internal.String())
	selectQuery = selectQuery.Order("organization_user.last_access_at DESC")
	selectQuery = selectQuery.Limit(1)

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

func (r *OrganizationUserBunRepository) GetOrganizationUserByIdentity(params organizationrepo.GetOrganizationUserByIdentityParams) (*organization.OrganizationUser, error) {
	var organizationUser *OrganizationUserTable = new(OrganizationUserTable)
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(organizationUser)
	selectQuery = selectQuery.Relation("Role.RolePermissions.Permission").Relation("User.Credentials").Relation("User.Data")
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

func (r *OrganizationUserBunRepository) PaginateOrganizationUsersBy(params organizationrepo.PaginateOrganizationUsersParams) (*core.PaginationOutput[organization.OrganizationUser], error) {
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
	selectQuery = r.applyFilters(selectQuery, params.Filters)
	countBeforePagination, err := selectQuery.Count(context.Background())
	if err != nil {
		return nil, err
	}

	selectQuery = coredatabase.ApplyPagination(selectQuery, params.Pagination)
	err = selectQuery.Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return &core.PaginationOutput[organization.OrganizationUser]{
				Data:    []organization.OrganizationUser{},
				Page:    page,
				HasMore: false,
				Total:   0,
			}, nil
		}

		return nil, err
	}

	var organizationUserEntities []organization.OrganizationUser = make([]organization.OrganizationUser, 0)
	for _, organizationUser := range organizationUsers {
		organizationUserEntities = append(organizationUserEntities, *organizationUser.ToEntity())
	}

	return &core.PaginationOutput[organization.OrganizationUser]{
		Data:    organizationUserEntities,
		Page:    page,
		HasMore: core.HasMorePages(page, countBeforePagination, perPage),
		Total:   countBeforePagination,
	}, nil
}

func (r *OrganizationUserBunRepository) StoreOrganizationUser(params organizationrepo.StoreOrganizationUserParams) (*organization.OrganizationUser, error) {
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
		LastAccessAt:           params.OrganizationUser.LastAccessAt,
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
		if err != nil {
			return nil, err
		}
	}

	return params.OrganizationUser, nil
}

func (r *OrganizationUserBunRepository) UpdateOrganizationUser(params organizationrepo.UpdateOrganizationUserParams) error {
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
		LastAccessAt:           params.OrganizationUser.LastAccessAt,
	}

	_, err := tx.NewUpdate().Model(organizationUserTable).Where("user_internal_id = ? AND organization_internal_id = ?", params.OrganizationUser.User.Identity.Internal.String(), params.OrganizationUser.OrganizationIdentity.Internal.String()).Exec(context.Background())
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

func (r *OrganizationUserBunRepository) DeleteOrganizationUser(params organizationrepo.DeleteOrganizationUserParams) error {
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
