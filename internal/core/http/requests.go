package core_http

import "github.com/gabrielmrtt/taski/internal/core"

type HttpRequest[INPUT core.ServiceInput] interface {
	ToInput() INPUT
}
