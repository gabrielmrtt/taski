package project

import (
	"strings"

	"github.com/gabrielmrtt/taski/internal/core"
	"golang.org/x/net/html"
)

type ProjectDocumentTitle struct {
	Value string
}

func NewProjectDocumentTitle(value string) (ProjectDocumentTitle, error) {
	t := ProjectDocumentTitle{Value: value}
	if err := t.Validate(); err != nil {
		return ProjectDocumentTitle{}, err
	}
	return t, nil
}

func (t ProjectDocumentTitle) Validate() error {
	if t.Value == "" {
		field := core.InvalidInputErrorField{
			Field: "title",
			Error: "title cannot be empty",
		}
		return core.NewInvalidInputError("title cannot be empty", []core.InvalidInputErrorField{field})
	}

	if len(t.Value) > 255 {
		field := core.InvalidInputErrorField{
			Field: "title",
			Error: "title must be less than 255 characters",
		}
		return core.NewInvalidInputError("title must be less than 255 characters", []core.InvalidInputErrorField{field})
	}

	return nil
}

type ProjectDocumentContent struct {
	Value string
}

func NewProjectDocumentContent(value string) (ProjectDocumentContent, error) {
	c := ProjectDocumentContent{Value: value}
	if err := c.Validate(); err != nil {
		return ProjectDocumentContent{}, err
	}
	return c, nil
}

func (c ProjectDocumentContent) Validate() error {
	if c.Value == "" {
		field := core.InvalidInputErrorField{
			Field: "content",
			Error: "content cannot be empty",
		}
		return core.NewInvalidInputError("content cannot be empty", []core.InvalidInputErrorField{field})
	}

	doc, err := html.Parse(strings.NewReader(c.Value))
	if err != nil {
		field := core.InvalidInputErrorField{
			Field: "content",
			Error: "content must be valid HTML",
		}
		return core.NewInvalidInputError("content must be valid HTML", []core.InvalidInputErrorField{field})
	}

	if doc == nil || (doc.FirstChild == nil && doc.LastChild == nil) {
		field := core.InvalidInputErrorField{
			Field: "content",
			Error: "content must contain valid HTML elements",
		}
		return core.NewInvalidInputError("content must contain valid HTML elements", []core.InvalidInputErrorField{field})
	}

	return nil
}
