package projectservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	project "github.com/gabrielmrtt/taski/internal/project"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
	userrepo "github.com/gabrielmrtt/taski/internal/user/repository"
	workspacerepo "github.com/gabrielmrtt/taski/internal/workspace/repository"
)

type CreateProjectService struct {
	ProjectRepository     projectrepo.ProjectRepository
	ProjectUserRepository projectrepo.ProjectUserRepository
	UserRepository        userrepo.UserRepository
	WorkspaceRepository   workspacerepo.WorkspaceRepository
	TransactionRepository core.TransactionRepository
}

func NewCreateProjectService(
	projectRepository projectrepo.ProjectRepository,
	projectUserRepository projectrepo.ProjectUserRepository,
	userRepository userrepo.UserRepository,
	workspaceRepository workspacerepo.WorkspaceRepository,
	transactionRepository core.TransactionRepository,
) *CreateProjectService {
	return &CreateProjectService{
		ProjectRepository:     projectRepository,
		ProjectUserRepository: projectUserRepository,
		UserRepository:        userRepository,
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
	PriorityLevel        project.ProjectPriorityLevels
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

func (s *CreateProjectService) Execute(input CreateProjectInput) (*project.ProjectDto, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return nil, err
	}

	s.ProjectRepository.SetTransaction(tx)
	s.ProjectUserRepository.SetTransaction(tx)
	s.UserRepository.SetTransaction(tx)

	user, err := s.UserRepository.GetUserByIdentity(userrepo.GetUserByIdentityParams{
		UserIdentity: input.UserCreatorIdentity,
	})
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, core.NewNotFoundError("user creator not found")
	}

	wrk, err := s.WorkspaceRepository.GetWorkspaceByIdentity(workspacerepo.GetWorkspaceByIdentityParams{
		WorkspaceIdentity:    input.WorkspaceIdentity,
		OrganizationIdentity: &input.OrganizationIdentity,
	})
	if err != nil {
		return nil, err
	}

	if wrk == nil {
		return nil, core.NewNotFoundError("workspace not found")
	}

	prj, err := project.NewProject(project.NewProjectInput{
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

	prj, err = s.ProjectRepository.StoreProject(projectrepo.StoreProjectParams{Project: prj})
	if err != nil {
		return nil, err
	}

	projectUser, err := project.NewProjectUser(project.NewProjectUserInput{
		ProjectIdentity: prj.Identity,
		User:            *user,
		Status:          project.ProjectUserStatusActive,
	})
	if err != nil {
		return nil, err
	}

	_, err = s.ProjectUserRepository.StoreProjectUser(projectrepo.StoreProjectUserParams{ProjectUser: projectUser})
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return project.ProjectToDto(prj), nil
}
