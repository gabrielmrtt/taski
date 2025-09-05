package user_http_requests

import user_services "github.com/gabrielmrtt/taski/internal/user/services"

type RecoverUserPasswordRequest struct {
	Token                string `json:"token"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

func (r *RecoverUserPasswordRequest) Validate() error {
	return nil
}

func (r *RecoverUserPasswordRequest) ToInput() user_services.RecoverUserPasswordInput {
	return user_services.RecoverUserPasswordInput{
		PasswordRecoveryToken: r.Token,
		Password:              r.Password,
	}
}
