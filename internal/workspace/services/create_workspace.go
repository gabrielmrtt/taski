package workspace_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	workspace_core "github.com/gabrielmrtt/taski/internal/workspace"
	workspace_repositories "github.com/gabrielmrtt/taski/internal/workspace/repositories"
)

type CreateWorkspaceService struct {
	WorkspaceRepository   workspace_repositories.WorkspaceRepository
	TransactionRepository core.TransactionRepository
}

func NewCreateWorkspaceService(
	workspaceRepository workspace_repositories.WorkspaceRepository,
	transactionRepository core.TransactionRepository,
) *CreateWorkspaceService {
	return &CreateWorkspaceService{
		WorkspaceRepository:   workspaceRepository,
		TransactionRepository: transactionRepository,
	}
}

type CreateWorkspaceInput struct {
	Name                 string
	Description          string
	Color                string
	OrganizationIdentity core.Identity
	UserCreatorIdentity  core.Identity
}

func (i CreateWorkspaceInput) Validate() error {
	var fields []core.InvalidInputErrorField

	if _, err := core.NewName(i.Name); err != nil {
		fields = append(fields, core.InvalidInputErrorField{
			Field: "name",
			Error: err.Error(),
		})
	}

	if _, err := core.NewDescription(i.Description); err != nil {
		fields = append(fields, core.InvalidInputErrorField{
			Field: "description",
			Error: err.Error(),
		})
	}

	if _, err := core.NewColor(i.Color); err != nil {
		fields = append(fields, core.InvalidInputErrorField{
			Field: "color",
			Error: err.Error(),
		})
	}

	if len(fields) > 0 {
		return core.NewInvalidInputError("invalid input", fields)
	}

	return nil
}

func (s *CreateWorkspaceService) Execute(input CreateWorkspaceInput) (*workspace_core.WorkspaceDto, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return nil, err
	}

	s.WorkspaceRepository.SetTransaction(tx)

	workspace, err := workspace_core.NewWorkspace(workspace_core.NewWorkspaceInput{
		Name:                 input.Name,
		Description:          input.Description,
		Color:                input.Color,
		OrganizationIdentity: input.OrganizationIdentity,
		UserCreatorIdentity:  &input.UserCreatorIdentity,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	workspace, err = s.WorkspaceRepository.StoreWorkspace(workspace_repositories.StoreWorkspaceParams{Workspace: workspace})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return workspace_core.WorkspaceToDto(workspace), nil
}
