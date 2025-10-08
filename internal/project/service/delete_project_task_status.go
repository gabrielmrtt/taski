package projectservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
)

type DeleteProjectTaskStatusService struct {
	ProjectRepository           projectrepo.ProjectRepository
	ProjectTaskStatusRepository projectrepo.ProjectTaskStatusRepository
	TransactionRepository       core.TransactionRepository
}

func NewDeleteProjectTaskStatusService(
	projectRepository projectrepo.ProjectRepository,
	projectTaskStatusRepository projectrepo.ProjectTaskStatusRepository,
	transactionRepository core.TransactionRepository,
) *DeleteProjectTaskStatusService {
	return &DeleteProjectTaskStatusService{
		ProjectRepository:           projectRepository,
		ProjectTaskStatusRepository: projectTaskStatusRepository,
		TransactionRepository:       transactionRepository,
	}
}

type DeleteProjectTaskStatusInput struct {
	OrganizationIdentity      core.Identity
	ProjectIdentity           core.Identity
	ProjectTaskStatusIdentity core.Identity
}

func (i DeleteProjectTaskStatusInput) Validate() error {
	return nil
}

func (s *DeleteProjectTaskStatusService) Execute(input DeleteProjectTaskStatusInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return err
	}

	s.ProjectTaskStatusRepository.SetTransaction(tx)
	s.ProjectRepository.SetTransaction(tx)

	prj, err := s.ProjectRepository.GetProjectByIdentity(projectrepo.GetProjectByIdentityParams{
		ProjectIdentity:      input.ProjectIdentity,
		OrganizationIdentity: &input.OrganizationIdentity,
	})
	if err != nil {
		return err
	}

	if prj == nil {
		return core.NewNotFoundError("project not found")
	}

	projectTaskStatus, err := s.ProjectTaskStatusRepository.GetProjectTaskStatusByIdentity(projectrepo.GetProjectTaskStatusByIdentityParams{
		ProjectTaskStatusIdentity: &input.ProjectTaskStatusIdentity,
		ProjectIdentity:           &prj.Identity,
	})
	if err != nil {
		return err
	}

	if projectTaskStatus == nil {
		return core.NewNotFoundError("project task status not found")
	}

	projectTaskStatus.Delete()

	err = s.ProjectTaskStatusRepository.UpdateProjectTaskStatus(projectrepo.UpdateProjectTaskStatusParams{
		ProjectTaskStatus: projectTaskStatus,
	})
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
