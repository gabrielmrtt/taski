package projectservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	project "github.com/gabrielmrtt/taski/internal/project"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
)

type UpdateProjectService struct {
	ProjectRepository     projectrepo.ProjectRepository
	TransactionRepository core.TransactionRepository
}

func NewUpdateProjectService(
	projectRepository projectrepo.ProjectRepository,
	transactionRepository core.TransactionRepository,
) *UpdateProjectService {
	return &UpdateProjectService{
		ProjectRepository:     projectRepository,
		TransactionRepository: transactionRepository,
	}
}

type UpdateProjectInput struct {
	OrganizationIdentity core.Identity
	WorkspaceIdentity    core.Identity
	ProjectIdentity      core.Identity
	UserEditorIdentity   core.Identity
	Name                 *string
	Description          *string
	Color                *string
	Status               *project.ProjectStatuses
	PriorityLevel        *project.ProjectPriorityLevels
	StartAt              *int64
	EndAt                *int64
}

func (i UpdateProjectInput) Validate() error {
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

func (s *UpdateProjectService) Execute(input UpdateProjectInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return err
	}

	s.ProjectRepository.SetTransaction(tx)

	project, err := s.ProjectRepository.GetProjectByIdentity(projectrepo.GetProjectByIdentityParams{
		ProjectIdentity:      input.ProjectIdentity,
		WorkspaceIdentity:    &input.WorkspaceIdentity,
		OrganizationIdentity: &input.OrganizationIdentity,
	})
	if err != nil {
		return err
	}

	if project == nil {
		return core.NewNotFoundError("project not found")
	}

	if input.Name != nil {
		err = project.ChangeName(*input.Name, &input.UserEditorIdentity)
		if err != nil {
			return err
		}
	}

	if input.Description != nil {
		err = project.ChangeDescription(*input.Description, &input.UserEditorIdentity)
		if err != nil {
			return err
		}
	}

	if input.Color != nil {
		err = project.ChangeColor(*input.Color, &input.UserEditorIdentity)
		if err != nil {
			return err
		}
	}

	if input.Status != nil {
		err = project.ChangeStatus(*input.Status, &input.UserEditorIdentity)
		if err != nil {
			return err
		}
	}

	if input.PriorityLevel != nil {
		err = project.ChangePriorityLevel(*input.PriorityLevel, &input.UserEditorIdentity)
		if err != nil {
			return err
		}
	}

	if input.StartAt != nil {
		err = project.ChangeStartAt(*input.StartAt, &input.UserEditorIdentity)
		if err != nil {
			return err
		}
	}

	if input.EndAt != nil {
		err = project.ChangeEndAt(*input.EndAt, &input.UserEditorIdentity)
		if err != nil {
			return err
		}
	}

	err = s.ProjectRepository.UpdateProject(projectrepo.UpdateProjectParams{Project: project})
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
