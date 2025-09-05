package core_database_postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gabrielmrtt/taski/config"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func connectToPostgres() *bun.DB {
	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(
			fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
				config.Instance.PostgresUsername,
				config.Instance.PostgresPassword,
				config.Instance.PostgresHost,
				config.Instance.PostgresPort,
				config.Instance.PostgresName,
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

var DB *bun.DB = connectToPostgres()
