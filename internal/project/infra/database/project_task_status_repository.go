package projectdatabase

import (
	"context"
	"database/sql"

	"github.com/gabrielmrtt/taski/internal/core"
	coredatabase "github.com/gabrielmrtt/taski/internal/core/database"
	"github.com/gabrielmrtt/taski/internal/project"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ProjectTaskStatusTable struct {
	bun.BaseModel `bun:"table:project_task_status,alias:project_task_status"`

	InternalId               string `bun:"internal_id,pk,notnull,type:uuid"`
	PublicId                 string `bun:"public_id,notnull,type:varchar(510)"`
	Name                     string `bun:"name,notnull,type:varchar(255)"`
	Color                    string `bun:"color,notnull,type:varchar(7)"`
	StatusOrder              *int8  `bun:"status_order,type:int8"`
	ShouldSetTaskToCompleted bool   `bun:"should_set_task_to_completed,notnull,type:boolean"`
	IsDefault                bool   `bun:"is_default,notnull,type:boolean"`
	ProjectInternalId        string `bun:"project_internal_id,notnull,type:uuid"`
	DeletedAt                *int64 `bun:"deleted_at,type:bigint"`

	Project *ProjectTable `bun:"rel:has-one,join:project_internal_id=internal_id"`
}

func (p *ProjectTaskStatusTable) ToEntity() *project.ProjectTaskStatus {
	return &project.ProjectTaskStatus{
		Identity:                 core.NewIdentityFromInternal(uuid.MustParse(p.InternalId), project.ProjectTaskStatusIdentityPrefix),
		ProjectIdentity:          core.NewIdentityFromInternal(uuid.MustParse(p.ProjectInternalId), project.ProjectIdentityPrefix),
		Name:                     p.Name,
		Color:                    p.Color,
		Order:                    p.StatusOrder,
		ShouldSetTaskToCompleted: p.ShouldSetTaskToCompleted,
		IsDefault:                p.IsDefault,
	}
}

type ProjectTaskStatusBunRepository struct {
	db *bun.DB
	tx *coredatabase.TransactionBun
}

func NewProjectTaskStatusBunRepository(connection *bun.DB) *ProjectTaskStatusBunRepository {
	return &ProjectTaskStatusBunRepository{db: connection, tx: nil}
}

func (r *ProjectTaskStatusBunRepository) SetTransaction(tx core.Transaction) error {
	if r.tx != nil && !r.tx.IsClosed() {
		return nil
	}

	r.tx = tx.(*coredatabase.TransactionBun)
	return nil
}

func (r *ProjectTaskStatusBunRepository) applyFilters(selectQuery *bun.SelectQuery, filters projectrepo.ProjectTaskStatusFilters) *bun.SelectQuery {
	if filters.ProjectIdentity != nil {
		selectQuery = selectQuery.Where("project_internal_id = ?", filters.ProjectIdentity.Internal.String())
	}

	if filters.IsDefault != nil {
		selectQuery = selectQuery.Where("is_default = ?", filters.IsDefault)
	}

	if filters.ShouldSetTaskToCompleted != nil {
		selectQuery = selectQuery.Where("should_set_task_to_completed = ?", filters.ShouldSetTaskToCompleted)
	}

	if filters.Name != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "name", filters.Name)
	}

	if filters.Order != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "status_order", filters.Order)
	}

	return selectQuery
}

func (r *ProjectTaskStatusBunRepository) GetLastTaskStatusOrder(params projectrepo.GetLastTaskStatusOrderParams) (int8, error) {
	var projectTaskStatus *ProjectTaskStatusTable = new(ProjectTaskStatusTable)
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(projectTaskStatus)

	if params.ProjectIdentity != nil {
		selectQuery = selectQuery.Where("project_internal_id = ?", params.ProjectIdentity.Internal.String())
	}

	selectQuery = selectQuery.Order("status_order DESC NULLS LAST")
	selectQuery = selectQuery.Limit(1)

	err := selectQuery.Scan(context.Background())
	if err != nil {
		return 0, err
	}

	if projectTaskStatus.InternalId == "" {
		return 0, nil
	}

	if projectTaskStatus.StatusOrder == nil {
		return 0, nil
	}

	return *projectTaskStatus.StatusOrder, nil
}

