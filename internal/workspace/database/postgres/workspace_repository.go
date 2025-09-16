package workspace_database_postgres

import (
	"context"
	"database/sql"

	"github.com/gabrielmrtt/taski/internal/core"
	core_database_postgres "github.com/gabrielmrtt/taski/internal/core/database/postgres"
	organization_core "github.com/gabrielmrtt/taski/internal/organization"
	organization_database_postgres "github.com/gabrielmrtt/taski/internal/organization/database/postgres"
	user_core "github.com/gabrielmrtt/taski/internal/user"
	workspace_core "github.com/gabrielmrtt/taski/internal/workspace"
	workspace_repositories "github.com/gabrielmrtt/taski/internal/workspace/repositories"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type WorkspaceTable struct {
	bun.BaseModel `bun:"table:workspace,alias:workspace"`

	InternalId             string  `bun:"internal_id,pk,notnull,type:uuid"`
	PublicId               string  `bun:"public_id,notnull,type:varchar(510)"`
	Name                   string  `bun:"name,notnull,type:varchar(255)"`
	Description            string  `bun:"description,type:varchar(510)"`
	Color                  string  `bun:"color,notnull,type:varchar(7)"`
	Status                 string  `bun:"status,notnull,type:varchar(100)"`
	OrganizationInternalId string  `bun:"organization_internal_id,notnull,type:uuid"`
	UserCreatorInternalId  string  `bun:"user_creator_internal_id,notnull,type:uuid"`
	UserEditorInternalId   *string `bun:"user_editor_internal_id,type:uuid"`
	CreatedAt              int64   `bun:"created_at,notnull,type:bigint"`
	UpdatedAt              *int64  `bun:"updated_at,type:bigint"`
	DeletedAt              *int64  `bun:"deleted_at,type:bigint"`

	Organization *organization_database_postgres.OrganizationTable `bun:"rel:has-one,join:organization_internal_id=internal_id"`
}

func (w *WorkspaceTable) ToEntity() *workspace_core.Workspace {
	var userCreatorIdentity core.Identity = core.NewIdentityFromInternal(uuid.MustParse(w.UserCreatorInternalId), user_core.UserIdentityPrefix)
	var userEditorIdentity *core.Identity = nil
	if w.UserEditorInternalId != nil {
		identity := core.NewIdentityFromInternal(uuid.MustParse(*w.UserEditorInternalId), user_core.UserIdentityPrefix)
		userEditorIdentity = &identity
	}

	return &workspace_core.Workspace{
		Identity:             core.NewIdentityFromInternal(uuid.MustParse(w.InternalId), workspace_core.WorkspaceIdentityPrefix),
		Name:                 w.Name,
		Description:          w.Description,
		Color:                w.Color,
		Status:               workspace_core.WorkspaceStatuses(w.Status),
		OrganizationIdentity: core.NewIdentityFromInternal(uuid.MustParse(w.OrganizationInternalId), organization_core.OrganizationIdentityPrefix),
		UserCreatorIdentity:  &userCreatorIdentity,
		UserEditorIdentity:   userEditorIdentity,
		Timestamps: core.Timestamps{
			CreatedAt: &w.CreatedAt,
			UpdatedAt: w.UpdatedAt,
		},
		DeletedAt: w.DeletedAt,
	}
}

type WorkspaceRepository struct {
	db *bun.DB
	tx *core_database_postgres.TransactionPostgres
}

func NewWorkspacePostgresRepository() *WorkspaceRepository {
	return &WorkspaceRepository{db: core_database_postgres.DB, tx: nil}
}

func (r *WorkspaceRepository) SetTransaction(tx core.Transaction) error {
	r.tx = tx.(*core_database_postgres.TransactionPostgres)
	return nil
}

func (r *WorkspaceRepository) applyFilters(selectQuery *bun.SelectQuery, filters workspace_repositories.WorkspaceFilters) *bun.SelectQuery {
	selectQuery = selectQuery.Where("organization_internal_id = ?", filters.OrganizationIdentity.Internal.String())

	if filters.Name != nil {
		selectQuery = core_database_postgres.ApplyComparableFilter(selectQuery, "name", filters.Name)
	}

	if filters.Description != nil {
		selectQuery = core_database_postgres.ApplyComparableFilter(selectQuery, "description", filters.Description)
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

func (r *WorkspaceRepository) GetWorkspaceByIdentity(params workspace_repositories.GetWorkspaceByIdentityParams) (*workspace_core.Workspace, error) {
	var workspace *WorkspaceTable = new(WorkspaceTable)
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(workspace)
	selectQuery = selectQuery.Where("internal_id = ?", params.WorkspaceIdentity.Internal.String())
	err := selectQuery.Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	if workspace.InternalId == "" {
		return nil, nil
	}

	return workspace.ToEntity(), nil
}

func (r *WorkspaceRepository) PaginateWorkspacesBy(params workspace_repositories.PaginateWorkspacesParams) (*core.PaginationOutput[workspace_core.Workspace], error) {
	var workspaces []WorkspaceTable = make([]WorkspaceTable, 0)
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

	selectQuery = selectQuery.Model(&workspaces)
	selectQuery = r.applyFilters(selectQuery, params.Filters)
	countBeforePagination, err := selectQuery.Count(context.Background())
	if err != nil {
		return nil, err
	}

	selectQuery = core_database_postgres.ApplyPagination(selectQuery, params.Pagination)
	err = selectQuery.Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return &core.PaginationOutput[workspace_core.Workspace]{
				Data:    []workspace_core.Workspace{},
				Page:    page,
				HasMore: false,
				Total:   0,
			}, nil
		}
	}

	var workspaceEntities []workspace_core.Workspace = make([]workspace_core.Workspace, 0)
	for _, workspace := range workspaces {
		workspaceEntities = append(workspaceEntities, *workspace.ToEntity())
	}

	return &core.PaginationOutput[workspace_core.Workspace]{
		Data:    workspaceEntities,
		Page:    page,
		HasMore: core.HasMorePages(page, countBeforePagination, perPage),
		Total:   countBeforePagination,
	}, nil
}

func (r *WorkspaceRepository) StoreWorkspace(params workspace_repositories.StoreWorkspaceParams) (*workspace_core.Workspace, error) {
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
	if params.Workspace.UserEditorIdentity != nil {
		identity := params.Workspace.UserEditorIdentity.Internal.String()
		userEditorInternalId = &identity
	}

	workspaceTable := &WorkspaceTable{
		InternalId:             params.Workspace.Identity.Internal.String(),
		PublicId:               params.Workspace.Identity.Public,
		Name:                   params.Workspace.Name,
		Description:            params.Workspace.Description,
		Color:                  params.Workspace.Color,
		Status:                 string(params.Workspace.Status),
		OrganizationInternalId: params.Workspace.OrganizationIdentity.Internal.String(),
		UserCreatorInternalId:  params.Workspace.UserCreatorIdentity.Internal.String(),
		UserEditorInternalId:   userEditorInternalId,
		CreatedAt:              *params.Workspace.Timestamps.CreatedAt,
		UpdatedAt:              params.Workspace.Timestamps.UpdatedAt,
		DeletedAt:              params.Workspace.DeletedAt,
	}

	_, err := tx.NewInsert().Model(workspaceTable).Exec(context.Background())
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

	return params.Workspace, nil
}

func (r *WorkspaceRepository) UpdateWorkspace(params workspace_repositories.UpdateWorkspaceParams) error {
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
	if params.Workspace.UserEditorIdentity != nil {
		identity := params.Workspace.UserEditorIdentity.Internal.String()
		userEditorInternalId = &identity
	}

	workspaceTable := &WorkspaceTable{
		InternalId:             params.Workspace.Identity.Internal.String(),
		PublicId:               params.Workspace.Identity.Public,
		Name:                   params.Workspace.Name,
		Description:            params.Workspace.Description,
		Color:                  params.Workspace.Color,
		Status:                 string(params.Workspace.Status),
		OrganizationInternalId: params.Workspace.OrganizationIdentity.Internal.String(),
		UserCreatorInternalId:  params.Workspace.UserCreatorIdentity.Internal.String(),
		UserEditorInternalId:   userEditorInternalId,
		UpdatedAt:              params.Workspace.Timestamps.UpdatedAt,
		DeletedAt:              params.Workspace.DeletedAt,
	}

	_, err := tx.NewUpdate().Model(workspaceTable).Where("internal_id = ?", params.Workspace.Identity.Internal.String()).Exec(context.Background())
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

func (r *WorkspaceRepository) DeleteWorkspace(params workspace_repositories.DeleteWorkspaceParams) error {
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

	_, err := tx.NewDelete().Model(&WorkspaceTable{}).Where("internal_id = ?", params.WorkspaceIdentity.Internal.String()).Exec(context.Background())
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
