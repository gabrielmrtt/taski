package user_http_requests

import user_services "github.com/gabrielmrtt/taski/internal/user/services"

type ChangeUserPasswordRequest struct {
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
}

func (r *ChangeUserPasswordRequest) ToInput() user_services.ChangeUserPasswordInput {
	return user_services.ChangeUserPasswordInput{
		Password: r.Password,
	}
}
