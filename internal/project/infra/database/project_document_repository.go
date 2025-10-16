package projectdatabase

import (
	"context"
	"database/sql"

	"github.com/gabrielmrtt/taski/internal/core"
	coredatabase "github.com/gabrielmrtt/taski/internal/core/database"
	"github.com/gabrielmrtt/taski/internal/project"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
	"github.com/gabrielmrtt/taski/internal/storage"
	storagedatabase "github.com/gabrielmrtt/taski/internal/storage/infra/database"
	"github.com/gabrielmrtt/taski/internal/user"
	userdatabase "github.com/gabrielmrtt/taski/internal/user/infra/database"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ProjectDocumentVersionManagerTable struct {
	bun.BaseModel `bun:"table:project_document_version_manager,alias:project_document_version_manager"`

	InternalId        string `bun:"internal_id,pk,notnull,type:uuid"`
	PublicId          string `bun:"public_id,notnull,type:varchar(510)"`
	ProjectInternalId string `bun:"project_internal_id,notnull,type:uuid"`

	Project       *ProjectTable                `bun:"rel:has-one,join:project_internal_id=internal_id"`
	LatestVersion *ProjectDocumentVersionTable `bun:"rel:has-one,join:internal_id=project_document_version_manager_internal_id"`
}

func (p *ProjectDocumentVersionManagerTable) ToEntity() *project.ProjectDocumentVersionManager {
	var latestVersion *project.ProjectDocumentVersion
	if p.LatestVersion != nil {
		latestVersion = p.LatestVersion.ToEntity()
	}

	return &project.ProjectDocumentVersionManager{
		Identity:        core.NewIdentityFromInternal(uuid.MustParse(p.InternalId), project.ProjectDocumentVersionManagerIdentityPrefix),
		ProjectIdentity: core.NewIdentityFromInternal(uuid.MustParse(p.ProjectInternalId), project.ProjectIdentityPrefix),
		LatestVersion:   latestVersion,
	}
}

type ProjectDocumentVersionTable struct {
	bun.BaseModel `bun:"table:project_document_version,alias:project_document_version"`

	InternalId                              string  `bun:"internal_id,pk,notnull,type:uuid"`
	PublicId                                string  `bun:"public_id,notnull,type:varchar(510)"`
	ProjectDocumentVersionManagerInternalId string  `bun:"project_document_version_manager_internal_id,notnull,type:uuid"`
	Version                                 string  `bun:"version,notnull,type:varchar(255)"`
	Title                                   string  `bun:"title,notnull,type:varchar(255)"`
	Content                                 string  `bun:"content,notnull,type:text"`
	UserCreatorInternalId                   string  `bun:"user_creator_internal_id,notnull,type:uuid"`
	UserEditorInternalId                    *string `bun:"user_editor_internal_id,type:uuid"`
	Latest                                  bool    `bun:"latest,notnull,type:boolean"`
	CreatedAt                               int64   `bun:"created_at,notnull,type:bigint"`
	UpdatedAt                               *int64  `bun:"updated_at,type:bigint"`

	ProjectDocumentVersionManager *ProjectDocumentVersionManagerTable `bun:"rel:has-one,join:project_document_version_manager_internal_id=internal_id"`
	ProjectDocumentFiles          []*ProjectDocumentFileTable         `bun:"rel:has-many,join:internal_id=project_document_version_internal_id"`
	UserCreator                   *userdatabase.UserTable             `bun:"rel:has-one,join:user_creator_internal_id=internal_id"`
	UserEditor                    *userdatabase.UserTable             `bun:"rel:has-one,join:user_editor_internal_id=internal_id"`
}

type ProjectDocumentFileTable struct {
	bun.BaseModel `bun:"table:project_document_file,alias:project_document_file"`

	InternalId                       string `bun:"internal_id,pk,notnull,type:uuid"`
	ProjectDocumentVersionInternalId string `bun:"project_document_version_internal_id,notnull,type:uuid"`
	FileInternalId                   string `bun:"file_internal_id,notnull,type:uuid"`

	ProjectDocumentVersion *ProjectDocumentVersionTable       `bun:"rel:has-one,join:project_document_version_internal_id=internal_id"`
	File                   *storagedatabase.UploadedFileTable `bun:"rel:has-one,join:file_internal_id=internal_id"`
}

func (p *ProjectDocumentFileTable) ToEntity() *project.ProjectDocumentFile {
	return &project.ProjectDocumentFile{
		Identity:     core.NewIdentityWithoutPublicFromInternal(uuid.MustParse(p.InternalId)),
		FileIdentity: core.NewIdentityFromInternal(uuid.MustParse(p.FileInternalId), storage.UploadedFileIdentityPrefix),
	}
}

func (p *ProjectDocumentVersionTable) ToEntity() *project.ProjectDocumentVersion {
	var userCreatorIdentity core.Identity = core.NewIdentityFromInternal(uuid.MustParse(p.UserCreatorInternalId), user.UserIdentityPrefix)
	var userEditorIdentity *core.Identity = nil

	if p.UserEditorInternalId != nil {
		identity := core.NewIdentityFromInternal(uuid.MustParse(*p.UserEditorInternalId), user.UserIdentityPrefix)
		userEditorIdentity = &identity
	}

	var files []project.ProjectDocumentFile = make([]project.ProjectDocumentFile, len(p.ProjectDocumentFiles))
	for i, file := range p.ProjectDocumentFiles {
		files[i] = *file.ToEntity()
	}

	return &project.ProjectDocumentVersion{
		Identity:                              core.NewIdentityFromInternal(uuid.MustParse(p.InternalId), project.ProjectDocumentVersionIdentityPrefix),
		ProjectDocumentVersionManagerIdentity: core.NewIdentityFromInternal(uuid.MustParse(p.ProjectDocumentVersionManagerInternalId), project.ProjectDocumentVersionManagerIdentityPrefix),
		Version:                               p.Version,
		Document: project.ProjectDocument{
			Identity: core.NewIdentityWithoutPublic(),
			Title:    p.Title,
			Content:  p.Content,
			Files:    files,
		},
		UserCreatorIdentity: &userCreatorIdentity,
		UserEditorIdentity:  userEditorIdentity,
		Latest:              p.Latest,
		Timestamps: core.Timestamps{
			CreatedAt: &p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		},
	}
}

type ProjectDocumentBunRepository struct {
	db *bun.DB
	tx *coredatabase.TransactionBun
}

func NewProjectDocumentBunRepository(connection *bun.DB) *ProjectDocumentBunRepository {
	return &ProjectDocumentBunRepository{db: connection, tx: nil}
}

func (r *ProjectDocumentBunRepository) SetTransaction(tx core.Transaction) error {
	if r.tx != nil && !r.tx.IsClosed() {
		return nil
	}

	r.tx = tx.(*coredatabase.TransactionBun)
	return nil
}

func (r *ProjectDocumentBunRepository) applyProjectDocumentVersionManagerFilters(selectQuery *bun.SelectQuery, filters projectrepo.ProjectDocumentVersionManagerFilters) *bun.SelectQuery {
	if filters.ProjectIdentity != nil {
		selectQuery = selectQuery.Where("project_internal_id = ?", filters.ProjectIdentity.Internal.String())
	}

	if filters.Title != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "title", filters.Title)
	}

	if filters.CreatedAt != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "created_at", filters.CreatedAt)
	}

	return selectQuery
}

