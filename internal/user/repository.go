package user_core

import (
	"github.com/gabrielmrtt/taski/internal/core"
)

type UserFilters struct {
	Email       *core.ComparableFilter[string]
	Status      *core.ComparableFilter[UserStatuses]
	Name        *core.ComparableFilter[string]
	DisplayName *core.ComparableFilter[string]
	CreatedAt   *core.ComparableFilter[int64]
	UpdatedAt   *core.ComparableFilter[int64]
	DeletedAt   *core.ComparableFilter[int64]
}

type GetUserByIdentityParams struct {
	UserIdentity core.Identity
}

type GetUserByEmailParams struct {
	Email string
}

type PaginateUsersParams struct {
	Filters    UserFilters
	Pagination *core.PaginationInput
}

type StoreUserParams struct {
	User *User
}

type UpdateUserParams struct {
	User *User
}

type DeleteUserParams struct {
	UserIdentity core.Identity
}

type UserRepository interface {
	SetTransaction(tx core.Transaction) error

	GetUserByIdentity(params GetUserByIdentityParams) (*User, error)
	GetUserByEmail(params GetUserByEmailParams) (*User, error)
	PaginateUsersBy(params PaginateUsersParams) (*core.PaginationOutput[User], error)

	StoreUser(params StoreUserParams) (*User, error)
	UpdateUser(params UpdateUserParams) error
	DeleteUser(params DeleteUserParams) error
}

type GetUserRegistrationByTokenParams struct {
	Token string
}

type StoreUserRegistrationParams struct {
	UserRegistration *UserRegistration
}

type UpdateUserRegistrationParams struct {
	UserRegistration *UserRegistration
}

type DeleteUserRegistrationParams struct {
	UserRegistrationIdentity core.Identity
}

type UserRegistrationRepository interface {
	SetTransaction(tx core.Transaction) error

	GetUserRegistrationByToken(params GetUserRegistrationByTokenParams) (*UserRegistration, error)

	StoreUserRegistration(params StoreUserRegistrationParams) (*UserRegistration, error)
	UpdateUserRegistration(params UpdateUserRegistrationParams) error
	DeleteUserRegistration(params DeleteUserRegistrationParams) error
}

type GetPasswordRecoveryByTokenParams struct {
	Token string
}

type StorePasswordRecoveryParams struct {
	PasswordRecovery *PasswordRecovery
}

type UpdatePasswordRecoveryParams struct {
	PasswordRecovery *PasswordRecovery
}

type DeletePasswordRecoveryParams struct {
	PasswordRecoveryIdentity core.Identity
}

type PasswordRecoveryRepository interface {
	SetTransaction(tx core.Transaction) error

	GetPasswordRecoveryByToken(params GetPasswordRecoveryByTokenParams) (*PasswordRecovery, error)

	StorePasswordRecovery(params StorePasswordRecoveryParams) (*PasswordRecovery, error)
	UpdatePasswordRecovery(params UpdatePasswordRecoveryParams) error
	DeletePasswordRecovery(params DeletePasswordRecoveryParams) error
}
