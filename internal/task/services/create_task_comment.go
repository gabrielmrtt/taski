package taskservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
	storagerepo "github.com/gabrielmrtt/taski/internal/storage/repository"
	storageservice "github.com/gabrielmrtt/taski/internal/storage/service"
	"github.com/gabrielmrtt/taski/internal/task"
	taskrepo "github.com/gabrielmrtt/taski/internal/task/repository"
)

type CreateTaskCommentService struct {
	TaskRepository         taskrepo.TaskRepository
	TaskCommentRepository  taskrepo.TaskCommentRepository
	UploadedFileRepository storagerepo.UploadedFileRepository
	StorageRepository      storagerepo.StorageRepository
	ProjectUserRepository  projectrepo.ProjectUserRepository
	TransactionRepository  core.TransactionRepository
}

func NewCreateTaskCommentService(
	taskRepository taskrepo.TaskRepository,
	taskCommentRepository taskrepo.TaskCommentRepository,
	uploadedFileRepository storagerepo.UploadedFileRepository,
	storageRepository storagerepo.StorageRepository,
	projectUserRepository projectrepo.ProjectUserRepository,
	transactionRepository core.TransactionRepository,
) *CreateTaskCommentService {
	return &CreateTaskCommentService{
		TaskRepository:         taskRepository,
		TaskCommentRepository:  taskCommentRepository,
		UploadedFileRepository: uploadedFileRepository,
		StorageRepository:      storageRepository,
		ProjectUserRepository:  projectUserRepository,
		TransactionRepository:  transactionRepository,
	}
}

type CreateTaskCommentInput struct {
	TaskIdentity   core.Identity
	Content        string
	Files          []core.FileInput
	AuthorIdentity core.Identity
}

func (i CreateTaskCommentInput) Validate() error {
	var fields []core.InvalidInputErrorField

	if _, err := core.NewDescription(i.Content); err != nil {
		fields = append(fields, core.InvalidInputErrorField{
			Field: "content",
			Error: err.Error(),
		})
	}

	if len(fields) > 0 {
		return core.NewInvalidInputError("invalid input", fields)
	}

	return nil
}

func (s *CreateTaskCommentService) Execute(input CreateTaskCommentInput) (*task.TaskCommentDto, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return nil, err
	}

	s.TaskRepository.SetTransaction(tx)
	s.TaskCommentRepository.SetTransaction(tx)
	s.UploadedFileRepository.SetTransaction(tx)

	tsk, err := s.TaskRepository.GetTaskByIdentity(taskrepo.GetTaskByIdentityParams{
		TaskIdentity: input.TaskIdentity,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if tsk == nil {
		tx.Rollback()
		return nil, core.NewNotFoundError("task not found")
	}

	usr, err := s.ProjectUserRepository.GetProjectUserByIdentity(projectrepo.GetProjectUserByIdentityParams{
		ProjectIdentity: tsk.ProjectIdentity,
		UserIdentity:    input.AuthorIdentity,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if usr == nil {
		tx.Rollback()
		return nil, core.NewNotFoundError("project user not found")
	}

	comment, err := task.NewTaskComment(task.NewTaskCommentInput{
		TaskIdentity: tsk.Identity,
		Content:      input.Content,
		Author:       &usr.User,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var filePath string = "tasks/" + input.TaskIdentity.Internal.String() + "/comments/" + comment.Identity.Internal.String() + "/files/"

	uploadFileService := storageservice.NewUploadFileService(
		s.UploadedFileRepository,
		s.StorageRepository,
	)

	for _, fileInput := range input.Files {
		uploadedFile, err := uploadFileService.Execute(storageservice.UploadFileInput{
			File:       fileInput,
			Directory:  filePath,
			UploadedBy: usr.User.Identity,
		})
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		comment.AddFile(task.TaskCommentFile{
			Identity:     uploadedFile.Identity,
			FileIdentity: uploadedFile.Identity,
		})
	}

	comment, err = s.TaskCommentRepository.StoreTaskComment(taskrepo.StoreTaskCommentParams{
		TaskComment: comment,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return task.TaskCommentToDto(comment), nil
}
