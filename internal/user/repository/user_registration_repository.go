package userrepo

import (
	"github.com/gabrielmrtt/taski/internal/core"
	user "github.com/gabrielmrtt/taski/internal/user"
)

type GetUserRegistrationByTokenParams struct {
	Token string
}

type StoreUserRegistrationParams struct {
	UserRegistration *user.UserRegistration
}

type UpdateUserRegistrationParams struct {
	UserRegistration *user.UserRegistration
}

type DeleteUserRegistrationParams struct {
	UserRegistrationIdentity core.Identity
}

type UserRegistrationRepository interface {
	SetTransaction(tx core.Transaction) error

	GetUserRegistrationByToken(params GetUserRegistrationByTokenParams) (*user.UserRegistration, error)

	StoreUserRegistration(params StoreUserRegistrationParams) (*user.UserRegistration, error)
	UpdateUserRegistration(params UpdateUserRegistrationParams) error
	DeleteUserRegistration(params DeleteUserRegistrationParams) error
}
