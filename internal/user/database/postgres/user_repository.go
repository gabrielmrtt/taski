package user_database_postgres

import (
	"context"

	"github.com/gabrielmrtt/taski/internal/core"
	core_database_postgres "github.com/gabrielmrtt/taski/internal/core/database/postgres"
	user_core "github.com/gabrielmrtt/taski/internal/user"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type UserTable struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	InternalId string `bun:"internal_id,pk,notnull,type:uuid"`
	PublicId   string `bun:"public_id,notnull,type:varchar(510)"`
	Status     string `bun:"status,notnull,type:varchar(100)"`
	CreatedAt  int64  `bun:"created_at,notnull,type:bigint"`
	UpdatedAt  *int64 `bun:"updated_at,type:bigint"`
	DeletedAt  *int64 `bun:"deleted_at,type:bigint"`

	UserCredentials *UserCredentialsTable `bun:"rel:has-one,join:internal_id=user_internal_id"`
	UserData        *UserDataTable        `bun:"rel:has-one,join:internal_id=user_internal_id"`
}

func (u *UserTable) ToEntity() *user_core.User {
	var userCredentials *user_core.UserCredentials
	var userData *user_core.UserData

	if u.UserCredentials != nil {
		userCredentials = &user_core.UserCredentials{
			Name:        u.UserCredentials.Name,
			Email:       u.UserCredentials.Email,
			Password:    u.UserCredentials.Password,
			PhoneNumber: u.UserCredentials.PhoneNumber,
		}
	}

	if u.UserData != nil {
		var profilePictureIdentity *core.Identity

		if u.UserData.ProfilePictureInternalId != nil {
			identity := core.NewIdentityFromInternal(uuid.MustParse(*u.UserData.ProfilePictureInternalId), "file")
			profilePictureIdentity = &identity
		}

		userData = &user_core.UserData{
			DisplayName:            u.UserData.DisplayName,
			About:                  u.UserData.About,
			ProfilePictureIdentity: profilePictureIdentity,
		}
	}

	return &user_core.User{
		Identity:    core.NewIdentityFromInternal(uuid.MustParse(u.InternalId), "user"),
		Status:      user_core.UserStatuses(u.Status),
		Credentials: userCredentials,
		Data:        userData,
		Timestamps:  core.Timestamps{CreatedAt: &u.CreatedAt, UpdatedAt: u.UpdatedAt},
		DeletedAt:   u.DeletedAt,
	}
}

type UserCredentialsTable struct {
	bun.BaseModel `bun:"table:user_credentials,alias:uc"`

	UserInternalId string  `bun:"user_internal_id,pk,notnull,type:uuid"`
	Name           string  `bun:"name,notnull,type:varchar(255)"`
	Email          string  `bun:"email,notnull,type:varchar(255)"`
	Password       string  `bun:"password,notnull,type:varchar(510)"`
	PhoneNumber    *string `bun:"phone_number,type:varchar(30)"`
}

type UserDataTable struct {
	bun.BaseModel `bun:"table:user_data,alias:ud"`

	UserInternalId           string  `bun:"user_internal_id,pk,notnull,type:uuid"`
	DisplayName              string  `bun:"display_name,notnull,type:varchar(255)"`
	About                    *string `bun:"about,type:varchar(510)"`
	ProfilePictureInternalId *string `bun:"profile_picture_internal_id,type:uuid"`
}

type UserPostgresRepository struct {
	db *bun.DB
	tx *core_database_postgres.TransactionPostgres
}

func NewUserPostgresRepository() *UserPostgresRepository {
	return &UserPostgresRepository{db: core_database_postgres.DB, tx: nil}
}

func (r *UserPostgresRepository) SetTransaction(tx core.Transaction) error {
	r.tx = tx.(*core_database_postgres.TransactionPostgres)
	return nil
}

