package utils

import (
	"4d63.com/optional"

	"github.com/Southclaws/storyden/api/src/infra/db/model"
)

func SetIfPresent[T any](m model.Mutation, field string, value optional.Optional[T]) {
	if v, ok := value.Get(); ok {
		m.SetField(field, v)
	}
}
