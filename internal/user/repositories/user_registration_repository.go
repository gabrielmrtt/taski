package user_repositories

import (
	"github.com/gabrielmrtt/taski/internal/core"
	user_core "github.com/gabrielmrtt/taski/internal/user"
)

type GetUserRegistrationByTokenParams struct {
	Token string
}

type StoreUserRegistrationParams struct {
	UserRegistration *user_core.UserRegistration
}

type UpdateUserRegistrationParams struct {
	UserRegistration *user_core.UserRegistration
}

type DeleteUserRegistrationParams struct {
	UserRegistrationIdentity core.Identity
}

type UserRegistrationRepository interface {
	SetTransaction(tx core.Transaction) error

	GetUserRegistrationByToken(params GetUserRegistrationByTokenParams) (*user_core.UserRegistration, error)

	StoreUserRegistration(params StoreUserRegistrationParams) (*user_core.UserRegistration, error)
	UpdateUserRegistration(params UpdateUserRegistrationParams) error
	DeleteUserRegistration(params DeleteUserRegistrationParams) error
}
