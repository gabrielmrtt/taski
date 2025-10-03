package userhttprequests

import userservice "github.com/gabrielmrtt/taski/internal/user/service"

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *UserLoginRequest) ToInput() userservice.UserLoginInput {
	return userservice.UserLoginInput{
		Email:    r.Email,
		Password: r.Password,
	}
}
