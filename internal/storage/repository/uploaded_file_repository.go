package storagerepo

import (
	"github.com/gabrielmrtt/taski/internal/core"
	storage "github.com/gabrielmrtt/taski/internal/storage"
)

type GetUploadedFileByIdentityParams struct {
	FileIdentity core.Identity
}

type StoreUploadedFileParams struct {
	UploadedFile *storage.UploadedFile
}

type DeleteUploadedFileParams struct {
	FileIdentity core.Identity
}

type UploadedFileRepository interface {
	SetTransaction(tx core.Transaction) error

	GetUploadedFileByIdentity(params GetUploadedFileByIdentityParams) (*storage.UploadedFile, error)
	StoreUploadedFile(params StoreUploadedFileParams) (*storage.UploadedFile, error)
	DeleteUploadedFile(params DeleteUploadedFileParams) error
}
