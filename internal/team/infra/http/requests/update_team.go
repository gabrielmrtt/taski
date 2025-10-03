package teamhttprequests

import (
	"github.com/gabrielmrtt/taski/internal/core"
	teamservice "github.com/gabrielmrtt/taski/internal/team/service"
)

type UpdateTeamRequest struct {
	Name        *string   `json:"name"`
	Description *string   `json:"description"`
	Members     *[]string `json:"members"`
}

func (r *UpdateTeamRequest) ToInput() teamservice.UpdateTeamInput {
	var users []core.Identity = make([]core.Identity, 0)
	for _, user := range *r.Members {
		users = append(users, core.NewIdentityFromPublic(user))
	}

	return teamservice.UpdateTeamInput{
		Name:        r.Name,
		Description: r.Description,
		Members:     users,
	}
}
