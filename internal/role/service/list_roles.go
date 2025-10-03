package roleservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/role"
	rolerepo "github.com/gabrielmrtt/taski/internal/role/repository"
)

type ListRolesService struct {
	RoleRepository rolerepo.RoleRepository
}

func NewListRolesService(
	roleRepository rolerepo.RoleRepository,
) *ListRolesService {
	return &ListRolesService{
		RoleRepository: roleRepository,
	}
}

type ListRolesInput struct {
	OrganizationIdentity core.Identity
	Filters              rolerepo.RoleFilters
	SortInput            *core.SortInput
	Pagination           *core.PaginationInput
	RelationsInput       core.RelationsInput
}

func (i ListRolesInput) Validate() error {
	return nil
}

func (s *ListRolesService) Execute(input ListRolesInput) (*core.PaginationOutput[role.RoleDto], error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	input.Filters.OrganizationIdentity = input.OrganizationIdentity

	roles, err := s.RoleRepository.PaginateRolesBy(rolerepo.PaginateRolesParams{
		Filters:    input.Filters,
		SortInput:  input.SortInput,
		Pagination: input.Pagination,
	})
	if err != nil {
		return nil, err
	}

	var rolesDto []role.RoleDto
	for _, rol := range roles.Data {
		rolesDto = append(rolesDto, *role.RoleToDto(&rol))
	}

	paginationOutput := core.PaginationOutput[role.RoleDto]{
		Data:    rolesDto,
		Page:    roles.Page,
		HasMore: roles.HasMore,
		Total:   roles.Total,
	}

	return &paginationOutput, nil
}
