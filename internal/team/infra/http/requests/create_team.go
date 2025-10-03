package teamhttprequests

import (
	"github.com/gabrielmrtt/taski/internal/core"
	teamservice "github.com/gabrielmrtt/taski/internal/team/service"
)

type CreateTeamRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Members     []string `json:"members"`
}

func (r *CreateTeamRequest) ToInput() teamservice.CreateTeamInput {
	var userIdentities []core.Identity = make([]core.Identity, 0)
	for _, user := range r.Members {
		userIdentities = append(userIdentities, core.NewIdentityFromPublic(user))
	}

	return teamservice.CreateTeamInput{
		Name:        r.Name,
		Description: r.Description,
		Members:     userIdentities,
	}
}
