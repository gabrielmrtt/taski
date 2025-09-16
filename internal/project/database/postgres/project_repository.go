package project_database_postgres

import (
	"context"
	"database/sql"

	"github.com/gabrielmrtt/taski/internal/core"
	core_database_postgres "github.com/gabrielmrtt/taski/internal/core/database/postgres"
	project_core "github.com/gabrielmrtt/taski/internal/project"
	project_repositories "github.com/gabrielmrtt/taski/internal/project/repositories"
	user_core "github.com/gabrielmrtt/taski/internal/user"
	user_database_postgres "github.com/gabrielmrtt/taski/internal/user/database/postgres"
	workspace_core "github.com/gabrielmrtt/taski/internal/workspace"
	workspace_database_postgres "github.com/gabrielmrtt/taski/internal/workspace/database/postgres"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ProjectTable struct {
	bun.BaseModel `bun:"table:project,alias:project"`

	InternalId            string  `bun:"internal_id,pk,notnull,type:uuid"`
	PublicId              string  `bun:"public_id,notnull,type:varchar(510)"`
	Name                  string  `bun:"name,notnull,type:varchar(255)"`
	Description           string  `bun:"description,type:varchar(510)"`
	Status                string  `bun:"status,notnull,type:varchar(100)"`
	Color                 string  `bun:"color,notnull,type:varchar(7)"`
	PriorityLevel         int     `bun:"priority_level,notnull,type:int"`
	StartAt               *int64  `bun:"start_at,type:bigint"`
	EndAt                 *int64  `bun:"end_at,type:bigint"`
	WorkspaceInternalId   string  `bun:"workspace_internal_id,notnull,type:uuid"`
	UserCreatorInternalId string  `bun:"user_creator_internal_id,notnull,type:uuid"`
	UserEditorInternalId  *string `bun:"user_editor_internal_id,type:uuid"`
	CreatedAt             int64   `bun:"created_at,notnull,type:bigint"`
	UpdatedAt             *int64  `bun:"updated_at,type:bigint"`
	DeletedAt             *int64  `bun:"deleted_at,type:bigint"`

	Workspace *workspace_database_postgres.WorkspaceTable `bun:"rel:has-one,join:workspace_internal_id=internal_id"`
	Creator   *user_database_postgres.UserTable           `bun:"rel:has-one,join:user_creator_internal_id=internal_id"`
	Editor    *user_database_postgres.UserTable           `bun:"rel:has-one,join:user_editor_internal_id=internal_id"`
}

func (p *ProjectTable) ToEntity() *project_core.Project {
	var userCreatorIdentity core.Identity = core.NewIdentityFromInternal(uuid.MustParse(p.UserCreatorInternalId), user_core.UserIdentityPrefix)
	var userEditorIdentity *core.Identity = nil

	if p.UserEditorInternalId != nil {
		identity := core.NewIdentityFromInternal(uuid.MustParse(*p.UserEditorInternalId), user_core.UserIdentityPrefix)
		userEditorIdentity = &identity
	}

	return &project_core.Project{
		Identity:            core.NewIdentityFromInternal(uuid.MustParse(p.InternalId), project_core.ProjectIdentityPrefix),
		WorkspaceIdentity:   core.NewIdentityFromInternal(uuid.MustParse(p.WorkspaceInternalId), workspace_core.WorkspaceIdentityPrefix),
		Name:                p.Name,
		Description:         p.Description,
		Status:              project_core.ProjectStatuses(p.Status),
		Color:               p.Color,
		PriorityLevel:       project_core.ProjectPriorityLevels(p.PriorityLevel),
		UserCreatorIdentity: &userCreatorIdentity,
		UserEditorIdentity:  userEditorIdentity,
		StartAt:             p.StartAt,
		EndAt:               p.EndAt,
		Timestamps: core.Timestamps{
			CreatedAt: &p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		},
		DeletedAt: p.DeletedAt,
	}
}

type ProjectPostgresRepository struct {
	db *bun.DB
	tx *core_database_postgres.TransactionPostgres
}

func NewProjectPostgresRepository() *ProjectPostgresRepository {
	return &ProjectPostgresRepository{db: core_database_postgres.DB, tx: nil}
}

func (r *ProjectPostgresRepository) SetTransaction(tx core.Transaction) error {
	r.tx = tx.(*core_database_postgres.TransactionPostgres)
	return nil
}

func (r *ProjectPostgresRepository) applyFilters(selectQuery *bun.SelectQuery, filters project_repositories.ProjectFilters) *bun.SelectQuery {
	selectQuery = selectQuery.Where("workspace_internal_id = ?", filters.WorkspaceIdentity.Internal.String())

	if filters.Name != nil {
		selectQuery = core_database_postgres.ApplyComparableFilter(selectQuery, "name", filters.Name)
	}

	if filters.Description != nil {
		selectQuery = core_database_postgres.ApplyComparableFilter(selectQuery, "description", filters.Description)
	}

	if filters.Color != nil {
		selectQuery = core_database_postgres.ApplyComparableFilter(selectQuery, "color", filters.Color)
	}

	if filters.PriorityLevel != nil {
		selectQuery = core_database_postgres.ApplyComparableFilter(selectQuery, "priority_level", filters.PriorityLevel)
	}

	if filters.Status != nil {
		selectQuery = core_database_postgres.ApplyComparableFilter(selectQuery, "status", filters.Status)
	}

	if filters.CreatedAt != nil {
		selectQuery = core_database_postgres.ApplyComparableFilter(selectQuery, "created_at", filters.CreatedAt)
	}

	if filters.UpdatedAt != nil {
		selectQuery = core_database_postgres.ApplyComparableFilter(selectQuery, "updated_at", filters.UpdatedAt)
	}

	if filters.DeletedAt != nil {
		selectQuery = core_database_postgres.ApplyComparableFilter(selectQuery, "deleted_at", filters.DeletedAt)
	}

	return selectQuery
}

func (r *ProjectPostgresRepository) GetProjectByIdentity(params project_repositories.GetProjectByIdentityParams) (*project_core.Project, error) {
	var project *ProjectTable = new(ProjectTable)
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(project)
	selectQuery = core_database_postgres.ApplyRelations(selectQuery, params.RelationsInput)
	selectQuery = selectQuery.Where("internal_id = ?", params.ProjectIdentity.Internal.String())

	err := selectQuery.Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	if project.InternalId == "" {
		return nil, nil
	}

	return project.ToEntity(), nil
}

func (r *ProjectPostgresRepository) PaginateProjectsBy(params project_repositories.PaginateProjectsParams) (*core.PaginationOutput[project_core.Project], error) {
	var projects []ProjectTable = make([]ProjectTable, 0)
	var selectQuery *bun.SelectQuery
	var perPage int = 10
	var page int = 1

	if params.Pagination.PerPage != nil {
		perPage = *params.Pagination.PerPage
	}

	if params.Pagination.Page != nil {
		page = *params.Pagination.Page
	}

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(&projects)
	selectQuery = core_database_postgres.ApplyRelations(selectQuery, params.RelationsInput)
	selectQuery = r.applyFilters(selectQuery, params.Filters)
	countBeforePagination, err := selectQuery.Count(context.Background())
	if err != nil {
		return nil, err
	}

	selectQuery = core_database_postgres.ApplyPagination(selectQuery, params.Pagination)
	err = selectQuery.Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return &core.PaginationOutput[project_core.Project]{
				Data:    []project_core.Project{},
				Page:    page,
				HasMore: false,
				Total:   0,
			}, nil
		}

		return nil, err
	}

	var projectEntities []project_core.Project = make([]project_core.Project, 0)
	for _, project := range projects {
		projectEntities = append(projectEntities, *project.ToEntity())
	}

	return &core.PaginationOutput[project_core.Project]{
		Data:    projectEntities,
		Page:    page,
		HasMore: core.HasMorePages(page, countBeforePagination, perPage),
		Total:   countBeforePagination,
	}, nil
}

