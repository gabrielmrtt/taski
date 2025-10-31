package projectservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/project"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
)

type GetProjectDocumentVersionService struct {
	ProjectRepository         projectrepo.ProjectRepository
	ProjectDocumentRepository projectrepo.ProjectDocumentRepository
}

func NewGetProjectDocumentVersionService(projectRepository projectrepo.ProjectRepository, projectDocumentRepository projectrepo.ProjectDocumentRepository) *GetProjectDocumentVersionService {
	return &GetProjectDocumentVersionService{
		ProjectRepository:         projectRepository,
		ProjectDocumentRepository: projectDocumentRepository,
	}
}

type GetProjectDocumentVersionInput struct {
	ProjectIdentity                       core.Identity
	ProjectDocumentVersionManagerIdentity core.Identity
	ProjectDocumentVersionIdentity        core.Identity
	RelationsInput                        core.RelationsInput
}

func (i GetProjectDocumentVersionInput) Validate() error {
	return nil
}

func (s *GetProjectDocumentVersionService) Execute(input GetProjectDocumentVersionInput) (*project.ProjectDocumentVersionDto, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	prj, err := s.ProjectRepository.GetProjectByIdentity(projectrepo.GetProjectByIdentityParams{
		ProjectIdentity: input.ProjectIdentity,
	})
	if err != nil {
		return nil, err
	}

	if prj == nil {
		return nil, core.NewNotFoundError("project not found")
	}

	projectDocumentVersion, err := s.ProjectDocumentRepository.GetProjectDocumentVersionBy(projectrepo.GetProjectDocumentVersionByParams{
		ProjectDocumentVersionManagerIdentity: &input.ProjectDocumentVersionManagerIdentity,
		ProjectDocumentVersionIdentity:        input.ProjectDocumentVersionIdentity,
		RelationsInput:                        input.RelationsInput,
	})
	if err != nil {
		return nil, err
	}

	if projectDocumentVersion == nil {
		return nil, core.NewNotFoundError("project document version not found")
	}

	return project.ProjectDocumentVersionToDto(projectDocumentVersion), nil
}
