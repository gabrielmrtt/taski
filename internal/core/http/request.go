package corehttp

import (
	"strings"

	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/pkg/stringutils"
)

type HttpRequest[INPUT core.ServiceInput] interface {
	ToInput() INPUT
}

func GetRelationsInput(relations *string) core.RelationsInput {
	var relationsInput core.RelationsInput = make([]string, 0)
	if relations != nil {
		relationsInput = strings.Split(strings.TrimSpace(*relations), ",")

		for i, relation := range relationsInput {
			var relationPath []string = make([]string, 0)
			parts := strings.Split(relation, ".")
			for _, part := range parts {
				pascalCasePart := stringutils.CamelCaseToPascalCase(part)
				relationPath = append(relationPath, pascalCasePart)
			}
			relationsInput[i] = strings.Join(relationPath, ".")
		}
	}

	return relationsInput
}
