package userhttprequests

import userservice "github.com/gabrielmrtt/taski/internal/user/service"

type RecoverUserPasswordRequest struct {
	Token                string `json:"token"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
}

func (r *RecoverUserPasswordRequest) ToInput() userservice.RecoverUserPasswordInput {
	return userservice.RecoverUserPasswordInput{
		PasswordRecoveryToken: r.Token,
		Password:              r.Password,
	}
}
