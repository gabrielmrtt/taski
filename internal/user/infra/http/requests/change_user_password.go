package userhttprequests

import userservice "github.com/gabrielmrtt/taski/internal/user/service"

type ChangeUserPasswordRequest struct {
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
}

func (r *ChangeUserPasswordRequest) ToInput() userservice.ChangeUserPasswordInput {
	return userservice.ChangeUserPasswordInput{
		Password: r.Password,
	}
}
