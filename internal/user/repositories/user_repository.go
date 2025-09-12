package user_repositories

import (
	"github.com/gabrielmrtt/taski/internal/core"
	user_core "github.com/gabrielmrtt/taski/internal/user"
)

type UserFilters struct {
	Email       *core.ComparableFilter[string]
	Status      *core.ComparableFilter[user_core.UserStatuses]
	Name        *core.ComparableFilter[string]
	DisplayName *core.ComparableFilter[string]
	CreatedAt   *core.ComparableFilter[int64]
	UpdatedAt   *core.ComparableFilter[int64]
	DeletedAt   *core.ComparableFilter[int64]
}

type GetUserByIdentityParams struct {
	UserIdentity   core.Identity
	RelationsInput core.RelationsInput
}

type GetUserByEmailParams struct {
	Email          string
	RelationsInput core.RelationsInput
}

type PaginateUsersParams struct {
	Filters        UserFilters
	Pagination     *core.PaginationInput
	SortInput      *core.SortInput
	RelationsInput core.RelationsInput
}

type StoreUserParams struct {
	User *user_core.User
}

type UpdateUserParams struct {
	User *user_core.User
}

type DeleteUserParams struct {
	UserIdentity core.Identity
}

type UserRepository interface {
	SetTransaction(tx core.Transaction) error

	GetUserByIdentity(params GetUserByIdentityParams) (*user_core.User, error)
	GetUserByEmail(params GetUserByEmailParams) (*user_core.User, error)
	PaginateUsersBy(params PaginateUsersParams) (*core.PaginationOutput[user_core.User], error)

	StoreUser(params StoreUserParams) (*user_core.User, error)
	UpdateUser(params UpdateUserParams) error
	DeleteUser(params DeleteUserParams) error
}
