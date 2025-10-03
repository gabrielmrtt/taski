package user

import (
	"github.com/gabrielmrtt/taski/pkg/datetimeutils"
)

type UserDto struct {
	Id          string              `json:"id"`
	Status      string              `json:"status"`
	Credentials *UserCredentialsDto `json:"credentials,omitempty"`
	Data        *UserDataDto        `json:"data,omitempty"`

	CreatedAt string  `json:"createdAt"`
	UpdatedAt *string `json:"updatedAt"`
}

type UserCredentialsDto struct {
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	PhoneNumber *string `json:"phoneNumber"`
}

type UserDataDto struct {
	DisplayName          string  `json:"displayName"`
	About                *string `json:"about"`
	ProfilePictureFileId *string `json:"profilePictureFileId"`
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
		var profilePicturePublicId *string = nil
		if user.Data.ProfilePictureIdentity != nil {
			profilePicturePublicId = &user.Data.ProfilePictureIdentity.Public
		}

		userDataDto = &UserDataDto{
			DisplayName:          user.Data.DisplayName,
			About:                user.Data.About,
			ProfilePictureFileId: profilePicturePublicId,
		}
	}

	var createdAt string = datetimeutils.EpochToRFC3339(*user.Timestamps.CreatedAt)
	var updatedAt *string = nil
	if user.Timestamps.UpdatedAt != nil {
		updatedAtString := datetimeutils.EpochToRFC3339(*user.Timestamps.UpdatedAt)
		updatedAt = &updatedAtString
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
