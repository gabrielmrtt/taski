package storageservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	storagerepo "github.com/gabrielmrtt/taski/internal/storage/repository"
)

type GetFileContentByIdentityService struct {
	FileRepository    storagerepo.UploadedFileRepository
	StorageRepository storagerepo.StorageRepository
}

func NewGetFileContentByIdentityService(
	fileRepository storagerepo.UploadedFileRepository,
	storageRepository storagerepo.StorageRepository,
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
	file, err := s.FileRepository.GetUploadedFileByIdentity(storagerepo.GetUploadedFileByIdentityParams{FileIdentity: input.FileIdentity})
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
