package workspaceservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	workspacerepo "github.com/gabrielmrtt/taski/internal/workspace/repository"
)

type DeleteWorkspaceService struct {
	WorkspaceRepository   workspacerepo.WorkspaceRepository
	TransactionRepository core.TransactionRepository
}

func NewDeleteWorkspaceService(
	workspaceRepository workspacerepo.WorkspaceRepository,
	transactionRepository core.TransactionRepository,
) *DeleteWorkspaceService {
	return &DeleteWorkspaceService{
		WorkspaceRepository:   workspaceRepository,
		TransactionRepository: transactionRepository,
	}
}

type DeleteWorkspaceInput struct {
	OrganizationIdentity core.Identity
	WorkspaceIdentity    core.Identity
}

func (i DeleteWorkspaceInput) Validate() error {
	return nil
}

func (s *DeleteWorkspaceService) Execute(input DeleteWorkspaceInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return err
	}

	s.WorkspaceRepository.SetTransaction(tx)

	workspace, err := s.WorkspaceRepository.GetWorkspaceByIdentity(workspacerepo.GetWorkspaceByIdentityParams{WorkspaceIdentity: input.WorkspaceIdentity, OrganizationIdentity: &input.OrganizationIdentity})
	if err != nil {
		return err
	}

	if workspace == nil {
		return core.NewNotFoundError("workspace not found")
	}

	workspace.Delete()

	err = s.WorkspaceRepository.UpdateWorkspace(workspacerepo.UpdateWorkspaceParams{Workspace: workspace})
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