func (r *ProjectDocumentBunRepository) applyProjectDocumentVersionFilters(selectQuery *bun.SelectQuery, filters projectrepo.ProjectDocumentVersionFilters) *bun.SelectQuery {
	if filters.ProjectDocumentVersionManagerIdentity != nil {
		selectQuery = selectQuery.Where("project_document_version_manager_internal_id = ?", filters.ProjectDocumentVersionManagerIdentity.Internal.String())
	}

	if filters.Version != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "version", filters.Version)
	}

	if filters.CreatedAt != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "created_at", filters.CreatedAt)
	}

	return selectQuery
}

func (r *ProjectDocumentBunRepository) GetProjectDocumentVersionManagerBy(params projectrepo.GetProjectDocumentVersionManagerByParams) (*project.ProjectDocumentVersionManager, error) {
	var projectDocumentVersionManager *ProjectDocumentVersionManagerTable = new(ProjectDocumentVersionManagerTable)
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(projectDocumentVersionManager)
	selectQuery = selectQuery.Relation("LatestVersion.ProjectDocumentFiles", func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where("latest = ?", true)
	})
	selectQuery = selectQuery.Where("internal_id = ?", params.ProjectDocumentVersionManagerIdentity.Internal.String())
	err := selectQuery.Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	if projectDocumentVersionManager.InternalId == "" {
		return nil, nil
	}

	return projectDocumentVersionManager.ToEntity(), nil
}

