package settings

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/internal/ent"
)

type Value[T ~string | ~int | ~float64 | ~uint32 | bool] struct {
	value T // value must remain as field 0 for simple reflection code
	key   string
}

// Get a setting value.
func (s Value[T]) Get() (v T) { return s.value }

// Set a value persistently using a settings repository.
func (s Value[T]) Set(ctx context.Context, r Repository, v T) error {
	if err := r.SetValue(ctx, s.key, fmt.Sprint(v)); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	return nil
}

func fromEnt(raw []*ent.Setting) (*Settings, error) {
	s := Settings{}

	keys := map[string]reflect.Type{}

	rt := reflect.TypeOf(s)
	rv := reflect.ValueOf(&s).Elem()

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)

		name := f.Name
		keys[name] = f.Type
	}

	for _, entry := range raw {
		f := rv.FieldByName(entry.ID)

		fv, err := bindEntry(entry, f)
		if err != nil {
			return nil, err
		}

		f.Set(fv)
	}

	return &s, nil
}

func bindEntry(entry *ent.Setting, f reflect.Value) (v reflect.Value, err error) {
	k := f.Type().Field(0).Type.Kind()
	switch k {
	case reflect.String:
		v = reflect.ValueOf(Value[string]{key: entry.ID, value: entry.Value})

	case reflect.Uint32:
		u64, err := strconv.ParseUint(entry.Value, 10, 32)
		if err != nil {
			panic(err)
		}

		v = reflect.ValueOf(Value[uint32]{key: entry.ID, value: uint32(u64)})

	default:
		fault.Newf("cannot auto bind type: '%s'", k.String())
	}

	return
}

func get[T any](ctx context.Context, r Repository) (v T, err error) {
	raw, err := r.GetValue(ctx, "reflect the key name from the struct")
	if err != nil {
		return
	}

	// output variable, upcasted to `any` in order to perform a type switch.
	out := any(v)

	// Within each block, we know the concrete type of the underlying value. So,
	// we can simply assign the upcasted output variable to the decoded raw data
	// and then downcast (via type assertion) this back into an actual `T` type.
	switch out.(type) {
	case string:
		out = raw

	case int, int32, int64:
		out, err = strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return
		}

	default:
		panic("idk 🤷")
	}

	return out.(T), nil
}
