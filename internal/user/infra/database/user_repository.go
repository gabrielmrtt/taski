package userdatabase

import (
	"context"
	"database/sql"

	"github.com/gabrielmrtt/taski/internal/core"
	coredatabase "github.com/gabrielmrtt/taski/internal/core/database"
	storage "github.com/gabrielmrtt/taski/internal/storage"
	user "github.com/gabrielmrtt/taski/internal/user"
	userrepo "github.com/gabrielmrtt/taski/internal/user/repository"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type UserTable struct {
	bun.BaseModel `bun:"table:users,alias:users"`

	InternalId string `bun:"internal_id,pk,notnull,type:uuid"`
	PublicId   string `bun:"public_id,notnull,type:varchar(510)"`
	Status     string `bun:"status,notnull,type:varchar(100)"`
	CreatedAt  int64  `bun:"created_at,notnull,type:bigint"`
	UpdatedAt  *int64 `bun:"updated_at,type:bigint"`
	DeletedAt  *int64 `bun:"deleted_at,type:bigint"`

	Credentials *UserCredentialsTable `bun:"rel:has-one,join:internal_id=user_internal_id"`
	Data        *UserDataTable        `bun:"rel:has-one,join:internal_id=user_internal_id"`
}

func (u *UserTable) ToEntity() *user.User {
	var userCredentials *user.UserCredentials
	var userData *user.UserData

	if u.Credentials != nil {
		userCredentials = &user.UserCredentials{
			Name:        u.Credentials.Name,
			Email:       u.Credentials.Email,
			Password:    u.Credentials.Password,
			PhoneNumber: u.Credentials.PhoneNumber,
		}
	}

	if u.Data != nil {
		var profilePictureIdentity *core.Identity

		if u.Data.ProfilePictureInternalId != nil {
			identity := core.NewIdentityFromInternal(uuid.MustParse(*u.Data.ProfilePictureInternalId), storage.UploadedFileIdentityPrefix)
			profilePictureIdentity = &identity
		}

		userData = &user.UserData{
			DisplayName:            u.Data.DisplayName,
			About:                  u.Data.About,
			ProfilePictureIdentity: profilePictureIdentity,
		}
	}

	return &user.User{
		Identity:    core.NewIdentityFromInternal(uuid.MustParse(u.InternalId), user.UserIdentityPrefix),
		Status:      user.UserStatuses(u.Status),
		Credentials: userCredentials,
		Data:        userData,
		Timestamps:  core.Timestamps{CreatedAt: &u.CreatedAt, UpdatedAt: u.UpdatedAt},
		DeletedAt:   u.DeletedAt,
	}
}

type UserCredentialsTable struct {
	bun.BaseModel `bun:"table:user_credentials,alias:user_credentials"`

	UserInternalId string  `bun:"user_internal_id,pk,notnull,type:uuid"`
	Name           string  `bun:"name,notnull,type:varchar(255)"`
	Email          string  `bun:"email,notnull,type:varchar(255)"`
	Password       string  `bun:"password,notnull,type:varchar(510)"`
	PhoneNumber    *string `bun:"phone_number,type:varchar(30)"`
}

type UserDataTable struct {
	bun.BaseModel `bun:"table:user_data,alias:user_data"`

	UserInternalId           string  `bun:"user_internal_id,pk,notnull,type:uuid"`
	DisplayName              string  `bun:"display_name,notnull,type:varchar(255)"`
	About                    *string `bun:"about,type:varchar(510)"`
	ProfilePictureInternalId *string `bun:"profile_picture_internal_id,type:uuid"`
}

type UserBunRepository struct {
	db *bun.DB
	tx *coredatabase.TransactionBun
}

func NewUserBunRepository(connection *bun.DB) *UserBunRepository {
	return &UserBunRepository{db: connection, tx: nil}
}

func (r *UserBunRepository) SetTransaction(tx core.Transaction) error {
	r.tx = tx.(*coredatabase.TransactionBun)
	return nil
}

func (r *UserBunRepository) applyFilters(selectQuery *bun.SelectQuery, filters userrepo.UserFilters) *bun.SelectQuery {
	if filters.Email != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "credentials.email", filters.Email)
	}

	if filters.Status != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "status", filters.Status)
	}

	if filters.Name != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "credentials.name", filters.Name)
	}

	if filters.DisplayName != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "data.display_name", filters.DisplayName)
	}

	if filters.CreatedAt != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "users.created_at", filters.CreatedAt)
	}

	if filters.UpdatedAt != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "users.updated_at", filters.UpdatedAt)
	}

	if filters.DeletedAt != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "users.deleted_at", filters.DeletedAt)
	}

	return selectQuery
}

func (r *UserBunRepository) GetUserByIdentity(params userrepo.GetUserByIdentityParams) (*user.User, error) {
	var user *UserTable = new(UserTable)
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(user)
	selectQuery = selectQuery.Relation("Credentials").Relation("Data")
	selectQuery = selectQuery.Where("users.internal_id = ?", params.UserIdentity.Internal)
	err := selectQuery.Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	if user.InternalId == "" {
		return nil, nil
	}

	return user.ToEntity(), nil
}

func (r *UserBunRepository) GetUserByEmail(params userrepo.GetUserByEmailParams) (*user.User, error) {
	var user *UserTable = new(UserTable)
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(user)
	selectQuery = selectQuery.Relation("Credentials").Relation("Data")
	selectQuery = selectQuery.Where("credentials.email = ?", params.Email)
	err := selectQuery.Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	if user.InternalId == "" {
		return nil, nil
	}

	return user.ToEntity(), nil
}

