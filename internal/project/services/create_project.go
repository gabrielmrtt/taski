package project_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	project_core "github.com/gabrielmrtt/taski/internal/project"
	project_repositories "github.com/gabrielmrtt/taski/internal/project/repositories"
	workspace_repositories "github.com/gabrielmrtt/taski/internal/workspace/repositories"
)

type CreateProjectService struct {
	ProjectRepository     project_repositories.ProjectRepository
	WorkspaceRepository   workspace_repositories.WorkspaceRepository
	TransactionRepository core.TransactionRepository
}

func NewCreateProjectService(
	projectRepository project_repositories.ProjectRepository,
	workspaceRepository workspace_repositories.WorkspaceRepository,
	transactionRepository core.TransactionRepository,
) *CreateProjectService {
	return &CreateProjectService{
		ProjectRepository:     projectRepository,
		WorkspaceRepository:   workspaceRepository,
		TransactionRepository: transactionRepository,
	}
}

type CreateProjectInput struct {
	WorkspaceIdentity    core.Identity
	OrganizationIdentity core.Identity
	UserCreatorIdentity  core.Identity
	Name                 string
	Description          string
	Color                string
	PriorityLevel        project_core.ProjectPriorityLevels
	StartAt              *int64
	EndAt                *int64
}

func (i CreateProjectInput) Validate() error {
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

func (s *CreateProjectService) Execute(input CreateProjectInput) (*project_core.ProjectDto, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return nil, err
	}

	s.ProjectRepository.SetTransaction(tx)

	workspace, err := s.WorkspaceRepository.GetWorkspaceByIdentity(workspace_repositories.GetWorkspaceByIdentityParams{
		WorkspaceIdentity:    input.WorkspaceIdentity,
		OrganizationIdentity: &input.OrganizationIdentity,
	})
	if err != nil {
		return nil, err
	}

	if workspace == nil {
		return nil, core.NewNotFoundError("workspace not found")
	}

	project, err := project_core.NewProject(project_core.NewProjectInput{
		Name:                input.Name,
		Description:         input.Description,
		Color:               input.Color,
		WorkspaceIdentity:   input.WorkspaceIdentity,
		PriorityLevel:       input.PriorityLevel,
		StartAt:             input.StartAt,
		EndAt:               input.EndAt,
		UserCreatorIdentity: &input.UserCreatorIdentity,
	})
	if err != nil {
		return nil, err
	}

	project, err = s.ProjectRepository.StoreProject(project_repositories.StoreProjectParams{Project: project})
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return project_core.ProjectToDto(project), nil
}
