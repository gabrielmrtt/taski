package core

type Transaction interface {
	Commit() error
	Rollback() error
	IsClosed() bool
}

type TransactionRepository interface {
	BeginTransaction() (Transaction, error)
}
