package datagraph

import (
	"net/url"
	"path"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"
)

// Ref is a non-hydrated type to express a result type from semdex operations
// such as searching or recommendations. It can be hydrated into an Item using
// the Kind field to select a relevant resource querier to find the full object.
type Ref struct {
	ID        xid.ID
	Kind      Kind
	Relevance float64
}

func (r *Ref) GetID() xid.ID {
	return r.ID
}

func (r *Ref) GetKind() Kind {
	return r.Kind
}

type RefList []*Ref

func (a RefList) Len() int           { return len(a) }
func (a RefList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a RefList) Less(i, j int) bool { return a[i].Relevance > a[j].Relevance }

func NewRef(i Item) *Ref {
	return &Ref{
		ID:   i.GetID(),
		Kind: i.GetKind(),
	}
}

func NewRefFromSDR(u url.URL) (*Ref, error) {
	if u.Scheme != RefScheme {
		return nil, fault.Wrap(ErrInvalidReferenceScheme, ftag.With(ftag.InvalidArgument))
	}

	resourcePath, identifier := path.Split(u.Opaque)
	resource := strings.Trim(resourcePath, "/")

	id, err := xid.FromString(identifier)
	if err != nil {
		return nil, err
	}

	k, err := NewKind(resource)
	if err != nil {
		return nil, fault.Wrap(err, ftag.With(ftag.InvalidArgument))
	}

	return &Ref{
		ID:   id,
		Kind: k,
	}, nil
}
