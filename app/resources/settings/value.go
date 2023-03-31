package settings

import (
	"context"
	"reflect"

	"github.com/Southclaws/storyden/internal/ent"
	"github.com/kr/pretty"
)

type Value[T any] struct {
	value T
}

// func (s Value[T]) Get(ctx context.Context, r Repository) (v T, err error) {
// 	raw, _ := r.GetValue(ctx, "reflect the key name from the struct")
// 	// cast the raw value to T somehow
// 	return nil, nil
// }

func (s Value[T]) Set(ctx context.Context, r Repository, v T) error {
	r.SetValue(ctx, "reflect the key name somehow", "cast v to a string")
	return nil
}

func fromEnt(raw []*ent.Setting) (*Settings, error) {
	s := Settings{}

	keys := map[string]reflect.Type{}

	rt := reflect.TypeOf(s)
	rv := reflect.ValueOf(s)

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)

		name := f.Name
		keys[name] = f.Type
	}

	for _, entry := range raw {
		f := rv.FieldByName(entry.ID)

		f.SetString(entry.Value)
	}

	pretty.Println(keys)

	return &s, nil
}
