package taskservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	storagerepo "github.com/gabrielmrtt/taski/internal/storage/repository"
	storageservice "github.com/gabrielmrtt/taski/internal/storage/service"
	taskrepo "github.com/gabrielmrtt/taski/internal/task/repository"
)

type DeleteTaskCommentService struct {
	TaskCommentRepository  taskrepo.TaskCommentRepository
	TaskRepository         taskrepo.TaskRepository
	UploadedFileRepository storagerepo.UploadedFileRepository
	StorageRepository      storagerepo.StorageRepository
	TransactionRepository  core.TransactionRepository
}

func NewDeleteTaskCommentService(
	taskCommentRepository taskrepo.TaskCommentRepository,
	taskRepository taskrepo.TaskRepository,
	uploadedFileRepository storagerepo.UploadedFileRepository,
	storageRepository storagerepo.StorageRepository,
	transactionRepository core.TransactionRepository,
) *DeleteTaskCommentService {
	return &DeleteTaskCommentService{
		TaskCommentRepository:  taskCommentRepository,
		TaskRepository:         taskRepository,
		UploadedFileRepository: uploadedFileRepository,
		StorageRepository:      storageRepository,
		TransactionRepository:  transactionRepository,
	}
}

type DeleteTaskCommentInput struct {
	TaskIdentity        core.Identity
	TaskCommentIdentity core.Identity
}

func (i DeleteTaskCommentInput) Validate() error { return nil }

func (s *DeleteTaskCommentService) Execute(input DeleteTaskCommentInput) error {
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

	deleteFileService := storageservice.NewDeleteFileByIdentityService(s.UploadedFileRepository, s.StorageRepository)

	for _, file := range comment.Files {
		err = deleteFileService.Execute(file.FileIdentity)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = s.TaskCommentRepository.DeleteTaskComment(taskrepo.DeleteTaskCommentParams{
		TaskCommentIdentity: input.TaskCommentIdentity,
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
