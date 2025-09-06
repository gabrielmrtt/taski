package storage_core

import "github.com/gabrielmrtt/taski/internal/core"

type GetUploadedFileByIdentityParams struct {
	Identity core.Identity
}

type UploadedFileRepository interface {
	SetTransaction(tx core.Transaction) error

	GetUploadedFileByIdentity(params GetUploadedFileByIdentityParams) (*UploadedFile, error)
	StoreUploadedFile(uploadedFile *UploadedFile) (*UploadedFile, error)
	DeleteUploadedFile(identity core.Identity) error
}

type StorageRepository interface {
	GetFile(dir string, filename string) ([]byte, error)
	StoreFile(dir string, filename string, file []byte) error
	DeleteFile(dir string, filename string) error
}
