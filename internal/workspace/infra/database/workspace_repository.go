package workspacedatabase

import (
	"context"
	"database/sql"

	"github.com/gabrielmrtt/taski/internal/core"
	coredatabase "github.com/gabrielmrtt/taski/internal/core/database"
	"github.com/gabrielmrtt/taski/internal/organization"
	organizationdatabase "github.com/gabrielmrtt/taski/internal/organization/infra/database"
	"github.com/gabrielmrtt/taski/internal/user"
	userdatabase "github.com/gabrielmrtt/taski/internal/user/infra/database"
	"github.com/gabrielmrtt/taski/internal/workspace"
	workspacerepo "github.com/gabrielmrtt/taski/internal/workspace/repository"
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

	Organization *organizationdatabase.OrganizationTable `bun:"rel:has-one,join:organization_internal_id=internal_id"`
	Creator      *userdatabase.UserTable                 `bun:"rel:has-one,join:user_creator_internal_id=internal_id"`
	Editor       *userdatabase.UserTable                 `bun:"rel:has-one,join:user_editor_internal_id=internal_id"`
}

func (w *WorkspaceTable) ToEntity() *workspace.Workspace {
	var userCreatorIdentity core.Identity = core.NewIdentityFromInternal(uuid.MustParse(w.UserCreatorInternalId), user.UserIdentityPrefix)
	var userEditorIdentity *core.Identity = nil
	if w.UserEditorInternalId != nil {
		identity := core.NewIdentityFromInternal(uuid.MustParse(*w.UserEditorInternalId), user.UserIdentityPrefix)
		userEditorIdentity = &identity
	}

	var creator *user.User = nil
	if w.Creator != nil {
		creator = w.Creator.ToEntity()
	}

	var editor *user.User = nil
	if w.Editor != nil {
		editor = w.Editor.ToEntity()
	}

	var org *organization.Organization = nil
	if w.Organization != nil {
		org = w.Organization.ToEntity()
	}

	return &workspace.Workspace{
		Identity:             core.NewIdentityFromInternal(uuid.MustParse(w.InternalId), workspace.WorkspaceIdentityPrefix),
		Name:                 w.Name,
		Description:          w.Description,
		Color:                w.Color,
		Status:               workspace.WorkspaceStatuses(w.Status),
		OrganizationIdentity: core.NewIdentityFromInternal(uuid.MustParse(w.OrganizationInternalId), organization.OrganizationIdentityPrefix),
		UserCreatorIdentity:  &userCreatorIdentity,
		UserEditorIdentity:   userEditorIdentity,
		Creator:              creator,
		Editor:               editor,
		Organization:         org,
		Timestamps: core.Timestamps{
			CreatedAt: &w.CreatedAt,
			UpdatedAt: w.UpdatedAt,
		},
		DeletedAt: w.DeletedAt,
	}
}

type WorkspaceRepository struct {
	db *bun.DB
	tx *coredatabase.TransactionBun
}

func NewWorkspaceBunRepository(connection *bun.DB) *WorkspaceRepository {
	return &WorkspaceRepository{db: connection, tx: nil}
}

func (r *WorkspaceRepository) SetTransaction(tx core.Transaction) error {
	if r.tx != nil && !r.tx.IsClosed() {
		return nil
	}

	r.tx = tx.(*coredatabase.TransactionBun)
	return nil
}

func (r *WorkspaceRepository) applyFilters(selectQuery *bun.SelectQuery, filters workspacerepo.WorkspaceFilters) *bun.SelectQuery {
	if filters.OrganizationIdentity != nil {
		selectQuery = selectQuery.Where("organization_internal_id = ?", filters.OrganizationIdentity.Internal.String())
	}

	if filters.AuthenticatedUserIdentity != nil {
		selectQuery = selectQuery.Where("workspace.internal_id IN (SELECT workspace_internal_id FROM workspace_user WHERE user_internal_id = ? AND workspace_user.status = ?)", filters.AuthenticatedUserIdentity.Internal.String(), workspace.WorkspaceUserStatusActive)
	}

	if filters.Name != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "name", filters.Name)
	}

	if filters.Description != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "description", filters.Description)
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

func (r *WorkspaceRepository) GetWorkspaceByIdentity(params workspacerepo.GetWorkspaceByIdentityParams) (*workspace.Workspace, error) {
	var workspace *WorkspaceTable = new(WorkspaceTable)
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(workspace)
	selectQuery = coredatabase.ApplyRelations(selectQuery, params.RelationsInput)
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

func (r *WorkspaceRepository) PaginateWorkspacesBy(params workspacerepo.PaginateWorkspacesParams) (*core.PaginationOutput[workspace.Workspace], error) {
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
	selectQuery = selectQuery.Relation("Creator").Relation("Editor").Relation("Organization")
	selectQuery = r.applyFilters(selectQuery, params.Filters)

	if !params.ShowDeleted {
		selectQuery = selectQuery.Where("deleted_at IS NULL")
	}

	countBeforePagination, err := selectQuery.Count(context.Background())
	if err != nil {
		return nil, err
	}

	selectQuery = coredatabase.ApplySort(selectQuery, params.SortInput)
	selectQuery = coredatabase.ApplyPagination(selectQuery, params.Pagination)
	err = selectQuery.Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return &core.PaginationOutput[workspace.Workspace]{
				Data:    []workspace.Workspace{},
				Page:    page,
				HasMore: false,
				Total:   0,
			}, nil
		}
	}

	var workspaceEntities []workspace.Workspace = make([]workspace.Workspace, 0)
	for _, workspace := range workspaces {
		workspaceEntities = append(workspaceEntities, *workspace.ToEntity())
	}

	return &core.PaginationOutput[workspace.Workspace]{
		Data:    workspaceEntities,
		Page:    page,
		HasMore: core.HasMorePages(page, countBeforePagination, perPage),
		Total:   countBeforePagination,
	}, nil
}

func (r *WorkspaceRepository) StoreWorkspace(params workspacerepo.StoreWorkspaceParams) (*workspace.Workspace, error) {
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

func (r *WorkspaceRepository) UpdateWorkspace(params workspacerepo.UpdateWorkspaceParams) error {
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

func (r *WorkspaceRepository) DeleteWorkspace(params workspacerepo.DeleteWorkspaceParams) error {
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
