package auth

import "github.com/gabrielmrtt/taski/internal/user"

type UserAuthDto struct {
	Token                      string        `json:"token"`
	LastAccessedOrganizationId *string       `json:"lastAccessedOrganizationId"`
	User                       *user.UserDto `json:"user"`
}

func UserAuthToDto(usr *user.User, token string, lastAccessedOrganizationId *string) *UserAuthDto {
	return &UserAuthDto{
		Token:                      token,
		LastAccessedOrganizationId: lastAccessedOrganizationId,
		User:                       user.UserToDto(usr),
	}
}
