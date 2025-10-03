package coredatabase

import (
	"context"

	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/uptrace/bun"
)

type TransactionBun struct {
	Tx     *bun.Tx
	closed bool
}

func (t *TransactionBun) Commit() error {
	t.closed = true
	return t.Tx.Commit()
}

func (t *TransactionBun) Rollback() error {
	t.closed = true
	return t.Tx.Rollback()
}

func (t *TransactionBun) IsClosed() bool {
	return t.closed
}

type TransactionBunRepository struct {
	db *bun.DB
}

func NewTransactionBunRepository(connection *bun.DB) *TransactionBunRepository {
	return &TransactionBunRepository{db: connection}
}

func (r *TransactionBunRepository) BeginTransaction() (core.Transaction, error) {
	tx, err := r.db.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return &TransactionBun{Tx: &tx, closed: false}, nil
}
