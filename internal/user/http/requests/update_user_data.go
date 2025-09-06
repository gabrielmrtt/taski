package user_http_requests

import (
	"io"
	"mime/multipart"

	"github.com/gabrielmrtt/taski/internal/core"
	user_services "github.com/gabrielmrtt/taski/internal/user/services"
)

type UpdateUserDataRequest struct {
	DisplayName    *string               `json:"display_name"`
	About          *string               `json:"about"`
	ProfilePicture *multipart.FileHeader `form:"profile_picture"`
}

func (r *UpdateUserDataRequest) ToInput() user_services.UpdateUserDataInput {
	var profilePicture *core.FileInput

	if r.ProfilePicture != nil {
		file, err := r.ProfilePicture.Open()
		if err != nil {
			return user_services.UpdateUserDataInput{}
		}

		file.Seek(0, io.SeekStart)

		content, err := io.ReadAll(file)
		if err != nil {
			return user_services.UpdateUserDataInput{}
		}

		profilePicture = &core.FileInput{
			FileName:     r.ProfilePicture.Filename,
			FileContent:  content,
			FileMimeType: r.ProfilePicture.Header.Get("Content-Type"),
		}
	}
	return user_services.UpdateUserDataInput{
		DisplayName:    r.DisplayName,
		About:          r.About,
		ProfilePicture: profilePicture,
	}
}
