package team_http_requests

import (
	"github.com/gabrielmrtt/taski/internal/core"
	team_services "github.com/gabrielmrtt/taski/internal/team/services"
)

type CreateTeamRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Members     []string `json:"members"`
}

func (r *CreateTeamRequest) ToInput() team_services.CreateTeamInput {
	var userIdentities []core.Identity = make([]core.Identity, 0)
	for _, user := range r.Members {
		userIdentities = append(userIdentities, core.NewIdentityFromPublic(user))
	}

	return team_services.CreateTeamInput{
		Name:        r.Name,
		Description: r.Description,
		Members:     userIdentities,
	}
}
