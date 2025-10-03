package organization

const OrganizationIdentityPrefix = "org"

type OrganizationStatuses string

const (
	OrganizationStatusActive   OrganizationStatuses = "active"
	OrganizationStatusInactive OrganizationStatuses = "inactive"
)

type OrganizationUserStatuses string

const (
	OrganizationUserStatusActive   OrganizationUserStatuses = "active"
	OrganizationUserStatusInactive OrganizationUserStatuses = "inactive"
	OrganizationUserStatusInvited  OrganizationUserStatuses = "invited"
	OrganizationUserStatusRefused  OrganizationUserStatuses = "refused"
)