func (r *ProjectDocumentBunRepository) GetProjectDocumentVersionBy(params projectrepo.GetProjectDocumentVersionByParams) (*project.ProjectDocumentVersion, error) {
	var projectDocumentVersion *ProjectDocumentVersionTable = new(ProjectDocumentVersionTable)
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(projectDocumentVersion)
	selectQuery = selectQuery.Where("internal_id = ?", params.ProjectDocumentVersionIdentity.Internal.String())
	err := selectQuery.Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	if projectDocumentVersion.InternalId == "" {
		return nil, nil
	}

	return projectDocumentVersion.ToEntity(), nil
}

func (r *ProjectDocumentBunRepository) ListProjectDocumentVersionsByProjectDocumentVersionManagerIdentity(params projectrepo.ListProjectDocumentVersionsByProjectDocumentVersionManagerIdentityParams) ([]project.ProjectDocumentVersion, error) {
	var projectDocumentVersions []ProjectDocumentVersionTable = make([]ProjectDocumentVersionTable, 0)
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(&projectDocumentVersions)
	selectQuery = selectQuery.Where("project_document_version_manager_internal_id = ?", params.ProjectDocumentVersionManagerIdentity.Internal.String())
	err := selectQuery.Scan(context.Background())
	if err != nil {
		return nil, err
	}

	var projectDocumentVersionsEntities []project.ProjectDocumentVersion = make([]project.ProjectDocumentVersion, 0)
	for _, projectDocumentVersion := range projectDocumentVersions {
		projectDocumentVersionsEntities = append(projectDocumentVersionsEntities, *projectDocumentVersion.ToEntity())
	}

	return projectDocumentVersionsEntities, nil
}

func (r *ProjectDocumentBunRepository) PaginateProjectDocumentVersionManagersBy(params projectrepo.PaginateProjectDocumentVersionManagersByParams) (*core.PaginationOutput[project.ProjectDocumentVersionManager], error) {
	var projectDocumentVersionManagers []ProjectDocumentVersionManagerTable = make([]ProjectDocumentVersionManagerTable, 0)
	var selectQuery *bun.SelectQuery
	var perPage int = 10
	var page int = 1

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	if params.Pagination.PerPage != nil {
		perPage = *params.Pagination.PerPage
	}

	if params.Pagination.Page != nil {
		page = *params.Pagination.Page
	}

	selectQuery = selectQuery.Model(&projectDocumentVersionManagers)
	selectQuery = selectQuery.Relation("LatestVersion", func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where("latest = ?", true)
	})
	selectQuery = r.applyProjectDocumentVersionManagerFilters(selectQuery, params.Filters)
	countBeforePagination, err := selectQuery.Count(context.Background())
	if err != nil {
		return nil, err
	}

	selectQuery = coredatabase.ApplySort(selectQuery, params.SortInput)
	selectQuery = coredatabase.ApplyPagination(selectQuery, params.Pagination)
	err = selectQuery.Scan(context.Background())
	if err != nil {
		return nil, err
	}

	var projectDocumentVersionManagersEntities []project.ProjectDocumentVersionManager = make([]project.ProjectDocumentVersionManager, 0)
	for _, projectDocumentVersionManager := range projectDocumentVersionManagers {
		projectDocumentVersionManagersEntities = append(projectDocumentVersionManagersEntities, *projectDocumentVersionManager.ToEntity())
	}

	return &core.PaginationOutput[project.ProjectDocumentVersionManager]{
		Data:    projectDocumentVersionManagersEntities,
		Page:    page,
		HasMore: core.HasMorePages(page, countBeforePagination, perPage),
		Total:   countBeforePagination,
	}, nil
}

