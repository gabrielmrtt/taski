package user_database_postgres

import (
	"context"

	"github.com/gabrielmrtt/taski/internal/core"
	core_database_postgres "github.com/gabrielmrtt/taski/internal/core/database/postgres"
	user_core "github.com/gabrielmrtt/taski/internal/user"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type PasswordRecoveryTable struct {
	bun.BaseModel `bun:"table:password_recovery,alias:pr"`

	InternalId     string `bun:"internal_id,pk,notnull,type:uuid"`
	UserInternalId string `bun:"user_internal_id,notnull,type:uuid"`
	Token          string `bun:"token,notnull,type:varchar(510)"`
	Status         string `bun:"status,notnull,type:varchar(100)"`
	RecoveredAt    *int64 `bun:"recovered_at,type:bigint"`
	ExpiresAt      int64  `bun:"expires_at,notnull,type:bigint"`
	RequestedAt    int64  `bun:"requested_at,notnull,type:bigint"`
}

func (p *PasswordRecoveryTable) ToEntity() *user_core.PasswordRecovery {
	return &user_core.PasswordRecovery{
		Identity:     core.NewIdentityFromInternal(uuid.MustParse(p.InternalId), "password_recovery"),
		UserIdentity: core.NewIdentityFromInternal(uuid.MustParse(p.UserInternalId), "user"),
		Token:        p.Token,
		Status:       user_core.PasswordRecoveryStatuses(p.Status),
		RecoveredAt:  p.RecoveredAt,
		ExpiresAt:    p.ExpiresAt,
		RequestedAt:  p.RequestedAt,
	}
}

type PasswordRecoveryPostgresRepository struct {
	db *bun.DB
	tx *core_database_postgres.TransactionPostgres
}

func NewPasswordRecoveryPostgresRepository() *PasswordRecoveryPostgresRepository {
	return &PasswordRecoveryPostgresRepository{db: core_database_postgres.DB, tx: nil}
}

func (r *PasswordRecoveryPostgresRepository) SetTransaction(tx core.Transaction) error {
	r.tx = tx.(*core_database_postgres.TransactionPostgres)
	return nil
}

func (r *PasswordRecoveryPostgresRepository) GetPasswordRecoveryByToken(params user_core.GetPasswordRecoveryByTokenParams) (*user_core.PasswordRecovery, error) {
	var passwordRecovery PasswordRecoveryTable
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(&passwordRecovery).Where("token = ?", params.Token)

	err := selectQuery.Scan(context.Background())

	if err != nil {
		return nil, err
	}

	return passwordRecovery.ToEntity(), nil
}

func (r *PasswordRecoveryPostgresRepository) StorePasswordRecovery(passwordRecovery *user_core.PasswordRecovery) (*user_core.PasswordRecovery, error) {
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

	passwordRecoveryTable := &PasswordRecoveryTable{
		InternalId:     passwordRecovery.Identity.Internal.String(),
		UserInternalId: passwordRecovery.UserIdentity.Internal.String(),
		Token:          passwordRecovery.Token,
		Status:         string(passwordRecovery.Status),
		RecoveredAt:    passwordRecovery.RecoveredAt,
		ExpiresAt:      passwordRecovery.ExpiresAt,
		RequestedAt:    passwordRecovery.RequestedAt,
	}

	_, err := tx.NewInsert().Model(passwordRecoveryTable).Exec(context.Background())

	if err != nil {
		return nil, err
	}

	return passwordRecoveryTable.ToEntity(), nil
}

func (r *PasswordRecoveryPostgresRepository) UpdatePasswordRecovery(passwordRecovery *user_core.PasswordRecovery) error {
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

	passwordRecoveryTable := &PasswordRecoveryTable{
		InternalId:     passwordRecovery.Identity.Internal.String(),
		UserInternalId: passwordRecovery.UserIdentity.Internal.String(),
		Token:          passwordRecovery.Token,
		Status:         string(passwordRecovery.Status),
		RecoveredAt:    passwordRecovery.RecoveredAt,
		ExpiresAt:      passwordRecovery.ExpiresAt,
		RequestedAt:    passwordRecovery.RequestedAt,
	}

	_, err := tx.NewUpdate().Model(passwordRecoveryTable).Where("internal_id = ?", passwordRecovery.Identity.Internal.String()).Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func (r *PasswordRecoveryPostgresRepository) DeletePasswordRecovery(passwordRecoveryIdentity core.Identity) error {
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

	_, err := tx.NewDelete().Model(&PasswordRecoveryTable{}).Where("internal_id = ?", passwordRecoveryIdentity.Internal.String()).Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}
