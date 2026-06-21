package robot_session

import (
	"github.com/Southclaws/storyden/app/resources/robot/robot_querier"
	"github.com/Southclaws/storyden/internal/ent"
)

type Repository struct {
	db           *ent.Client
	robotQuerier *robot_querier.Querier
}

func New(db *ent.Client, robotQuerier *robot_querier.Querier) *Repository {
	return &Repository{
		db:           db,
		robotQuerier: robotQuerier,
	}
}
