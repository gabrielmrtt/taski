package projecthttp

import (
	"net/http"

	authhttpmiddlewares "github.com/gabrielmrtt/taski/internal/auth/infra/http/middlewares"
	"github.com/gabrielmrtt/taski/internal/core"
	corehttp "github.com/gabrielmrtt/taski/internal/core/http"
	organizationhttpmiddlewares "github.com/gabrielmrtt/taski/internal/organization/infra/http/middlewares"
	"github.com/gabrielmrtt/taski/internal/project"
	projecthttpmiddlewares "github.com/gabrielmrtt/taski/internal/project/infra/http/middlewares"
	projecthttprequests "github.com/gabrielmrtt/taski/internal/project/infra/http/requests"
	projectservice "github.com/gabrielmrtt/taski/internal/project/service"
	"github.com/gin-gonic/gin"
)

type ProjectDocumentHandler struct {
	GetProjectDocumentVersionService    *projectservice.GetProjectDocumentVersionService
	ListProjectDocumentsService         *projectservice.ListProjectDocumentsService
	ListProjectDocumentVersionsService  *projectservice.ListProjectDocumentVersionsService
	CreateProjectDocumentService        *projectservice.CreateProjectDocumentService
	UpdateProjectDocumentService        *projectservice.UpdateProjectDocumentService
	DeleteProjectDocumentService        *projectservice.DeleteProjectDocumentService
	DeleteProjectDocumentVersionService *projectservice.DeleteProjectDocumentVersionService
}

func NewProjectDocumentHandler(
	getProjectDocumentVersionService *projectservice.GetProjectDocumentVersionService,
	listProjectDocumentsService *projectservice.ListProjectDocumentsService,
	listProjectDocumentVersionsService *projectservice.ListProjectDocumentVersionsService,
	createProjectDocumentService *projectservice.CreateProjectDocumentService,
	updateProjectDocumentService *projectservice.UpdateProjectDocumentService,
	deleteProjectDocumentService *projectservice.DeleteProjectDocumentService,
	deleteProjectDocumentVersionService *projectservice.DeleteProjectDocumentVersionService,
) *ProjectDocumentHandler {
	return &ProjectDocumentHandler{
		GetProjectDocumentVersionService:    getProjectDocumentVersionService,
		ListProjectDocumentsService:         listProjectDocumentsService,
		ListProjectDocumentVersionsService:  listProjectDocumentVersionsService,
		CreateProjectDocumentService:        createProjectDocumentService,
		UpdateProjectDocumentService:        updateProjectDocumentService,
		DeleteProjectDocumentService:        deleteProjectDocumentService,
		DeleteProjectDocumentVersionService: deleteProjectDocumentVersionService,
	}
}

type ListProjectDocumentsResponse = corehttp.HttpSuccessResponseWithData[core.PaginationOutput[project.ProjectDocumentVersionDto]]

