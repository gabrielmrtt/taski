package user_http_requests

import user_services "github.com/gabrielmrtt/taski/internal/user/services"

type ForgotUserPasswordRequest struct {
	Email string `json:"email"`
}

func (r *ForgotUserPasswordRequest) ToInput() user_services.ForgotUserPasswordInput {
	return user_services.ForgotUserPasswordInput{
		Email: r.Email,
	}
}
