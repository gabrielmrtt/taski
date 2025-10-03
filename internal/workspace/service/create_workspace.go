package workspaceservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	userrepo "github.com/gabrielmrtt/taski/internal/user/repository"
	"github.com/gabrielmrtt/taski/internal/workspace"
	workspacerepo "github.com/gabrielmrtt/taski/internal/workspace/repository"
)

type CreateWorkspaceService struct {
	WorkspaceRepository     workspacerepo.WorkspaceRepository
	UserRepository          userrepo.UserRepository
	WorkspaceUserRepository workspacerepo.WorkspaceUserRepository
	TransactionRepository   core.TransactionRepository
}

func NewCreateWorkspaceService(
	workspaceRepository workspacerepo.WorkspaceRepository,
	userRepository userrepo.UserRepository,
	workspaceUserRepository workspacerepo.WorkspaceUserRepository,
	transactionRepository core.TransactionRepository,
) *CreateWorkspaceService {
	return &CreateWorkspaceService{
		WorkspaceRepository:     workspaceRepository,
		UserRepository:          userRepository,
		WorkspaceUserRepository: workspaceUserRepository,
		TransactionRepository:   transactionRepository,
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

func (s *CreateWorkspaceService) Execute(input CreateWorkspaceInput) (*workspace.WorkspaceDto, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return nil, err
	}

	s.WorkspaceRepository.SetTransaction(tx)
	s.UserRepository.SetTransaction(tx)
	s.WorkspaceUserRepository.SetTransaction(tx)

	wrk, err := workspace.NewWorkspace(workspace.NewWorkspaceInput{
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

	wrk, err = s.WorkspaceRepository.StoreWorkspace(workspacerepo.StoreWorkspaceParams{Workspace: wrk})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	user, err := s.UserRepository.GetUserByIdentity(userrepo.GetUserByIdentityParams{
		UserIdentity: input.UserCreatorIdentity,
	})
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, core.NewNotFoundError("user not found")
	}

	workspaceUser, err := workspace.NewWorkspaceUser(workspace.NewWorkspaceUserInput{
		WorkspaceIdentity: wrk.Identity,
		User:              *user,
		Status:            workspace.WorkspaceUserStatusActive,
	})
	if err != nil {
		return nil, err
	}

	_, err = s.WorkspaceUserRepository.StoreWorkspaceUser(workspacerepo.StoreWorkspaceUserParams{
		WorkspaceUser: workspaceUser,
	})
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return workspace.WorkspaceToDto(wrk), nil
}
