package projectservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
	storagerepo "github.com/gabrielmrtt/taski/internal/storage/repository"
	storageservice "github.com/gabrielmrtt/taski/internal/storage/service"
)

type DeleteProjectDocumentService struct {
	ProjectRepository         projectrepo.ProjectRepository
	ProjectDocumentRepository projectrepo.ProjectDocumentRepository
	UploadedFileRepository    storagerepo.UploadedFileRepository
	StorageRepository         storagerepo.StorageRepository
	TransactionRepository     core.TransactionRepository
}

func NewDeleteProjectDocumentService(
	projectRepository projectrepo.ProjectRepository,
	projectDocumentRepository projectrepo.ProjectDocumentRepository,
	uploadedFileRepository storagerepo.UploadedFileRepository,
	storageRepository storagerepo.StorageRepository,
	transactionRepository core.TransactionRepository,
) *DeleteProjectDocumentService {
	return &DeleteProjectDocumentService{
		ProjectRepository:         projectRepository,
		ProjectDocumentRepository: projectDocumentRepository,
		UploadedFileRepository:    uploadedFileRepository,
		StorageRepository:         storageRepository,
		TransactionRepository:     transactionRepository,
	}
}

type DeleteProjectDocumentInput struct {
	ProjectIdentity                core.Identity
	ProjectDocumentVersionIdentity core.Identity
}

func (i DeleteProjectDocumentInput) Validate() error {
	return nil
}

func (s *DeleteProjectDocumentService) Execute(input DeleteProjectDocumentInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return err
	}

	s.ProjectRepository.SetTransaction(tx)
	s.ProjectDocumentRepository.SetTransaction(tx)

	prj, err := s.ProjectRepository.GetProjectByIdentity(projectrepo.GetProjectByIdentityParams{
		ProjectIdentity: input.ProjectIdentity,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	if prj == nil {
		tx.Rollback()
		return core.NewNotFoundError("project not found")
	}

	projectDocumentVersion, err := s.ProjectDocumentRepository.GetProjectDocumentVersionBy(projectrepo.GetProjectDocumentVersionByParams{
		ProjectDocumentVersionIdentity: input.ProjectDocumentVersionIdentity,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	if projectDocumentVersion == nil {
		tx.Rollback()
		return core.NewNotFoundError("project document version not found")
	}

	projectDocumentVersionManager, err := s.ProjectDocumentRepository.GetProjectDocumentVersionManagerBy(projectrepo.GetProjectDocumentVersionManagerByParams{
		ProjectDocumentVersionManagerIdentity: projectDocumentVersion.ProjectDocumentVersionManagerIdentity,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	if projectDocumentVersionManager == nil {
		tx.Rollback()
		return core.NewNotFoundError("project document version manager not found")
	}

	versions, err := s.ProjectDocumentRepository.ListProjectDocumentVersionsByProjectDocumentVersionManagerIdentity(projectrepo.ListProjectDocumentVersionsByProjectDocumentVersionManagerIdentityParams{
		ProjectDocumentVersionManagerIdentity: projectDocumentVersionManager.Identity,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	deleteFileService := storageservice.NewDeleteFileByIdentityService(s.UploadedFileRepository, s.StorageRepository)

	for _, version := range versions {
		err = s.ProjectDocumentRepository.DeleteProjectDocumentVersion(projectrepo.DeleteProjectDocumentVersionParams{ProjectDocumentVersionIdentity: version.Identity})
		if err != nil {
			tx.Rollback()
			return err
		}

		for _, file := range version.Document.Files {
			err = deleteFileService.Execute(file.FileIdentity)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	err = s.ProjectDocumentRepository.DeleteProjectDocumentVersionManager(projectrepo.DeleteProjectDocumentVersionManagerParams{ProjectDocumentVersionManagerIdentity: projectDocumentVersionManager.Identity})
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
