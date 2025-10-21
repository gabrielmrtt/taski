package task

import (
	"strings"

	"github.com/gabrielmrtt/taski/internal/core"
	"golang.org/x/net/html"
)

type TaskCommentContent struct {
	Value string
}

func NewTaskCommentContent(value string) (TaskCommentContent, error) {
	c := TaskCommentContent{Value: value}
	if err := c.Validate(); err != nil {
		return TaskCommentContent{}, err
	}
	return c, nil
}

func (c TaskCommentContent) Validate() error {
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