func (r *ProjectPostgresRepository) StoreProject(params project_repositories.StoreProjectParams) (*project_core.Project, error) {
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
	if params.Project.UserEditorIdentity != nil {
		identity := params.Project.UserEditorIdentity.Internal.String()
		userEditorInternalId = &identity
	}

	_, err := tx.NewInsert().Model(&ProjectTable{
		InternalId:            params.Project.Identity.Internal.String(),
		PublicId:              params.Project.Identity.Public,
		Name:                  params.Project.Name,
		Description:           params.Project.Description,
		Status:                string(params.Project.Status),
		Color:                 params.Project.Color,
		PriorityLevel:         int(params.Project.PriorityLevel),
		StartAt:               params.Project.StartAt,
		EndAt:                 params.Project.EndAt,
		WorkspaceInternalId:   params.Project.WorkspaceIdentity.Internal.String(),
		UserCreatorInternalId: params.Project.UserCreatorIdentity.Internal.String(),
		UserEditorInternalId:  userEditorInternalId,
		CreatedAt:             *params.Project.Timestamps.CreatedAt,
		UpdatedAt:             params.Project.Timestamps.UpdatedAt,
		DeletedAt:             params.Project.DeletedAt,
	}).Exec(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	if shouldCommit {
		err = tx.Commit()
		if err != nil {
			return nil, err
		}
	}

	return params.Project, nil
}

func (r *ProjectPostgresRepository) UpdateProject(params project_repositories.UpdateProjectParams) error {
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
	if params.Project.UserEditorIdentity != nil {
		identity := params.Project.UserEditorIdentity.Internal.String()
		userEditorInternalId = &identity
	}

	projectTable := &ProjectTable{
		InternalId:            params.Project.Identity.Internal.String(),
		PublicId:              params.Project.Identity.Public,
		Name:                  params.Project.Name,
		Description:           params.Project.Description,
		Status:                string(params.Project.Status),
		Color:                 params.Project.Color,
		PriorityLevel:         int(params.Project.PriorityLevel),
		StartAt:               params.Project.StartAt,
		EndAt:                 params.Project.EndAt,
		WorkspaceInternalId:   params.Project.WorkspaceIdentity.Internal.String(),
		UserCreatorInternalId: params.Project.UserCreatorIdentity.Internal.String(),
		UserEditorInternalId:  userEditorInternalId,
		UpdatedAt:             params.Project.Timestamps.UpdatedAt,
		DeletedAt:             params.Project.DeletedAt,
	}

	_, err := tx.NewUpdate().Model(projectTable).Where("internal_id = ?", params.Project.Identity.Internal.String()).Exec(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
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

func (r *ProjectPostgresRepository) DeleteProject(params project_repositories.DeleteProjectParams) error {
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

	_, err := tx.NewDelete().Model(&ProjectTable{}).Where("internal_id = ?", params.ProjectIdentity.Internal.String()).Exec(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
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
