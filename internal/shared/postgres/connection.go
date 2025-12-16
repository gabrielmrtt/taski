package shreadpostgres

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/gabrielmrtt/taski/config"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var lock = &sync.Mutex{}

// -- Default Connection --

var postgresConnection *bun.DB = nil

func createPostgresConnection() *bun.DB {
	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(
			fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
				config.GetInstance().PostgresUsername,
				config.GetInstance().PostgresPassword,
				config.GetInstance().PostgresHost,
				config.GetInstance().PostgresPort,
				config.GetInstance().PostgresName,
			),
		),
	))

	sqldb.SetMaxOpenConns(25)
	sqldb.SetMaxIdleConns(10)
	sqldb.SetConnMaxLifetime(10 * time.Minute)
	sqldb.SetConnMaxIdleTime(10 * time.Minute)

	if err := sqldb.Ping(); err != nil {
		panic(fmt.Errorf("failed to ping postgres: %w", err))
	}

	db := bun.NewDB(sqldb, pgdialect.New())

	return db
}

func GetPostgresConnection() *bun.DB {
	if postgresConnection == nil {
		lock.Lock()
		defer lock.Unlock()
		if postgresConnection == nil {
			postgresConnection = createPostgresConnection()
			return postgresConnection
		}
	}

	return postgresConnection
}
