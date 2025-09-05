package user_core

import (
	"time"

	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/pkg/datetimeutils"
	"github.com/gabrielmrtt/taski/pkg/hashutils"
	"github.com/gabrielmrtt/taski/pkg/stringutils"
	"github.com/google/uuid"
)

type UserStatuses string

const (
	UserStatusActive     UserStatuses = "active"
	UserStatusUnverified UserStatuses = "unverified"
	UserStatusInactive   UserStatuses = "inactive"
)

type User struct {
	Identity    core.Identity
	Status      UserStatuses
	Credentials *UserCredentials
	Data        *UserData
	Timestamps  core.Timestamps
	DeletedAt   *int64
}

type UserCredentials struct {
	Name        string
	Email       string
	Password    string
	PhoneNumber *string
}

type UserData struct {
	DisplayName            string
	About                  *string
	ProfilePictureIdentity *core.Identity
}

type NewUserInput struct {
	Name        string
	Email       string
	Password    string
	PhoneNumber *string
}

func NewUser(input NewUserInput) (*User, error) {
	emailValueObject, err := NewEmail(input.Email)
	if err != nil {
		return nil, err
	}

	phoneNumberValueObject := (*PhoneNumber)(nil)

	if input.PhoneNumber != nil {
		vo, err := NewPhoneNumber(*input.PhoneNumber)
		if err != nil {
			return nil, err
		}

		phoneNumberValueObject = &vo
	}

	passwordValueObject, err := NewPassword(input.Password)
	if err != nil {
		return nil, err
	}

	nameValueObject, err := core.NewName(input.Name)
	if err != nil {
		return nil, err
	}

	identity := core.NewIdentity("user")

	now := datetimeutils.EpochNow()

	timestamps := core.Timestamps{
		CreatedAt: &now,
		UpdatedAt: nil,
	}

	hashedPassword, err := hashutils.HashPassword(passwordValueObject.Value)
	if err != nil {
		return nil, err
	}

	user := &User{
		Identity: identity,
		Status:   UserStatusUnverified,
		Credentials: &UserCredentials{
			Name:        nameValueObject.Value,
			Password:    hashedPassword,
			Email:       emailValueObject.Value,
			PhoneNumber: nil,
		},
		Data: &UserData{
			DisplayName:            input.Name,
			About:                  nil,
			ProfilePictureIdentity: nil,
		},
		Timestamps: timestamps,
		DeletedAt:  nil,
	}

	if phoneNumberValueObject != nil {
		user.Credentials.PhoneNumber = &phoneNumberValueObject.Value
	}

	return user, nil
}

func (u *User) ChangeCredentialsName(name string) error {
	if u.Credentials == nil {
		return core.NewInternalError("user credentials not found")
	}

	nameValueObject, err := core.NewName(name)
	if err != nil {
		return err
	}

	u.Credentials.Name = nameValueObject.Value
	now := datetimeutils.EpochNow()

	u.Timestamps.UpdatedAt = &now
	return nil
}

func (u *User) ChangeCredentialsEmail(email string) error {
	if u.Credentials == nil {
		return core.NewInternalError("user credentials not found")
	}

	emailValueObject, err := NewEmail(email)
	if err != nil {
		return err
	}

	now := datetimeutils.EpochNow()

	u.Credentials.Email = emailValueObject.Value
	u.Timestamps.UpdatedAt = &now

	return nil
}

func (u *User) ChangeCredentialsPassword(password string) error {
	if u.Credentials == nil {
		return core.NewInternalError("user credentials not found")
	}

	passwordValueObject, err := NewPassword(password)
	if err != nil {
		return err
	}

	hashedPassword, err := hashutils.HashPassword(passwordValueObject.Value)
	if err != nil {
		return err
	}

	now := datetimeutils.EpochNow()

	u.Credentials.Password = hashedPassword
	u.Timestamps.UpdatedAt = &now

	return nil
}

func (u *User) ChangeCredentialsPhoneNumber(phoneNumber string) error {
	if u.Credentials == nil {
		return core.NewInternalError("user credentials not found")
	}

	phoneNumberValueObject, err := NewPhoneNumber(phoneNumber)
	if err != nil {
		return err
	}

	now := datetimeutils.EpochNow()

	u.Credentials.PhoneNumber = &phoneNumberValueObject.Value
	u.Timestamps.UpdatedAt = &now

	return nil
}

func (u *User) ChangeUserDataDisplayName(displayName string) error {
	if u.Data == nil {
		return core.NewInternalError("user data not found")
	}

	displayNameValueObject, err := core.NewName(displayName)
	if err != nil {
		return err
	}

	u.Data.DisplayName = displayNameValueObject.Value
	now := datetimeutils.EpochNow()

	u.Timestamps.UpdatedAt = &now

	return nil
}

