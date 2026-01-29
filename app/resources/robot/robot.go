package robot

import (
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/robot/robot_ref"
	"github.com/Southclaws/storyden/internal/ent"
)

type Robot struct {
	robot_ref.Robot

	Author account.Account
}

func Map(in *ent.Robot) (*Robot, error) {
	author, err := account.MapRef(in.Edges.Author)
	if err != nil {
		return nil, err
	}

	return &Robot{
		Robot: robot_ref.Robot{
			ID:          robot_ref.ID(in.ID),
			CreatedAt:   in.CreatedAt,
			UpdatedAt:   in.UpdatedAt,
			Name:        in.Name,
			Description: in.Description,
			Playbook:    in.Playbook,
			Metadata:    in.Metadata,
			AuthorID:    account.AccountID(in.AuthorID),
			Tools:       in.Tools,
		},
		Author: *author,
	}, nil
}
