package user_database_postgres

import (
	"context"
	"database/sql"

	"github.com/gabrielmrtt/taski/internal/core"
	core_database_postgres "github.com/gabrielmrtt/taski/internal/core/database/postgres"
	user_core "github.com/gabrielmrtt/taski/internal/user"
	user_repositories "github.com/gabrielmrtt/taski/internal/user/repositories"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type PasswordRecoveryTable struct {
	bun.BaseModel `bun:"table:password_recovery"`

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
		Identity:     core.NewIdentityWithoutPublicFromInternal(uuid.MustParse(p.InternalId)),
		UserIdentity: core.NewIdentityFromInternal(uuid.MustParse(p.UserInternalId), user_core.UserIdentityPrefix),
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

func (r *PasswordRecoveryPostgresRepository) GetPasswordRecoveryByToken(params user_repositories.GetPasswordRecoveryByTokenParams) (*user_core.PasswordRecovery, error) {
	var passwordRecovery *PasswordRecoveryTable = new(PasswordRecoveryTable)
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(passwordRecovery).Where("token = ?", params.Token)
	err := selectQuery.Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	if passwordRecovery.InternalId == "" {
		return nil, nil
	}

	return passwordRecovery.ToEntity(), nil
}

func (r *PasswordRecoveryPostgresRepository) StorePasswordRecovery(params user_repositories.StorePasswordRecoveryParams) (*user_core.PasswordRecovery, error) {
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

	passwordRecoveryTable := &PasswordRecoveryTable{
		InternalId:     params.PasswordRecovery.Identity.Internal.String(),
		UserInternalId: params.PasswordRecovery.UserIdentity.Internal.String(),
		Token:          params.PasswordRecovery.Token,
		Status:         string(params.PasswordRecovery.Status),
		RecoveredAt:    params.PasswordRecovery.RecoveredAt,
		ExpiresAt:      params.PasswordRecovery.ExpiresAt,
		RequestedAt:    params.PasswordRecovery.RequestedAt,
	}

	_, err := tx.NewInsert().Model(passwordRecoveryTable).Exec(context.Background())
	if err != nil {
		return nil, err
	}

	if shouldCommit {
		err = tx.Commit()
		if err != nil {
			return nil, err
		}
	}

	return params.PasswordRecovery, nil
}

func (r *PasswordRecoveryPostgresRepository) UpdatePasswordRecovery(params user_repositories.UpdatePasswordRecoveryParams) error {
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

	passwordRecoveryTable := &PasswordRecoveryTable{
		InternalId:     params.PasswordRecovery.Identity.Internal.String(),
		UserInternalId: params.PasswordRecovery.UserIdentity.Internal.String(),
		Token:          params.PasswordRecovery.Token,
		Status:         string(params.PasswordRecovery.Status),
		RecoveredAt:    params.PasswordRecovery.RecoveredAt,
		ExpiresAt:      params.PasswordRecovery.ExpiresAt,
		RequestedAt:    params.PasswordRecovery.RequestedAt,
	}

	_, err := tx.NewUpdate().Model(passwordRecoveryTable).Where("internal_id = ?", params.PasswordRecovery.Identity.Internal.String()).Exec(context.Background())
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

func (r *PasswordRecoveryPostgresRepository) DeletePasswordRecovery(params user_repositories.DeletePasswordRecoveryParams) error {
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

	_, err := tx.NewDelete().Model(&PasswordRecoveryTable{}).Where("internal_id = ?", params.PasswordRecoveryIdentity.Internal.String()).Exec(context.Background())
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
