package workspace_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	workspace_core "github.com/gabrielmrtt/taski/internal/workspace"
	workspace_repositories "github.com/gabrielmrtt/taski/internal/workspace/repositories"
)

type UpdateWorkspaceService struct {
	WorkspaceRepository   workspace_repositories.WorkspaceRepository
	TransactionRepository core.TransactionRepository
}

func NewUpdateWorkspaceService(
	workspaceRepository workspace_repositories.WorkspaceRepository,
	transactionRepository core.TransactionRepository,
) *UpdateWorkspaceService {
	return &UpdateWorkspaceService{
		WorkspaceRepository:   workspaceRepository,
		TransactionRepository: transactionRepository,
	}
}

type UpdateWorkspaceInput struct {
	OrganizationIdentity core.Identity
	WorkspaceIdentity    core.Identity
	Name                 *string
	Description          *string
	Color                *string
	Status               *workspace_core.WorkspaceStatuses
	UserEditorIdentity   core.Identity
}

func (i UpdateWorkspaceInput) Validate() error {
	var fields []core.InvalidInputErrorField

	if i.Name != nil {
		_, err := core.NewName(*i.Name)
		if err != nil {
			fields = append(fields, core.InvalidInputErrorField{
				Field: "name",
				Error: err.Error(),
			})
		}
	}

	if i.Description != nil {
		_, err := core.NewDescription(*i.Description)
		if err != nil {
			fields = append(fields, core.InvalidInputErrorField{
				Field: "description",
				Error: err.Error(),
			})
		}
	}

	if i.Color != nil {
		_, err := core.NewColor(*i.Color)
		if err != nil {
			fields = append(fields, core.InvalidInputErrorField{
				Field: "color",
				Error: err.Error(),
			})
		}
	}

	if len(fields) > 0 {
		return core.NewInvalidInputError("invalid input", fields)
	}

	return nil
}

func (s *UpdateWorkspaceService) Execute(input UpdateWorkspaceInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return err
	}

	s.WorkspaceRepository.SetTransaction(tx)

	workspace, err := s.WorkspaceRepository.GetWorkspaceByIdentity(workspace_repositories.GetWorkspaceByIdentityParams{
		WorkspaceIdentity:    input.WorkspaceIdentity,
		OrganizationIdentity: &input.OrganizationIdentity,
	})
	if err != nil {
		return err
	}

	if workspace == nil {
		return core.NewNotFoundError("workspace not found")
	}

	if input.Name != nil {
		err = workspace.ChangeName(*input.Name, &input.UserEditorIdentity)
	}

	if input.Description != nil {
		err = workspace.ChangeDescription(*input.Description, &input.UserEditorIdentity)
	}

	if input.Color != nil {
		err = workspace.ChangeColor(*input.Color, &input.UserEditorIdentity)
	}

	if input.Status != nil {
		err = workspace.ChangeStatus(*input.Status, &input.UserEditorIdentity)
	}

	err = s.WorkspaceRepository.UpdateWorkspace(workspace_repositories.UpdateWorkspaceParams{Workspace: workspace})
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
