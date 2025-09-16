package workspace_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	workspace_core "github.com/gabrielmrtt/taski/internal/workspace"
	workspace_repositories "github.com/gabrielmrtt/taski/internal/workspace/repositories"
)

type GetWorkspaceService struct {
	WorkspaceRepository workspace_repositories.WorkspaceRepository
}

func NewGetWorkspaceService(
	workspaceRepository workspace_repositories.WorkspaceRepository,
) *GetWorkspaceService {
	return &GetWorkspaceService{
		WorkspaceRepository: workspaceRepository,
	}
}

type GetWorkspaceInput struct {
	OrganizationIdentity core.Identity
	WorkspaceIdentity    core.Identity
}

func (i GetWorkspaceInput) Validate() error {
	return nil
}

func (s *GetWorkspaceService) Execute(input GetWorkspaceInput) (*workspace_core.WorkspaceDto, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	workspace, err := s.WorkspaceRepository.GetWorkspaceByIdentity(workspace_repositories.GetWorkspaceByIdentityParams{WorkspaceIdentity: input.WorkspaceIdentity, OrganizationIdentity: &input.OrganizationIdentity})
	if err != nil {
		return nil, err
	}

	if workspace == nil {
		return nil, core.NewNotFoundError("workspace not found")
	}

	return workspace_core.WorkspaceToDto(workspace), nil
}
