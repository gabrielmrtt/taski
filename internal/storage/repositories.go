package storage_core

import "github.com/gabrielmrtt/taski/internal/core"

type GetUploadedFileByIdentityParams struct {
	FileIdentity core.Identity
}

type StoreUploadedFileParams struct {
	UploadedFile *UploadedFile
}

type DeleteUploadedFileParams struct {
	FileIdentity core.Identity
}

type UploadedFileRepository interface {
	SetTransaction(tx core.Transaction) error

	GetUploadedFileByIdentity(params GetUploadedFileByIdentityParams) (*UploadedFile, error)
	StoreUploadedFile(params StoreUploadedFileParams) (*UploadedFile, error)
	DeleteUploadedFile(params DeleteUploadedFileParams) error
}

type StorageRepository interface {
	GetFile(dir string, filename string) ([]byte, error)
	StoreFile(dir string, filename string, file []byte) error
	DeleteFile(dir string, filename string) error
}
