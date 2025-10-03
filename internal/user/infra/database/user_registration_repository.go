package userdatabase

import (
	"context"
	"database/sql"

	"github.com/gabrielmrtt/taski/internal/core"
	coredatabase "github.com/gabrielmrtt/taski/internal/core/database"
	user "github.com/gabrielmrtt/taski/internal/user"
	userrepo "github.com/gabrielmrtt/taski/internal/user/repository"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type UserRegistrationTable struct {
	bun.BaseModel `bun:"table:user_registration"`

	InternalId     string `bun:"internal_id,pk,notnull,type:uuid"`
	UserInternalId string `bun:"user_internal_id,notnull,type:uuid"`
	Token          string `bun:"token,notnull,type:varchar(510)"`
	Status         string `bun:"status,notnull,type:varchar(100)"`
	ExpiresAt      int64  `bun:"expires_at,notnull,type:bigint"`
	RegisteredAt   int64  `bun:"registered_at,notnull,type:bigint"`
	VerifiedAt     *int64 `bun:"verified_at,type:bigint"`
}

func (u *UserRegistrationTable) ToEntity() *user.UserRegistration {
	return &user.UserRegistration{
		Identity:     core.NewIdentityWithoutPublicFromInternal(uuid.MustParse(u.InternalId)),
		UserIdentity: core.NewIdentityFromInternal(uuid.MustParse(u.UserInternalId), user.UserIdentityPrefix),
		Token:        u.Token,
		Status:       user.UserRegistrationStatuses(u.Status),
		ExpiresAt:    u.ExpiresAt,
		RegisteredAt: u.RegisteredAt,
		VerifiedAt:   u.VerifiedAt,
	}
}

type UserRegistrationBunRepository struct {
	db *bun.DB
	tx *coredatabase.TransactionBun
}

func NewUserRegistrationBunRepository(connection *bun.DB) *UserRegistrationBunRepository {
	return &UserRegistrationBunRepository{db: connection, tx: nil}
}

func (r *UserRegistrationBunRepository) SetTransaction(tx core.Transaction) error {
	r.tx = tx.(*coredatabase.TransactionBun)
	return nil
}

func (r *UserRegistrationBunRepository) GetUserRegistrationByToken(params userrepo.GetUserRegistrationByTokenParams) (*user.UserRegistration, error) {
	var userRegistration *UserRegistrationTable = new(UserRegistrationTable)
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(userRegistration).Where("token = ?", params.Token)
	err := selectQuery.Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	if userRegistration.InternalId == "" {
		return nil, nil
	}

	return userRegistration.ToEntity(), nil
}

func (r *UserRegistrationBunRepository) StoreUserRegistration(params userrepo.StoreUserRegistrationParams) (*user.UserRegistration, error) {
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

	userRegistrationTable := &UserRegistrationTable{
		InternalId:     params.UserRegistration.Identity.Internal.String(),
		UserInternalId: params.UserRegistration.UserIdentity.Internal.String(),
		Token:          params.UserRegistration.Token,
		Status:         string(params.UserRegistration.Status),
		ExpiresAt:      params.UserRegistration.ExpiresAt,
		RegisteredAt:   params.UserRegistration.RegisteredAt,
		VerifiedAt:     params.UserRegistration.VerifiedAt,
	}

	_, err := tx.NewInsert().Model(userRegistrationTable).Exec(context.Background())
	if err != nil {
		return nil, err
	}

	if shouldCommit {
		err = tx.Commit()

		if err != nil {
			return nil, err
		}
	}

	return params.UserRegistration, nil
}

func (r *UserRegistrationBunRepository) UpdateUserRegistration(params userrepo.UpdateUserRegistrationParams) error {
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

	userRegistrationTable := &UserRegistrationTable{
		InternalId:     params.UserRegistration.Identity.Internal.String(),
		UserInternalId: params.UserRegistration.UserIdentity.Internal.String(),
		Token:          params.UserRegistration.Token,
		Status:         string(params.UserRegistration.Status),
		ExpiresAt:      params.UserRegistration.ExpiresAt,
		RegisteredAt:   params.UserRegistration.RegisteredAt,
		VerifiedAt:     params.UserRegistration.VerifiedAt,
	}

	_, err := tx.NewUpdate().Model(userRegistrationTable).Where("internal_id = ?", params.UserRegistration.Identity.Internal.String()).Exec(context.Background())
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

func (r *UserRegistrationBunRepository) DeleteUserRegistration(params userrepo.DeleteUserRegistrationParams) error {
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

	_, err := tx.NewDelete().Model(&UserRegistrationTable{}).Where("internal_id = ?", params.UserRegistrationIdentity.Internal.String()).Exec(context.Background())
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
