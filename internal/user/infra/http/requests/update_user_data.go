package userhttprequests

import (
	"io"
	"mime/multipart"

	"github.com/gabrielmrtt/taski/internal/core"
	userservice "github.com/gabrielmrtt/taski/internal/user/service"
)

type UpdateUserDataRequest struct {
	DisplayName    *string               `json:"displayName"`
	About          *string               `json:"about"`
	ProfilePicture *multipart.FileHeader `form:"profilePicture"`
}

func (r *UpdateUserDataRequest) ToInput() userservice.UpdateUserDataInput {
	var profilePicture *core.FileInput

	if r.ProfilePicture != nil {
		file, err := r.ProfilePicture.Open()
		if err != nil {
			return userservice.UpdateUserDataInput{}
		}

		file.Seek(0, io.SeekStart)

		content, err := io.ReadAll(file)
		if err != nil {
			return userservice.UpdateUserDataInput{}
		}

		profilePicture = &core.FileInput{
			FileName:     r.ProfilePicture.Filename,
			FileContent:  content,
			FileMimeType: r.ProfilePicture.Header.Get("Content-Type"),
		}
	}
	return userservice.UpdateUserDataInput{
		DisplayName:    r.DisplayName,
		About:          r.About,
		ProfilePicture: profilePicture,
	}
}
