package user

const UserIdentityPrefix = "usr"

type UserStatuses string

const (
	UserStatusActive     UserStatuses = "active"
	UserStatusUnverified UserStatuses = "unverified"
	UserStatusInactive   UserStatuses = "inactive"
)

type PasswordRecoveryStatuses string

const (
	PasswordRecoveryStatusNew     PasswordRecoveryStatuses = "new"
	PasswordRecoveryStatusUsed    PasswordRecoveryStatuses = "used"
	PasswordRecoveryStatusExpired PasswordRecoveryStatuses = "expired"
)

type UserRegistrationStatuses string

const (
	UserRegistrationStatusPending  UserRegistrationStatuses = "pending"
	UserRegistrationStatusVerified UserRegistrationStatuses = "verified"
	UserRegistrationStatusExpired  UserRegistrationStatuses = "expired"
)
