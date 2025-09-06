package storage_core

import (
	"slices"

	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/pkg/datetimeutils"
)

type UploadedFile struct {
	Identity               core.Identity
	File                   *string
	FileDirectory          *string
	FileMimeType           *string
	FileExtension          *string
	UserUploadedByIdentity core.Identity
	UploadedAt             int64
}

type NewUploadedFileInput struct {
	File                   *string
	FileDirectory          *string
	FileMimeType           *string
	FileExtension          *string
	UserUploadedByIdentity core.Identity
}

func NewUploadedFile(input NewUploadedFileInput) (*UploadedFile, error) {
	return &UploadedFile{
		Identity:               core.NewIdentity("file"),
		File:                   input.File,
		FileDirectory:          input.FileDirectory,
		FileMimeType:           input.FileMimeType,
		FileExtension:          input.FileExtension,
		UserUploadedByIdentity: input.UserUploadedByIdentity,
		UploadedAt:             datetimeutils.EpochNow(),
	}, nil
}

func GetSupportedImageMimeTypes() []string {
	return []string{
		"image/png",
		"image/jpeg",
	}
}

func (u *UploadedFile) IsImage() bool {
	return slices.Contains(GetSupportedImageMimeTypes(), *u.FileMimeType)
}

func GetSupportedVideoMimeTypes() []string {
	return []string{
		"video/mp4",
		"video/mpeg",
		"video/ogg",
		"video/webm",
	}
}

func (u *UploadedFile) IsVideo() bool {
	return slices.Contains(GetSupportedVideoMimeTypes(), *u.FileMimeType)
}

func (u *UploadedFile) IsPdf() bool {
	return *u.FileMimeType == "application/pdf"
}
