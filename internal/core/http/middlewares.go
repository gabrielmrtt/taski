package corehttp

import "github.com/uptrace/bun"

type MiddlewareOptions struct {
	DbConnection *bun.DB
}
