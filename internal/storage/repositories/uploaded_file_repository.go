package storage_repositories

import (
	"github.com/gabrielmrtt/taski/internal/core"
	storage_core "github.com/gabrielmrtt/taski/internal/storage"
)

type GetUploadedFileByIdentityParams struct {
	FileIdentity core.Identity
}

type StoreUploadedFileParams struct {
	UploadedFile *storage_core.UploadedFile
}

type DeleteUploadedFileParams struct {
	FileIdentity core.Identity
}

type UploadedFileRepository interface {
	SetTransaction(tx core.Transaction) error

	GetUploadedFileByIdentity(params GetUploadedFileByIdentityParams) (*storage_core.UploadedFile, error)
	StoreUploadedFile(params StoreUploadedFileParams) (*storage_core.UploadedFile, error)
	DeleteUploadedFile(params DeleteUploadedFileParams) error
}
