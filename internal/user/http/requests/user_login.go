package user_http_requests

import user_services "github.com/gabrielmrtt/taski/internal/user/services"

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *UserLoginRequest) ToInput() user_services.UserLoginInput {
	return user_services.UserLoginInput{
		Email:    r.Email,
		Password: r.Password,
	}
}
