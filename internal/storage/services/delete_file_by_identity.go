package storage_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	storage_core "github.com/gabrielmrtt/taski/internal/storage"
)

type DeleteFileByIdentityService struct {
	UploadedFileRepository storage_core.UploadedFileRepository
	StorageRepository      storage_core.StorageRepository
}

func NewDeleteFileByIdentityService(uploadedFileRepository storage_core.UploadedFileRepository, storageRepository storage_core.StorageRepository) *DeleteFileByIdentityService {
	return &DeleteFileByIdentityService{uploadedFileRepository, storageRepository}
}

func (e *DeleteFileByIdentityService) Execute(identity core.Identity) error {
	uploadedFile, err := e.UploadedFileRepository.GetUploadedFileByIdentity(storage_core.GetUploadedFileByIdentityParams{FileIdentity: identity})
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

	err = e.UploadedFileRepository.DeleteUploadedFile(storage_core.DeleteUploadedFileParams{FileIdentity: identity})
	if err != nil {
		return err
	}

	return nil
}
