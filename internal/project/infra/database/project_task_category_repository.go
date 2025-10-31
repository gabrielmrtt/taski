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

type ProjectTaskCategoryTable struct {
	bun.BaseModel `bun:"table:project_task_category,alias:project_task_category"`

	InternalId        string `bun:"internal_id,pk,notnull,type:uuid"`
	PublicId          string `bun:"public_id,notnull,type:varchar(510)"`
	Name              string `bun:"name,notnull,type:varchar(255)"`
	Color             string `bun:"color,notnull,type:varchar(7)"`
	DeletedAt         *int64 `bun:"deleted_at,type:bigint"`
	ProjectInternalId string `bun:"project_internal_id,notnull,type:uuid"`

	Project *ProjectTable `bun:"rel:has-one,join:project_internal_id=internal_id"`
}

func (p *ProjectTaskCategoryTable) ToEntity() *project.ProjectTaskCategory {
	return &project.ProjectTaskCategory{
		Identity:        core.NewIdentityFromInternal(uuid.MustParse(p.InternalId), project.ProjectTaskCategoryIdentityPrefix),
		ProjectIdentity: core.NewIdentityFromInternal(uuid.MustParse(p.ProjectInternalId), project.ProjectIdentityPrefix),
		Name:            p.Name,
		Color:           p.Color,
	}
}

type ProjectTaskCategoryBunRepository struct {
	db *bun.DB
	tx *coredatabase.TransactionBun
}

func NewProjectTaskCategoryBunRepository(connection *bun.DB) *ProjectTaskCategoryBunRepository {
	return &ProjectTaskCategoryBunRepository{db: connection, tx: nil}
}

func (r *ProjectTaskCategoryBunRepository) SetTransaction(tx core.Transaction) error {
	if r.tx != nil && !r.tx.IsClosed() {
		return nil
	}

	r.tx = tx.(*coredatabase.TransactionBun)
	return nil
}

func (r *ProjectTaskCategoryBunRepository) applyFilters(selectQuery *bun.SelectQuery, filters projectrepo.ProjectTaskCategoryFilters) *bun.SelectQuery {
	if filters.ProjectIdentity != nil {
		selectQuery = selectQuery.Where("project_internal_id = ?", filters.ProjectIdentity.Internal.String())
	}

	if filters.Name != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "name", filters.Name)
	}

	return selectQuery
}

func (r *ProjectTaskCategoryBunRepository) GetProjectTaskCategoryByIdentity(params projectrepo.GetProjectTaskCategoryByIdentityParams) (*project.ProjectTaskCategory, error) {
	var projectTaskCategory *ProjectTaskCategoryTable = new(ProjectTaskCategoryTable)
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(projectTaskCategory)
	selectQuery = coredatabase.ApplyRelations(selectQuery, params.RelationsInput)
	selectQuery = selectQuery.Where("internal_id = ?", params.ProjectTaskCategoryIdentity.Internal.String())
	err := selectQuery.Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	if projectTaskCategory.InternalId == "" {
		return nil, nil
	}

	return projectTaskCategory.ToEntity(), nil
}

func (r *ProjectTaskCategoryBunRepository) PaginateProjectTaskCategoryBy(params projectrepo.PaginateProjectTaskCategoryParams) (*core.PaginationOutput[project.ProjectTaskCategory], error) {
	var projectTaskCategories []ProjectTaskCategoryTable = make([]ProjectTaskCategoryTable, 0)
	var selectQuery *bun.SelectQuery
	var perPage int = 10
	var page int = 1

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(&projectTaskCategories)
	selectQuery = coredatabase.ApplyRelations(selectQuery, params.RelationsInput)
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
			return &core.PaginationOutput[project.ProjectTaskCategory]{
				Data:    []project.ProjectTaskCategory{},
				Page:    page,
				HasMore: false,
				Total:   0,
			}, nil
		}

		return nil, err
	}

	var projectTaskCategoryEntities []project.ProjectTaskCategory = make([]project.ProjectTaskCategory, 0)
	for _, projectTaskCategory := range projectTaskCategories {
		projectTaskCategoryEntities = append(projectTaskCategoryEntities, *projectTaskCategory.ToEntity())
	}

	return &core.PaginationOutput[project.ProjectTaskCategory]{
		Data:    projectTaskCategoryEntities,
		Page:    page,
		HasMore: core.HasMorePages(page, countBeforePagination, perPage),
		Total:   countBeforePagination,
	}, nil
}

func (r *ProjectTaskCategoryBunRepository) StoreProjectTaskCategory(params projectrepo.StoreProjectTaskCategoryParams) (*project.ProjectTaskCategory, error) {
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

	_, err := tx.NewInsert().Model(&ProjectTaskCategoryTable{
		InternalId:        params.ProjectTaskCategory.Identity.Internal.String(),
		PublicId:          params.ProjectTaskCategory.Identity.Public,
		Name:              params.ProjectTaskCategory.Name,
		Color:             params.ProjectTaskCategory.Color,
		ProjectInternalId: params.ProjectTaskCategory.ProjectIdentity.Internal.String(),
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

	return params.ProjectTaskCategory, nil
}

func (r *ProjectTaskCategoryBunRepository) UpdateProjectTaskCategory(params projectrepo.UpdateProjectTaskCategoryParams) error {
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

	_, err := tx.NewUpdate().Model(&ProjectTaskCategoryTable{
		ProjectInternalId: params.ProjectTaskCategory.ProjectIdentity.Internal.String(),
		Name:              params.ProjectTaskCategory.Name,
		Color:             params.ProjectTaskCategory.Color,
		DeletedAt:         params.ProjectTaskCategory.DeletedAt,
	}).Where("internal_id = ?", params.ProjectTaskCategory.Identity.Internal.String()).Exec(context.Background())
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

func (r *ProjectTaskCategoryBunRepository) DeleteProjectTaskCategory(params projectrepo.DeleteProjectTaskCategoryParams) error {
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

	_, err := tx.NewDelete().Model(&ProjectTaskCategoryTable{}).Where("internal_id = ?", params.ProjectTaskCategoryIdentity.Internal.String()).Exec(context.Background())
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
