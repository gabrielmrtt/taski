package teamhttp

import (
	"net/http"

	authhttpmiddlewares "github.com/gabrielmrtt/taski/internal/auth/infra/http/middlewares"
	"github.com/gabrielmrtt/taski/internal/core"
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	organizationhttpmiddlewares "github.com/gabrielmrtt/taski/internal/organization/infra/http/middlewares"
	"github.com/gabrielmrtt/taski/internal/team"
	teamhttprequests "github.com/gabrielmrtt/taski/internal/team/infra/http/requests"
	teamservice "github.com/gabrielmrtt/taski/internal/team/service"
	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	ListTeamsService  *teamservice.ListTeamsService
	GetTeamService    *teamservice.GetTeamService
	CreateTeamService *teamservice.CreateTeamService
	UpdateTeamService *teamservice.UpdateTeamService
	DeleteTeamService *teamservice.DeleteTeamService
}

func NewTeamHandler(
	listTeamsService *teamservice.ListTeamsService,
	getTeamService *teamservice.GetTeamService,
	createTeamService *teamservice.CreateTeamService,
	updateTeamService *teamservice.UpdateTeamService,
	deleteTeamService *teamservice.DeleteTeamService,
) *TeamHandler {
	return &TeamHandler{
		ListTeamsService:  listTeamsService,
		GetTeamService:    getTeamService,
		CreateTeamService: createTeamService,
		UpdateTeamService: updateTeamService,
		DeleteTeamService: deleteTeamService,
	}
}

type ListTeamsResponse = corehttp.HttpSuccessResponseWithData[core.PaginationOutput[team.TeamDto]]

// ListTeams godoc
// @Summary List teams in an organization
// @Description Lists all existing teams in an organization.
// @Tags Team
// @Accept json
// @Param request query teamhttprequests.ListTeamsRequest true "Query parameters"
// @Produce json
// @Success 200 {object} ListTeamsResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /team [get]
func (c *TeamHandler) ListTeams(ctx *gin.Context) {
	var request teamhttprequests.ListTeamsRequest
	var organizationIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserLastAccessedOrganizationIdentity(ctx)
	var input teamservice.ListTeamsInput

	if err := request.FromQuery(ctx); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	input = request.ToInput()
	input.Filters.OrganizationIdentity = *organizationIdentity

	response, err := c.ListTeamsService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

type GetTeamResponse = corehttp.HttpSuccessResponseWithData[team.TeamDto]

// GetTeam godoc
// @Summary Get a team in an organization
// @Description Returns a team in an organization by its ID.
// @Tags Team
// @Accept json
// @Param teamId path string true "Team ID"
// @Produce json
// @Success 200 {object} GetTeamResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /team/:teamId [get]
func (c *TeamHandler) GetTeam(ctx *gin.Context) {
	var teamIdentity core.Identity = core.NewIdentityFromPublic(ctx.Param("teamId"))
	var organizationIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserLastAccessedOrganizationIdentity(ctx)
	var input teamservice.GetTeamInput = teamservice.GetTeamInput{
		TeamIdentity:         teamIdentity,
		OrganizationIdentity: *organizationIdentity,
	}

	response, err := c.GetTeamService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

type CreateTeamResponse = corehttp.HttpSuccessResponseWithData[team.TeamDto]

// CreateTeam godoc
// @Summary Create a team in an organization
// @Description Creates a new team in an organization.
// @Tags Team
// @Accept json
// @Param request body teamhttprequests.CreateTeamRequest true "Request body"
// @Produce json
// @Success 200 {object} CreateTeamResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /organization/:organizationId/team [post]
func (c *TeamHandler) CreateTeam(ctx *gin.Context) {
	var request teamhttprequests.CreateTeamRequest
	var organizationIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserLastAccessedOrganizationIdentity(ctx)
	var authenticatedUserIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserIdentity(ctx)
	var input teamservice.CreateTeamInput

	if err := ctx.ShouldBindJSON(&request); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	input = request.ToInput()
	input.OrganizationIdentity = *organizationIdentity
	input.UserCreatorIdentity = *authenticatedUserIdentity

	response, err := c.CreateTeamService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

type UpdateTeamResponse = corehttp.EmptyHttpSuccessResponse

// UpdateTeam godoc
// @Summary Update a team in an organization
// @Description Updates an existing team in an organization.
// @Tags Team
// @Accept json
// @Param teamId path string true "Team ID"
// @Param request body teamhttprequests.UpdateTeamRequest true "Request body"
// @Produce json
// @Success 200 {object} UpdateTeamResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /team/:teamId [put]
func (c *TeamHandler) UpdateTeam(ctx *gin.Context) {
	var request teamhttprequests.UpdateTeamRequest
	var organizationIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserLastAccessedOrganizationIdentity(ctx)
	var teamIdentity core.Identity = core.NewIdentityFromPublic(ctx.Param("teamId"))
	var authenticatedUserIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserIdentity(ctx)
	var input teamservice.UpdateTeamInput

	if err := ctx.ShouldBindJSON(&request); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	input = request.ToInput()
	input.OrganizationIdentity = *organizationIdentity
	input.TeamIdentity = teamIdentity
	input.UserEditorIdentity = *authenticatedUserIdentity

	err := c.UpdateTeamService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

type DeleteTeamResponse = corehttp.EmptyHttpSuccessResponse

// DeleteTeam godoc
// @Summary Delete a team in an organization
// @Description Deletes an existing team in an organization.
// @Tags Team
// @Accept json
// @Param teamId path string true "Team ID"
// @Produce json
// @Success 200 {object} DeleteTeamResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /team/:teamId [delete]
func (c *TeamHandler) DeleteTeam(ctx *gin.Context) {
	var organizationIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserLastAccessedOrganizationIdentity(ctx)
	var teamIdentity core.Identity = core.NewIdentityFromPublic(ctx.Param("teamId"))
	var input teamservice.DeleteTeamInput = teamservice.DeleteTeamInput{
		TeamIdentity:         teamIdentity,
		OrganizationIdentity: *organizationIdentity,
	}

	err := c.DeleteTeamService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

func (c *TeamHandler) ConfigureRoutes(options corehttp.ConfigureRoutesOptions) *gin.RouterGroup {
	middlewareOptions := corehttp.MiddlewareOptions{
		DbConnection: options.DbConnection,
	}

	g := options.RouterGroup.Group("/team")
	{
		g.Use(authhttpmiddlewares.AuthMiddleware(middlewareOptions))

		g.GET("", organizationhttpmiddlewares.UserMustHavePermission("teams:view", middlewareOptions), c.ListTeams)
		g.GET("/:teamId", organizationhttpmiddlewares.UserMustHavePermission("teams:view", middlewareOptions), c.GetTeam)
		g.POST("", organizationhttpmiddlewares.UserMustHavePermission("teams:create", middlewareOptions), c.CreateTeam)
		g.PUT("/:teamId", organizationhttpmiddlewares.UserMustHavePermission("teams:update", middlewareOptions), c.UpdateTeam)
		g.DELETE("/:teamId", organizationhttpmiddlewares.UserMustHavePermission("teams:delete", middlewareOptions), c.DeleteTeam)
	}

	return g
}
