package corehttp

import (
	"strings"

	"github.com/gabrielmrtt/taski/internal/core"
	"github.com/gabrielmrtt/taski/pkg/stringutils"
)

type HttpRequest[INPUT core.ServiceInput] interface {
	ToInput() INPUT
}

func GetRelationsInput(relations string) core.RelationsInput {
	var relationsInput core.RelationsInput = make([]string, 0)
	if relations != "" {
		relationsInput = strings.Split(strings.TrimSpace(relations), ",")

		for i, relation := range relationsInput {
			relationsInput[i] = stringutils.CamelCaseToPascalCase(relation)
		}
	}

	return relationsInput
}