func (r *UserBunRepository) PaginateUsersBy(params userrepo.PaginateUsersParams) (*core.PaginationOutput[user.User], error) {
	var users []UserTable = make([]UserTable, 0)
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
	selectQuery = selectQuery.Relation("Credentials").Relation("Data")
	selectQuery = r.applyFilters(selectQuery, params.Filters)
	selectQuery = coredatabase.ApplySort(selectQuery, *params.SortInput)
	countBeforePagination, err := selectQuery.Count(context.Background())
	if err != nil {
		return nil, err
	}

	selectQuery = coredatabase.ApplyPagination(selectQuery, params.Pagination)
	err = selectQuery.Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return &core.PaginationOutput[user.User]{
				Data:    []user.User{},
				Page:    page,
				HasMore: false,
				Total:   0,
			}, nil
		}

		return nil, err
	}

	var userEntities []user.User = make([]user.User, 0)
	for _, user := range users {
		userEntities = append(userEntities, *user.ToEntity())
	}

	return &core.PaginationOutput[user.User]{
		Data:    userEntities,
		Page:    page,
		HasMore: core.HasMorePages(page, countBeforePagination, perPage),
		Total:   countBeforePagination,
	}, nil
}

func (r *UserBunRepository) StoreUser(params userrepo.StoreUserParams) (*user.User, error) {
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

	userTable := &UserTable{
		InternalId: params.User.Identity.Internal.String(),
		PublicId:   params.User.Identity.Public,
		Status:     string(params.User.Status),
		CreatedAt:  *params.User.Timestamps.CreatedAt,
		UpdatedAt:  params.User.Timestamps.UpdatedAt,
		DeletedAt:  params.User.DeletedAt,
	}

	_, err := tx.NewInsert().Model(userTable).Exec(context.Background())
	if err != nil {
		return nil, err
	}

	userCredentialsTable := &UserCredentialsTable{
		UserInternalId: userTable.InternalId,
		Name:           params.User.Credentials.Name,
		Email:          params.User.Credentials.Email,
		Password:       params.User.Credentials.Password,
		PhoneNumber:    params.User.Credentials.PhoneNumber,
	}

	_, err = tx.NewInsert().Model(userCredentialsTable).Exec(context.Background())
	if err != nil {
		return nil, err
	}

	var profilePictureInternalId *string
	if params.User.Data.ProfilePictureIdentity != nil {
		internalId := params.User.Data.ProfilePictureIdentity.Internal.String()
		profilePictureInternalId = &internalId
	}

	userDataTable := &UserDataTable{
		UserInternalId:           userTable.InternalId,
		DisplayName:              params.User.Data.DisplayName,
		About:                    params.User.Data.About,
		ProfilePictureInternalId: profilePictureInternalId,
	}

	_, err = tx.NewInsert().Model(userDataTable).Exec(context.Background())
	if err != nil {
		return nil, err
	}

	if shouldCommit {
		err = tx.Commit()
		if err != nil {
			return nil, err
		}
	}

	return params.User, nil
}

func (r *UserBunRepository) UpdateUser(params userrepo.UpdateUserParams) error {
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

	userTable := &UserTable{
		InternalId: params.User.Identity.Internal.String(),
		PublicId:   params.User.Identity.Public,
		Status:     string(params.User.Status),
		CreatedAt:  *params.User.Timestamps.CreatedAt,
		UpdatedAt:  params.User.Timestamps.UpdatedAt,
		DeletedAt:  params.User.DeletedAt,
	}

	_, err := tx.NewUpdate().Model(userTable).Where("internal_id = ?", params.User.Identity.Internal.String()).Exec(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}

		return err
	}

	if params.User.Credentials != nil {
		userCredentialsTable := &UserCredentialsTable{
			UserInternalId: userTable.InternalId,
			Name:           params.User.Credentials.Name,
			Email:          params.User.Credentials.Email,
			Password:       params.User.Credentials.Password,
			PhoneNumber:    params.User.Credentials.PhoneNumber,
		}

		_, err = tx.NewUpdate().Model(userCredentialsTable).Where("user_internal_id = ?", params.User.Identity.Internal.String()).Exec(context.Background())
		if err != nil {
			if err == sql.ErrNoRows {
				return nil
			}

			return err
		}
	}

	if params.User.Data != nil {
		var profilePictureInternalId *string
		if params.User.Data.ProfilePictureIdentity != nil {
			internalId := params.User.Data.ProfilePictureIdentity.Internal.String()
			profilePictureInternalId = &internalId
		}

		userDataTable := &UserDataTable{
			UserInternalId:           userTable.InternalId,
			DisplayName:              params.User.Data.DisplayName,
			About:                    params.User.Data.About,
			ProfilePictureInternalId: profilePictureInternalId,
		}

		_, err = tx.NewUpdate().Model(userDataTable).Where("user_internal_id = ?", params.User.Identity.Internal.String()).Exec(context.Background())
		if err != nil {
			if err == sql.ErrNoRows {
				return nil
			}

			return err
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

func (r *UserBunRepository) DeleteUser(params userrepo.DeleteUserParams) error {
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

	_, err := tx.NewDelete().Model(&UserTable{}).Where("internal_id = ?", params.UserIdentity.Internal.String()).Exec(context.Background())
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
