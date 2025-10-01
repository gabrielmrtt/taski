package role_core

const RoleIdentityPrefix = "rol"

type PermissionSlugs string

const (
	OrganizationsView        PermissionSlugs = "organizations:view"
	OrganizationsUpdate      PermissionSlugs = "organizations:update"
	OrganizationsDelete      PermissionSlugs = "organizations:delete"
	OrganizationsUsersView   PermissionSlugs = "organizations:users:view"
	OrganizationsUsersCreate PermissionSlugs = "organizations:users:create"
	OrganizationsUsersUpdate PermissionSlugs = "organizations:users:update"
	OrganizationsUsersDelete PermissionSlugs = "organizations:users:delete"
	RolesView                PermissionSlugs = "roles:view"
	RolesCreate              PermissionSlugs = "roles:create"
	RolesUpdate              PermissionSlugs = "roles:update"
	RolesDelete              PermissionSlugs = "roles:delete"
	ProjectsView             PermissionSlugs = "projects:view"
	ProjectsCreate           PermissionSlugs = "projects:create"
	ProjectsUpdate           PermissionSlugs = "projects:update"
	ProjectsDelete           PermissionSlugs = "projects:delete"
	WorkspacesView           PermissionSlugs = "workspaces:view"
	WorkspacesCreate         PermissionSlugs = "workspaces:create"
	WorkspacesUpdate         PermissionSlugs = "workspaces:update"
	WorkspacesDelete         PermissionSlugs = "workspaces:delete"
	TeamsView                PermissionSlugs = "teams:view"
	TeamsCreate              PermissionSlugs = "teams:create"
	TeamsUpdate              PermissionSlugs = "teams:update"
	TeamsDelete              PermissionSlugs = "teams:delete"
)

type PermissionSlugsArrayItem struct {
	Name        string
	Slug        PermissionSlugs
	Description string
}

var PermissionSlugsArray = []PermissionSlugsArrayItem{
	{Name: "Organizations View", Slug: OrganizationsView, Description: "Allow users to view organizations"},
	{Name: "Organizations Update", Slug: OrganizationsUpdate, Description: "Allow users to update organizations"},
	{Name: "Organizations Delete", Slug: OrganizationsDelete, Description: "Allow users to delete organizations"},
	{Name: "Organizations Users View", Slug: OrganizationsUsersView, Description: "Allow users to view organizations users"},
	{Name: "Organizations Users Create", Slug: OrganizationsUsersCreate, Description: "Allow users to create organizations users"},
	{Name: "Organizations Users Update", Slug: OrganizationsUsersUpdate, Description: "Allow users to update organizations users"},
	{Name: "Organizations Users Delete", Slug: OrganizationsUsersDelete, Description: "Allow users to delete organizations users"},
	{Name: "Roles View", Slug: RolesView, Description: "Allow users to view roles"},
	{Name: "Roles Create", Slug: RolesCreate, Description: "Allow users to create roles"},
	{Name: "Roles Update", Slug: RolesUpdate, Description: "Allow users to update roles"},
	{Name: "Roles Delete", Slug: RolesDelete, Description: "Allow users to delete roles"},
	{Name: "Projects View", Slug: ProjectsView, Description: "Allow users to view projects"},
	{Name: "Projects Create", Slug: ProjectsCreate, Description: "Allow users to create projects"},
	{Name: "Projects Update", Slug: ProjectsUpdate, Description: "Allow users to update projects"},
	{Name: "Projects Delete", Slug: ProjectsDelete, Description: "Allow users to delete projects"},
	{Name: "Workspaces View", Slug: WorkspacesView, Description: "Allow users to view workspaces"},
	{Name: "Workspaces Create", Slug: WorkspacesCreate, Description: "Allow users to create workspaces"},
	{Name: "Workspaces Update", Slug: WorkspacesUpdate, Description: "Allow users to update workspaces"},
	{Name: "Workspaces Delete", Slug: WorkspacesDelete, Description: "Allow users to delete workspaces"},
	{Name: "Teams View", Slug: TeamsView, Description: "Allow users to view teams"},
	{Name: "Teams Create", Slug: TeamsCreate, Description: "Allow users to create teams"},
	{Name: "Teams Update", Slug: TeamsUpdate, Description: "Allow users to update teams"},
	{Name: "Teams Delete", Slug: TeamsDelete, Description: "Allow users to delete teams"},
}

type DefaultRoleSlugs string

const (
	DefaultRoleSlug DefaultRoleSlugs = "default"
	AdminRoleSlug   DefaultRoleSlugs = "admin"
)

type DefaultRoleSlugsArrayItem struct {
	Name        string
	Slug        DefaultRoleSlugs
	Description string
	Permissions []PermissionSlugs
}

var DefaultRoleSlugsArray = []DefaultRoleSlugsArrayItem{
	{
		Name:        "Default",
		Slug:        DefaultRoleSlug,
		Description: "Default role",
		Permissions: []PermissionSlugs{
			OrganizationsView,
			RolesView,
			ProjectsView,
			ProjectsCreate,
			ProjectsUpdate,
			ProjectsDelete,
			WorkspacesView,
		},
	},
	{
		Name:        "Admin",
		Slug:        AdminRoleSlug,
		Description: "Admin role",
		Permissions: []PermissionSlugs{
			OrganizationsView,
			OrganizationsUpdate,
			OrganizationsDelete,
			OrganizationsUsersView,
			OrganizationsUsersCreate,
			OrganizationsUsersUpdate,
			OrganizationsUsersDelete,
			RolesView,
			RolesCreate,
			RolesUpdate,
			RolesDelete,
			ProjectsView,
			ProjectsCreate,
			ProjectsUpdate,
			ProjectsDelete,
			WorkspacesView,
			WorkspacesCreate,
			WorkspacesUpdate,
			WorkspacesDelete,
			TeamsView,
			TeamsCreate,
			TeamsUpdate,
			TeamsDelete,
		},
	},
}
