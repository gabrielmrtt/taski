package team_http_requests

import (
	"github.com/gabrielmrtt/taski/internal/core"
	team_services "github.com/gabrielmrtt/taski/internal/team/services"
)

type UpdateTeamRequest struct {
	Name        *string   `json:"name"`
	Description *string   `json:"description"`
	Members     *[]string `json:"members"`
}

func (r *UpdateTeamRequest) ToInput() team_services.UpdateTeamInput {
	var users []core.Identity = make([]core.Identity, 0)
	for _, user := range *r.Members {
		users = append(users, core.NewIdentityFromPublic(user))
	}

	return team_services.UpdateTeamInput{
		Name:        r.Name,
		Description: r.Description,
		Members:     users,
	}
}