func (r *ProjectTaskStatusBunRepository) GetProjectTaskStatusByIdentity(params projectrepo.GetProjectTaskStatusByIdentityParams) (*project.ProjectTaskStatus, error) {
	var projectTaskStatus *ProjectTaskStatusTable = new(ProjectTaskStatusTable)
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(projectTaskStatus)

	if params.ProjectIdentity != nil {
		selectQuery = selectQuery.Where("project_internal_id = ?", params.ProjectIdentity.Internal.String())
	}

	if params.ProjectTaskStatusIdentity != nil {
		selectQuery = selectQuery.Where("internal_id = ?", params.ProjectTaskStatusIdentity.Internal.String())
	}

	if params.IsDefault != nil {
		selectQuery = selectQuery.Where("is_default = ?", params.IsDefault)
	}

	if params.ShouldSetTaskToCompleted != nil {
		selectQuery = selectQuery.Where("should_set_task_to_completed = ?", params.ShouldSetTaskToCompleted)
	}

	err := selectQuery.Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	if projectTaskStatus.InternalId == "" {
		return nil, nil
	}

	return projectTaskStatus.ToEntity(), nil
}

func (r *ProjectTaskStatusBunRepository) ListProjectTaskStatusesBy(params projectrepo.ListProjectTaskStatusesByParams) ([]project.ProjectTaskStatus, error) {
	var projectTaskStatuses []ProjectTaskStatusTable = make([]ProjectTaskStatusTable, 0)
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(&projectTaskStatuses)
	selectQuery = r.applyFilters(selectQuery, params.Filters)
	err := selectQuery.Scan(context.Background())
	if err != nil {
		return nil, err
	}

	var projectTaskStatusEntities []project.ProjectTaskStatus = make([]project.ProjectTaskStatus, 0)
	for _, projectTaskStatus := range projectTaskStatuses {
		projectTaskStatusEntities = append(projectTaskStatusEntities, *projectTaskStatus.ToEntity())
	}

	return projectTaskStatusEntities, nil
}

func (r *ProjectTaskStatusBunRepository) PaginateProjectTaskStatusesBy(params projectrepo.PaginateProjectTaskStatusesParams) (*core.PaginationOutput[project.ProjectTaskStatus], error) {
	var projectTaskStatuses []ProjectTaskStatusTable = make([]ProjectTaskStatusTable, 0)
	var selectQuery *bun.SelectQuery
	var perPage int = 10
	var page int = 1

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(&projectTaskStatuses)
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
			return nil, nil
		}
	}

	var projectTaskStatusEntities []project.ProjectTaskStatus = make([]project.ProjectTaskStatus, 0)
	for _, projectTaskStatus := range projectTaskStatuses {
		projectTaskStatusEntities = append(projectTaskStatusEntities, *projectTaskStatus.ToEntity())
	}

	return &core.PaginationOutput[project.ProjectTaskStatus]{
		Data:    projectTaskStatusEntities,
		Page:    page,
		HasMore: core.HasMorePages(page, countBeforePagination, perPage),
		Total:   countBeforePagination,
	}, nil
}

func (r *ProjectTaskStatusBunRepository) StoreProjectTaskStatus(params projectrepo.StoreProjectTaskStatusParams) (*project.ProjectTaskStatus, error) {
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

	_, err := tx.NewInsert().Model(&ProjectTaskStatusTable{
		InternalId:               params.ProjectTaskStatus.Identity.Internal.String(),
		PublicId:                 params.ProjectTaskStatus.Identity.Public,
		Name:                     params.ProjectTaskStatus.Name,
		Color:                    params.ProjectTaskStatus.Color,
		StatusOrder:              params.ProjectTaskStatus.Order,
		ShouldSetTaskToCompleted: params.ProjectTaskStatus.ShouldSetTaskToCompleted,
		IsDefault:                params.ProjectTaskStatus.IsDefault,
		ProjectInternalId:        params.ProjectTaskStatus.ProjectIdentity.Internal.String(),
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

	return params.ProjectTaskStatus, nil
}

func (r *ProjectTaskStatusBunRepository) UpdateProjectTaskStatus(params projectrepo.UpdateProjectTaskStatusParams) error {
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

	_, err := tx.NewUpdate().Model(&ProjectTaskStatusTable{
		ProjectInternalId:        params.ProjectTaskStatus.ProjectIdentity.Internal.String(),
		Name:                     params.ProjectTaskStatus.Name,
		Color:                    params.ProjectTaskStatus.Color,
		StatusOrder:              params.ProjectTaskStatus.Order,
		ShouldSetTaskToCompleted: params.ProjectTaskStatus.ShouldSetTaskToCompleted,
		IsDefault:                params.ProjectTaskStatus.IsDefault,
	}).Where("internal_id = ?", params.ProjectTaskStatus.Identity.Internal.String()).Exec(context.Background())
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

func (r *ProjectTaskStatusBunRepository) DeleteProjectTaskStatus(params projectrepo.DeleteProjectTaskStatusParams) error {
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

	_, err := tx.NewDelete().Model(&ProjectTaskStatusTable{}).Where("internal_id = ?", params.ProjectTaskStatusIdentity.Internal.String()).Exec(context.Background())
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
