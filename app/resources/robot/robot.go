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

	ref, err := robot_ref.Map(in)
	if err != nil {
		return nil, err
	}

	return &Robot{
		Robot:  *ref,
		Author: *author,
	}, nil
}
