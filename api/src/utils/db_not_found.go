package utils

import "github.com/Southclaws/storyden/api/src/infra/db/model"

func NotFoundOrError[T *any](t T, err error) (T, error) {
	if model.IsNotFound(err) {
		return nil, nil
	}
	return t, err
}
