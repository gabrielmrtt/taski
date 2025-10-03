package workspaceservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/workspace"
	workspacerepo "github.com/gabrielmrtt/taski/internal/workspace/repository"
)

type GetWorkspaceService struct {
	WorkspaceRepository workspacerepo.WorkspaceRepository
}

func NewGetWorkspaceService(
	workspaceRepository workspacerepo.WorkspaceRepository,
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

func (s *GetWorkspaceService) Execute(input GetWorkspaceInput) (*workspace.WorkspaceDto, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	wrk, err := s.WorkspaceRepository.GetWorkspaceByIdentity(workspacerepo.GetWorkspaceByIdentityParams{WorkspaceIdentity: input.WorkspaceIdentity, OrganizationIdentity: &input.OrganizationIdentity})
	if err != nil {
		return nil, err
	}

	if wrk == nil {
		return nil, core.NewNotFoundError("workspace not found")
	}

	return workspace.WorkspaceToDto(wrk), nil
}
