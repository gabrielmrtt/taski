package userrepo

import (
	"github.com/gabrielmrtt/taski/internal/core"
	user "github.com/gabrielmrtt/taski/internal/user"
)

type UserFilters struct {
	Email       *core.ComparableFilter[string]
	Status      *core.ComparableFilter[user.UserStatuses]
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
	SortInput  *core.SortInput
}

type StoreUserParams struct {
	User *user.User
}

type UpdateUserParams struct {
	User *user.User
}

type DeleteUserParams struct {
	UserIdentity core.Identity
}

type UserRepository interface {
	SetTransaction(tx core.Transaction) error

	GetUserByIdentity(params GetUserByIdentityParams) (*user.User, error)
	GetUserByEmail(params GetUserByEmailParams) (*user.User, error)
	PaginateUsersBy(params PaginateUsersParams) (*core.PaginationOutput[user.User], error)

	StoreUser(params StoreUserParams) (*user.User, error)
	UpdateUser(params UpdateUserParams) error
	DeleteUser(params DeleteUserParams) error
}
