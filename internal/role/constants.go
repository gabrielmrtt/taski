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
		},
	},
}
