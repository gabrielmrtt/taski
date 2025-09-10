package role_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	role_core "github.com/gabrielmrtt/taski/internal/role"
)

type ListRolesService struct {
	RoleRepository role_core.RoleRepository
}

func NewListRolesService(
	roleRepository role_core.RoleRepository,
) *ListRolesService {
	return &ListRolesService{
		RoleRepository: roleRepository,
	}
}

type ListRolesInput struct {
	OrganizationIdentity core.Identity
	Filters              role_core.RoleFilters
	SortInput            *core.SortInput
	Pagination           *core.PaginationInput
	LoggedUserIdentity   core.Identity
}

func (i ListRolesInput) Validate() error {
	return nil
}

func (s *ListRolesService) Execute(input ListRolesInput) (*core.PaginationOutput[role_core.RoleDto], error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	organizationHasUser, err := s.RoleRepository.CheckIfOrganizatonHasUser(input.OrganizationIdentity, input.LoggedUserIdentity)
	if err != nil {
		return nil, err
	}

	if !organizationHasUser {
		return nil, core.NewUnauthorizedError("user is not part of the organization")
	}

	input.Filters.OrganizationIdentity = input.OrganizationIdentity

	roles, err := s.RoleRepository.PaginateRolesBy(role_core.PaginateRolesParams{
		Filters:    input.Filters,
		SortInput:  input.SortInput,
		Pagination: input.Pagination,
	})
	if err != nil {
		return nil, err
	}

	var rolesDto []role_core.RoleDto
	for _, role := range roles.Data {
		rolesDto = append(rolesDto, *role_core.RoleToDto(&role))
	}

	paginationOutput := core.PaginationOutput[role_core.RoleDto]{
		Data:    rolesDto,
		Page:    roles.Page,
		HasMore: roles.HasMore,
		Total:   roles.Total,
	}

	return &paginationOutput, nil
}
