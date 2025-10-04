package corehttp

import (
	"net/http"
	"time"

	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gin-gonic/gin"
)

type HttpSuccessResponseWithData[T any] struct {
	Status    int    `json:"status"`
	Resource  string `json:"resource"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	Data      *T     `json:"data,omitempty"`
}

func NewHttpSuccessResponseWithData[T any](gc *gin.Context, status int, data *T) {
	resource := gc.Request.URL.Path
	message := "success"
	timestamp := time.Now().Format(time.RFC3339)
	res := &HttpSuccessResponseWithData[T]{
		Status:    status,
		Resource:  resource,
		Message:   message,
		Timestamp: timestamp,
		Data:      data,
	}

	gc.JSON(status, res)
}

type EmptyHttpSuccessResponse struct {
	Status    int    `json:"status"`
	Resource  string `json:"resource"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

func NewEmptyHttpSuccessResponse(gc *gin.Context, status int) {
	resource := gc.Request.URL.Path
	message := "success"
	timestamp := time.Now().Format(time.RFC3339)
	res := &EmptyHttpSuccessResponse{
		Status:    status,
		Resource:  resource,
		Message:   message,
		Timestamp: timestamp,
	}

	gc.JSON(status, res)
}

type HttpErrorResponse struct {
	Status    int          `json:"status"`
	Resource  string       `json:"resource"`
	Message   string       `json:"message"`
	Timestamp string       `json:"timestamp"`
	Errors    *interface{} `json:"errors,omitempty"`
}

func NewHttpErrorResponse(gc *gin.Context, err error) {
	resource := gc.Request.URL.Path
	timestamp := time.Now().Format(time.RFC3339)
	status := http.StatusInternalServerError

	message := "Request unsuccessful"
	var errors interface{}
	errors = nil

	switch e := err.(type) {
	case *core.AlreadyExistsError:
		status = http.StatusConflict
		message = e.Error()
		errors = nil
	case *core.InvalidInputError:
		status = http.StatusBadRequest
		message = e.Error()
		errors = err.(*core.InvalidInputError).Fields
	case *core.UnauthorizedError:
		status = http.StatusForbidden
		message = e.Error()
		errors = nil
	case *core.UnauthenticatedError:
		status = http.StatusUnauthorized
		message = e.Error()
		errors = nil
	case *core.InternalError:
		status = http.StatusInternalServerError
		message = e.Error()
		errors = nil
	case *core.NotFoundError:
		status = http.StatusNotFound
		message = e.Error()
		errors = nil
	default:
		status = http.StatusInternalServerError
		message = e.Error()
		errors = nil
	}

	res := &HttpErrorResponse{
		Status:    status,
		Resource:  resource,
		Message:   message,
		Timestamp: timestamp,
		Errors:    &errors,
	}

	gc.JSON(status, res)
}
