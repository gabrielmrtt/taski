package user_http_requests

import user_services "github.com/gabrielmrtt/taski/internal/user/services"

type ChangeUserPasswordRequest struct {
	Password             string `json:"password" validate:"required,strong_password"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
}

func (r *ChangeUserPasswordRequest) Validate() error {
	return nil
}

func (r *ChangeUserPasswordRequest) ToInput() user_services.ChangeUserPasswordInput {
	return user_services.ChangeUserPasswordInput{
		Password: r.Password,
	}
}