func (r *ProjectDocumentBunRepository) PaginateProjectDocumentVersionsBy(params projectrepo.PaginateProjectDocumentVersionsByParams) (*core.PaginationOutput[project.ProjectDocumentVersion], error) {
	var projectDocumentVersions []ProjectDocumentVersionTable = make([]ProjectDocumentVersionTable, 0)
	var selectQuery *bun.SelectQuery
	var perPage int = 10
	var page int = 1

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	if params.Pagination.PerPage != nil {
		perPage = *params.Pagination.PerPage
	}

	if params.Pagination.Page != nil {
		page = *params.Pagination.Page
	}

	selectQuery = selectQuery.Model(&projectDocumentVersions)
	selectQuery = selectQuery.Relation("ProjectDocumentFiles")
	selectQuery = r.applyProjectDocumentVersionFilters(selectQuery, params.Filters)
	countBeforePagination, err := selectQuery.Count(context.Background())
	if err != nil {
		return nil, err
	}

	selectQuery = coredatabase.ApplySort(selectQuery, params.SortInput)
	selectQuery = coredatabase.ApplyPagination(selectQuery, params.Pagination)
	err = selectQuery.Scan(context.Background())
	if err != nil {
		return nil, err
	}

	var projectDocumentVersionsEntities []project.ProjectDocumentVersion = make([]project.ProjectDocumentVersion, 0)
	for _, projectDocumentVersion := range projectDocumentVersions {
		projectDocumentVersionsEntities = append(projectDocumentVersionsEntities, *projectDocumentVersion.ToEntity())
	}

	return &core.PaginationOutput[project.ProjectDocumentVersion]{
		Data:    projectDocumentVersionsEntities,
		Page:    page,
		HasMore: core.HasMorePages(page, countBeforePagination, perPage),
		Total:   countBeforePagination,
	}, nil
}

func (r *ProjectDocumentBunRepository) StoreProjectDocumentVersionManager(params projectrepo.StoreProjectDocumentVersionManagerParams) (*project.ProjectDocumentVersionManager, error) {
	var tx bun.Tx
	var shouldCommit bool = false

	if r.tx != nil && !r.tx.IsClosed() {
		tx = *r.tx.Tx
	} else {
		var err error
		tx, err = r.db.BeginTx(context.Background(), nil)
		shouldCommit = true

		if err != nil {
			return nil, err
		}
	}

	projectDocumentVersionManagerTable := &ProjectDocumentVersionManagerTable{
		InternalId:        params.ProjectDocumentVersionManager.Identity.Internal.String(),
		PublicId:          params.ProjectDocumentVersionManager.Identity.Public,
		ProjectInternalId: params.ProjectDocumentVersionManager.ProjectIdentity.Internal.String(),
	}

	_, err := tx.NewInsert().Model(projectDocumentVersionManagerTable).Exec(context.Background())
	if err != nil {
		return nil, err
	}

	if shouldCommit {
		err = tx.Commit()
		if err != nil {
			return nil, err
		}
	}

	return params.ProjectDocumentVersionManager, nil
}

