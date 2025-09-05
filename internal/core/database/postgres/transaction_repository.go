package core_database_postgres

import (
	"context"

	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/uptrace/bun"
)

type TransactionPostgres struct {
	Tx     *bun.Tx
	closed bool
}

func (t *TransactionPostgres) Commit() error {
	t.closed = true
	return t.Tx.Commit()
}

func (t *TransactionPostgres) Rollback() error {
	t.closed = true
	return t.Tx.Rollback()
}

func (t *TransactionPostgres) IsClosed() bool {
	return t.closed
}

type TransactionPostgresRepository struct{}

func NewTransactionPostgresRepository() *TransactionPostgresRepository {
	return &TransactionPostgresRepository{}
}

func (r *TransactionPostgresRepository) BeginTransaction() (core.Transaction, error) {
	tx, err := DB.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return &TransactionPostgres{Tx: &tx, closed: false}, nil
}