func (u *User) ChangeUserDataAbout(about string) error {
	if u.Data == nil {
		return core.NewInternalError("user data not found")
	}

	aboutValueObject, err := core.NewDescription(about)
	if err != nil {
		return err
	}

	u.Data.About = &aboutValueObject.Value
	now := datetimeutils.EpochNow()

	u.Timestamps.UpdatedAt = &now

	return nil
}

func (u *User) ChangeUserDataProfilePicture(profilePictureIdentity *core.Identity) error {
	if u.Data == nil {
		return core.NewInternalError("user data not found")
	}

	u.Data.ProfilePictureIdentity = profilePictureIdentity
	now := datetimeutils.EpochNow()

	u.Timestamps.UpdatedAt = &now

	return nil
}

func (u *User) Activate() {
	u.Status = UserStatusActive
	now := datetimeutils.EpochNow()

	u.Timestamps.UpdatedAt = &now
}

func (u *User) Deactivate() {
	u.Status = UserStatusInactive
	now := datetimeutils.EpochNow()

	u.Timestamps.UpdatedAt = &now
}

func (u *User) Delete() {
	now := datetimeutils.EpochNow()

	u.DeletedAt = &now
}

func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

func (u *User) IsInactive() bool {
	return u.Status == UserStatusInactive
}

func (u *User) IsUnverified() bool {
	return u.Status == UserStatusUnverified
}

func (u *User) IsDeleted() bool {
	return u.DeletedAt != nil
}

func (u *User) CheckPassword(password string) bool {
	if u.Credentials == nil {
		return false
	}
	return hashutils.ComparePassword(password, u.Credentials.Password)
}

type PasswordRecoveryStatuses string

const (
	PasswordRecoveryStatusNew     PasswordRecoveryStatuses = "new"
	PasswordRecoveryStatusUsed    PasswordRecoveryStatuses = "used"
	PasswordRecoveryStatusExpired PasswordRecoveryStatuses = "expired"
)

type PasswordRecovery struct {
	Identity     core.Identity
	UserIdentity core.Identity
	Token        string
	Status       PasswordRecoveryStatuses
	RecoveredAt  *int64
	ExpiresAt    int64
	RequestedAt  int64
}

func NewPasswordRecovery(userIdentity core.Identity, expiresIn time.Duration) (*PasswordRecovery, error) {
	now := datetimeutils.EpochNow()
	token := stringutils.GenerateRandomString(32, stringutils.RandomStringOptionsAlphanumeric)
	expiresAt := now + int64(expiresIn.Seconds())

	passwordRecovery := &PasswordRecovery{
		Identity: core.Identity{
			Internal: uuid.New(),
		},
		UserIdentity: userIdentity,
		Token:        token,
		Status:       PasswordRecoveryStatusNew,
		ExpiresAt:    expiresAt,
		RequestedAt:  now,
		RecoveredAt:  nil,
	}

	return passwordRecovery, nil
}

func (p *PasswordRecovery) IsUsed() bool {
	return p.Status == PasswordRecoveryStatusUsed
}

func (p *PasswordRecovery) IsExpired() bool {
	return p.Status == PasswordRecoveryStatusExpired || p.ExpiresAt < datetimeutils.EpochNow()
}

func (p *PasswordRecovery) Use() {
	p.Status = PasswordRecoveryStatusUsed
	now := datetimeutils.EpochNow()
	p.RecoveredAt = &now
}

type UserRegistrationStatuses string

const (
	UserRegistrationStatusPending  UserRegistrationStatuses = "pending"
	UserRegistrationStatusVerified UserRegistrationStatuses = "verified"
	UserRegistrationStatusExpired  UserRegistrationStatuses = "expired"
)

type UserRegistration struct {
	Identity     core.Identity
	UserIdentity core.Identity
	Token        string
	Status       UserRegistrationStatuses
	ExpiresAt    int64
	RegisteredAt int64
	VerifiedAt   *int64
}

func NewUserRegistration(userIdentity core.Identity, expiresIn time.Duration) (*UserRegistration, error) {
	now := datetimeutils.EpochNow()

	token := stringutils.GenerateRandomString(32, stringutils.RandomStringOptionsAlphanumeric)

	expiresAt := now + int64(expiresIn.Seconds())

	userRegistration := &UserRegistration{
		Identity: core.Identity{
			Internal: uuid.New(),
		},
		UserIdentity: userIdentity,
		Token:        token,
		Status:       UserRegistrationStatusPending,
		ExpiresAt:    expiresAt,
		RegisteredAt: now,
		VerifiedAt:   nil,
	}

	return userRegistration, nil
}

func (u *UserRegistration) Verify() {
	now := datetimeutils.EpochNow()

	u.Status = UserRegistrationStatusVerified
	u.VerifiedAt = &now
}

func (u *UserRegistration) IsExpired() bool {
	return u.Status == UserRegistrationStatusExpired || u.ExpiresAt < datetimeutils.EpochNow()
}
