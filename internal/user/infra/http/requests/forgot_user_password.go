package userhttprequests

import userservice "github.com/gabrielmrtt/taski/internal/user/service"

type ForgotUserPasswordRequest struct {
	Email string `json:"email"`
}

func (r *ForgotUserPasswordRequest) ToInput() userservice.ForgotUserPasswordInput {
	return userservice.ForgotUserPasswordInput{
		Email: r.Email,
	}
}
