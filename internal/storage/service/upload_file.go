package storageservice

import (
	"path/filepath"
	"strings"

	"github.com/gabrielmrtt/taski/internal/core"
	storage "github.com/gabrielmrtt/taski/internal/storage"
	storagerepo "github.com/gabrielmrtt/taski/internal/storage/repository"
)

type UploadFileService struct {
	UploadedFileRepository storagerepo.UploadedFileRepository
	StorageRepository      storagerepo.StorageRepository
}

func NewUploadFileService(uploadedFileRepository storagerepo.UploadedFileRepository, storageRepository storagerepo.StorageRepository) *UploadFileService {
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

func (e *UploadFileService) Execute(input UploadFileInput) (*storage.UploadedFile, error) {
	extension := filepath.Ext(input.File.FileName)
	extension = strings.TrimPrefix(extension, ".")

	uploadedFile, err := storage.NewUploadedFile(storage.NewUploadedFileInput{
		File:                   &input.File.FileName,
		FileDirectory:          &input.Directory,
		FileMimeType:           &input.File.FileMimeType,
		FileExtension:          &extension,
		UserUploadedByIdentity: input.UploadedBy,
	})

	if err != nil {
		return nil, err
	}

	uploadedFile, err = e.UploadedFileRepository.StoreUploadedFile(storagerepo.StoreUploadedFileParams{UploadedFile: uploadedFile})
	if err != nil {
		return nil, err
	}

	err = e.StorageRepository.StoreFile(input.Directory, input.File.FileName, input.File.FileContent)
	if err != nil {
		return nil, err
	}

	return uploadedFile, nil
}
