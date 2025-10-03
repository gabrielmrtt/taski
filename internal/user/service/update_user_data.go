package userservice

import (
	"slices"
	"strings"

	"github.com/gabrielmrtt/taski/internal/core"
	storage "github.com/gabrielmrtt/taski/internal/storage"
	storagerepo "github.com/gabrielmrtt/taski/internal/storage/repository"
	storageservice "github.com/gabrielmrtt/taski/internal/storage/service"
	userrepo "github.com/gabrielmrtt/taski/internal/user/repository"
)

type UpdateUserDataService struct {
	UserRepository         userrepo.UserRepository
	TransactionRepository  core.TransactionRepository
	UploadedFileRepository storagerepo.UploadedFileRepository
	StorageRepository      storagerepo.StorageRepository
}

func NewUpdateUserDataService(
	userRepository userrepo.UserRepository,
	transactionRepository core.TransactionRepository,
	uploadedFileRepository storagerepo.UploadedFileRepository,
	storageRepository storagerepo.StorageRepository,
) *UpdateUserDataService {
	return &UpdateUserDataService{
		UserRepository:         userRepository,
		TransactionRepository:  transactionRepository,
		UploadedFileRepository: uploadedFileRepository,
		StorageRepository:      storageRepository,
	}
}

type UpdateUserDataInput struct {
	UserIdentity   core.Identity
	DisplayName    *string
	About          *string
	ProfilePicture *core.FileInput
}

func (i UpdateUserDataInput) Validate() error {
	var fields []core.InvalidInputErrorField

	if i.DisplayName != nil {
		_, err := core.NewName(*i.DisplayName)
		if err != nil {
			fields = append(fields, core.InvalidInputErrorField{
				Field: "display_name",
				Error: err.Error(),
			})
		}
	}

	if i.About != nil {
		_, err := core.NewDescription(*i.About)
		if err != nil {
			fields = append(fields, core.InvalidInputErrorField{
				Field: "about",
				Error: err.Error(),
			})
		}
	}

	if i.ProfilePicture != nil {
		acceptedMimeTypes := storage.GetSupportedImageMimeTypes()

		if !slices.Contains(acceptedMimeTypes, i.ProfilePicture.FileMimeType) {
			fields = append(fields, core.InvalidInputErrorField{
				Field: "profile_picture",
				Error: "invalid file type. supported mime types are: " + strings.Join(acceptedMimeTypes, ", "),
			})
		}
	}

	if len(fields) > 0 {
		return core.NewInvalidInputError("invalid input", fields)
	}

	return nil
}

func (s *UpdateUserDataService) Execute(input UpdateUserDataInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return core.NewInternalError(err.Error())
	}

	s.UserRepository.SetTransaction(tx)

	usr, err := s.UserRepository.GetUserByIdentity(userrepo.GetUserByIdentityParams{UserIdentity: input.UserIdentity})
	if err != nil {
		tx.Rollback()
		return core.NewInternalError(err.Error())
	}

	if usr == nil {
		tx.Rollback()
		return core.NewNotFoundError("user not found")
	}

	if input.DisplayName != nil {
		err = usr.ChangeUserDataDisplayName(*input.DisplayName)
		if err != nil {
			tx.Rollback()
			return core.NewInternalError(err.Error())
		}
	}

	if input.About != nil {
		err = usr.ChangeUserDataAbout(*input.About)
		if err != nil {
			tx.Rollback()
			return core.NewInternalError(err.Error())
		}
	}

	if input.ProfilePicture != nil {
		uploadedFile, err := storageservice.NewUploadFileService(s.UploadedFileRepository, s.StorageRepository).Execute(storageservice.UploadFileInput{
			File:       *input.ProfilePicture,
			Directory:  "users/" + input.UserIdentity.Internal.String() + "/profile_picture",
			UploadedBy: input.UserIdentity,
		})
		if err != nil {
			tx.Rollback()
			return core.NewInternalError(err.Error())
		}

		if uploadedFile == nil {
			tx.Rollback()
			return core.NewInternalError("failed to upload file")
		}

		usr.ChangeUserDataProfilePicture(&uploadedFile.Identity)
	} else {
		storageservice.NewDeleteFileByIdentityService(s.UploadedFileRepository, s.StorageRepository).Execute(input.UserIdentity)
	}

	err = s.UserRepository.UpdateUser(userrepo.UpdateUserParams{User: usr})
	if err != nil {
		tx.Rollback()
		return core.NewInternalError(err.Error())
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return core.NewInternalError(err.Error())
	}

	return nil
}
