package role_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	role_core "github.com/gabrielmrtt/taski/internal/role"
	role_repositories "github.com/gabrielmrtt/taski/internal/role/repositories"
)

type ListRolesService struct {
	RoleRepository role_repositories.RoleRepository
}

func NewListRolesService(
	roleRepository role_repositories.RoleRepository,
) *ListRolesService {
	return &ListRolesService{
		RoleRepository: roleRepository,
	}
}

type ListRolesInput struct {
	OrganizationIdentity core.Identity
	Filters              role_repositories.RoleFilters
	SortInput            *core.SortInput
	Pagination           *core.PaginationInput
	RelationsInput       core.RelationsInput
}

func (i ListRolesInput) Validate() error {
	return nil
}

func (s *ListRolesService) Execute(input ListRolesInput) (*core.PaginationOutput[role_core.RoleDto], error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	input.Filters.OrganizationIdentity = input.OrganizationIdentity

	roles, err := s.RoleRepository.PaginateRolesBy(role_repositories.PaginateRolesParams{
		Filters:        input.Filters,
		SortInput:      input.SortInput,
		Pagination:     input.Pagination,
		RelationsInput: input.RelationsInput,
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
