package user_http_requests

import user_services "github.com/gabrielmrtt/taski/internal/user/services"

type VerifyUserRegistrationRequest struct {
	Token string `json:"token"`
}

func (r *VerifyUserRegistrationRequest) Validate() error {
	return nil
}

func (r *VerifyUserRegistrationRequest) ToInput() user_services.VerifyUserRegistrationInput {
	return user_services.VerifyUserRegistrationInput{
		Token: r.Token,
	}
}
