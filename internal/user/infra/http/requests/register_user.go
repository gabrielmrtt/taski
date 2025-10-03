package userhttprequests

import userservice "github.com/gabrielmrtt/taski/internal/user/service"

type RegisterUserRequest struct {
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	Password    string  `json:"password"`
	PhoneNumber *string `json:"phoneNumber"`
}

func (r *RegisterUserRequest) ToInput() userservice.RegisterUserInput {
	return userservice.RegisterUserInput{
		Name:        r.Name,
		Email:       r.Email,
		Password:    r.Password,
		PhoneNumber: r.PhoneNumber,
	}
}
