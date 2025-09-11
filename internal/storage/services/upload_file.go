package storage_services

import (
	"path/filepath"
	"strings"

	"github.com/gabrielmrtt/taski/internal/core"
	storage_core "github.com/gabrielmrtt/taski/internal/storage"
	storage_repositories "github.com/gabrielmrtt/taski/internal/storage/repositories"
)

type UploadFileService struct {
	UploadedFileRepository storage_repositories.UploadedFileRepository
	StorageRepository      storage_repositories.StorageRepository
}

func NewUploadFileService(uploadedFileRepository storage_repositories.UploadedFileRepository, storageRepository storage_repositories.StorageRepository) *UploadFileService {
	return &UploadFileService{
		UploadedFileRepository: uploadedFileRepository,
		StorageRepository:      storageRepository,
	}
}

type UploadFileInput struct {
	File       core.FileInput
	Directory  string
	UploadedBy core.Identity
}

func (e *UploadFileService) Execute(input UploadFileInput) (*storage_core.UploadedFile, error) {
	extension := filepath.Ext(input.File.FileName)
	extension = strings.TrimPrefix(extension, ".")

	uploadedFile, err := storage_core.NewUploadedFile(storage_core.NewUploadedFileInput{
		File:                   &input.File.FileName,
		FileDirectory:          &input.Directory,
		FileMimeType:           &input.File.FileMimeType,
		FileExtension:          &extension,
		UserUploadedByIdentity: input.UploadedBy,
	})

	if err != nil {
		return nil, err
	}

	uploadedFile, err = e.UploadedFileRepository.StoreUploadedFile(storage_repositories.StoreUploadedFileParams{UploadedFile: uploadedFile})
	if err != nil {
		return nil, err
	}

	err = e.StorageRepository.StoreFile(input.Directory, input.File.FileName, input.File.FileContent)
	if err != nil {
		return nil, err
	}

	return uploadedFile, nil
}