func (r *ProjectDocumentBunRepository) DeleteProjectDocumentVersionManager(params projectrepo.DeleteProjectDocumentVersionManagerParams) error {
	var tx bun.Tx
	var shouldCommit bool = false

	if r.tx != nil && !r.tx.IsClosed() {
		tx = *r.tx.Tx
	} else {
		var err error
		tx, err = r.db.BeginTx(context.Background(), nil)
		shouldCommit = true

		if err != nil {
			return err
		}
	}

	_, err := tx.NewDelete().Model(&ProjectDocumentVersionManagerTable{}).Where("internal_id = ?", params.ProjectDocumentVersionManagerIdentity.Internal.String()).Exec(context.Background())
	if err != nil {
		return err
	}

	if shouldCommit {
		err = tx.Commit()
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ProjectDocumentBunRepository) StoreProjectDocumentVersion(params projectrepo.StoreProjectDocumentVersionParams) (*project.ProjectDocumentVersion, error) {
	var tx bun.Tx
	var shouldCommit bool = false

	if r.tx != nil && !r.tx.IsClosed() {
		tx = *r.tx.Tx
	} else {
		var err error
		tx, err = r.db.BeginTx(context.Background(), nil)
		shouldCommit = true

		if err != nil {
			return nil, err
		}
	}

	var userEditorInternalId *string
	if params.ProjectDocumentVersion.UserEditorIdentity != nil {
		identity := params.ProjectDocumentVersion.UserEditorIdentity.Internal.String()
		userEditorInternalId = &identity
	}

	projectDocumentVersionTable := &ProjectDocumentVersionTable{
		InternalId:                              params.ProjectDocumentVersion.Identity.Internal.String(),
		PublicId:                                params.ProjectDocumentVersion.Identity.Public,
		ProjectDocumentVersionManagerInternalId: params.ProjectDocumentVersion.ProjectDocumentVersionManagerIdentity.Internal.String(),
		Version:                                 params.ProjectDocumentVersion.Version,
		Title:                                   params.ProjectDocumentVersion.Document.Title,
		Content:                                 params.ProjectDocumentVersion.Document.Content,
		UserCreatorInternalId:                   params.ProjectDocumentVersion.UserCreatorIdentity.Internal.String(),
		UserEditorInternalId:                    userEditorInternalId,
		Latest:                                  params.ProjectDocumentVersion.Latest,
		CreatedAt:                               *params.ProjectDocumentVersion.Timestamps.CreatedAt,
		UpdatedAt:                               params.ProjectDocumentVersion.Timestamps.UpdatedAt,
	}

	_, err := tx.NewInsert().Model(projectDocumentVersionTable).Exec(context.Background())
	if err != nil {
		return nil, err
	}

	for _, file := range params.ProjectDocumentVersion.Document.Files {
		projectDocumentFileTable := &ProjectDocumentFileTable{
			InternalId:                       file.Identity.Internal.String(),
			ProjectDocumentVersionInternalId: params.ProjectDocumentVersion.Identity.Internal.String(),
			FileInternalId:                   file.FileIdentity.Internal.String(),
		}

		_, err := tx.NewInsert().Model(projectDocumentFileTable).Exec(context.Background())
		if err != nil {
			return nil, err
		}
	}

	if shouldCommit {
		err = tx.Commit()
		if err != nil {
			return nil, err
		}
	}

	return params.ProjectDocumentVersion, nil
}

func (r *ProjectDocumentBunRepository) UpdateProjectDocumentVersion(params projectrepo.UpdateProjectDocumentVersionParams) error {
	var tx bun.Tx
	var shouldCommit bool = false

	if r.tx != nil && !r.tx.IsClosed() {
		tx = *r.tx.Tx
	} else {
		var err error
		tx, err = r.db.BeginTx(context.Background(), nil)
		shouldCommit = true

		if err != nil {
			return err
		}
	}

	var userEditorInternalId *string
	if params.ProjectDocumentVersion.UserEditorIdentity != nil {
		identity := params.ProjectDocumentVersion.UserEditorIdentity.Internal.String()
		userEditorInternalId = &identity
	}

	projectDocumentVersionTable := &ProjectDocumentVersionTable{
		InternalId:                              params.ProjectDocumentVersion.Identity.Internal.String(),
		PublicId:                                params.ProjectDocumentVersion.Identity.Public,
		ProjectDocumentVersionManagerInternalId: params.ProjectDocumentVersion.ProjectDocumentVersionManagerIdentity.Internal.String(),
		Version:                                 params.ProjectDocumentVersion.Version,
		Title:                                   params.ProjectDocumentVersion.Document.Title,
		Content:                                 params.ProjectDocumentVersion.Document.Content,
		UserEditorInternalId:                    userEditorInternalId,
		UpdatedAt:                               params.ProjectDocumentVersion.Timestamps.UpdatedAt,
	}

	_, err := tx.NewUpdate().Model(projectDocumentVersionTable).Where("internal_id = ?", params.ProjectDocumentVersion.Identity.Internal.String()).Exec(context.Background())
	if err != nil {
		return err
	}

	_, err = tx.NewDelete().Model(&ProjectDocumentFileTable{}).Where("project_document_version_internal_id = ?", params.ProjectDocumentVersion.Identity.Internal.String()).Exec(context.Background())
	if err != nil {
		return err
	}

	for _, file := range params.ProjectDocumentVersion.Document.Files {
		projectDocumentFileTable := &ProjectDocumentFileTable{
			InternalId:                       file.Identity.Internal.String(),
			ProjectDocumentVersionInternalId: params.ProjectDocumentVersion.Identity.Internal.String(),
			FileInternalId:                   file.FileIdentity.Internal.String(),
		}

		_, err := tx.NewInsert().Model(projectDocumentFileTable).Exec(context.Background())
		if err != nil {
			return err
		}
	}

	if shouldCommit {
		err = tx.Commit()
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ProjectDocumentBunRepository) DeleteProjectDocumentVersion(params projectrepo.DeleteProjectDocumentVersionParams) error {
	var tx bun.Tx
	var shouldCommit bool = false

	if r.tx != nil && !r.tx.IsClosed() {
		tx = *r.tx.Tx
	} else {
		var err error
		tx, err = r.db.BeginTx(context.Background(), nil)
		shouldCommit = true

		if err != nil {
			return err
		}
	}

	_, err := tx.NewDelete().Model(&ProjectDocumentVersionTable{}).Where("internal_id = ?", params.ProjectDocumentVersionIdentity.Internal.String()).Exec(context.Background())
	if err != nil {
		return err
	}

	if shouldCommit {
		err = tx.Commit()
		if err != nil {
			return err
		}
	}

	return nil
}
