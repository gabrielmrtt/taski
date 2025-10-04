package authhttp

import (
	"fmt"
	"net/http"

	"github.com/gabrielmrtt/taski/internal/auth"
	authhttpmiddlewares "github.com/gabrielmrtt/taski/internal/auth/infra/http/middlewares"
	authhttprequests "github.com/gabrielmrtt/taski/internal/auth/infra/http/requests"
	authservice "github.com/gabrielmrtt/taski/internal/auth/service"
	"github.com/gabrielmrtt/taski/internal/core"
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	UserLoginService          *authservice.UserLoginService
	AccessOrganizationService *authservice.AccessOrganizationService
}

func NewAuthHandler(
	userLoginService *authservice.UserLoginService,
	accessOrganizationService *authservice.AccessOrganizationService,
) *AuthHandler {
	return &AuthHandler{
		UserLoginService:          userLoginService,
		AccessOrganizationService: accessOrganizationService,
	}
}

type LoginResponse = corehttp.HttpSuccessResponseWithData[auth.UserAuthDto]

// Login godoc
// @Summary Login
// @Description Authenticates an user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body authhttprequests.UserLoginRequest true "Request body"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /auth/login [post]
func (c *AuthHandler) Login(ctx *gin.Context) {
	var request authhttprequests.UserLoginRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	response, err := c.UserLoginService.Execute(authservice.UserLoginInput{
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

type AccessOrganizationResponse = corehttp.EmptyHttpSuccessResponse

// AccessOrganization godoc
// @Summary Access an organization
// @Description Access an organization
// @Tags Auth
// @Accept json
// @Produce json
// @Param organizationId path string true "Organization ID"
// @Success 200 {object} corehttp.EmptyHttpSuccessResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /auth/organization/:organizationId/access [patch]
func (c *AuthHandler) AccessOrganization(ctx *gin.Context) {
	authenticatedUserIdentity := authhttpmiddlewares.GetAuthenticatedUserIdentity(ctx)
	organizationIdentity := ctx.Param("organizationId")

	input := authservice.AccessOrganizationInput{
		LoggedUserIdentity:   *authenticatedUserIdentity,
		OrganizationIdentity: core.NewIdentityFromPublic(organizationIdentity),
	}

	token, err := c.AccessOrganizationService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	ctx.Header("Authorization", fmt.Sprintf("Bearer %s", *token))
	corehttp.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

func (c *AuthHandler) ConfigureRoutes(options corehttp.ConfigureRoutesOptions) *gin.RouterGroup {
	middlewareOptions := corehttp.MiddlewareOptions{
		DbConnection: options.DbConnection,
	}

	g := options.RouterGroup.Group("/auth")
	{
		g.POST("/login", c.Login)
		g.PATCH("/organization/:organizationId/access", authhttpmiddlewares.AuthMiddleware(middlewareOptions), c.AccessOrganization)
	}

	return g
}
