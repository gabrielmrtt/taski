package user_repositories

import (
	"github.com/gabrielmrtt/taski/internal/core"
	user_core "github.com/gabrielmrtt/taski/internal/user"
)

type GetPasswordRecoveryByTokenParams struct {
	Token string
}

type StorePasswordRecoveryParams struct {
	PasswordRecovery *user_core.PasswordRecovery
}

type UpdatePasswordRecoveryParams struct {
	PasswordRecovery *user_core.PasswordRecovery
}

type DeletePasswordRecoveryParams struct {
	PasswordRecoveryIdentity core.Identity
}

type PasswordRecoveryRepository interface {
	SetTransaction(tx core.Transaction) error

	GetPasswordRecoveryByToken(params GetPasswordRecoveryByTokenParams) (*user_core.PasswordRecovery, error)

	StorePasswordRecovery(params StorePasswordRecoveryParams) (*user_core.PasswordRecovery, error)
	UpdatePasswordRecovery(params UpdatePasswordRecoveryParams) error
	DeletePasswordRecovery(params DeletePasswordRecoveryParams) error
}
