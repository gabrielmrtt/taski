package storage_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	storage_repositories "github.com/gabrielmrtt/taski/internal/storage/repositories"
)

type DeleteFileByIdentityService struct {
	UploadedFileRepository storage_repositories.UploadedFileRepository
	StorageRepository      storage_repositories.StorageRepository
}

func NewDeleteFileByIdentityService(uploadedFileRepository storage_repositories.UploadedFileRepository, storageRepository storage_repositories.StorageRepository) *DeleteFileByIdentityService {
	return &DeleteFileByIdentityService{uploadedFileRepository, storageRepository}
}

func (e *DeleteFileByIdentityService) Execute(identity core.Identity) error {
	uploadedFile, err := e.UploadedFileRepository.GetUploadedFileByIdentity(storage_repositories.GetUploadedFileByIdentityParams{FileIdentity: identity})
	if err != nil {
		return err
	}

	if uploadedFile == nil {
		return core.NewNotFoundError("file not found")
	}

	err = e.StorageRepository.DeleteFile(*uploadedFile.FileDirectory, *uploadedFile.File)
	if err != nil {
		return err
	}

	err = e.UploadedFileRepository.DeleteUploadedFile(storage_repositories.DeleteUploadedFileParams{FileIdentity: identity})
	if err != nil {
		return err
	}

	return nil
}
