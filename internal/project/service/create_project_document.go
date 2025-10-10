package projectservice

import (
	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/internal/project"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
	storagerepo "github.com/gabrielmrtt/taski/internal/storage/repository"
	storageservice "github.com/gabrielmrtt/taski/internal/storage/service"
)

type CreateProjectDocumentService struct {
	ProjectRepository         projectrepo.ProjectRepository
	ProjectDocumentRepository projectrepo.ProjectDocumentRepository
	UploadedFileRepository    storagerepo.UploadedFileRepository
	StorageRepository         storagerepo.StorageRepository
	TransactionRepository     core.TransactionRepository
}

func NewCreateProjectDocumentService(
	projectRepository projectrepo.ProjectRepository,
	projectDocumentRepository projectrepo.ProjectDocumentRepository,
	uploadedFileRepository storagerepo.UploadedFileRepository,
	storageRepository storagerepo.StorageRepository,
	transactionRepository core.TransactionRepository,
) *CreateProjectDocumentService {
	return &CreateProjectDocumentService{
		ProjectRepository:         projectRepository,
		ProjectDocumentRepository: projectDocumentRepository,
		UploadedFileRepository:    uploadedFileRepository,
		StorageRepository:         storageRepository,
		TransactionRepository:     transactionRepository,
	}
}

type CreateProjectDocumentInput struct {
	ProjectIdentity     core.Identity
	Title               string
	Content             string
	Version             string
	Files               []core.FileInput
	UserCreatorIdentity core.Identity
}

func (i CreateProjectDocumentInput) Validate() error {
	var fields []core.InvalidInputErrorField

	if _, err := project.NewProjectDocumentTitle(i.Title); err != nil {
		fields = append(fields, core.InvalidInputErrorField{
			Field: "title",
			Error: err.Error(),
		})
	}

	if _, err := project.NewProjectDocumentContent(i.Content); err != nil {
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

func (s *CreateProjectDocumentService) Execute(input CreateProjectDocumentInput) (*project.ProjectDocumentVersionDto, error) {
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

	projectDocumentVersionManager := &project.ProjectDocumentVersionManager{
		Identity:        core.NewIdentity(project.ProjectDocumentVersionManagerIdentityPrefix),
		ProjectIdentity: input.ProjectIdentity,
		LatestVersion:   nil,
	}

	uploadFileService := storageservice.NewUploadFileService(s.UploadedFileRepository, s.StorageRepository)
	filePath := "projects/" + input.ProjectIdentity.Internal.String() + "/documents/" + projectDocumentVersionManager.Identity.Internal.String() + "/versions/" + input.Version + "/files/"

	var files []project.ProjectDocumentFile = make([]project.ProjectDocumentFile, len(input.Files))
	for i, file := range input.Files {
		uploadedFile, err := uploadFileService.Execute(storageservice.UploadFileInput{
			File:       file,
			Directory:  filePath,
			UploadedBy: input.UserCreatorIdentity,
		})
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		files[i] = project.ProjectDocumentFile{
			Identity:     core.NewIdentity(project.ProjectDocumentVersionIdentityPrefix),
			FileIdentity: uploadedFile.Identity,
		}
	}

	projectDocumentVersion, err := project.NewProjectDocument(project.NewProjectDocumentInput{
		ProjectIdentity:                       input.ProjectIdentity,
		ProjectDocumentVersionManagerIdentity: projectDocumentVersionManager.Identity,
		Title:                                 input.Title,
		Content:                               input.Content,
		Version:                               input.Version,
		Files:                                 files,
		UserCreatorIdentity:                   &input.UserCreatorIdentity,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	projectDocumentVersionManager.LatestVersion = projectDocumentVersion

	_, err = s.ProjectDocumentRepository.StoreProjectDocumentVersionManager(projectrepo.StoreProjectDocumentVersionManagerParams{
		ProjectDocumentVersionManager: projectDocumentVersionManager,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = s.ProjectDocumentRepository.StoreProjectDocumentVersion(projectrepo.StoreProjectDocumentVersionParams{
		ProjectDocumentVersion: projectDocumentVersion,
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

	return project.ProjectDocumentVersionToDto(projectDocumentVersion), nil
}
