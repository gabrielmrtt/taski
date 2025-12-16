package taskdatabase

import (
	"context"
	"database/sql"

	"github.com/gabrielmrtt/taski/internal/core"
	coredatabase "github.com/gabrielmrtt/taski/internal/core/database"
	"github.com/gabrielmrtt/taski/internal/task"
	taskrepo "github.com/gabrielmrtt/taski/internal/task/repository"
	userdatabase "github.com/gabrielmrtt/taski/internal/user/infra/database"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type TaskActionTable struct {
	bun.BaseModel `bun:"table:task_action,alias:task_action"`

	InternalId     string `bun:"internal_id,pk,notnull,type:uuid"`
	PublicId       string `bun:"public_id,notnull,type:varchar(510)"`
	Type           string `bun:"type,notnull,type:varchar(100)"`
	TaskInternalId string `bun:"task_internal_id,notnull,type:uuid"`
	UserInternalId string `bun:"user_internal_id,notnull,type:uuid"`
	CreatedAt      int64  `bun:"created_at,notnull,type:bigint"`

	Task *TaskTable              `bun:"rel:has-one,join:task_internal_id=internal_id"`
	User *userdatabase.UserTable `bun:"rel:has-one,join:user_internal_id=internal_id"`
}

func (t *TaskActionTable) ToEntity() *task.TaskAction {
	return &task.TaskAction{
		Identity:     core.NewIdentityFromInternal(uuid.MustParse(t.InternalId), task.TaskActionIdentityPrefix),
		Type:         task.TaskActionType(t.Type),
		TaskIdentity: core.NewIdentityFromInternal(uuid.MustParse(t.TaskInternalId), task.TaskIdentityPrefix),
		User:         t.User.ToEntity(),
		CreatedAt:    core.DateTime{Value: t.CreatedAt},
	}
}

type TaskActionBunRepository struct {
	db *bun.DB
	tx *coredatabase.TransactionBun
}

func NewTaskActionBunRepository(connection *bun.DB) *TaskActionBunRepository {
	return &TaskActionBunRepository{db: connection, tx: nil}
}

func (r *TaskActionBunRepository) SetTransaction(tx core.Transaction) error {
	if r.tx != nil && !r.tx.IsClosed() {
		return nil
	}

	r.tx = tx.(*coredatabase.TransactionBun)
	return nil
}

func (r *TaskActionBunRepository) applyFilters(query *bun.SelectQuery, filters taskrepo.TaskActionFilters) *bun.SelectQuery {
	if filters.TaskIdentity != nil {
		query = query.Where("task_action.task_internal_id = ?", filters.TaskIdentity.Internal.String())
	}

	if filters.CreatedAt != nil {
		query = coredatabase.ApplyComparableFilter(query, "task_action.created_at", filters.CreatedAt)
	}

	if filters.Type != nil {
		query = coredatabase.ApplyComparableFilter(query, "task_action.type", filters.Type)
	}

	return query
}

func (r *TaskActionBunRepository) PaginateTaskActionsBy(params taskrepo.PaginateTaskActionsParams) (*core.PaginationOutput[task.TaskAction], error) {
	var taskActions []*TaskActionTable = make([]*TaskActionTable, 0)
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

	selectQuery = selectQuery.Model(&taskActions)
	selectQuery = selectQuery.Relation("Task").Relation("User")
	selectQuery = coredatabase.ApplyRelations(selectQuery, params.RelationsInput)
	selectQuery = r.applyFilters(selectQuery, params.Filters)
	countBeforePagination, err := selectQuery.Count(context.Background())
	if err != nil {
		return nil, err
	}

	selectQuery = coredatabase.ApplySort(selectQuery, params.SortInput)
	selectQuery = coredatabase.ApplyPagination(selectQuery, params.Pagination)
	err = selectQuery.Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
	}

	var taskActionEntities []task.TaskAction = make([]task.TaskAction, 0)
	for _, taskAction := range taskActions {
		taskActionEntities = append(taskActionEntities, *taskAction.ToEntity())
	}

	return &core.PaginationOutput[task.TaskAction]{
		Data:    taskActionEntities,
		Page:    page,
		HasMore: core.HasMorePages(page, countBeforePagination, perPage),
		Total:   countBeforePagination,
	}, nil
}

func (r *TaskActionBunRepository) StoreTaskAction(params taskrepo.StoreTaskActionParams) (*task.TaskAction, error) {
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

	var createdAt *int64 = nil
	if params.TaskAction.CreatedAt.Value != 0 {
		createdAt = &params.TaskAction.CreatedAt.Value
	}

	taskActionTable := &TaskActionTable{
		InternalId:     params.TaskAction.Identity.Internal.String(),
		PublicId:       params.TaskAction.Identity.Public,
		Type:           string(params.TaskAction.Type),
		TaskInternalId: params.TaskAction.TaskIdentity.Internal.String(),
		UserInternalId: params.TaskAction.User.Identity.Internal.String(),
		CreatedAt:      *createdAt,
	}

	_, err := tx.NewInsert().Model(taskActionTable).Exec(context.Background())
	if err != nil {
		return nil, err
	}

	if shouldCommit {
		err = tx.Commit()
		if err != nil {
			return nil, err
		}
	}

	return params.TaskAction, nil
}

func (r *TaskActionBunRepository) DeleteTaskAction(params taskrepo.DeleteTaskActionParams) error {
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

	_, err := tx.NewDelete().Model(&TaskActionTable{}).Where("internal_id = ?", params.TaskActionIdentity.Internal.String()).Exec(context.Background())
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