func applyFilters(selectQuery *bun.SelectQuery, filters user_core.UserFilters) *bun.SelectQuery {
	if filters.Email != nil {
		selectQuery = core_database_postgres.ApplyComparableFilter(selectQuery, "email", filters.Email)
	}

	if filters.Status != nil {
		selectQuery = core_database_postgres.ApplyComparableFilter(selectQuery, "status", filters.Status)
	}

	if filters.Name != nil {
		selectQuery = core_database_postgres.ApplyComparableFilter(selectQuery, "name", filters.Name)
	}

	if filters.DisplayName != nil {
		selectQuery = core_database_postgres.ApplyComparableFilter(selectQuery, "display_name", filters.DisplayName)
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

func (r *UserPostgresRepository) GetUserByIdentity(params user_core.GetUserByIdentityParams) (*user_core.User, error) {
	var user UserTable
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(&user).Where("internal_id = ?", params.Identity.Internal)

	err := selectQuery.Scan(context.Background())

	if err != nil {
		return nil, err
	}

	return user.ToEntity(), nil
}

func (r *UserPostgresRepository) GetUserByEmail(params user_core.GetUserByEmailParams) (*user_core.User, error) {
	var user UserTable
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(&user).Join("JOIN user_credentials uc ON users.internal_id = uc.user_internal_id").Where("uc.email = ?", params.Email)

	err := selectQuery.Scan(context.Background())

	if err != nil {
		return nil, err
	}

	return user.ToEntity(), nil
}

func (r *UserPostgresRepository) ListUsersBy(params user_core.ListUsersParams) (*[]user_core.User, error) {
	var users []UserTable
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(&users)
	selectQuery = applyFilters(selectQuery, params.Filters)

	err := selectQuery.Scan(context.Background())

	if err != nil {
		return nil, err
	}

	var userEntities []user_core.User

	for _, user := range users {
		userEntities = append(userEntities, *user.ToEntity())
	}

	return &userEntities, nil
}

func (r *UserPostgresRepository) PaginateUsersBy(params user_core.PaginateUsersParams) (*core.PaginationOutput[user_core.User], error) {
	var users []UserTable
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

	selectQuery = selectQuery.Model(&users)
	selectQuery = applyFilters(selectQuery, params.Filters)

	countBeforePagination, err := selectQuery.Count(context.Background())

	if err != nil {
		return nil, err
	}

	selectQuery = core_database_postgres.ApplyPagination(selectQuery, params.Pagination)

	err = selectQuery.Scan(context.Background(), &users)

	if err != nil {
		return nil, err
	}

	var userEntities []user_core.User

	for _, user := range users {
		userEntities = append(userEntities, *user.ToEntity())
	}

	return &core.PaginationOutput[user_core.User]{
		Data:    userEntities,
		Page:    page,
		HasMore: core.HasMorePages(page, countBeforePagination, perPage),
		Total:   countBeforePagination,
	}, nil
}

func (r *UserPostgresRepository) StoreUser(user *user_core.User) (*user_core.User, error) {
	var tx bun.Tx

	if r.tx != nil && !r.tx.IsClosed() {
		tx = *r.tx.Tx
	} else {
		var err error
		tx, err = r.db.BeginTx(context.Background(), nil)

		if err != nil {
			return nil, err
		}
	}

	userTable := &UserTable{
		InternalId: user.Identity.Internal.String(),
		PublicId:   user.Identity.Public,
		Status:     string(user.Status),
		CreatedAt:  *user.Timestamps.CreatedAt,
		UpdatedAt:  user.Timestamps.UpdatedAt,
		DeletedAt:  user.DeletedAt,
	}

	_, err := tx.NewInsert().Model(userTable).Exec(context.Background())

	if err != nil {
		return nil, err
	}

	userCredentialsTable := &UserCredentialsTable{
		UserInternalId: userTable.InternalId,
		Name:           user.Credentials.Name,
		Email:          user.Credentials.Email,
		Password:       user.Credentials.Password,
		PhoneNumber:    user.Credentials.PhoneNumber,
	}

	_, err = tx.NewInsert().Model(userCredentialsTable).Exec(context.Background())

	if err != nil {
		return nil, err
	}

	var profilePictureInternalId *string

	if user.Data.ProfilePictureIdentity != nil {
		internalId := user.Data.ProfilePictureIdentity.Internal.String()
		profilePictureInternalId = &internalId
	}

	userDataTable := &UserDataTable{
		UserInternalId:           userTable.InternalId,
		DisplayName:              user.Data.DisplayName,
		About:                    user.Data.About,
		ProfilePictureInternalId: profilePictureInternalId,
	}

	_, err = tx.NewInsert().Model(userDataTable).Exec(context.Background())

	if err != nil {
		return nil, err
	}

	return userTable.ToEntity(), nil
}

func (r *UserPostgresRepository) UpdateUser(user *user_core.User) error {
	var tx bun.Tx

	if r.tx != nil && !r.tx.IsClosed() {
		tx = *r.tx.Tx
	} else {
		var err error
		tx, err = r.db.BeginTx(context.Background(), nil)

		if err != nil {
			return err
		}
	}

	userTable := &UserTable{
		InternalId: user.Identity.Internal.String(),
		PublicId:   user.Identity.Public,
		Status:     string(user.Status),
		UpdatedAt:  user.Timestamps.UpdatedAt,
		DeletedAt:  user.DeletedAt,
	}

	_, err := tx.NewUpdate().Model(userTable).Where("internal_id = ?", user.Identity.Internal.String()).Exec(context.Background())

	if err != nil {
		return err
	}

	if user.Credentials != nil {
		userCredentialsTable := &UserCredentialsTable{
			UserInternalId: userTable.InternalId,
			Name:           user.Credentials.Name,
			Email:          user.Credentials.Email,
			Password:       user.Credentials.Password,
			PhoneNumber:    user.Credentials.PhoneNumber,
		}

		_, err = tx.NewUpdate().Model(userCredentialsTable).Where("user_internal_id = ?", user.Identity.Internal.String()).Exec(context.Background())

		if err != nil {
			return err
		}
	}

	if user.Data != nil {
		var profilePictureInternalId *string

		if user.Data.ProfilePictureIdentity != nil {
			internalId := user.Data.ProfilePictureIdentity.Internal.String()
			profilePictureInternalId = &internalId
		}

		userDataTable := &UserDataTable{
			UserInternalId:           userTable.InternalId,
			DisplayName:              user.Data.DisplayName,
			About:                    user.Data.About,
			ProfilePictureInternalId: profilePictureInternalId,
		}

		_, err = tx.NewUpdate().Model(userDataTable).Where("user_internal_id = ?", user.Identity.Internal.String()).Exec(context.Background())

		if err != nil {
			return err
		}
	}

	return nil
}

func (r *UserPostgresRepository) DeleteUser(userIdentity core.Identity) error {
	var tx bun.Tx

	if r.tx != nil && !r.tx.IsClosed() {
		tx = *r.tx.Tx
	} else {
		var err error
		tx, err = r.db.BeginTx(context.Background(), nil)

		if err != nil {
			return err
		}
	}

	_, err := tx.NewDelete().Model(&UserTable{}).Where("internal_id = ?", userIdentity.Internal.String()).Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}
