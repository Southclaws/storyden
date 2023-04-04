package settings

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/internal/ent"
)

type Value[T ~string | ~int | ~float64 | ~uint32 | bool] struct {
	value T // value must remain as field 0 for simple reflection code
	key   string
}

// Get a setting value.
func (s *Value[T]) Get() (v T) { return s.value }

// Set a value persistently using a settings repository.
func (s *Value[T]) Set(ctx context.Context, r Repository, v T) error {
	if err := r.SetValue(ctx, s.key, fmt.Sprint(v)); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	s.value = v
	return nil
}

// fromEnt takes all the rows from a query on the settings table and spits back
// a populated settings struct. It achieves this with some basic reflection.
//
// Each setting is stored in the DB as a string, which keeps things simple for
// editing ad-hoc via SQL and database tools. This function will iterate through
// all possible settings from the `Settings` struct and check if each one is
// present in the list of raw ent settings rows. Those that are are bound to the
// struct using `bindEntry` which performs most of the actual reflection logic.
func fromEnt(raw []*ent.Setting) (*Settings, error) {
	s := Settings{}

	rawmap := lo.FromEntries(dt.Map(raw, func(in *ent.Setting) lo.Entry[string, *ent.Setting] {
		return lo.Entry[string, *ent.Setting]{Key: in.ID, Value: in}
	}))

	rt := reflect.TypeOf(s)
	rv := reflect.ValueOf(&s).Elem()

	for i := 0; i < rt.NumField(); i++ {
		ft := rt.Field(i)
		fv := rv.Field(i)

		if entry, ok := rawmap[ft.Name]; ok {
			nv, err := bindEntry(fv, entry.ID, &entry.Value)
			if err != nil {
				return nil, err
			}

			fv.Set(nv)
		} else {
			nv, err := bindEntry(fv, ft.Name, nil)
			if err != nil {
				return nil, err
			}

			fv.Set(nv)
		}
	}

	return &s, nil
}

// bindEntry takes some field `f` as well as a key and a value in string form
// which represent a single setting row from the database (or elsewhere) and
// switches on the field type in order to use the correct decoding method.
//
// One added detail to note here is that each value of the struct isn't a simple
// scalar type but a generic type `Value[T]` so the return value must be the
// reflected value of this type, not the underlying type.
//
// Also note that the value parameter is a pointer that may be nil. This is just
// to allow use of this for empty/default values without lots of duplicate code.
func bindEntry(f reflect.Value, key string, value *string) (v reflect.Value, err error) {
	k := f.Type().Field(0).Type.Kind()
	switch k {
	case reflect.String:
		var s string
		if value != nil {
			s = *value
		}

		v = reflect.ValueOf(Value[string]{key: key, value: s})

	case reflect.Bool:
		var b bool
		if value != nil {
			b, _ = strconv.ParseBool(*value)
		}

		v = reflect.ValueOf(Value[bool]{key: key, value: b})

	case reflect.Uint32:
		var u64 uint64
		if value != nil {
			u64, err = strconv.ParseUint(*value, 10, 32)
			if err != nil {
				panic(err)
			}
		}

		v = reflect.ValueOf(Value[uint32]{key: key, value: uint32(u64)})

	default:
		err = fault.Newf("cannot auto bind type: '%s'", k.String())
	}

	return
}
