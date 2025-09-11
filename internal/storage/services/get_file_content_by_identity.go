package storage_services

import (
	"github.com/gabrielmrtt/taski/internal/core"
	storage_repositories "github.com/gabrielmrtt/taski/internal/storage/repositories"
)

type GetFileContentByIdentityService struct {
	FileRepository    storage_repositories.UploadedFileRepository
	StorageRepository storage_repositories.StorageRepository
}

func NewGetFileContentByIdentityService(
	fileRepository storage_repositories.UploadedFileRepository,
	storageRepository storage_repositories.StorageRepository,
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
	file, err := s.FileRepository.GetUploadedFileByIdentity(storage_repositories.GetUploadedFileByIdentityParams{FileIdentity: input.FileIdentity})
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
