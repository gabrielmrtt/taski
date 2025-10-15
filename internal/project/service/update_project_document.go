package projectservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/project"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
	storagerepo "github.com/gabrielmrtt/taski/internal/storage/repository"
	storageservice "github.com/gabrielmrtt/taski/internal/storage/service"
)

type UpdateProjectDocumentService struct {
	ProjectRepository         projectrepo.ProjectRepository
	ProjectDocumentRepository projectrepo.ProjectDocumentRepository
	UploadedFileRepository    storagerepo.UploadedFileRepository
	StorageRepository         storagerepo.StorageRepository
	TransactionRepository     core.TransactionRepository
}

func NewUpdateProjectDocumentService(
	projectRepository projectrepo.ProjectRepository,
	projectDocumentRepository projectrepo.ProjectDocumentRepository,
	uploadedFileRepository storagerepo.UploadedFileRepository,
	storageRepository storagerepo.StorageRepository,
	transactionRepository core.TransactionRepository,
) *UpdateProjectDocumentService {
	return &UpdateProjectDocumentService{
		ProjectRepository:         projectRepository,
		ProjectDocumentRepository: projectDocumentRepository,
		UploadedFileRepository:    uploadedFileRepository,
		StorageRepository:         storageRepository,
		TransactionRepository:     transactionRepository,
	}
}

type UpdateProjectDocumentInput struct {
	ProjectIdentity                       core.Identity
	ProjectDocumentVersionManagerIdentity core.Identity
	ProjectDocumentVersionIdentity        core.Identity
	Version                               *string
	Title                                 *string
	Content                               *string
	Files                                 []core.FileInput
	UserEditorIdentity                    core.Identity
}

func (i UpdateProjectDocumentInput) Validate() error {
	var fields []core.InvalidInputErrorField

	if i.Title != nil {
		_, err := project.NewProjectDocumentTitle(*i.Title)
		if err != nil {
			fields = append(fields, core.InvalidInputErrorField{
				Field: "title",
				Error: err.Error(),
			})
		}
	}

	if i.Content != nil {
		_, err := project.NewProjectDocumentContent(*i.Content)
		if err != nil {
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

func (s *UpdateProjectDocumentService) Execute(input UpdateProjectDocumentInput) (*project.ProjectDocumentVersionDto, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	tx, err := s.TransactionRepository.BeginTransaction()
	if err != nil {
		return nil, err
	}

	s.ProjectRepository.SetTransaction(tx)
	s.ProjectDocumentRepository.SetTransaction(tx)

	prj, err := s.ProjectRepository.GetProjectByIdentity(projectrepo.GetProjectByIdentityParams{
		ProjectIdentity: input.ProjectIdentity,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if prj == nil {
		tx.Rollback()
		return nil, core.NewNotFoundError("project not found")
	}

	projectDocumentVersion, err := s.ProjectDocumentRepository.GetProjectDocumentVersionBy(projectrepo.GetProjectDocumentVersionByParams{
		ProjectDocumentVersionManagerIdentity: &input.ProjectDocumentVersionManagerIdentity,
		ProjectDocumentVersionIdentity:        input.ProjectDocumentVersionIdentity,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if projectDocumentVersion == nil {
		tx.Rollback()
		return nil, core.NewNotFoundError("project document version not found")
	}

	if input.Version != nil {
		projectDocumentVersion = projectDocumentVersion.NewVersion(*input.Version)
	}

	if input.Title != nil {
		err = projectDocumentVersion.ChangeTitle(*input.Title, &input.UserEditorIdentity)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if input.Content != nil {
		err = projectDocumentVersion.ChangeContent(*input.Content, &input.UserEditorIdentity)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if len(input.Files) > 0 {
		uploadFileService := storageservice.NewUploadFileService(s.UploadedFileRepository, s.StorageRepository)
		deleteFileService := storageservice.NewDeleteFileByIdentityService(s.UploadedFileRepository, s.StorageRepository)

		if input.Version == nil {
			for _, file := range projectDocumentVersion.Document.Files {
				err = deleteFileService.Execute(file.FileIdentity)
				if err != nil {
					tx.Rollback()
					return nil, err
				}
			}
		}

		projectDocumentVersion.ClearAllFiles()

		filePath := "projects/" + input.ProjectIdentity.Internal.String() + "/documents/" + projectDocumentVersion.Identity.Internal.String() + "/versions/" + projectDocumentVersion.Identity.Internal.String() + "/files"

		for _, file := range input.Files {
			uploadedFile, err := uploadFileService.Execute(storageservice.UploadFileInput{
				File:       file,
				Directory:  filePath,
				UploadedBy: input.UserEditorIdentity,
			})
			if err != nil {
				tx.Rollback()
				return nil, err
			}

			projectDocumentVersion.AddFile(project.ProjectDocumentFile{
				Identity:     core.NewIdentity(project.ProjectDocumentVersionIdentityPrefix),
				FileIdentity: uploadedFile.Identity,
			})
		}
	}

	if input.Version != nil {
		_, err = s.ProjectDocumentRepository.StoreProjectDocumentVersion(projectrepo.StoreProjectDocumentVersionParams{
			ProjectDocumentVersion: projectDocumentVersion,
		})
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	} else {
		err = s.ProjectDocumentRepository.UpdateProjectDocumentVersion(projectrepo.UpdateProjectDocumentVersionParams{
			ProjectDocumentVersion: projectDocumentVersion,
		})
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return project.ProjectDocumentVersionToDto(projectDocumentVersion), nil
}
