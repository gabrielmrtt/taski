package taskdatabase

import (
	"context"
	"database/sql"

	"github.com/gabrielmrtt/taski/internal/core"
	coredatabase "github.com/gabrielmrtt/taski/internal/core/database"
	"github.com/gabrielmrtt/taski/internal/storage"
	storagedatabase "github.com/gabrielmrtt/taski/internal/storage/infra/database"
	"github.com/gabrielmrtt/taski/internal/task"
	taskrepo "github.com/gabrielmrtt/taski/internal/task/repository"
	userdatabase "github.com/gabrielmrtt/taski/internal/user/infra/database"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type TaskCommentFileTable struct {
	bun.BaseModel `bun:"table:task_comment_file,alias:task_comment_file"`

	InternalId            string `bun:"internal_id,pk,notnull,type:uuid"`
	TaskCommentInternalId string `bun:"task_comment_internal_id,notnull,type:uuid"`
	FileInternalId        string `bun:"file_internal_id,notnull,type:uuid"`

	TaskComment *TaskCommentTable                  `bun:"rel:has-one,join:task_comment_internal_id=internal_id"`
	File        *storagedatabase.UploadedFileTable `bun:"rel:has-one,join:file_internal_id=internal_id"`
}

func (t *TaskCommentFileTable) ToEntity() *task.TaskCommentFile {
	return &task.TaskCommentFile{
		Identity:     core.NewIdentityWithoutPublicFromInternal(uuid.MustParse(t.InternalId)),
		FileIdentity: core.NewIdentityFromInternal(uuid.MustParse(t.FileInternalId), storage.UploadedFileIdentityPrefix),
	}
}

type TaskCommentTable struct {
	bun.BaseModel `bun:"table:task_comment,alias:task_comment"`

	InternalId           string `bun:"internal_id,pk,notnull,type:uuid"`
	PublicId             string `bun:"public_id,notnull,type:varchar(510)"`
	Content              string `bun:"content,notnull,type:text"`
	TaskInternalId       string `bun:"task_internal_id,notnull,type:uuid"`
	UserAuthorInternalId string `bun:"user_author_internal_id,notnull,type:uuid"`
	CreatedAt            int64  `bun:"created_at,notnull,type:bigint"`
	UpdatedAt            *int64 `bun:"updated_at,type:bigint"`

	Task   *TaskTable              `bun:"rel:has-one,join:task_internal_id=internal_id"`
	Author *userdatabase.UserTable `bun:"rel:has-one,join:user_author_internal_id=internal_id"`
	Files  []*TaskCommentFileTable `bun:"rel:has-many,join:internal_id=task_comment_internal_id"`
}

func (t *TaskCommentTable) ToEntity() *task.TaskComment {
	var files []task.TaskCommentFile = make([]task.TaskCommentFile, len(t.Files))
	for i, file := range t.Files {
		files[i] = *file.ToEntity()
	}

	createdAt := core.DateTime{Value: t.CreatedAt}
	var updatedAt *core.DateTime = nil
	if t.UpdatedAt != nil {
		updatedAt = &core.DateTime{Value: *t.UpdatedAt}
	}

	return &task.TaskComment{
		Identity:     core.NewIdentityFromInternal(uuid.MustParse(t.InternalId), task.TaskCommentIdentityPrefix),
		Content:      t.Content,
		Files:        files,
		TaskIdentity: core.NewIdentityFromInternal(uuid.MustParse(t.TaskInternalId), task.TaskIdentityPrefix),
		Author:       t.Author.ToEntity(),
		Timestamps: core.Timestamps{
			CreatedAt: &createdAt,
			UpdatedAt: updatedAt,
		},
	}
}

type TaskCommentBunRepository struct {
	db *bun.DB
	tx *coredatabase.TransactionBun
}

func NewTaskCommentBunRepository(connection *bun.DB) *TaskCommentBunRepository {
	return &TaskCommentBunRepository{db: connection, tx: nil}
}

func (r *TaskCommentBunRepository) SetTransaction(tx core.Transaction) error {
	if r.tx != nil && !r.tx.IsClosed() {
		return nil
	}

	r.tx = tx.(*coredatabase.TransactionBun)
	return nil
}

func (r *TaskCommentBunRepository) applyFilters(query *bun.SelectQuery, filters taskrepo.TaskCommentFilters) *bun.SelectQuery {
	if filters.TaskIdentity != nil {
		query = query.Where("task_comment.task_internal_id = ?", filters.TaskIdentity.Internal.String())
	}

	if filters.AuthorIdentity != nil {
		query = query.Where("task_comment.user_author_internal_id = ?", filters.AuthorIdentity.Internal.String())
	}

	if filters.CreatedAt != nil {
		query = coredatabase.ApplyComparableFilter(query, "task_comment.created_at", filters.CreatedAt)
	}

	if filters.UpdatedAt != nil {
		query = coredatabase.ApplyComparableFilter(query, "task_comment.updated_at", filters.UpdatedAt)
	}

	return query
}

func (r *TaskCommentBunRepository) GetTaskCommentByIdentity(params taskrepo.GetTaskCommentByIdentityParams) (*task.TaskComment, error) {
	var taskComment *TaskCommentTable = new(TaskCommentTable)
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(taskComment)
	selectQuery = selectQuery.Relation("Author").Relation("Files.File")
	selectQuery = coredatabase.ApplyRelations(selectQuery, params.RelationsInput)
	selectQuery = selectQuery.Where("task_comment.internal_id = ?", params.TaskCommentIdentity.Internal.String())
	err := selectQuery.Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	if taskComment.InternalId == "" {
		return nil, nil
	}

	return taskComment.ToEntity(), nil
}

func (r *TaskCommentBunRepository) PaginateTaskCommentsBy(params taskrepo.PaginateTaskCommentsParams) (*core.PaginationOutput[task.TaskComment], error) {
	var taskComments []*TaskCommentTable = make([]*TaskCommentTable, 0)
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(&taskComments)
	selectQuery = selectQuery.Relation("Author").Relation("Files.File")
	selectQuery = coredatabase.ApplyRelations(selectQuery, params.RelationsInput)
	selectQuery = r.applyFilters(selectQuery, params.Filters)
	selectQuery = coredatabase.ApplyPagination(selectQuery, params.Pagination)
	selectQuery = coredatabase.ApplySort(selectQuery, params.SortInput)

	countBeforePagination, err := selectQuery.Count(context.Background())
	if err != nil {
		return nil, err
	}

	err = selectQuery.Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return &core.PaginationOutput[task.TaskComment]{
				Data:    []task.TaskComment{},
				Page:    *params.Pagination.Page,
				HasMore: false,
				Total:   0,
			}, nil
		}

		return nil, err
	}

	var taskCommentEntities []task.TaskComment = make([]task.TaskComment, 0)
	for _, taskComment := range taskComments {
		taskCommentEntities = append(taskCommentEntities, *taskComment.ToEntity())
	}

	return &core.PaginationOutput[task.TaskComment]{
		Data:    taskCommentEntities,
		Page:    *params.Pagination.Page,
		HasMore: core.HasMorePages(*params.Pagination.Page, countBeforePagination, *params.Pagination.PerPage),
		Total:   countBeforePagination,
	}, nil
}

