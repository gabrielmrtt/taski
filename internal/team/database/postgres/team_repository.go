package team_database_postgres

import (
	"context"
	"database/sql"

	"github.com/gabrielmrtt/taski/internal/core"
	core_database_postgres "github.com/gabrielmrtt/taski/internal/core/database/postgres"
	organization_core "github.com/gabrielmrtt/taski/internal/organization"
	team_core "github.com/gabrielmrtt/taski/internal/team"
	team_repositories "github.com/gabrielmrtt/taski/internal/team/repositories"
	user_core "github.com/gabrielmrtt/taski/internal/user"
	user_database_postgres "github.com/gabrielmrtt/taski/internal/user/database/postgres"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type TeamTable struct {
	bun.BaseModel `bun:"table:team,alias:team"`

	InternalId             string  `bun:"internal_id,pk,notnull,type:uuid"`
	PublicId               string  `bun:"public_id,notnull,type:varchar(510)"`
	Name                   string  `bun:"name,notnull,type:varchar(255)"`
	Description            string  `bun:"description,type:varchar(510)"`
	Status                 string  `bun:"status,notnull,type:varchar(100)"`
	OrganizationInternalId string  `bun:"organization_internal_id,notnull,type:uuid"`
	UserCreatorInternalId  string  `bun:"user_creator_internal_id,notnull,type:uuid"`
	UserEditorInternalId   *string `bun:"user_editor_internal_id,type:uuid"`
	CreatedAt              int64   `bun:"created_at,notnull,type:bigint"`
	UpdatedAt              *int64  `bun:"updated_at,type:bigint"`

	Members []*TeamUserTable `bun:"rel:has-many,join:internal_id=team_internal_id"`
}

type TeamUserTable struct {
	bun.BaseModel `bun:"table:team_user,alias:team_user"`

	TeamInternalId string `bun:"team_internal_id,pk,notnull,type:uuid"`
	UserInternalId string `bun:"user_internal_id,pk,notnull,type:uuid"`

	Team *TeamTable                        `bun:"rel:has-one,join:team_internal_id=internal_id"`
	User *user_database_postgres.UserTable `bun:"rel:has-one,join:user_internal_id=internal_id"`
}

func (t *TeamTable) ToEntity() *team_core.Team {
	var userCreatorIdentity *core.Identity = nil
	var userEditorIdentity *core.Identity = nil

	if t.UserCreatorInternalId != "" {
		identity := core.NewIdentityFromInternal(uuid.MustParse(t.UserCreatorInternalId), user_core.UserIdentityPrefix)
		userCreatorIdentity = &identity
	}

	if t.UserEditorInternalId != nil {
		identity := core.NewIdentityFromInternal(uuid.MustParse(*t.UserEditorInternalId), user_core.UserIdentityPrefix)
		userEditorIdentity = &identity
	}

	var members []team_core.TeamUser = make([]team_core.TeamUser, 0)
	for _, user := range t.Members {
		members = append(members, *user.ToEntity())
	}

	return &team_core.Team{
		Identity:             core.NewIdentityFromInternal(uuid.MustParse(t.InternalId), team_core.TeamIdentityPrefix),
		OrganizationIdentity: core.NewIdentityFromInternal(uuid.MustParse(t.OrganizationInternalId), organization_core.OrganizationIdentityPrefix),
		Name:                 t.Name,
		Description:          t.Description,
		Status:               team_core.TeamStatuses(t.Status),
		UserCreatorIdentity:  userCreatorIdentity,
		UserEditorIdentity:   userEditorIdentity,
		Timestamps:           core.Timestamps{CreatedAt: &t.CreatedAt, UpdatedAt: t.UpdatedAt},
		Members:              members,
	}
}

func (t *TeamUserTable) ToEntity() *team_core.TeamUser {
	return &team_core.TeamUser{
		TeamIdentity: core.NewIdentityFromInternal(uuid.MustParse(t.TeamInternalId), team_core.TeamIdentityPrefix),
		User:         *t.User.ToEntity(),
	}
}

type TeamPostgresRepository struct {
	db *bun.DB
	tx *core_database_postgres.TransactionPostgres
}

func NewTeamPostgresRepository() *TeamPostgresRepository {
	return &TeamPostgresRepository{db: core_database_postgres.DB, tx: nil}
}

func (r *TeamPostgresRepository) applyFilters(selectQuery *bun.SelectQuery, filters team_repositories.TeamFilters) *bun.SelectQuery {
	selectQuery = selectQuery.Where("organization_internal_id = ?", filters.OrganizationIdentity.Internal.String())

	if filters.Name != nil {
		selectQuery = core_database_postgres.ApplyComparableFilter(selectQuery, "name", filters.Name)
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

	return selectQuery
}

func (r *TeamPostgresRepository) SetTransaction(tx core.Transaction) error {
	r.tx = tx.(*core_database_postgres.TransactionPostgres)
	return nil
}

func (r *TeamPostgresRepository) GetTeamByIdentity(params team_repositories.GetTeamByIdentityParams) (*team_core.Team, error) {
	var team *TeamTable = new(TeamTable)
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(team)
	selectQuery = selectQuery.Relation("Members.User").Relation("Members.User.Credentials").Relation("Members.User.Data")
	selectQuery = selectQuery.Where("internal_id = ?", params.TeamIdentity.Internal.String())
	err := selectQuery.Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	if team.InternalId == "" {
		return nil, nil
	}

	return team.ToEntity(), nil
}

func (r *TeamPostgresRepository) PaginateTeamsBy(params team_repositories.PaginateTeamsParams) (*core.PaginationOutput[team_core.Team], error) {
	var teams []*TeamTable = make([]*TeamTable, 0)
	var selectQuery *bun.SelectQuery
	var perPage int = 10
	var page int = 1

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(&teams)
	selectQuery = selectQuery.Relation("Members.User").Relation("Members.User.Credentials").Relation("Members.User.Data")
	selectQuery = r.applyFilters(selectQuery, params.Filters)
	selectQuery = core_database_postgres.ApplySort(selectQuery, *params.SortInput)
	countBeforePagination, err := selectQuery.Count(context.Background())
	if err != nil {
		return nil, err
	}

	selectQuery = core_database_postgres.ApplyPagination(selectQuery, params.Pagination)
	err = selectQuery.Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return &core.PaginationOutput[team_core.Team]{
				Data:    []team_core.Team{},
				Page:    page,
				HasMore: false,
				Total:   0,
			}, nil
		}
	}

	var teamEntities []team_core.Team = make([]team_core.Team, 0)
	for _, team := range teams {
		teamEntities = append(teamEntities, *team.ToEntity())
	}

	return &core.PaginationOutput[team_core.Team]{
		Data:    teamEntities,
		Page:    page,
		HasMore: core.HasMorePages(page, countBeforePagination, perPage),
		Total:   countBeforePagination,
	}, nil
}

