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
	var completedAt *core.DateTime = nil
	if s.CompletedAt != nil {
		completedAt = &core.DateTime{Value: *s.CompletedAt}
	}

	return &task.SubTask{
		Identity:    core.NewIdentityFromInternal(uuid.MustParse(s.InternalId), task.SubTaskIdentityPrefix),
		Name:        s.Name,
		CompletedAt: completedAt,
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

	var dueDate *core.DateTime = nil
	if t.DueDate != nil {
		dueDate = &core.DateTime{Value: *t.DueDate}
	}

	var completedAt *core.DateTime = nil
	if t.CompletedAt != nil {
		completedAt = &core.DateTime{Value: *t.CompletedAt}
	}

	var createdAt *core.DateTime = nil
	if t.CreatedAt != 0 {
		createdAt = &core.DateTime{Value: t.CreatedAt}
	}

	var updatedAt *core.DateTime = nil
	if t.UpdatedAt != nil {
		updatedAt = &core.DateTime{Value: *t.UpdatedAt}
	}

	var deletedAt *core.DateTime = nil
	if t.DeletedAt != nil {
		deletedAt = &core.DateTime{Value: *t.DeletedAt}
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
		DueDate:                 dueDate,
		CompletedAt:             completedAt,
		SubTasks:                subTasks,
		ChildrenTasks:           childrenTasks,
		Users:                   users,
		UserCompletedByIdentity: userCompletedIdentity,
		UserCreatorIdentity:     userCreatorIdentity,
		UserEditorIdentity:      userEditorIdentity,
		Timestamps: core.Timestamps{
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		},
		DeletedAt: deletedAt,
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
	if filters.OrganizationIdentity != nil {
		args := []interface{}{filters.OrganizationIdentity.Internal.String()}

		projectQueryString := `
			SELECT project.internal_id FROM project 
			WHERE project.workspace_internal_id IN (
				SELECT workspace.internal_id FROM workspace 
				WHERE workspace.organization_internal_id = ?
			)
		`

		if filters.AuthenticatedUserIdentity != nil {
			projectQueryString = projectQueryString + ` AND project.internal_id IN (
				SELECT project_user.project_internal_id FROM project_user
				WHERE project_user.user_internal_id = ? AND project_user.status = ?
			)`
			args = append(args, filters.AuthenticatedUserIdentity.Internal.String(), project.ProjectUserStatusActive)
		}

		selectQuery = selectQuery.Where("task.project_internal_id IN ("+projectQueryString+")", args...)
	}

	if filters.ProjectIdentity != nil {
		selectQuery = selectQuery.Where("task.project_internal_id = ?", filters.ProjectIdentity.Internal.String())
	}

	if filters.TaskStatusIdentity != nil {
		selectQuery = selectQuery.Where("task.project_task_status_internal_id = ?", filters.TaskStatusIdentity.Internal.String())
	}

	if filters.TaskCategoryIdentity != nil {
		selectQuery = selectQuery.Where("task.project_task_category_internal_id = ?", filters.TaskCategoryIdentity.Internal.String())
	}

	if filters.ParentTaskIdentity != nil {
		selectQuery = selectQuery.Where("task.parent_task_internal_id = ?", filters.ParentTaskIdentity.Internal.String())
	}

	if filters.Name != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "task.name", filters.Name)
	}

	if filters.CompletedAt != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "task.completed_at", filters.CompletedAt)
	}

	if filters.CreatedAt != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "task.created_at", filters.CreatedAt)
	}

	if filters.UpdatedAt != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "task.updated_at", filters.UpdatedAt)
	}

	if filters.DueDate != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "task.due_date", filters.DueDate)
	}

	if filters.Type != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "task.type", filters.Type)
	}

	if filters.Priority != nil {
		selectQuery = coredatabase.ApplyComparableFilter(selectQuery, "task.priority_level", filters.Priority)
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
	selectQuery = selectQuery.Relation("ProjectTaskStatus").Relation("ProjectTaskCategory").Relation("SubTasks")
	selectQuery = coredatabase.ApplyRelations(selectQuery, params.RelationsInput)
	selectQuery = selectQuery.Where("task.internal_id = ?", params.TaskIdentity.Internal.String())

	if params.OrganizationIdentity != nil {
		selectQuery = selectQuery.Where(`
			task.project_internal_id IN (
				SELECT project.internal_id FROM project 
				WHERE project.workspace_internal_id IN (
					SELECT workspace.internal_id FROM workspace 
					WHERE workspace.organization_internal_id = ?
				)
			)`, params.OrganizationIdentity.Internal.String())
	}

	if params.ProjectIdentity != nil {
		selectQuery = selectQuery.Where("task.project_internal_id = ?", params.ProjectIdentity.Internal.String())
	}

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

func (r *TaskBunRepository) PaginateTasksBy(params taskrepo.PaginateTasksParams) (*core.PaginationOutput[task.Task], error) {
	var tasks []*TaskTable = make([]*TaskTable, 0)
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

	selectQuery = selectQuery.Model(&tasks)
	selectQuery = selectQuery.Relation("ProjectTaskStatus").Relation("ProjectTaskCategory").Relation("SubTasks")
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

	var taskEntities []task.Task = make([]task.Task, 0)
	for _, task := range tasks {
		taskEntities = append(taskEntities, *task.ToEntity())
	}

	return &core.PaginationOutput[task.Task]{
		Data:    taskEntities,
		Page:    page,
		HasMore: core.HasMorePages(page, countBeforePagination, perPage),
		Total:   countBeforePagination,
	}, nil
}

func (r *TaskBunRepository) GetTasksByParentTaskIdentity(params taskrepo.GetTasksByParentTaskIdentityParams) ([]*task.Task, error) {
	var tasks []*TaskTable = make([]*TaskTable, 0)
	var selectQuery *bun.SelectQuery

	if r.tx != nil && !r.tx.IsClosed() {
		selectQuery = r.tx.Tx.NewSelect()
	} else {
		selectQuery = r.db.NewSelect()
	}

	selectQuery = selectQuery.Model(&tasks)
	selectQuery = selectQuery.Relation("ProjectTaskStatus").Relation("ProjectTaskCategory").Relation("SubTasks")
	selectQuery = coredatabase.ApplyRelations(selectQuery, params.RelationsInput)
	selectQuery = selectQuery.Where("task.parent_task_internal_id = ?", params.ParentTaskIdentity.Internal.String())
	err := selectQuery.Scan(context.Background())
	if err != nil {
		return nil, err
	}

	var taskEntities []*task.Task = make([]*task.Task, 0)
	for _, task := range tasks {
		taskEntities = append(taskEntities, task.ToEntity())
	}

	return taskEntities, nil
}

func (r *TaskBunRepository) AddSubTask(params taskrepo.AddSubTaskParams) error {
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

	var completedAt *int64 = nil
	if params.SubTask.CompletedAt != nil {
		completedAt = &params.SubTask.CompletedAt.Value
	}

	var subTaskTable *SubTaskTable = &SubTaskTable{
		InternalId:     params.SubTask.Identity.Internal.String(),
		PublicId:       params.SubTask.Identity.Public,
		Name:           params.SubTask.Name,
		CompletedAt:    completedAt,
		TaskInternalId: params.Task.Identity.Internal.String(),
	}

	_, err := tx.NewInsert().Model(subTaskTable).Exec(context.Background())
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

func (r *TaskBunRepository) UpdateSubTask(params taskrepo.UpdateSubTaskParams) error {
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

	var completedAt *int64 = nil
	if params.SubTask.CompletedAt != nil {
		completedAt = &params.SubTask.CompletedAt.Value
	}

	_, err := tx.NewUpdate().Model(&SubTaskTable{
		InternalId:  params.SubTask.Identity.Internal.String(),
		PublicId:    params.SubTask.Identity.Public,
		Name:        params.SubTask.Name,
		CompletedAt: completedAt,
	}).Where("sub_task.internal_id = ?", params.SubTask.Identity.Internal.String()).Exec(context.Background())
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

func (r *TaskBunRepository) RemoveSubTask(params taskrepo.RemoveSubTaskParams) error {
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

	_, err := tx.NewDelete().Model(&SubTaskTable{}).Where("sub_task.internal_id = ?", params.SubTask.Identity.Internal.String()).Exec(context.Background())
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

func (r *TaskBunRepository) StoreTask(params taskrepo.StoreTaskParams) (*task.Task, error) {
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

	var parentTaskInternalId *string = nil
	if params.Task.ParentTaskIdentity != nil {
		internalId := params.Task.ParentTaskIdentity.Internal.String()
		parentTaskInternalId = &internalId
	}

	var userEditorInternalId *string = nil
	if params.Task.UserEditorIdentity != nil {
		internalId := params.Task.UserEditorIdentity.Internal.String()
		userEditorInternalId = &internalId
	}

	var dueDate *int64 = nil
	if params.Task.DueDate != nil {
		dueDate = &params.Task.DueDate.Value
	}

	var completedAt *int64 = nil
	if params.Task.CompletedAt != nil {
		completedAt = &params.Task.CompletedAt.Value
	}

	var createdAt *int64 = nil
	if params.Task.Timestamps.CreatedAt != nil {
		createdAt = &params.Task.Timestamps.CreatedAt.Value
	}

	var updatedAt *int64 = nil
	if params.Task.Timestamps.UpdatedAt != nil {
		updatedAt = &params.Task.Timestamps.UpdatedAt.Value
	}

	var deletedAt *int64 = nil
	if params.Task.DeletedAt != nil {
		deletedAt = &params.Task.DeletedAt.Value
	}

	_, err := tx.NewInsert().Model(&TaskTable{
		InternalId:                    params.Task.Identity.Internal.String(),
		PublicId:                      params.Task.Identity.Public,
		Name:                          params.Task.Name,
		Description:                   params.Task.Description,
		EstimatedMinutes:              *params.Task.EstimatedMinutes,
		PriorityLevel:                 int8(params.Task.PriorityLevel),
		DueDate:                       dueDate,
		CompletedAt:                   completedAt,
		Type:                          string(params.Task.Type),
		ProjectTaskStatusInternalId:   params.Task.Status.Identity.Internal.String(),
		ProjectTaskCategoryInternalId: params.Task.Category.Identity.Internal.String(),
		ParentTaskInternalId:          parentTaskInternalId,
		ProjectInternalId:             params.Task.ProjectIdentity.Internal.String(),
		UserCreatorInternalId:         params.Task.UserCreatorIdentity.Internal.String(),
		UserEditorInternalId:          userEditorInternalId,
		CreatedAt:                     *createdAt,
		UpdatedAt:                     updatedAt,
		DeletedAt:                     deletedAt,
	}).Exec(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	for _, subTask := range params.Task.SubTasks {
		err = r.AddSubTask(taskrepo.AddSubTaskParams{
			Task:    params.Task,
			SubTask: subTask,
		})
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

	return params.Task, nil
}

func (r *TaskBunRepository) UpdateTask(params taskrepo.UpdateTaskParams) error {
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

	var parentTaskInternalId *string = nil
	if params.Task.ParentTaskIdentity != nil {
		internalId := params.Task.ParentTaskIdentity.Internal.String()
		parentTaskInternalId = &internalId
	}

	var userEditorInternalId *string = nil
	if params.Task.UserEditorIdentity != nil {
		internalId := params.Task.UserEditorIdentity.Internal.String()
		userEditorInternalId = &internalId
	}

	var dueDate *int64 = nil
	if params.Task.DueDate != nil {
		dueDate = &params.Task.DueDate.Value
	}

	var completedAt *int64 = nil
	if params.Task.CompletedAt != nil {
		completedAt = &params.Task.CompletedAt.Value
	}

	var updatedAt *int64 = nil
	if params.Task.Timestamps.UpdatedAt != nil {
		updatedAt = &params.Task.Timestamps.UpdatedAt.Value
	}

	var deletedAt *int64 = nil
	if params.Task.DeletedAt != nil {
		deletedAt = &params.Task.DeletedAt.Value
	}

	taskTable := &TaskTable{
		InternalId:                    params.Task.Identity.Internal.String(),
		PublicId:                      params.Task.Identity.Public,
		Name:                          params.Task.Name,
		Description:                   params.Task.Description,
		EstimatedMinutes:              *params.Task.EstimatedMinutes,
		PriorityLevel:                 int8(params.Task.PriorityLevel),
		DueDate:                       dueDate,
		CompletedAt:                   completedAt,
		Type:                          string(params.Task.Type),
		ProjectTaskStatusInternalId:   params.Task.Status.Identity.Internal.String(),
		ProjectTaskCategoryInternalId: params.Task.Category.Identity.Internal.String(),
		ParentTaskInternalId:          parentTaskInternalId,
		ProjectInternalId:             params.Task.ProjectIdentity.Internal.String(),
		UserCreatorInternalId:         params.Task.UserCreatorIdentity.Internal.String(),
		UserEditorInternalId:          userEditorInternalId,
		UpdatedAt:                     updatedAt,
		DeletedAt:                     deletedAt,
	}

	_, err := tx.NewUpdate().Model(taskTable).Where("task.internal_id = ?", params.Task.Identity.Internal.String()).Exec(context.Background())
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

func (r *TaskBunRepository) DeleteTask(params taskrepo.DeleteTaskParams) error {
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

	_, err := tx.NewDelete().Model(&TaskTable{}).Where("task.internal_id = ?", params.TaskIdentity.Internal.String()).Exec(context.Background())
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