func (r *TaskCommentBunRepository) StoreTaskComment(params taskrepo.StoreTaskCommentParams) (*task.TaskComment, error) {
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
	if params.TaskComment.Timestamps.CreatedAt != nil {
		createdAt = &params.TaskComment.Timestamps.CreatedAt.Value
	}

	var updatedAt *int64 = nil
	if params.TaskComment.Timestamps.UpdatedAt != nil {
		updatedAt = &params.TaskComment.Timestamps.UpdatedAt.Value
	}

	taskCommentTable := &TaskCommentTable{
		InternalId:           params.TaskComment.Identity.Internal.String(),
		PublicId:             params.TaskComment.Identity.Public,
		Content:              params.TaskComment.Content,
		TaskInternalId:       params.TaskComment.TaskIdentity.Internal.String(),
		UserAuthorInternalId: params.TaskComment.Author.Identity.Internal.String(),
		CreatedAt:            *createdAt,
		UpdatedAt:            updatedAt,
	}

	_, err := tx.NewInsert().Model(taskCommentTable).Exec(context.Background())
	if err != nil {
		return nil, err
	}

	for _, file := range params.TaskComment.Files {
		taskCommentFileTable := &TaskCommentFileTable{
			InternalId:            file.Identity.Internal.String(),
			TaskCommentInternalId: taskCommentTable.InternalId,
			FileInternalId:        file.FileIdentity.Internal.String(),
		}

		_, err = tx.NewInsert().Model(taskCommentFileTable).Exec(context.Background())
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

	return params.TaskComment, nil
}

func (r *TaskCommentBunRepository) UpdateTaskComment(params taskrepo.UpdateTaskCommentParams) error {
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

	var createdAt *int64 = nil
	if params.TaskComment.Timestamps.CreatedAt != nil {
		createdAt = &params.TaskComment.Timestamps.CreatedAt.Value
	}

	var updatedAt *int64 = nil
	if params.TaskComment.Timestamps.UpdatedAt != nil {
		updatedAt = &params.TaskComment.Timestamps.UpdatedAt.Value
	}

	taskCommentTable := &TaskCommentTable{
		InternalId:           params.TaskComment.Identity.Internal.String(),
		Content:              params.TaskComment.Content,
		TaskInternalId:       params.TaskComment.TaskIdentity.Internal.String(),
		UserAuthorInternalId: params.TaskComment.Author.Identity.Internal.String(),
		CreatedAt:            *createdAt,
		UpdatedAt:            updatedAt,
	}

	_, err := tx.NewUpdate().Model(taskCommentTable).Where("task_comment.internal_id = ?", params.TaskComment.Identity.Internal.String()).Exec(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}

		return err
	}

	_, err = tx.NewDelete().Model(&TaskCommentFileTable{}).Where("task_comment_file.task_comment_internal_id = ?", params.TaskComment.Identity.Internal.String()).Exec(context.Background())
	if err != nil {
		return err
	}

	for _, file := range params.TaskComment.Files {
		taskCommentFileTable := &TaskCommentFileTable{
			InternalId:            file.Identity.Internal.String(),
			TaskCommentInternalId: taskCommentTable.InternalId,
			FileInternalId:        file.FileIdentity.Internal.String(),
		}

		_, err = tx.NewInsert().Model(taskCommentFileTable).Exec(context.Background())
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

func (r *TaskCommentBunRepository) DeleteTaskComment(params taskrepo.DeleteTaskCommentParams) error {
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

	_, err := tx.NewDelete().Model(&TaskCommentFileTable{}).Where("task_comment_file.task_comment_internal_id = ?", params.TaskCommentIdentity.Internal.String()).Exec(context.Background())
	if err != nil {
		return err
	}

	_, err = tx.NewDelete().Model(&TaskCommentTable{}).Where("task_comment.internal_id = ?", params.TaskCommentIdentity.Internal.String()).Exec(context.Background())
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
