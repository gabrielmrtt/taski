package userrepo

import (
	"github.com/gabrielmrtt/taski/internal/core"
	user "github.com/gabrielmrtt/taski/internal/user"
)

type GetPasswordRecoveryByTokenParams struct {
	Token string
}

type StorePasswordRecoveryParams struct {
	PasswordRecovery *user.PasswordRecovery
}

type UpdatePasswordRecoveryParams struct {
	PasswordRecovery *user.PasswordRecovery
}

type DeletePasswordRecoveryParams struct {
	PasswordRecoveryIdentity core.Identity
}

type PasswordRecoveryRepository interface {
	SetTransaction(tx core.Transaction) error

	GetPasswordRecoveryByToken(params GetPasswordRecoveryByTokenParams) (*user.PasswordRecovery, error)

	StorePasswordRecovery(params StorePasswordRecoveryParams) (*user.PasswordRecovery, error)
	UpdatePasswordRecovery(params UpdatePasswordRecoveryParams) error
	DeletePasswordRecovery(params DeletePasswordRecoveryParams) error
}
