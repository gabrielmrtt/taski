package projectdatabase

import (
	"context"
	"database/sql"

	"github.com/gabrielmrtt/taski/internal/core"
	coredatabase "github.com/gabrielmrtt/taski/internal/core/database"
	project "github.com/gabrielmrtt/taski/internal/project"
	projectrepo "github.com/gabrielmrtt/taski/internal/project/repository"
	userdatabase "github.com/gabrielmrtt/taski/internal/user/infra/database"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ProjectUserTable struct {
	bun.BaseModel `bun:"table:project_user,alias:project_user"`

	ProjectInternalId string `bun:"project_internal_id,pk,notnull,type:uuid"`
	UserInternalId    string `bun:"user_internal_id,pk,notnull,type:uuid"`
	Status            string `bun:"status,notnull,type:varchar(100)"`

	Project *ProjectTable           `bun:"rel:has-one,join:project_internal_id=internal_id"`
	User    *userdatabase.UserTable `bun:"rel:has-one,join:user_internal_id=internal_id"`
}

func (p *ProjectUserTable) ToEntity() *project.ProjectUser {
	return &project.ProjectUser{
		ProjectIdentity: core.NewIdentityFromInternal(uuid.MustParse(p.ProjectInternalId), project.ProjectIdentityPrefix),
		User:            *p.User.ToEntity(),
		Status:          project.ProjectUserStatuses(p.Status),
	}
}

type ProjectUserBunRepository struct {
	db *bun.DB
	tx *coredatabase.TransactionBun
}

func NewProjectUserBunRepository(connection *bun.DB) *ProjectUserBunRepository {
	return &ProjectUserBunRepository{db: connection, tx: nil}
}

func (r *ProjectUserBunRepository) SetTransaction(tx core.Transaction) error {
	if r.tx != nil && !r.tx.IsClosed() {
		return nil
	}

	r.tx = tx.(*coredatabase.TransactionBun)
	return nil
}

func (r *ProjectUserBunRepository) applyFilters(selectQuery *bun.SelectQuery, filters projectrepo.ProjectUserFilters) *bun.SelectQuery {
	selectQuery = selectQuery.Where("project_internal_id = ?", filters.ProjectIdentity.Internal.String())

	if filters.UserIdentity != nil {
		selectQuery = selectQuery.Where("user_internal_id = ?", filters.UserIdentity.Internal.String())
	}

	if filters.Status != nil {
		selectQuery = selectQuery.Where("status = ?", filters.Status)
	}

	return selectQuery
}

func (r *ProjectUserBunRepository) GetProjectUserByIdentity(params projectrepo.GetProjectUserByIdentityParams) (*project.ProjectUser, error) {
	var projectUser *ProjectUserTable = new(ProjectUserTable)
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(projectUser)
	selectQuery = selectQuery.Relation("User")
	selectQuery = coredatabase.ApplyRelations(selectQuery, params.RelationsInput)
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

func (r *ProjectUserBunRepository) GetProjectUsersByUserIdentity(params projectrepo.GetProjectUsersByUserIdentityParams) ([]project.ProjectUser, error) {
	var projectUsers []ProjectUserTable = make([]ProjectUserTable, 0)
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(&projectUsers)
	selectQuery = coredatabase.ApplyRelations(selectQuery, params.RelationsInput)
	selectQuery = selectQuery.Where("user_internal_id = ?", params.UserIdentity.Internal.String())

	err := selectQuery.Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return []project.ProjectUser{}, nil
		}
	}

	var projectUserEntities []project.ProjectUser = make([]project.ProjectUser, 0)
	for _, projectUser := range projectUsers {
		projectUserEntities = append(projectUserEntities, *projectUser.ToEntity())
	}

	return projectUserEntities, nil
}

func (r *ProjectUserBunRepository) StoreProjectUser(params projectrepo.StoreProjectUserParams) (*project.ProjectUser, error) {
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
		UserInternalId:    params.ProjectUser.User.Identity.Internal.String(),
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

func (r *ProjectUserBunRepository) UpdateProjectUser(params projectrepo.UpdateProjectUserParams) error {
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
		UserInternalId:    params.ProjectUser.User.Identity.Internal.String(),
		Status:            string(params.ProjectUser.Status),
	}).Where("project_internal_id = ? and user_internal_id = ?", params.ProjectUser.ProjectIdentity.Internal.String(), params.ProjectUser.User.Identity.Internal.String()).Exec(context.Background())
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

func (r *ProjectUserBunRepository) DeleteProjectUser(params projectrepo.DeleteProjectUserParams) error {
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

func (r *ProjectUserBunRepository) DeleteAllByUserIdentity(params projectrepo.DeleteAllByUserIdentityParams) error {
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
