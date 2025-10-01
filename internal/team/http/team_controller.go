package team_http

import (
	"net/http"

	"github.com/gabrielmrtt/taski/internal/core"
	core_http "github.com/gabrielmrtt/taski/internal/core/http"
	organization_http_middlewares "github.com/gabrielmrtt/taski/internal/organization/http/middlewares"
	team_core "github.com/gabrielmrtt/taski/internal/team"
	team_http_requests "github.com/gabrielmrtt/taski/internal/team/http/requests"
	team_services "github.com/gabrielmrtt/taski/internal/team/services"
	user_http_middlewares "github.com/gabrielmrtt/taski/internal/user/http/middlewares"
	"github.com/gin-gonic/gin"
)

type TeamController struct {
	ListTeamsService  *team_services.ListTeamsService
	GetTeamService    *team_services.GetTeamService
	CreateTeamService *team_services.CreateTeamService
	UpdateTeamService *team_services.UpdateTeamService
	DeleteTeamService *team_services.DeleteTeamService
}

func NewTeamController(
	listTeamsService *team_services.ListTeamsService,
	getTeamService *team_services.GetTeamService,
	createTeamService *team_services.CreateTeamService,
	updateTeamService *team_services.UpdateTeamService,
	deleteTeamService *team_services.DeleteTeamService,
) *TeamController {
	return &TeamController{
		ListTeamsService:  listTeamsService,
		GetTeamService:    getTeamService,
		CreateTeamService: createTeamService,
		UpdateTeamService: updateTeamService,
		DeleteTeamService: deleteTeamService,
	}
}

type ListTeamsResponse = core_http.HttpSuccessResponseWithData[core.PaginationOutput[team_core.TeamDto]]

