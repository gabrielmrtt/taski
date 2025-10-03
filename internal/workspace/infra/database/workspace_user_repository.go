package workspacedatabase

import (
	"context"
	"database/sql"

	"github.com/gabrielmrtt/taski/internal/core"
	coredatabase "github.com/gabrielmrtt/taski/internal/core/database"
	userdatabase "github.com/gabrielmrtt/taski/internal/user/infra/database"
	"github.com/gabrielmrtt/taski/internal/workspace"
	workspacerepo "github.com/gabrielmrtt/taski/internal/workspace/repository"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type WorkspaceUserTable struct {
	bun.BaseModel `bun:"table:workspace_user,alias:workspace_user"`

	WorkspaceInternalId string `bun:"workspace_internal_id,pk,notnull,type:uuid"`
	UserInternalId      string `bun:"user_internal_id,pk,notnull,type:uuid"`
	Status              string `bun:"status,notnull,type:varchar(100)"`

	Workspace *WorkspaceTable         `bun:"rel:has-one,join:workspace_internal_id=internal_id"`
	User      *userdatabase.UserTable `bun:"rel:has-one,join:user_internal_id=internal_id"`
}

func (w *WorkspaceUserTable) ToEntity() *workspace.WorkspaceUser {
	return &workspace.WorkspaceUser{
		WorkspaceIdentity: core.NewIdentityFromInternal(uuid.MustParse(w.WorkspaceInternalId), workspace.WorkspaceIdentityPrefix),
		User:              *w.User.ToEntity(),
		Status:            workspace.WorkspaceUserStatuses(w.Status),
	}
}

type WorkspaceUserBunRepository struct {
	db *bun.DB
	tx *coredatabase.TransactionBun
}

func NewWorkspaceUserBunRepository(connection *bun.DB) *WorkspaceUserBunRepository {
	return &WorkspaceUserBunRepository{db: connection, tx: nil}
}

func (r *WorkspaceUserBunRepository) SetTransaction(tx core.Transaction) error {
	r.tx = tx.(*coredatabase.TransactionBun)
	return nil
}

func (r *WorkspaceUserBunRepository) applyFilters(selectQuery *bun.SelectQuery, filters workspacerepo.WorkspaceUserFilters) *bun.SelectQuery {
	selectQuery = selectQuery.Where("workspace_internal_id = ?", filters.WorkspaceIdentity.Internal.String())

	if filters.UserIdentity != nil {
		selectQuery = selectQuery.Where("user_internal_id = ?", filters.UserIdentity.Internal.String())
	}

	if filters.Status != nil {
		selectQuery = selectQuery.Where("status = ?", filters.Status)
	}

	return selectQuery
}

func (r *WorkspaceUserBunRepository) GetWorkspaceUserByIdentity(params workspacerepo.GetWorkspaceUserByIdentityParams) (*workspace.WorkspaceUser, error) {
	var workspaceUser *WorkspaceUserTable = new(WorkspaceUserTable)
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(workspaceUser)
	selectQuery = selectQuery.Relation("User")
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

func (r *WorkspaceUserBunRepository) GetWorkspaceUsersByUserIdentity(params workspacerepo.GetWorkspaceUsersByUserIdentityParams) ([]workspace.WorkspaceUser, error) {
	var workspaceUsers []WorkspaceUserTable = make([]WorkspaceUserTable, 0)
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(&workspaceUsers)
	selectQuery = selectQuery.Relation("User")
	selectQuery = selectQuery.Where("user_internal_id = ?", params.UserIdentity.Internal.String())

	err := selectQuery.Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return []workspace.WorkspaceUser{}, nil
		}

		return []workspace.WorkspaceUser{}, err
	}

	var workspaceUserEntities []workspace.WorkspaceUser = make([]workspace.WorkspaceUser, 0)
	for _, workspaceUser := range workspaceUsers {
		workspaceUserEntities = append(workspaceUserEntities, *workspaceUser.ToEntity())
	}

	return workspaceUserEntities, nil
}

func (r *WorkspaceUserBunRepository) StoreWorkspaceUser(params workspacerepo.StoreWorkspaceUserParams) (*workspace.WorkspaceUser, error) {
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
		UserInternalId:      params.WorkspaceUser.User.Identity.Internal.String(),
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

func (r *WorkspaceUserBunRepository) UpdateWorkspaceUser(params workspacerepo.UpdateWorkspaceUserParams) error {
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
		UserInternalId:      params.WorkspaceUser.User.Identity.Internal.String(),
		Status:              string(params.WorkspaceUser.Status),
	}).Where("workspace_internal_id = ? and user_internal_id = ?", params.WorkspaceUser.WorkspaceIdentity.Internal.String(), params.WorkspaceUser.User.Identity.Internal.String()).Exec(context.Background())
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

func (r *WorkspaceUserBunRepository) DeleteWorkspaceUser(params workspacerepo.DeleteWorkspaceUserParams) error {
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

func (r *WorkspaceUserBunRepository) DeleteAllByUserIdentity(params workspacerepo.DeleteAllByUserIdentityParams) error {
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
