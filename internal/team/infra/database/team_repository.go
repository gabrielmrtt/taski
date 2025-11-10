package teamdatabase

import (
	"context"
	"database/sql"

	"github.com/gabrielmrtt/taski/internal/core"
	coredatabase "github.com/gabrielmrtt/taski/internal/core/database"
	"github.com/gabrielmrtt/taski/internal/organization"
	organizationdatabase "github.com/gabrielmrtt/taski/internal/organization/infra/database"
	"github.com/gabrielmrtt/taski/internal/team"
	teamrepo "github.com/gabrielmrtt/taski/internal/team/repository"
	"github.com/gabrielmrtt/taski/internal/user"
	userdatabase "github.com/gabrielmrtt/taski/internal/user/infra/database"
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

	Members      []*TeamUserTable                        `bun:"rel:has-many,join:internal_id=team_internal_id"`
	Creator      *userdatabase.UserTable                 `bun:"rel:has-one,join:user_creator_internal_id=internal_id"`
	Editor       *userdatabase.UserTable                 `bun:"rel:has-one,join:user_editor_internal_id=internal_id"`
	Organization *organizationdatabase.OrganizationTable `bun:"rel:has-one,join:organization_internal_id=internal_id"`
}

type TeamUserTable struct {
	bun.BaseModel `bun:"table:team_user,alias:team_user"`

	TeamInternalId string `bun:"team_internal_id,pk,notnull,type:uuid"`
	UserInternalId string `bun:"user_internal_id,pk,notnull,type:uuid"`

	Team *TeamTable              `bun:"rel:has-one,join:team_internal_id=internal_id"`
	User *userdatabase.UserTable `bun:"rel:has-one,join:user_internal_id=internal_id"`
}

func (t *TeamTable) ToEntity() *team.Team {
	var userCreatorIdentity *core.Identity = nil
	var userEditorIdentity *core.Identity = nil

	if t.UserCreatorInternalId != "" {
		identity := core.NewIdentityFromInternal(uuid.MustParse(t.UserCreatorInternalId), user.UserIdentityPrefix)
		userCreatorIdentity = &identity
	}

	if t.UserEditorInternalId != nil {
		identity := core.NewIdentityFromInternal(uuid.MustParse(*t.UserEditorInternalId), user.UserIdentityPrefix)
		userEditorIdentity = &identity
	}

	var members []team.TeamUser = make([]team.TeamUser, 0)
	for _, user := range t.Members {
		members = append(members, *user.ToEntity())
	}

	var creator *user.User = nil
	if t.Creator != nil {
		creator = t.Creator.ToEntity()
	}

	var editor *user.User = nil
	if t.Editor != nil {
		editor = t.Editor.ToEntity()
	}

	var org *organization.Organization = nil
	if t.Organization != nil {
		org = t.Organization.ToEntity()
	}

	var createdAt *core.DateTime = nil
	if t.CreatedAt != 0 {
		createdAt = &core.DateTime{Value: t.CreatedAt}
	}

	var updatedAt *core.DateTime = nil
	if t.UpdatedAt != nil {
		updatedAt = &core.DateTime{Value: *t.UpdatedAt}
	}

	return &team.Team{
		Identity:             core.NewIdentityFromInternal(uuid.MustParse(t.InternalId), team.TeamIdentityPrefix),
		OrganizationIdentity: core.NewIdentityFromInternal(uuid.MustParse(t.OrganizationInternalId), organization.OrganizationIdentityPrefix),
		Name:                 t.Name,
		Description:          t.Description,
		Status:               team.TeamStatuses(t.Status),
		UserCreatorIdentity:  userCreatorIdentity,
		UserEditorIdentity:   userEditorIdentity,
		Creator:              creator,
		Editor:               editor,
		Organization:         org,
		Timestamps:           core.Timestamps{CreatedAt: createdAt, UpdatedAt: updatedAt},
		Members:              members,
	}
}

func (t *TeamUserTable) ToEntity() *team.TeamUser {
	return &team.TeamUser{
		TeamIdentity: core.NewIdentityFromInternal(uuid.MustParse(t.TeamInternalId), team.TeamIdentityPrefix),
		User:         *t.User.ToEntity(),
	}
}

type TeamBunRepository struct {
	db *bun.DB
	tx *coredatabase.TransactionBun
}

func NewTeamBunRepository(connection *bun.DB) *TeamBunRepository {
	return &TeamBunRepository{db: connection, tx: nil}
}

func (r *TeamBunRepository) applyFilters(selectQuery *bun.SelectQuery, filters teamrepo.TeamFilters) *bun.SelectQuery {
	selectQuery = selectQuery.Where("organization_internal_id = ?", filters.OrganizationIdentity.Internal.String())

	if filters.Name != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "name", filters.Name)
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

	return selectQuery
}

func (r *TeamBunRepository) SetTransaction(tx core.Transaction) error {
	if r.tx != nil && !r.tx.IsClosed() {
		return nil
	}

	r.tx = tx.(*coredatabase.TransactionBun)
	return nil
}

func (r *TeamBunRepository) GetTeamByIdentity(params teamrepo.GetTeamByIdentityParams) (*team.Team, error) {
	var team *TeamTable = new(TeamTable)
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(team)
	selectQuery = selectQuery.Relation("Members.User").Relation("Members.User.Credentials").Relation("Members.User.Data")
	selectQuery = coredatabase.ApplyRelations(selectQuery, params.RelationsInput)
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

func (r *TeamBunRepository) PaginateTeamsBy(params teamrepo.PaginateTeamsParams) (*core.PaginationOutput[team.Team], error) {
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
			return &core.PaginationOutput[team.Team]{
				Data:    []team.Team{},
				Page:    page,
				HasMore: false,
				Total:   0,
			}, nil
		}
	}

	var teamEntities []team.Team = make([]team.Team, 0)
	for _, team := range teams {
		teamEntities = append(teamEntities, *team.ToEntity())
	}

	return &core.PaginationOutput[team.Team]{
		Data:    teamEntities,
		Page:    page,
		HasMore: core.HasMorePages(page, countBeforePagination, perPage),
		Total:   countBeforePagination,
	}, nil
}

func (r *TeamBunRepository) StoreTeam(params teamrepo.StoreTeamParams) (*team.Team, error) {
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

	var createdAt *int64 = nil
	if params.Team.Timestamps.CreatedAt != nil {
		createdAt = &params.Team.Timestamps.CreatedAt.Value
	}

	var updatedAt *int64 = nil
	if params.Team.Timestamps.UpdatedAt != nil {
		updatedAt = &params.Team.Timestamps.UpdatedAt.Value
	}

	teamTable := &TeamTable{
		InternalId:             params.Team.Identity.Internal.String(),
		PublicId:               params.Team.Identity.Public,
		Name:                   params.Team.Name,
		Description:            params.Team.Description,
		Status:                 string(params.Team.Status),
		OrganizationInternalId: params.Team.OrganizationIdentity.Internal.String(),
		UserCreatorInternalId:  *userCreatorInternalId,
		CreatedAt:              *createdAt,
		UpdatedAt:              updatedAt,
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

func (r *TeamBunRepository) UpdateTeam(params teamrepo.UpdateTeamParams) error {
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

	var updatedAt *int64 = nil
	if params.Team.Timestamps.UpdatedAt != nil {
		updatedAt = &params.Team.Timestamps.UpdatedAt.Value
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
		UpdatedAt:              updatedAt,
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

func (r *TeamBunRepository) DeleteTeam(params teamrepo.DeleteTeamParams) error {
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
