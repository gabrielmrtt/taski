package user_core

import (
	"github.com/gabrielmrtt/taski/pkg/datetimeutils"
)

type UserDto struct {
	Id          string              `json:"id"`
	Status      string              `json:"status"`
	Credentials *UserCredentialsDto `json:"credentials"`
	Data        *UserDataDto        `json:"data"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UserCredentialsDto struct {
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	PhoneNumber *string `json:"phone_number"`
}

type UserDataDto struct {
	DisplayName            string  `json:"display_name"`
	About                  *string `json:"about"`
	ProfilePicturePublicId *string `json:"profile_picture"`
}

func UserToDto(user *User) *UserDto {
	var userCredentialsDto *UserCredentialsDto
	var userDataDto *UserDataDto

	if user.Credentials != nil {
		userCredentialsDto = &UserCredentialsDto{
			Name:        user.Credentials.Name,
			Email:       user.Credentials.Email,
			PhoneNumber: user.Credentials.PhoneNumber,
		}
	}

	if user.Data != nil {
		userDataDto = &UserDataDto{
			DisplayName:            user.Data.DisplayName,
			About:                  user.Data.About,
			ProfilePicturePublicId: &user.Data.ProfilePictureIdentity.Public,
		}
	}

	createdAt := datetimeutils.EpochToRFC3339(*user.Timestamps.CreatedAt)

	var updatedAt string

	if user.Timestamps.UpdatedAt != nil {
		updatedAt = datetimeutils.EpochToRFC3339(*user.Timestamps.UpdatedAt)
	}

	return &UserDto{
		Id:          user.Identity.Public,
		Status:      string(user.Status),
		Credentials: userCredentialsDto,
		Data:        userDataDto,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

type UserLoginDto struct {
	Token string   `json:"token"`
	User  *UserDto `json:"user"`
}

func UserLoginToDto(user *User, token string) *UserLoginDto {
	return &UserLoginDto{
		Token: token,
		User:  UserToDto(user),
	}
}
