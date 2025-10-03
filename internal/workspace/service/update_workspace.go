package workspaceservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/workspace"
	workspacerepo "github.com/gabrielmrtt/taski/internal/workspace/repository"
)

type UpdateWorkspaceService struct {
	WorkspaceRepository   workspacerepo.WorkspaceRepository
	TransactionRepository core.TransactionRepository
}

func NewUpdateWorkspaceService(
	workspaceRepository workspacerepo.WorkspaceRepository,
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
	Status               *workspace.WorkspaceStatuses
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

	wrk, err := s.WorkspaceRepository.GetWorkspaceByIdentity(workspacerepo.GetWorkspaceByIdentityParams{
		WorkspaceIdentity:    input.WorkspaceIdentity,
		OrganizationIdentity: &input.OrganizationIdentity,
	})
	if err != nil {
		return err
	}

	if wrk == nil {
		return core.NewNotFoundError("workspace not found")
	}

	if input.Name != nil {
		err = wrk.ChangeName(*input.Name, &input.UserEditorIdentity)
		if err != nil {
			return err
		}
	}

	if input.Description != nil {
		err = wrk.ChangeDescription(*input.Description, &input.UserEditorIdentity)
		if err != nil {
			return err
		}
	}

	if input.Color != nil {
		err = wrk.ChangeColor(*input.Color, &input.UserEditorIdentity)
		if err != nil {
			return err
		}
	}

	if input.Status != nil {
		err = wrk.ChangeStatus(*input.Status, &input.UserEditorIdentity)
		if err != nil {
			return err
		}
	}

	err = s.WorkspaceRepository.UpdateWorkspace(workspacerepo.UpdateWorkspaceParams{Workspace: wrk})
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
