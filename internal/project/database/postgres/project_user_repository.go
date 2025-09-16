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
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ProjectUserTable struct {
	bun.BaseModel `bun:"table:project_user,alias:project_user"`

	ProjectInternalId string `bun:"project_internal_id,pk,notnull,type:uuid"`
	UserInternalId    string `bun:"user_internal_id,pk,notnull,type:uuid"`
	Status            string `bun:"status,notnull,type:varchar(100)"`

	Project *ProjectTable                     `bun:"rel:has-one,join:project_internal_id=internal_id"`
	User    *user_database_postgres.UserTable `bun:"rel:has-one,join:user_internal_id=internal_id"`
}

func (p *ProjectUserTable) ToEntity() *project_core.ProjectUser {
	return &project_core.ProjectUser{
		ProjectIdentity: core.NewIdentityFromInternal(uuid.MustParse(p.ProjectInternalId), project_core.ProjectIdentityPrefix),
		UserIdentity:    core.NewIdentityFromInternal(uuid.MustParse(p.UserInternalId), user_core.UserIdentityPrefix),
		Status:          project_core.ProjectUserStatuses(p.Status),
	}
}

type ProjectUserPostgresRepository struct {
	db *bun.DB
	tx *core_database_postgres.TransactionPostgres
}

func NewProjectUserPostgresRepository() *ProjectUserPostgresRepository {
	return &ProjectUserPostgresRepository{db: core_database_postgres.DB, tx: nil}
}

func (r *ProjectUserPostgresRepository) SetTransaction(tx core.Transaction) error {
	r.tx = tx.(*core_database_postgres.TransactionPostgres)
	return nil
}

func (r *ProjectUserPostgresRepository) applyFilters(selectQuery *bun.SelectQuery, filters project_repositories.ProjectUserFilters) *bun.SelectQuery {
	selectQuery = selectQuery.Where("project_internal_id = ?", filters.ProjectIdentity.Internal.String())

	if filters.UserIdentity != nil {
		selectQuery = selectQuery.Where("user_internal_id = ?", filters.UserIdentity.Internal.String())
	}

	if filters.Status != nil {
		selectQuery = selectQuery.Where("status = ?", filters.Status)
	}

	return selectQuery
}

func (r *ProjectUserPostgresRepository) GetProjectUserByIdentity(params project_repositories.GetProjectUserByIdentityParams) (*project_core.ProjectUser, error) {
	var projectUser *ProjectUserTable = new(ProjectUserTable)
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(projectUser)
	selectQuery = core_database_postgres.ApplyRelations(selectQuery, params.RelationsInput)
	selectQuery = selectQuery.Where("project_internal_id = ? and user_internal_id = ?", params.ProjectIdentity.Internal.String(), params.UserIdentity.Internal.String())

	err := selectQuery.Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	if projectUser.ProjectInternalId == "" {
		return nil, nil
	}

	return projectUser.ToEntity(), nil
}

func (r *ProjectUserPostgresRepository) GetProjectUsersByUserIdentity(params project_repositories.GetProjectUsersByUserIdentityParams) ([]project_core.ProjectUser, error) {
	var projectUsers []ProjectUserTable = make([]ProjectUserTable, 0)
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(&projectUsers)
	selectQuery = core_database_postgres.ApplyRelations(selectQuery, *params.RelationsInput)
	selectQuery = selectQuery.Where("user_internal_id = ?", params.UserIdentity.Internal.String())

	err := selectQuery.Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return []project_core.ProjectUser{}, nil
		}
	}

	var projectUserEntities []project_core.ProjectUser = make([]project_core.ProjectUser, 0)
	for _, projectUser := range projectUsers {
		projectUserEntities = append(projectUserEntities, *projectUser.ToEntity())
	}

	return projectUserEntities, nil
}

func (r *ProjectUserPostgresRepository) StoreProjectUser(params project_repositories.StoreProjectUserParams) (*project_core.ProjectUser, error) {
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

	_, err := tx.NewInsert().Model(&ProjectUserTable{
		ProjectInternalId: params.ProjectUser.ProjectIdentity.Internal.String(),
		UserInternalId:    params.ProjectUser.UserIdentity.Internal.String(),
		Status:            string(params.ProjectUser.Status),
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

	return params.ProjectUser, nil
}

func (r *ProjectUserPostgresRepository) UpdateProjectUser(params project_repositories.UpdateProjectUserParams) error {
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

	_, err := tx.NewUpdate().Model(&ProjectUserTable{
		ProjectInternalId: params.ProjectUser.ProjectIdentity.Internal.String(),
		UserInternalId:    params.ProjectUser.UserIdentity.Internal.String(),
		Status:            string(params.ProjectUser.Status),
	}).Where("project_internal_id = ? and user_internal_id = ?", params.ProjectUser.ProjectIdentity.Internal.String(), params.ProjectUser.UserIdentity.Internal.String()).Exec(context.Background())
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

func (r *ProjectUserPostgresRepository) DeleteProjectUser(params project_repositories.DeleteProjectUserParams) error {
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

	_, err := tx.NewDelete().Model(&ProjectUserTable{}).Where("project_internal_id = ? and user_internal_id = ?", params.ProjectIdentity.Internal.String(), params.UserIdentity.Internal.String()).Exec(context.Background())
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

func (r *ProjectUserPostgresRepository) DeleteAllByUserIdentity(params project_repositories.DeleteAllByUserIdentityParams) error {
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

	_, err := tx.NewDelete().Model(&ProjectUserTable{}).Where("user_internal_id = ?", params.UserIdentity.Internal.String()).Exec(context.Background())
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
