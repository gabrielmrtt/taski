package core_http

type HttpRequest[INPUT any] interface {
	Validate() error
	ToInput() any
}