// ListTeams godoc
// @Summary List teams in an organization
// @Description Lists all existing teams in an organization.
// @Tags Team
// @Accept json
// @Param request query team_http_requests.ListTeamsRequest true "Query parameters"
// @Produce json
// @Success 200 {object} ListTeamsResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 401 {object} core_http.HttpErrorResponse
// @Failure 403 {object} core_http.HttpErrorResponse
// @Failure 404 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /organization/:organizationId/team [get]
func (c *TeamController) ListTeams(ctx *gin.Context) {
	var request team_http_requests.ListTeamsRequest

	organizationIdentity := organization_http_middlewares.GetOrganizationIdentityFromPath(ctx)

	if err := request.FromQuery(ctx); err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	input := request.ToInput()
	input.OrganizationIdentity = organizationIdentity

	response, err := c.ListTeamsService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

type GetTeamResponse = core_http.HttpSuccessResponseWithData[team_core.TeamDto]

// GetTeam godoc
// @Summary Get a team in an organization
// @Description Returns a team in an organization by its ID.
// @Tags Team
// @Accept json
// @Param teamId path string true "Team ID"
// @Produce json
// @Success 200 {object} GetTeamResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 401 {object} core_http.HttpErrorResponse
// @Failure 403 {object} core_http.HttpErrorResponse
// @Failure 404 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /organization/:organizationId/team/:teamId [get]
func (c *TeamController) GetTeam(ctx *gin.Context) {
	teamIdentity := core.NewIdentityFromPublic(ctx.Param("teamId"))
	organizationIdentity := organization_http_middlewares.GetOrganizationIdentityFromPath(ctx)

	input := team_services.GetTeamInput{
		TeamIdentity:         teamIdentity,
		OrganizationIdentity: organizationIdentity,
	}

	response, err := c.GetTeamService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

type CreateTeamResponse = core_http.HttpSuccessResponseWithData[team_core.TeamDto]

// CreateTeam godoc
// @Summary Create a team in an organization
// @Description Creates a new team in an organization.
// @Tags Team
// @Accept json
// @Param request body team_http_requests.CreateTeamRequest true "Request body"
// @Produce json
// @Success 200 {object} CreateTeamResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 401 {object} core_http.HttpErrorResponse
// @Failure 403 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /organization/:organizationId/team [post]
func (c *TeamController) CreateTeam(ctx *gin.Context) {
	var request team_http_requests.CreateTeamRequest

	organizationIdentity := organization_http_middlewares.GetOrganizationIdentityFromPath(ctx)
	authenticatedUserIdentity := user_http_middlewares.GetAuthenticatedUserIdentity(ctx)

	if err := ctx.ShouldBindJSON(&request); err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	input := request.ToInput()
	input.OrganizationIdentity = organizationIdentity
	input.UserCreatorIdentity = authenticatedUserIdentity

	response, err := c.CreateTeamService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

type UpdateTeamResponse = core_http.EmptyHttpSuccessResponse

// UpdateTeam godoc
// @Summary Update a team in an organization
// @Description Updates an existing team in an organization.
// @Tags Team
// @Accept json
// @Param teamId path string true "Team ID"
// @Param request body team_http_requests.UpdateTeamRequest true "Request body"
// @Produce json
// @Success 200 {object} UpdateTeamResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 401 {object} core_http.HttpErrorResponse
// @Failure 403 {object} core_http.HttpErrorResponse
// @Failure 404 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /organization/:organizationId/team/:teamId [put]
func (c *TeamController) UpdateTeam(ctx *gin.Context) {
	var request team_http_requests.UpdateTeamRequest

	organizationIdentity := organization_http_middlewares.GetOrganizationIdentityFromPath(ctx)
	teamIdentity := core.NewIdentityFromPublic(ctx.Param("teamId"))
	authenticatedUserIdentity := user_http_middlewares.GetAuthenticatedUserIdentity(ctx)

	if err := ctx.ShouldBindJSON(&request); err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	input := request.ToInput()
	input.OrganizationIdentity = organizationIdentity
	input.TeamIdentity = teamIdentity
	input.UserEditorIdentity = authenticatedUserIdentity

	err := c.UpdateTeamService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

type DeleteTeamResponse = core_http.EmptyHttpSuccessResponse

// DeleteTeam godoc
// @Summary Delete a team in an organization
// @Description Deletes an existing team in an organization.
// @Tags Team
// @Accept json
// @Param teamId path string true "Team ID"
// @Produce json
// @Success 200 {object} DeleteTeamResponse
// @Failure 400 {object} core_http.HttpErrorResponse
// @Failure 401 {object} core_http.HttpErrorResponse
// @Failure 403 {object} core_http.HttpErrorResponse
// @Failure 404 {object} core_http.HttpErrorResponse
// @Failure 500 {object} core_http.HttpErrorResponse
// @Router /organization/:organizationId/team/:teamId [delete]
func (c *TeamController) DeleteTeam(ctx *gin.Context) {
	organizationIdentity := organization_http_middlewares.GetOrganizationIdentityFromPath(ctx)
	teamIdentity := core.NewIdentityFromPublic(ctx.Param("teamId"))

	input := team_services.DeleteTeamInput{
		TeamIdentity:         teamIdentity,
		OrganizationIdentity: organizationIdentity,
	}

	err := c.DeleteTeamService.Execute(input)
	if err != nil {
		core_http.NewHttpErrorResponse(ctx, err)
		return
	}

	core_http.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

func (c *TeamController) ConfigureRoutes(g *gin.RouterGroup) *gin.RouterGroup {
	group := g.Group("/organization/:organizationId/team")
	{
		group.Use(user_http_middlewares.AuthMiddleware())

		group.GET("", organization_http_middlewares.UserMustHavePermission("teams:view"), c.ListTeams)
		group.GET("/:teamId", organization_http_middlewares.UserMustHavePermission("teams:view"), c.GetTeam)
		group.POST("", organization_http_middlewares.UserMustHavePermission("teams:create"), c.CreateTeam)
		group.PUT("/:teamId", organization_http_middlewares.UserMustHavePermission("teams:update"), c.UpdateTeam)
		group.DELETE("/:teamId", organization_http_middlewares.UserMustHavePermission("teams:delete"), c.DeleteTeam)
	}

	return group
}
