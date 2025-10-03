package coredatabase

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gabrielmrtt/taski/config"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/driver/sqliteshim"
)

var lock = &sync.Mutex{}

// -- Default Connection --

var postgresConnection *bun.DB = nil

func createPostgresConnection() *bun.DB {
	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(
			fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
				config.GetConfig().PostgresUsername,
				config.GetConfig().PostgresPassword,
				config.GetConfig().PostgresHost,
				config.GetConfig().PostgresPort,
				config.GetConfig().PostgresName,
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

// -- Test Connection --

var sqliteConnection *bun.DB = nil

func createSQLiteConnection() *bun.DB {
	sqldb, err := sql.Open(sqliteshim.ShimName, "file:test.db?cache=shared&mode=rwc")
	if err != nil {
		panic(fmt.Errorf("failed to open sqlite connection: %w", err))
	}

	sqldb.SetMaxOpenConns(25)
	sqldb.SetMaxIdleConns(10)
	sqldb.SetConnMaxLifetime(10 * time.Minute)
	sqldb.SetConnMaxIdleTime(10 * time.Minute)

	if err := sqldb.Ping(); err != nil {
		log.Fatal("failed to connect to sqlite: %w", err)
	}

	db := bun.NewDB(sqldb, sqlitedialect.New())
	return db
}

func GetSQLiteConnection() *bun.DB {
	if sqliteConnection == nil {
		lock.Lock()
		defer lock.Unlock()
		if sqliteConnection == nil {
			sqliteConnection = createSQLiteConnection()
			return sqliteConnection
		}
	}

	return sqliteConnection
}
