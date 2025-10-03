package storageservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	storagerepo "github.com/gabrielmrtt/taski/internal/storage/repository"
)

type DeleteFileByIdentityService struct {
	UploadedFileRepository storagerepo.UploadedFileRepository
	StorageRepository      storagerepo.StorageRepository
}

func NewDeleteFileByIdentityService(uploadedFileRepository storagerepo.UploadedFileRepository, storageRepository storagerepo.StorageRepository) *DeleteFileByIdentityService {
	return &DeleteFileByIdentityService{uploadedFileRepository, storageRepository}
}

func (e *DeleteFileByIdentityService) Execute(identity core.Identity) error {
	uploadedFile, err := e.UploadedFileRepository.GetUploadedFileByIdentity(storagerepo.GetUploadedFileByIdentityParams{FileIdentity: identity})
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

	err = e.UploadedFileRepository.DeleteUploadedFile(storagerepo.DeleteUploadedFileParams{FileIdentity: identity})
	if err != nil {
		return err
	}

	return nil
}
