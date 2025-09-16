package workspace_database_postgres

import (
	"context"
	"database/sql"

	"github.com/gabrielmrtt/taski/internal/core"
	core_database_postgres "github.com/gabrielmrtt/taski/internal/core/database/postgres"
	user_core "github.com/gabrielmrtt/taski/internal/user"
	user_database_postgres "github.com/gabrielmrtt/taski/internal/user/database/postgres"
	workspace_core "github.com/gabrielmrtt/taski/internal/workspace"
	workspace_repositories "github.com/gabrielmrtt/taski/internal/workspace/repositories"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type WorkspaceUserTable struct {
	bun.BaseModel `bun:"table:workspace_user,alias:workspace_user"`

	WorkspaceInternalId string `bun:"workspace_internal_id,pk,notnull,type:uuid"`
	UserInternalId      string `bun:"user_internal_id,pk,notnull,type:uuid"`
	Status              string `bun:"status,notnull,type:varchar(100)"`

	Workspace *WorkspaceTable                   `bun:"rel:has-one,join:workspace_internal_id=internal_id"`
	User      *user_database_postgres.UserTable `bun:"rel:has-one,join:user_internal_id=internal_id"`
}

func (w *WorkspaceUserTable) ToEntity() *workspace_core.WorkspaceUser {
	return &workspace_core.WorkspaceUser{
		WorkspaceIdentity: core.NewIdentityFromInternal(uuid.MustParse(w.WorkspaceInternalId), workspace_core.WorkspaceIdentityPrefix),
		UserIdentity:      core.NewIdentityFromInternal(uuid.MustParse(w.UserInternalId), user_core.UserIdentityPrefix),
		Status:            workspace_core.WorkspaceUserStatuses(w.Status),
	}
}

type WorkspaceUserPostgresRepository struct {
	db *bun.DB
	tx *core_database_postgres.TransactionPostgres
}

func NewWorkspaceUserPostgresRepository() *WorkspaceUserPostgresRepository {
	return &WorkspaceUserPostgresRepository{db: core_database_postgres.DB, tx: nil}
}

func (r *WorkspaceUserPostgresRepository) SetTransaction(tx core.Transaction) error {
	r.tx = tx.(*core_database_postgres.TransactionPostgres)
	return nil
}

func (r *WorkspaceUserPostgresRepository) applyFilters(selectQuery *bun.SelectQuery, filters workspace_repositories.WorkspaceUserFilters) *bun.SelectQuery {
	selectQuery = selectQuery.Where("workspace_internal_id = ?", filters.WorkspaceIdentity.Internal.String())

	if filters.UserIdentity != nil {
		selectQuery = selectQuery.Where("user_internal_id = ?", filters.UserIdentity.Internal.String())
	}

	if filters.Status != nil {
		selectQuery = selectQuery.Where("status = ?", filters.Status)
	}

	return selectQuery
}

func (r *WorkspaceUserPostgresRepository) GetWorkspaceUserByIdentity(params workspace_repositories.GetWorkspaceUserByIdentityParams) (*workspace_core.WorkspaceUser, error) {
	var workspaceUser *WorkspaceUserTable = new(WorkspaceUserTable)
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(workspaceUser)
	selectQuery = core_database_postgres.ApplyRelations(selectQuery, params.RelationsInput)
	selectQuery = selectQuery.Where("workspace_internal_id = ? and user_internal_id = ?", params.WorkspaceIdentity.Internal.String(), params.UserIdentity.Internal.String())

	err := selectQuery.Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	if workspaceUser.WorkspaceInternalId == "" {
		return nil, nil
	}

	return workspaceUser.ToEntity(), nil
}

func (r *WorkspaceUserPostgresRepository) GetWorkspaceUsersByUserIdentity(params workspace_repositories.GetWorkspaceUsersByUserIdentityParams) ([]workspace_core.WorkspaceUser, error) {
	var workspaceUsers []WorkspaceUserTable = make([]WorkspaceUserTable, 0)
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(&workspaceUsers)
	selectQuery = core_database_postgres.ApplyRelations(selectQuery, *params.RelationsInput)
	selectQuery = selectQuery.Where("user_internal_id = ?", params.UserIdentity.Internal.String())

	err := selectQuery.Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return []workspace_core.WorkspaceUser{}, nil
		}

		return []workspace_core.WorkspaceUser{}, err
	}

	var workspaceUserEntities []workspace_core.WorkspaceUser = make([]workspace_core.WorkspaceUser, 0)
	for _, workspaceUser := range workspaceUsers {
		workspaceUserEntities = append(workspaceUserEntities, *workspaceUser.ToEntity())
	}

	return workspaceUserEntities, nil
}

func (r *WorkspaceUserPostgresRepository) StoreWorkspaceUser(params workspace_repositories.StoreWorkspaceUserParams) (*workspace_core.WorkspaceUser, error) {
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

	_, err := tx.NewInsert().Model(&WorkspaceUserTable{
		WorkspaceInternalId: params.WorkspaceUser.WorkspaceIdentity.Internal.String(),
		UserInternalId:      params.WorkspaceUser.UserIdentity.Internal.String(),
		Status:              string(params.WorkspaceUser.Status),
	}).Exec(context.Background())
	if err != nil {
		return nil, err
	}

	if shouldCommit {
		err = tx.Commit()
		if err != nil {
			return nil, err
		}
	}

	return params.WorkspaceUser, nil
}

func (r *WorkspaceUserPostgresRepository) UpdateWorkspaceUser(params workspace_repositories.UpdateWorkspaceUserParams) error {
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

	_, err := tx.NewUpdate().Model(&WorkspaceUserTable{
		WorkspaceInternalId: params.WorkspaceUser.WorkspaceIdentity.Internal.String(),
		UserInternalId:      params.WorkspaceUser.UserIdentity.Internal.String(),
		Status:              string(params.WorkspaceUser.Status),
	}).Where("workspace_internal_id = ? and user_internal_id = ?", params.WorkspaceUser.WorkspaceIdentity.Internal.String(), params.WorkspaceUser.UserIdentity.Internal.String()).Exec(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}

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

func (r *WorkspaceUserPostgresRepository) DeleteWorkspaceUser(params workspace_repositories.DeleteWorkspaceUserParams) error {
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

	_, err := tx.NewDelete().Model(&WorkspaceUserTable{}).Where("workspace_internal_id = ? and user_internal_id = ?", params.WorkspaceIdentity.Internal.String(), params.UserIdentity.Internal.String()).Exec(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}

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

func (r *WorkspaceUserPostgresRepository) DeleteAllByUserIdentity(params workspace_repositories.DeleteAllByUserIdentityParams) error {
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

	_, err := tx.NewDelete().Model(&WorkspaceUserTable{}).Where("user_internal_id = ?", params.UserIdentity.Internal.String()).Exec(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}

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
