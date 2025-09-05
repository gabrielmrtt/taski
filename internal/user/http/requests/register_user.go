package user_http_requests

import user_services "github.com/gabrielmrtt/taski/internal/user/services"

type RegisterUserRequest struct {
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	Password    string  `json:"password"`
	PhoneNumber *string `json:"phone_number"`
}

func (r *RegisterUserRequest) ToInput() user_services.RegisterUserInput {
	return user_services.RegisterUserInput{
		Name:        r.Name,
		Email:       r.Email,
		Password:    r.Password,
		PhoneNumber: r.PhoneNumber,
	}
}
