package taskdatabase

import (
	"context"
	"database/sql"

	"github.com/gabrielmrtt/taski/internal/core"
	coredatabase "github.com/gabrielmrtt/taski/internal/core/database"
	"github.com/gabrielmrtt/taski/internal/project"
	projectdatabase "github.com/gabrielmrtt/taski/internal/project/infra/database"
	"github.com/gabrielmrtt/taski/internal/task"
	taskrepo "github.com/gabrielmrtt/taski/internal/task/repository"
	"github.com/gabrielmrtt/taski/internal/user"
	userdatabase "github.com/gabrielmrtt/taski/internal/user/infra/database"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type TaskUserTable struct {
	bun.BaseModel `bun:"table:task_user,alias:task_user"`

	TaskInternalId string `bun:"task_internal_id,pk,notnull,type:uuid"`
	UserInternalId string `bun:"user_internal_id,pk,notnull,type:uuid"`

	Task *TaskTable              `bun:"rel:has-one,join:task_internal_id=internal_id"`
	User *userdatabase.UserTable `bun:"rel:has-one,join:user_internal_id=internal_id"`
}

func (t *TaskUserTable) ToEntity() *task.TaskUser {
	return &task.TaskUser{
		User: t.User.ToEntity(),
	}
}

type SubTaskTable struct {
	bun.BaseModel `bun:"table:sub_task,alias:sub_task"`

	InternalId     string `bun:"internal_id,pk,notnull,type:uuid"`
	PublicId       string `bun:"public_id,notnull,type:varchar(510)"`
	Name           string `bun:"name,notnull,type:varchar(255)"`
	CompletedAt    *int64 `bun:"completed_at,type:bigint"`
	TaskInternalId string `bun:"task_internal_id,pk,notnull,type:uuid"`

	Task *TaskTable `bun:"rel:has-one,join:task_internal_id=internal_id"`
}

func (s *SubTaskTable) ToEntity() *task.SubTask {
	return &task.SubTask{
		Identity:    core.NewIdentityFromInternal(uuid.MustParse(s.InternalId), task.SubTaskIdentityPrefix),
		Name:        s.Name,
		CompletedAt: s.CompletedAt,
	}
}

type TaskTable struct {
	bun.BaseModel `bun:"table:task,alias:task"`

	InternalId                    string  `bun:"internal_id,pk,notnull,type:uuid"`
	PublicId                      string  `bun:"public_id,notnull,type:varchar(510)"`
	Name                          string  `bun:"name,notnull,type:varchar(255)"`
	Description                   string  `bun:"description,type:varchar(510)"`
	EstimatedMinutes              int16   `bun:"estimated_minutes,notnull,type:int16"`
	PriorityLevel                 int8    `bun:"priority_level,notnull,type:int8"`
	DueDate                       *int64  `bun:"due_date,type:bigint"`
	CompletedAt                   *int64  `bun:"completed_at,type:bigint"`
	Type                          string  `bun:"type,notnull,type:varchar(100)"`
	ProjectTaskStatusInternalId   string  `bun:"project_task_status_internal_id,notnull,type:uuid"`
	ProjectTaskCategoryInternalId string  `bun:"project_task_category_internal_id,notnull,type:uuid"`
	ParentTaskInternalId          *string `bun:"parent_task_internal_id,type:uuid"`
	ProjectInternalId             string  `bun:"project_internal_id,notnull,type:uuid"`
	UserCompletedInternalId       *string `bun:"user_completed_internal_id,type:uuid"`
	UserCreatorInternalId         string  `bun:"user_creator_internal_id,notnull,type:uuid"`
	UserEditorInternalId          *string `bun:"user_editor_internal_id,type:uuid"`
	CreatedAt                     int64   `bun:"created_at,notnull,type:bigint"`
	UpdatedAt                     *int64  `bun:"updated_at,type:bigint"`
	DeletedAt                     *int64  `bun:"deleted_at,type:bigint"`

	ProjectTaskStatus   *projectdatabase.ProjectTaskStatusTable   `bun:"rel:has-one,join:project_task_status_internal_id=internal_id"`
	ProjectTaskCategory *projectdatabase.ProjectTaskCategoryTable `bun:"rel:has-one,join:project_task_category_internal_id=internal_id"`
	ParentTask          *TaskTable                                `bun:"rel:has-one,join:parent_task_internal_id=internal_id"`
	ChildrenTasks       []*TaskTable                              `bun:"rel:has-many,join:internal_id=parent_task_internal_id"`
	Users               []*TaskUserTable                          `bun:"rel:has-many,join:internal_id=task_internal_id"`
	SubTasks            []*SubTaskTable                           `bun:"rel:has-many,join:internal_id=task_internal_id"`
	Project             *projectdatabase.ProjectTable             `bun:"rel:has-one,join:project_internal_id=internal_id"`
	UserCompleted       *userdatabase.UserTable                   `bun:"rel:has-one,join:user_completed_internal_id=internal_id"`
	UserCreator         *userdatabase.UserTable                   `bun:"rel:has-one,join:user_creator_internal_id=internal_id"`
	UserEditor          *userdatabase.UserTable                   `bun:"rel:has-one,join:user_editor_internal_id=internal_id"`
}

func (t *TaskTable) ToEntity() *task.Task {
	var userCreatorIdentity *core.Identity = nil
	if t.UserCreatorInternalId != "" {
		identity := core.NewIdentityFromInternal(uuid.MustParse(t.UserCreatorInternalId), user.UserIdentityPrefix)
		userCreatorIdentity = &identity
	}

	var userEditorIdentity *core.Identity = nil
	if t.UserEditorInternalId != nil {
		identity := core.NewIdentityFromInternal(uuid.MustParse(*t.UserEditorInternalId), user.UserIdentityPrefix)
		userEditorIdentity = &identity
	}

	var userCompletedIdentity *core.Identity = nil
	if t.UserCompletedInternalId != nil {
		identity := core.NewIdentityFromInternal(uuid.MustParse(*t.UserCompletedInternalId), user.UserIdentityPrefix)
		userCompletedIdentity = &identity
	}

	var projectTaskStatus *project.ProjectTaskStatus = nil
	if t.ProjectTaskStatusInternalId != "" {
		projectTaskStatus = t.ProjectTaskStatus.ToEntity()
	}

	var projectTaskCategory *project.ProjectTaskCategory = nil
	if t.ProjectTaskCategoryInternalId != "" {
		projectTaskCategory = t.ProjectTaskCategory.ToEntity()
	}

	var parentTaskIdentity *core.Identity = nil
	if t.ParentTaskInternalId != nil {
		identity := core.NewIdentityFromInternal(uuid.MustParse(*t.ParentTaskInternalId), task.TaskIdentityPrefix)
		parentTaskIdentity = &identity
	}

	var subTasks []*task.SubTask = make([]*task.SubTask, 0)
	for _, subTask := range t.SubTasks {
		subTasks = append(subTasks, subTask.ToEntity())
	}

	var users []*task.TaskUser = make([]*task.TaskUser, 0)
	for _, user := range t.Users {
		users = append(users, user.ToEntity())
	}

	var childrenTasks []*task.Task = make([]*task.Task, 0)
	for _, childTask := range t.ChildrenTasks {
		childrenTasks = append(childrenTasks, childTask.ToEntity())
	}

	return &task.Task{
		Identity:                core.NewIdentityFromInternal(uuid.MustParse(t.InternalId), task.TaskIdentityPrefix),
		ProjectIdentity:         core.NewIdentityFromInternal(uuid.MustParse(t.ProjectInternalId), project.ProjectIdentityPrefix),
		Status:                  projectTaskStatus,
		Category:                projectTaskCategory,
		ParentTaskIdentity:      parentTaskIdentity,
		Type:                    task.TaskType(t.Type),
		Name:                    t.Name,
		Description:             t.Description,
		EstimatedMinutes:        &t.EstimatedMinutes,
		PriorityLevel:           task.TaskPriorityLevels(t.PriorityLevel),
		DueDate:                 t.DueDate,
		CompletedAt:             t.CompletedAt,
		SubTasks:                subTasks,
		ChildrenTasks:           childrenTasks,
		Users:                   users,
		UserCompletedByIdentity: userCompletedIdentity,
		UserCreatorIdentity:     userCreatorIdentity,
		UserEditorIdentity:      userEditorIdentity,
		Timestamps: core.Timestamps{
			CreatedAt: &t.CreatedAt,
			UpdatedAt: t.UpdatedAt,
		},
		DeletedAt: t.DeletedAt,
	}
}

type TaskBunRepository struct {
	db *bun.DB
	tx *coredatabase.TransactionBun
}

func NewTaskBunRepository(connection *bun.DB) *TaskBunRepository {
	return &TaskBunRepository{db: connection, tx: nil}
}

func (r *TaskBunRepository) SetTransaction(tx core.Transaction) error {
	if r.tx != nil && !r.tx.IsClosed() {
		return nil
	}

	r.tx = tx.(*coredatabase.TransactionBun)
	return nil
}

func (r *TaskBunRepository) applyFilters(selectQuery *bun.SelectQuery, filters taskrepo.TaskFilters) *bun.SelectQuery {
	if filters.ProjectIdentity != nil {
		selectQuery = selectQuery.Where("project_internal_id = ?", filters.ProjectIdentity.Internal.String())
	}

	if filters.TaskStatusIdentity != nil {
		selectQuery = selectQuery.Where("project_task_status_internal_id = ?", filters.TaskStatusIdentity.Internal.String())
	}

	if filters.TaskCategoryIdentity != nil {
		selectQuery = selectQuery.Where("project_task_category_internal_id = ?", filters.TaskCategoryIdentity.Internal.String())
	}

	if filters.ParentTaskIdentity != nil {
		selectQuery = selectQuery.Where("parent_task_internal_id = ?", filters.ParentTaskIdentity.Internal.String())
	}

	if filters.Name != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "name", filters.Name)
	}

	if filters.CompletedAt != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "completed_at", filters.CompletedAt)
	}

	if filters.CreatedAt != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "created_at", filters.CreatedAt)
	}

	if filters.UpdatedAt != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "updated_at", filters.UpdatedAt)
	}

	if filters.DueDate != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "due_date", filters.DueDate)
	}

	if filters.Type != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "type", filters.Type)
	}

	if filters.Priority != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "priority_level", filters.Priority)
	}

	return selectQuery
}

func (r *TaskBunRepository) GetTaskByIdentity(params taskrepo.GetTaskByIdentityParams) (*task.Task, error) {
	var task *TaskTable = new(TaskTable)
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	}

	selectQuery = selectQuery.Model(task)
	selectQuery = selectQuery.Relation("ProjectTaskStatus").Relation("ProjectTaskCategory").Relation("ChildrenTasks").Relation("Users").Relation("SubTasks")
	selectQuery = selectQuery.Where("internal_id = ?", params.TaskIdentity.Internal.String())
	err := selectQuery.Scan(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	if task.InternalId == "" {
		return nil, nil
	}

	return task.ToEntity(), nil
}
