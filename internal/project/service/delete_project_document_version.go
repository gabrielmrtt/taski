package projectservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
	storagerepo "github.com/gabrielmrtt/taski/internal/storage/repository"
	storageservice "github.com/gabrielmrtt/taski/internal/storage/service"
)

type DeleteProjectDocumentVersionService struct {
	ProjectRepository         projectrepo.ProjectRepository
	ProjectDocumentRepository projectrepo.ProjectDocumentRepository
	UploadedFileRepository    storagerepo.UploadedFileRepository
	StorageRepository         storagerepo.StorageRepository
	TransactionRepository     core.TransactionRepository
}

func NewDeleteProjectDocumentVersionService(
	projectRepository projectrepo.ProjectRepository,
	projectDocumentRepository projectrepo.ProjectDocumentRepository,
	uploadedFileRepository storagerepo.UploadedFileRepository,
	storageRepository storagerepo.StorageRepository,
	transactionRepository core.TransactionRepository,
) *DeleteProjectDocumentVersionService {
	return &DeleteProjectDocumentVersionService{
		ProjectRepository:         projectRepository,
		ProjectDocumentRepository: projectDocumentRepository,
		UploadedFileRepository:    uploadedFileRepository,
		StorageRepository:         storageRepository,
		TransactionRepository:     transactionRepository,
	}
}

type DeleteProjectDocumentVersionInput struct {
	ProjectIdentity                core.Identity
	ProjectDocumentVersionIdentity core.Identity
}

func (i DeleteProjectDocumentVersionInput) Validate() error {
	return nil
}

func (s *DeleteProjectDocumentVersionService) Execute(input DeleteProjectDocumentVersionInput) error {
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

	sortBy := "created_at"
	sortDirection := core.SortDirectionDesc

	versions, err := s.ProjectDocumentRepository.ListProjectDocumentVersionsByProjectDocumentVersionManagerIdentity(projectrepo.ListProjectDocumentVersionsByProjectDocumentVersionManagerIdentityParams{
		ProjectDocumentVersionManagerIdentity: projectDocumentVersionManager.Identity,
		SortInput: core.SortInput{
			By:        &sortBy,
			Direction: &sortDirection,
		},
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	if len(versions) == 0 {
		err = s.ProjectDocumentRepository.DeleteProjectDocumentVersionManager(projectrepo.DeleteProjectDocumentVersionManagerParams{ProjectDocumentVersionManagerIdentity: projectDocumentVersionManager.Identity})
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = s.ProjectDocumentRepository.DeleteProjectDocumentVersion(projectrepo.DeleteProjectDocumentVersionParams{ProjectDocumentVersionIdentity: projectDocumentVersion.Identity})
	if err != nil {
		tx.Rollback()
		return err
	}

	deleteFileService := storageservice.NewDeleteFileByIdentityService(s.UploadedFileRepository, s.StorageRepository)

	for _, file := range projectDocumentVersion.Document.Files {
		err = deleteFileService.Execute(file.FileIdentity)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	latestVersion := versions[0]

	latestVersion.Latest = true

	err = s.ProjectDocumentRepository.UpdateProjectDocumentVersion(projectrepo.UpdateProjectDocumentVersionParams{
		ProjectDocumentVersion: &latestVersion,
	})
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
