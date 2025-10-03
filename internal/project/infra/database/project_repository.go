package projectdatabase

import (
	"context"
	"database/sql"

	"github.com/gabrielmrtt/taski/internal/core"
	coredatabase "github.com/gabrielmrtt/taski/internal/core/database"
	project "github.com/gabrielmrtt/taski/internal/project"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
	"github.com/gabrielmrtt/taski/internal/user"
	userdatabase "github.com/gabrielmrtt/taski/internal/user/infra/database"
	"github.com/gabrielmrtt/taski/internal/workspace"
	workspacedatabase "github.com/gabrielmrtt/taski/internal/workspace/infra/database"
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

	Workspace *workspacedatabase.WorkspaceTable `bun:"rel:has-one,join:workspace_internal_id=internal_id"`
	Creator   *userdatabase.UserTable           `bun:"rel:has-one,join:user_creator_internal_id=internal_id"`
	Editor    *userdatabase.UserTable           `bun:"rel:has-one,join:user_editor_internal_id=internal_id"`
}

func (p *ProjectTable) ToEntity() *project.Project {
	var userCreatorIdentity core.Identity = core.NewIdentityFromInternal(uuid.MustParse(p.UserCreatorInternalId), user.UserIdentityPrefix)
	var userEditorIdentity *core.Identity = nil

	if p.UserEditorInternalId != nil {
		identity := core.NewIdentityFromInternal(uuid.MustParse(*p.UserEditorInternalId), user.UserIdentityPrefix)
		userEditorIdentity = &identity
	}

	return &project.Project{
		Identity:            core.NewIdentityFromInternal(uuid.MustParse(p.InternalId), project.ProjectIdentityPrefix),
		WorkspaceIdentity:   core.NewIdentityFromInternal(uuid.MustParse(p.WorkspaceInternalId), workspace.WorkspaceIdentityPrefix),
		Name:                p.Name,
		Description:         p.Description,
		Status:              project.ProjectStatuses(p.Status),
		Color:               p.Color,
		PriorityLevel:       project.ProjectPriorityLevels(p.PriorityLevel),
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

type ProjectBunRepository struct {
	db *bun.DB
	tx *coredatabase.TransactionBun
}

func NewProjectBunRepository(connection *bun.DB) *ProjectBunRepository {
	return &ProjectBunRepository{db: connection, tx: nil}
}

func (r *ProjectBunRepository) SetTransaction(tx core.Transaction) error {
	r.tx = tx.(*coredatabase.TransactionBun)
	return nil
}

func (r *ProjectBunRepository) applyFilters(selectQuery *bun.SelectQuery, filters projectrepo.ProjectFilters) *bun.SelectQuery {
	selectQuery = selectQuery.Where("workspace_internal_id = ?", filters.WorkspaceIdentity.Internal.String())

	if filters.LoggedUserIdentity != nil {
		selectQuery = selectQuery.Where("internal_id IN (SELECT project_internal_id FROM project_user WHERE user_internal_id = ? AND status = ?)", filters.LoggedUserIdentity.Internal.String(), project.ProjectUserStatusActive)
	}

	if filters.Name != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "name", filters.Name)
	}

	if filters.Description != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "description", filters.Description)
	}

	if filters.Color != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "color", filters.Color)
	}

	if filters.PriorityLevel != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "priority_level", filters.PriorityLevel)
	}

	if filters.Status != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "status", filters.Status)
	}

	if filters.CreatedAt != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "created_at", filters.CreatedAt)
	}

	if filters.UpdatedAt != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "updated_at", filters.UpdatedAt)
	}

	if filters.DeletedAt != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "deleted_at", filters.DeletedAt)
	}

	return selectQuery
}

func (r *ProjectBunRepository) GetProjectByIdentity(params projectrepo.GetProjectByIdentityParams) (*project.Project, error) {
	var project *ProjectTable = new(ProjectTable)
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(project)
	selectQuery = coredatabase.ApplyRelations(selectQuery, params.RelationsInput)
	selectQuery = selectQuery.Where("internal_id = ?", params.ProjectIdentity.Internal.String())
	if params.WorkspaceIdentity != nil {
		selectQuery = selectQuery.Where("workspace_internal_id = ?", params.WorkspaceIdentity.Internal.String())
	}

	if params.OrganizationIdentity != nil {
		selectQuery = selectQuery.Where("workspace_internal_id IN (SELECT internal_id FROM workspace WHERE workspace.organization_internal_id = ?)", params.OrganizationIdentity.Internal.String())
	}

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

func (r *ProjectBunRepository) PaginateProjectsBy(params projectrepo.PaginateProjectsParams) (*core.PaginationOutput[project.Project], error) {
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
	selectQuery = coredatabase.ApplyRelations(selectQuery, params.RelationsInput)
	selectQuery = r.applyFilters(selectQuery, params.Filters)
	countBeforePagination, err := selectQuery.Count(context.Background())
	if err != nil {
		return nil, err
	}

	selectQuery = coredatabase.ApplyPagination(selectQuery, params.Pagination)
	err = selectQuery.Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return &core.PaginationOutput[project.Project]{
				Data:    []project.Project{},
				Page:    page,
				HasMore: false,
				Total:   0,
			}, nil
		}

		return nil, err
	}

	var projectEntities []project.Project = make([]project.Project, 0)
	for _, project := range projects {
		projectEntities = append(projectEntities, *project.ToEntity())
	}

	return &core.PaginationOutput[project.Project]{
		Data:    projectEntities,
		Page:    page,
		HasMore: core.HasMorePages(page, countBeforePagination, perPage),
		Total:   countBeforePagination,
	}, nil
}

func (r *ProjectBunRepository) StoreProject(params projectrepo.StoreProjectParams) (*project.Project, error) {
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

func (r *ProjectBunRepository) UpdateProject(params projectrepo.UpdateProjectParams) error {
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

func (r *ProjectBunRepository) DeleteProject(params projectrepo.DeleteProjectParams) error {
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
