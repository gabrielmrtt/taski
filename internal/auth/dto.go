package auth

import "github.com/gabrielmrtt/taski/internal/user"

type UserAuthDto struct {
	Token string        `json:"token"`
	User  *user.UserDto `json:"user"`
}

func UserAuthToDto(usr *user.User, token string) *UserAuthDto {
	return &UserAuthDto{
		Token: token,
		User:  user.UserToDto(usr),
	}
}
