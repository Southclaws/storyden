package deletable

import (
	"github.com/Southclaws/opt"
	"github.com/oapi-codegen/nullable"
)

// Value specifies a type which communicates whether a field for a resource
// should be either:
// - deleted (no value, delete = true)
// - updated to the new value (the held value)
// - left untouched (no value)
type Value[T any] struct {
	v      opt.Optional[T]
	delete bool
}

func New[T any](value nullable.Nullable[T]) Value[T] {
	if value.IsNull() {
		return Value[T]{v: opt.NewEmpty[T](), delete: true}
	}

	v, err := value.Get()
	if err != nil {
		return Value[T]{v: opt.NewEmpty[T]()}
	}

	return Value[T]{v: opt.New(v)}
}

func NewMap[T, R any](value nullable.Nullable[T], fn func(T) R) Value[R] {
	if value.IsNull() {
		return Value[R]{v: opt.NewEmpty[R](), delete: true}
	}

	v, err := value.Get()
	if err != nil {
		return Value[R]{v: opt.NewEmpty[R]()}
	}

	return Value[R]{v: opt.NewMap(v, fn)}
}

func NewMapErr[T, R any](value nullable.Nullable[T], fn func(T) (R, error)) (Value[R], error) {
	if value.IsNull() {
		return Value[R]{v: opt.NewEmpty[R](), delete: true}, nil
	}

	v, err := value.Get()
	if err != nil {
		return Value[R]{v: opt.NewEmpty[R]()}, nil
	}

	pv := opt.New(v)

	mv, err := opt.MapErr(pv, fn)
	if err != nil {
		return Value[R]{}, err
	}
	return Value[R]{v: mv}, nil
}

// Skip is an opt-out for cases where the inbound type is not a nullable type,
// but the receiving location is a deletable type. Not common, but a minor
// workaround to save designing more complex APIs at the service layer. Mostly
// those that involve `Partial` update struct types where it's not really worth
// creating a whole separate partial struct for creation vs updating a resource.
func Skip[T any](v opt.Optional[T]) Value[T] {
	return Value[T]{v: v}
}

func (v Value[T]) Get() (opt.Optional[T], bool) {
	return v.v, v.delete
}

func (v Value[T]) Call(
	setValueFn func(value T),
	deleteValueFn func(),
) {
	if v.delete {
		deleteValueFn()
	}

	v.v.Call(setValueFn)
}