// ListProjectDocuments godoc
// @Summary List project documents
// @Description Returns all project documents.
// @Tags Project Document
// @Accept json
// @Param request query projecthttprequests.ListProjectDocumentsRequest true "Query parameters"
// @Produce json
// @Success 200 {object} ListProjectDocumentsResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /project/:projectId/document [get]
func (c *ProjectDocumentHandler) ListProjectDocuments(ctx *gin.Context) {
	var request projecthttprequests.ListProjectDocumentsRequest
	var projectIdentity core.Identity = core.NewIdentityFromPublic(ctx.Param("projectId"))
	var input projectservice.ListProjectDocumentsInput

	if err := request.FromQuery(ctx); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	input = request.ToInput()
	input.Filters.ProjectIdentity = &projectIdentity

	response, err := c.ListProjectDocumentsService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

type ListProjectDocumentVersionsResponse = corehttp.HttpSuccessResponseWithData[core.PaginationOutput[project.ProjectDocumentVersionDto]]

// ListProjectDocumentVersions godoc
// @Summary List project document versions
// @Description Returns all project document versions.
// @Tags Project Document
// @Accept json
// @Param projectId path string true "Project ID"
// @Param documentVersionManagerId path string true "Document Version Manager ID"
// @Produce json
// @Success 200 {object} ListProjectDocumentVersionsResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /project/:projectId/document/:documentVersionManagerId [get]
func (c *ProjectDocumentHandler) ListProjectDocumentVersions(ctx *gin.Context) {
	var request projecthttprequests.ListProjectDocumentVersionsRequest
	var projectIdentity core.Identity = core.NewIdentityFromPublic(ctx.Param("projectId"))
	var documentVersionManagerIdentity core.Identity = core.NewIdentityFromPublic(ctx.Param("documentVersionManagerId"))
	var input projectservice.ListProjectDocumentVersionsInput
	if err := request.FromQuery(ctx); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	input = request.ToInput()
	input.Filters.ProjectIdentity = &projectIdentity
	input.Filters.ProjectDocumentVersionManagerIdentity = &documentVersionManagerIdentity

	response, err := c.ListProjectDocumentVersionsService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

type GetProjectDocumentVersionResponse = corehttp.HttpSuccessResponseWithData[project.ProjectDocumentVersionDto]

// GetProjectDocumentVersion godoc
// @Summary Get a project document version
// @Description Returns a project document version by its ID.
// @Tags Project Document
// @Accept json
// @Param projectId path string true "Project ID"
// @Param documentVersionManagerId path string true "Document Version Manager ID"
// @Param documentVersionId path string true "Document Version ID"
// @Produce json
// @Success 200 {object} GetProjectDocumentVersionResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /project/:projectId/document/:documentVersionManagerId/version/:documentVersionId [get]
func (c *ProjectDocumentHandler) GetProjectDocumentVersion(ctx *gin.Context) {
	var projectIdentity core.Identity = core.NewIdentityFromPublic(ctx.Param("projectId"))
	var documentVersionManagerIdentity core.Identity = core.NewIdentityFromPublic(ctx.Param("documentVersionManagerId"))
	var documentVersionIdentity core.Identity = core.NewIdentityFromPublic(ctx.Param("documentVersionId"))
	var request projecthttprequests.GetProjectDocumentVersionRequest

	if err := request.FromQuery(ctx); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	input := request.ToInput()
	input.ProjectIdentity = projectIdentity
	input.ProjectDocumentVersionManagerIdentity = documentVersionManagerIdentity
	input.ProjectDocumentVersionIdentity = documentVersionIdentity

	response, err := c.GetProjectDocumentVersionService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

type CreateProjectDocumentResponse = corehttp.HttpSuccessResponseWithData[project.ProjectDocumentVersionDto]

// CreateProjectDocument godoc
// @Summary Create a project document
// @Description Creates a new project document.
// @Tags Project Document
// @Accept json
// @Param projectId path string true "Project ID"
// @Param request body projecthttprequests.CreateProjectDocumentRequest true "Request body"
// @Produce json
// @Success 200 {object} CreateProjectDocumentResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /project/:projectId/document [post]
func (c *ProjectDocumentHandler) CreateProjectDocument(ctx *gin.Context) {
	var request projecthttprequests.CreateProjectDocumentRequest
	var authenticatedUserIdentity *core.Identity = authhttpmiddlewares.GetAuthenticatedUserIdentity(ctx)
	var projectIdentity core.Identity = core.NewIdentityFromPublic(ctx.Param("projectId"))
	var input projectservice.CreateProjectDocumentInput

	if err := ctx.ShouldBind(&request); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	input = request.ToInput()
	input.ProjectIdentity = projectIdentity
	input.UserCreatorIdentity = *authenticatedUserIdentity

	response, err := c.CreateProjectDocumentService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

type UpdateProjectDocumentResponse = corehttp.EmptyHttpSuccessResponse

// UpdateProjectDocument godoc
// @Summary Update a project document
// @Description Updates an existing project document.
// @Tags Project Document
// @Accept json
// @Param projectId path string true "Project ID"
// @Param documentVersionManagerId path string true "Document Version Manager ID"
// @Param documentVersionId path string true "Document Version ID"
// @Param request body projecthttprequests.UpdateProjectDocumentRequest true "Request body"
// @Produce json
// @Success 200 {object} UpdateProjectDocumentResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 404 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /project/:projectId/document/:documentVersionManagerId/version/:documentVersionId [put]
func (c *ProjectDocumentHandler) UpdateProjectDocument(ctx *gin.Context) {
	var request projecthttprequests.UpdateProjectDocumentRequest
	var projectIdentity core.Identity = core.NewIdentityFromPublic(ctx.Param("projectId"))
	var documentVersionManagerIdentity core.Identity = core.NewIdentityFromPublic(ctx.Param("documentVersionManagerId"))
	var documentVersionIdentity core.Identity = core.NewIdentityFromPublic(ctx.Param("documentVersionId"))
	var input projectservice.UpdateProjectDocumentInput

	if err := ctx.ShouldBindJSON(&request); err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	input = request.ToInput()
	input.ProjectIdentity = projectIdentity
	input.ProjectDocumentVersionManagerIdentity = documentVersionManagerIdentity
	input.ProjectDocumentVersionIdentity = documentVersionIdentity

	response, err := c.UpdateProjectDocumentService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewHttpSuccessResponseWithData(ctx, http.StatusOK, response)
}

type DeleteProjectDocumentResponse = corehttp.EmptyHttpSuccessResponse

// DeleteProjectDocument godoc
// @Summary Delete a project document
// @Description Deletes an existing project document.
// @Tags Project Document
// @Accept json
// @Param projectId path string true "Project ID"
// @Param documentVersionManagerId path string true "Document Version Manager ID"
// @Produce json
// @Success 200 {object} DeleteProjectDocumentResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /project/:projectId/document/:documentVersionManagerId [delete]
func (c *ProjectDocumentHandler) DeleteProjectDocument(ctx *gin.Context) {
	var projectIdentity core.Identity = core.NewIdentityFromPublic(ctx.Param("projectId"))
	var documentVersionManagerIdentity core.Identity = core.NewIdentityFromPublic(ctx.Param("documentVersionManagerId"))
	var input projectservice.DeleteProjectDocumentInput = projectservice.DeleteProjectDocumentInput{
		ProjectIdentity:                       projectIdentity,
		ProjectDocumentVersionManagerIdentity: documentVersionManagerIdentity,
	}

	err := c.DeleteProjectDocumentService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

type DeleteProjectDocumentVersionResponse = corehttp.EmptyHttpSuccessResponse

// DeleteProjectDocumentVersion godoc
// @Summary Delete a project document version
// @Description Deletes an existing project document version.
// @Tags Project Document
// @Accept json
// @Param projectId path string true "Project ID"
// @Param documentVersionManagerId path string true "Document Version Manager ID"
// @Param documentVersionId path string true "Document Version ID"
// @Produce json
// @Success 200 {object} DeleteProjectDocumentVersionResponse
// @Failure 400 {object} corehttp.HttpErrorResponse
// @Failure 401 {object} corehttp.HttpErrorResponse
// @Failure 403 {object} corehttp.HttpErrorResponse
// @Failure 500 {object} corehttp.HttpErrorResponse
// @Router /project/:projectId/document/:documentVersionManagerId/version/:documentVersionId [delete]
func (c *ProjectDocumentHandler) DeleteProjectDocumentVersion(ctx *gin.Context) {
	var projectIdentity core.Identity = core.NewIdentityFromPublic(ctx.Param("projectId"))
	var documentVersionManagerIdentity core.Identity = core.NewIdentityFromPublic(ctx.Param("documentVersionManagerId"))
	var documentVersionIdentity core.Identity = core.NewIdentityFromPublic(ctx.Param("documentVersionId"))
	var input projectservice.DeleteProjectDocumentVersionInput = projectservice.DeleteProjectDocumentVersionInput{
		ProjectIdentity:                       projectIdentity,
		ProjectDocumentVersionManagerIdentity: documentVersionManagerIdentity,
		ProjectDocumentVersionIdentity:        documentVersionIdentity,
	}

	err := c.DeleteProjectDocumentVersionService.Execute(input)
	if err != nil {
		corehttp.NewHttpErrorResponse(ctx, err)
		return
	}

	corehttp.NewEmptyHttpSuccessResponse(ctx, http.StatusOK)
}

func (c *ProjectDocumentHandler) ConfigureRoutes(options corehttp.ConfigureRoutesOptions) *gin.RouterGroup {
	middlewareOptions := corehttp.MiddlewareOptions{
		DbConnection: options.DbConnection,
	}

	g := options.RouterGroup.Group("/project/:projectId/document")
	{
		g.Use(authhttpmiddlewares.AuthMiddleware(middlewareOptions))
		g.Use(projecthttpmiddlewares.UserMustBeInProject(middlewareOptions))

		g.GET("", organizationhttpmiddlewares.UserMustHavePermission("projects:view", middlewareOptions), c.ListProjectDocuments)
		g.GET("/:documentVersionManagerId", organizationhttpmiddlewares.UserMustHavePermission("projects:view", middlewareOptions), c.ListProjectDocumentVersions)
		g.GET("/:documentVersionManagerId/version/:documentVersionId", organizationhttpmiddlewares.UserMustHavePermission("projects:view", middlewareOptions), c.GetProjectDocumentVersion)
		g.POST("", organizationhttpmiddlewares.UserMustHavePermission("projects:update", middlewareOptions), c.CreateProjectDocument)
		g.PUT("/:documentVersionManagerId/version/:documentVersionId", organizationhttpmiddlewares.UserMustHavePermission("projects:update", middlewareOptions), c.UpdateProjectDocument)
		g.DELETE("/:documentVersionManagerId", organizationhttpmiddlewares.UserMustHavePermission("projects:update", middlewareOptions), c.DeleteProjectDocument)
		g.DELETE("/:documentVersionManagerId/version/:documentVersionId", organizationhttpmiddlewares.UserMustHavePermission("projects:update", middlewareOptions), c.DeleteProjectDocumentVersion)
	}

	return g
}
