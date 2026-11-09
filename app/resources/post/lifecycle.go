package post

import ent_post "github.com/Southclaws/storyden/internal/ent/post"

//go:generate go run github.com/Southclaws/enumerator

type threadLifecycleStateEnum string

const (
	threadLifecycleStateOpen     threadLifecycleStateEnum = "open"
	threadLifecycleStateResolved threadLifecycleStateEnum = "resolved"
)

func ThreadLifecycleStateFromEnt(state ent_post.LifecycleState) ThreadLifecycleState {
	return ThreadLifecycleState{threadLifecycleStateEnum(state)}
}
