package taskservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	storagerepo "github.com/gabrielmrtt/taski/internal/storage/repository"
	storageservice "github.com/gabrielmrtt/taski/internal/storage/service"
	"github.com/gabrielmrtt/taski/internal/task"
	taskrepo "github.com/gabrielmrtt/taski/internal/task/repository"
)

type UpdateTaskCommentService struct {
	TaskCommentRepository  taskrepo.TaskCommentRepository
	TaskRepository         taskrepo.TaskRepository
	UploadedFileRepository storagerepo.UploadedFileRepository
	StorageRepository      storagerepo.StorageRepository
	TransactionRepository  core.TransactionRepository
}

func NewUpdateTaskCommentService(
	taskCommentRepository taskrepo.TaskCommentRepository,
	taskRepository taskrepo.TaskRepository,
	uploadedFileRepository storagerepo.UploadedFileRepository,
	storageRepository storagerepo.StorageRepository,
	transactionRepository core.TransactionRepository,
) *UpdateTaskCommentService {
	return &UpdateTaskCommentService{
		TaskCommentRepository:  taskCommentRepository,
		TaskRepository:         taskRepository,
		UploadedFileRepository: uploadedFileRepository,
		StorageRepository:      storageRepository,
		TransactionRepository:  transactionRepository,
	}
}

type UpdateTaskCommentInput struct {
	TaskIdentity        core.Identity
	TaskCommentIdentity core.Identity
	Content             *string
	Files               []core.FileInput
	UserEditorIdentity  core.Identity
}

func (i UpdateTaskCommentInput) Validate() error {
	var fields []core.InvalidInputErrorField

	if i.Content != nil {
		if _, err := core.NewDescription(*i.Content); err != nil {
			fields = append(fields, core.InvalidInputErrorField{
				Field: "content",
				Error: err.Error(),
			})
		}
	}

	if len(fields) > 0 {
		return core.NewInvalidInputError("invalid input", fields)
	}

	return nil
}

func (s *UpdateTaskCommentService) Execute(input UpdateTaskCommentInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return err
	}

	s.TaskCommentRepository.SetTransaction(tx)
	s.TaskRepository.SetTransaction(tx)
	s.UploadedFileRepository.SetTransaction(tx)

	comment, err := s.TaskCommentRepository.GetTaskCommentByIdentity(taskrepo.GetTaskCommentByIdentityParams{
		TaskCommentIdentity: input.TaskCommentIdentity,
		TaskIdentity:        &input.TaskIdentity,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	if comment == nil {
		tx.Rollback()
		return core.NewNotFoundError("task comment not found")
	}

	if input.Content != nil {
		err = comment.ChangeContent(*input.Content)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(input.Files) > 0 {
		deleteFileService := storageservice.NewDeleteFileByIdentityService(s.UploadedFileRepository, s.StorageRepository)
		uploadFileService := storageservice.NewUploadFileService(s.UploadedFileRepository, s.StorageRepository)

		for _, file := range comment.Files {
			err = deleteFileService.Execute(file.FileIdentity)
			if err != nil {
				tx.Rollback()
				return err
			}
		}

		comment.ClearAllFiles()

		filePath := "tasks/" + input.TaskIdentity.Internal.String() + "/comments/" + comment.Identity.Internal.String() + "/files"

		for _, file := range input.Files {
			uploadedFile, err := uploadFileService.Execute(storageservice.UploadFileInput{
				File:       file,
				Directory:  filePath,
				UploadedBy: input.UserEditorIdentity,
			})
			if err != nil {
				tx.Rollback()
				return err
			}

			comment.AddFile(task.TaskCommentFile{
				Identity:     uploadedFile.Identity,
				FileIdentity: uploadedFile.Identity,
			})
		}
	}

	err = s.TaskCommentRepository.UpdateTaskComment(taskrepo.UpdateTaskCommentParams{
		TaskComment: comment,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
