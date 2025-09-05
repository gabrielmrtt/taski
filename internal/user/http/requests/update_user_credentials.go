package user_http_requests

import user_services "github.com/gabrielmrtt/taski/internal/user/services"

type UpdateUserCredentialsRequest struct {
	Name        *string `json:"name"`
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phone_number"`
}

func (r *UpdateUserCredentialsRequest) Validate() error {
	return nil
}

func (r *UpdateUserCredentialsRequest) ToInput() user_services.UpdateUserCredentialsInput {
	return user_services.UpdateUserCredentialsInput{
		Name:        r.Name,
		Email:       r.Email,
		PhoneNumber: r.PhoneNumber,
	}
}