func (r *TeamPostgresRepository) StoreTeam(params team_repositories.StoreTeamParams) (*team_core.Team, error) {
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

	var userCreatorInternalId *string
	if params.Team.UserCreatorIdentity != nil {
		identity := params.Team.UserCreatorIdentity.Internal.String()
		userCreatorInternalId = &identity
	}

	teamTable := &TeamTable{
		InternalId:             params.Team.Identity.Internal.String(),
		PublicId:               params.Team.Identity.Public,
		Name:                   params.Team.Name,
		Description:            params.Team.Description,
		Status:                 string(params.Team.Status),
		OrganizationInternalId: params.Team.OrganizationIdentity.Internal.String(),
		UserCreatorInternalId:  *userCreatorInternalId,
		CreatedAt:              *params.Team.Timestamps.CreatedAt,
		UpdatedAt:              params.Team.Timestamps.UpdatedAt,
	}

	_, err := tx.NewInsert().Model(teamTable).Exec(context.Background())
	if err != nil {
		return nil, err
	}

	for _, user := range params.Team.Members {
		teamUserTable := &TeamUserTable{
			TeamInternalId: params.Team.Identity.Internal.String(),
			UserInternalId: user.User.Identity.Internal.String(),
		}

		_, err = tx.NewInsert().Model(teamUserTable).Exec(context.Background())
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

	return params.Team, nil
}

func (r *TeamPostgresRepository) UpdateTeam(params team_repositories.UpdateTeamParams) error {
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
	if params.Team.UserEditorIdentity != nil {
		identity := params.Team.UserEditorIdentity.Internal.String()
		userEditorInternalId = &identity
	}

	teamTable := &TeamTable{
		InternalId:             params.Team.Identity.Internal.String(),
		PublicId:               params.Team.Identity.Public,
		Name:                   params.Team.Name,
		Description:            params.Team.Description,
		Status:                 string(params.Team.Status),
		OrganizationInternalId: params.Team.OrganizationIdentity.Internal.String(),
		UserCreatorInternalId:  params.Team.UserCreatorIdentity.Internal.String(),
		UserEditorInternalId:   userEditorInternalId,
		UpdatedAt:              params.Team.Timestamps.UpdatedAt,
	}

	_, err := tx.NewUpdate().Model(teamTable).Where("internal_id = ?", params.Team.Identity.Internal.String()).Exec(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
	}

	_, err = tx.NewDelete().Model(&TeamUserTable{}).Where("team_internal_id = ?", params.Team.Identity.Internal.String()).Exec(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
	}

	for _, user := range params.Team.Members {
		teamUserTable := &TeamUserTable{
			TeamInternalId: params.Team.Identity.Internal.String(),
			UserInternalId: user.User.Identity.Internal.String(),
		}

		_, err = tx.NewInsert().Model(teamUserTable).Exec(context.Background())
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

func (r *TeamPostgresRepository) DeleteTeam(params team_repositories.DeleteTeamParams) error {
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

	_, err := tx.NewDelete().Model(&TeamTable{}).Where("internal_id = ?", params.TeamIdentity.Internal.String()).Exec(context.Background())
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
