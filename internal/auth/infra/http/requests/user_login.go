package authhttprequests

import authservice "github.com/gabrielmrtt/taski/internal/auth/service"

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *UserLoginRequest) ToInput() authservice.UserLoginInput {
	return authservice.UserLoginInput{
		Email:    r.Email,
		Password: r.Password,
	}
}
