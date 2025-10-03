package userhttprequests

import userservice "github.com/gabrielmrtt/taski/internal/user/service"

type VerifyUserRegistrationRequest struct {
	Token string `json:"token"`
}

func (r *VerifyUserRegistrationRequest) ToInput() userservice.VerifyUserRegistrationInput {
	return userservice.VerifyUserRegistrationInput{
		Token: r.Token,
	}
}
