package storage_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	storage_core "github.com/gabrielmrtt/taski/internal/storage"
)

type GetFileContentByIdentityService struct {
	FileRepository    storage_core.UploadedFileRepository
	StorageRepository storage_core.StorageRepository
}

func NewGetFileContentByIdentityService(
	fileRepository storage_core.UploadedFileRepository,
	storageRepository storage_core.StorageRepository,
) *GetFileContentByIdentityService {
	return &GetFileContentByIdentityService{
		FileRepository:    fileRepository,
		StorageRepository: storageRepository,
	}
}

type GetFileContentByIdentityInput struct {
	FileIdentity core.Identity
}

func (s *GetFileContentByIdentityService) Execute(input GetFileContentByIdentityInput) (*core.FileInput, error) {
	file, err := s.FileRepository.GetUploadedFileByIdentity(storage_core.GetUploadedFileByIdentityParams{
		Identity: input.FileIdentity,
	})
	if err != nil {
		return nil, err
	}

	if file == nil {
		return nil, core.NewNotFoundError("file not found")
	}

	fileContent, err := s.StorageRepository.GetFile(*file.FileDirectory, *file.File)
	if err != nil {
		return nil, err
	}

	if fileContent == nil {
		return nil, core.NewNotFoundError("storage not found")
	}

	return &core.FileInput{
		FileName:     *file.File,
		FileContent:  fileContent,
		FileMimeType: *file.FileMimeType,
	}, nil
}
