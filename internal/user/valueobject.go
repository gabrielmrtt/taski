package user

import (
	"errors"
	"regexp"

	"github.com/go-playground/validator/v10"
)

type Email struct {
	Value string
}

func NewEmail(value string) (Email, error) {
	e := Email{Value: value}

	if err := e.Validate(); err != nil {
		return Email{}, err
	}

	return e, nil
}

func (e Email) Validate() error {
	if e.Value == "" {
		return errors.New("email cannot be empty")
	}

	reg := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)

	if !reg.MatchString(e.Value) {
		return errors.New("invalid email format")
	}

	return nil
}

func (e Email) String() string {
	return e.Value
}

func (e Email) Equals(_e Email) bool {
	return e.Value == _e.Value
}

func IsValidEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()

	_, err := NewEmail(email)

	return err == nil
}

type Password struct {
	Value string
}

func NewPassword(value string) (Password, error) {
	p := Password{Value: value}

	if err := p.Validate(); err != nil {
		return Password{}, err
	}

	return p, nil
}

func (p Password) Validate() error {
	if len(p.Value) < 12 || len(p.Value) > 64 {
		return errors.New("password must be between 12 and 64 characters")
	}

	hasUppercase, _ := regexp.MatchString(`[A-Z]`, p.Value)
	hasLowercase, _ := regexp.MatchString(`[a-z]`, p.Value)
	hasNumber, _ := regexp.MatchString(`[0-9]`, p.Value)
	hasSpecial, _ := regexp.MatchString(`[^A-Za-z0-9]`, p.Value)

	if !hasUppercase || !hasLowercase || !hasNumber || !hasSpecial {
		return errors.New("password must contain at least one uppercase letter, one lowercase letter, one number, and one special character")
	}

	return nil
}

func (p Password) String() string {
	return p.Value
}

func (p Password) Equals(_p Password) bool {
	return p.Value == _p.Value
}

func IsValidPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	_, err := NewPassword(password)

	return err == nil
}

type PhoneNumber struct {
	Value string
}

func NewPhoneNumber(value string) (PhoneNumber, error) {
	p := PhoneNumber{Value: value}

	if err := p.Validate(); err != nil {
		return PhoneNumber{}, err
	}

	return p, nil
}

func (p PhoneNumber) Validate() error {
	reg := regexp.MustCompile(`\D`)
	digits := reg.ReplaceAllString(p.Value, "")

	length := len(digits)

	if length < 10 || length > 11 {
		return errors.New("invalid phone number: must have 10 or 11 digits")
	}

	return nil
}

func (p PhoneNumber) String() string {
	return p.Value
}

func (p PhoneNumber) Equals(_p PhoneNumber) bool {
	return p.Value == _p.Value
}
