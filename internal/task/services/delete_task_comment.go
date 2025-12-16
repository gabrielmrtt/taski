package taskservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
	storagerepo "github.com/gabrielmrtt/taski/internal/storage/repository"
	storageservice "github.com/gabrielmrtt/taski/internal/storage/service"
	"github.com/gabrielmrtt/taski/internal/task"
	taskrepo "github.com/gabrielmrtt/taski/internal/task/repository"
)

type DeleteTaskCommentService struct {
	TaskCommentRepository  taskrepo.TaskCommentRepository
	TaskRepository         taskrepo.TaskRepository
	UploadedFileRepository storagerepo.UploadedFileRepository
	StorageRepository      storagerepo.StorageRepository
	TaskActionRepository   taskrepo.TaskActionRepository
	ProjectUserRepository  projectrepo.ProjectUserRepository
	TransactionRepository  core.TransactionRepository
}

func NewDeleteTaskCommentService(
	taskCommentRepository taskrepo.TaskCommentRepository,
	taskRepository taskrepo.TaskRepository,
	projectUserRepository projectrepo.ProjectUserRepository,
	uploadedFileRepository storagerepo.UploadedFileRepository,
	storageRepository storagerepo.StorageRepository,
	taskActionRepository taskrepo.TaskActionRepository,
	transactionRepository core.TransactionRepository,
) *DeleteTaskCommentService {
	return &DeleteTaskCommentService{
		TaskCommentRepository:  taskCommentRepository,
		TaskRepository:         taskRepository,
		UploadedFileRepository: uploadedFileRepository,
		StorageRepository:      storageRepository,
		TaskActionRepository:   taskActionRepository,
		ProjectUserRepository:  projectUserRepository,
		TransactionRepository:  transactionRepository,
	}
}

type DeleteTaskCommentInput struct {
	TaskIdentity        core.Identity
	TaskCommentIdentity core.Identity
	UserDeleterIdentity core.Identity
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
	s.TaskActionRepository.SetTransaction(tx)
	s.ProjectUserRepository.SetTransaction(tx)

	tsk, err := s.TaskRepository.GetTaskByIdentity(taskrepo.GetTaskByIdentityParams{
		TaskIdentity: input.TaskIdentity,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	if tsk == nil {
		tx.Rollback()
		return core.NewNotFoundError("task not found")
	}

	userDeleter, err := s.ProjectUserRepository.GetProjectUserByIdentity(projectrepo.GetProjectUserByIdentityParams{
		ProjectIdentity: tsk.ProjectIdentity,
		UserIdentity:    input.UserDeleterIdentity,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	if userDeleter == nil {
		tx.Rollback()
		return core.NewNotFoundError("project user deleter not found")
	}

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

	taskAction := tsk.RegisterAction(task.TaskActionTypeDeleteComment, &userDeleter.User)
	_, err = s.TaskActionRepository.StoreTaskAction(taskrepo.StoreTaskActionParams{
		TaskAction: &taskAction,
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
