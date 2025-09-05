package user_database_postgres

import (
	"context"
	"database/sql"

	"github.com/gabrielmrtt/taski/internal/core"
	core_database_postgres "github.com/gabrielmrtt/taski/internal/core/database/postgres"
	user_core "github.com/gabrielmrtt/taski/internal/user"
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

func (u *UserRegistrationTable) ToEntity() *user_core.UserRegistration {
	return &user_core.UserRegistration{
		Identity:     core.NewIdentityFromInternal(uuid.MustParse(u.InternalId), "user_registration"),
		UserIdentity: core.NewIdentityFromInternal(uuid.MustParse(u.UserInternalId), "user"),
		Token:        u.Token,
		Status:       user_core.UserRegistrationStatuses(u.Status),
		ExpiresAt:    u.ExpiresAt,
		RegisteredAt: u.RegisteredAt,
		VerifiedAt:   u.VerifiedAt,
	}
}

type UserRegistrationPostgresRepository struct {
	db *bun.DB
	tx *core_database_postgres.TransactionPostgres
}

func NewUserRegistrationPostgresRepository() *UserRegistrationPostgresRepository {
	return &UserRegistrationPostgresRepository{db: core_database_postgres.DB, tx: nil}
}

func (r *UserRegistrationPostgresRepository) SetTransaction(tx core.Transaction) error {
	r.tx = tx.(*core_database_postgres.TransactionPostgres)
	return nil
}

func (r *UserRegistrationPostgresRepository) GetUserRegistrationByToken(params user_core.GetUserRegistrationByTokenParams) (*user_core.UserRegistration, error) {
	var userRegistration UserRegistrationTable
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(&userRegistration).Where("token = ?", params.Token)

	err := selectQuery.Scan(context.Background())

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return userRegistration.ToEntity(), nil
}

func (r *UserRegistrationPostgresRepository) StoreUserRegistration(userRegistration *user_core.UserRegistration) (*user_core.UserRegistration, error) {
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
		InternalId:     userRegistration.Identity.Internal.String(),
		UserInternalId: userRegistration.UserIdentity.Internal.String(),
		Token:          userRegistration.Token,
		Status:         string(userRegistration.Status),
		ExpiresAt:      userRegistration.ExpiresAt,
		RegisteredAt:   userRegistration.RegisteredAt,
		VerifiedAt:     userRegistration.VerifiedAt,
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

	return userRegistrationTable.ToEntity(), nil
}

func (r *UserRegistrationPostgresRepository) UpdateUserRegistration(userRegistration *user_core.UserRegistration) error {
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
		InternalId:     userRegistration.Identity.Internal.String(),
		UserInternalId: userRegistration.UserIdentity.Internal.String(),
		Token:          userRegistration.Token,
		Status:         string(userRegistration.Status),
		ExpiresAt:      userRegistration.ExpiresAt,
		RegisteredAt:   userRegistration.RegisteredAt,
		VerifiedAt:     userRegistration.VerifiedAt,
	}

	_, err := tx.NewUpdate().Model(userRegistrationTable).Where("internal_id = ?", userRegistration.Identity.Internal.String()).Exec(context.Background())

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

func (r *UserRegistrationPostgresRepository) DeleteUserRegistration(userRegistrationIdentity core.Identity) error {
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

	_, err := tx.NewDelete().Model(&UserRegistrationTable{}).Where("internal_id = ?", userRegistrationIdentity.Internal.String()).Exec(context.Background())

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
