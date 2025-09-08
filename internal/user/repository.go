package user_core

import (
	"github.com/gabrielmrtt/taski/internal/core"
)

type GetUserByIdentityParams struct {
	Identity core.Identity
	Include  map[string]any
}

type GetUserByEmailParams struct {
	Email   string
	Include map[string]any
}

type UserFilters struct {
	Email       *core.ComparableFilter[string]
	Status      *core.ComparableFilter[UserStatuses]
	Name        *core.ComparableFilter[string]
	DisplayName *core.ComparableFilter[string]
	CreatedAt   *core.ComparableFilter[int64]
	UpdatedAt   *core.ComparableFilter[int64]
	DeletedAt   *core.ComparableFilter[int64]
}

type ListUsersParams struct {
	Filters UserFilters
	Include map[string]any
}

type PaginateUsersParams struct {
	Filters    UserFilters
	Include    map[string]any
	Pagination *core.PaginationInput
}

type UserRepository interface {
	SetTransaction(tx core.Transaction) error

	GetUserByIdentity(params GetUserByIdentityParams) (*User, error)
	GetUserByEmail(params GetUserByEmailParams) (*User, error)
	ListUsersBy(params ListUsersParams) (*[]User, error)
	PaginateUsersBy(params PaginateUsersParams) (*core.PaginationOutput[User], error)

	StoreUser(user *User) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(userIdentity core.Identity) error
}

type GetUserRegistrationByTokenParams struct {
	Token   string
	Include map[string]any
}

type UserRegistrationRepository interface {
	SetTransaction(tx core.Transaction) error

	GetUserRegistrationByToken(params GetUserRegistrationByTokenParams) (*UserRegistration, error)

	StoreUserRegistration(userRegistration *UserRegistration) (*UserRegistration, error)
	UpdateUserRegistration(userRegistration *UserRegistration) error
	DeleteUserRegistration(userRegistrationIdentity core.Identity) error
}

type GetPasswordRecoveryByTokenParams struct {
	Token   string
	Include map[string]any
}

type PasswordRecoveryRepository interface {
	SetTransaction(tx core.Transaction) error

	GetPasswordRecoveryByToken(params GetPasswordRecoveryByTokenParams) (*PasswordRecovery, error)

	StorePasswordRecovery(passwordRecovery *PasswordRecovery) (*PasswordRecovery, error)
	UpdatePasswordRecovery(passwordRecovery *PasswordRecovery) error
	DeletePasswordRecovery(passwordRecoveryIdentity core.Identity) error
}
