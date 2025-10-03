package corehttp

import "github.com/gabrielmrtt/taski/internal/core"

type HttpRequest[INPUT core.ServiceInput] interface {
	ToInput() INPUT
}
