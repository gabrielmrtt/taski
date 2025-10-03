package userhttprequests

import userservice "github.com/gabrielmrtt/taski/internal/user/service"

type UpdateUserCredentialsRequest struct {
	Name        *string `json:"name"`
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phoneNumber"`
}

func (r *UpdateUserCredentialsRequest) ToInput() userservice.UpdateUserCredentialsInput {
	return userservice.UpdateUserCredentialsInput{
		Name:        r.Name,
		Email:       r.Email,
		PhoneNumber: r.PhoneNumber,
	}
}
