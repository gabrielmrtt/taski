package storage

import (
	"time"

	"github.com/gabrielmrtt/taski/internal/core"
)

type UploadedFile struct {
	Identity               core.Identity
	File                   *string
	FileDirectory          *string
	FileMimeType           *string
	FileExtension          *string
	UserUploadedByIdentity core.Identity
	UploadedAt             *time.Time
}
