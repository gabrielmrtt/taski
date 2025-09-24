package core

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{Message: message}
}

type AlreadyExistsError struct {
	Message string
}

func (e *AlreadyExistsError) Error() string {
	return e.Message
}

func NewAlreadyExistsError(message string) *AlreadyExistsError {
	return &AlreadyExistsError{Message: message}
}

type InvalidInputErrorField struct {
	Field string
	Error string
}

type InvalidInputError struct {
	Message string
	Fields  []InvalidInputErrorField
}

func NewInvalidInputError(message string, fields []InvalidInputErrorField) *InvalidInputError {
	return &InvalidInputError{Message: message, Fields: fields}
}

func (e *InvalidInputError) Error() string {
	return e.Message
}

type UnauthorizedError struct {
	Message string
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}

func NewUnauthorizedError(message string) *UnauthorizedError {
	return &UnauthorizedError{Message: message}
}

type UnauthenticatedError struct {
	Message string
}

func (e *UnauthenticatedError) Error() string {
	return e.Message
}

func NewUnauthenticatedError(message string) *UnauthenticatedError {
	return &UnauthenticatedError{Message: message}
}

type InternalError struct {
	Message string
}

func (e *InternalError) Error() string {
	return e.Message
}

func NewInternalError(message string) *InternalError {
	return &InternalError{Message: message}
}

type ConflictError struct {
	Message string
}

func (e *ConflictError) Error() string {
	return e.Message
}

func NewConflictError(message string) *ConflictError {
	return &ConflictError{Message: message}
}
