package tag_ref

import (
	"github.com/Southclaws/dt"
	"github.com/mazznoer/csscolorparser"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/internal/ent"
)

type ID xid.ID

func (id ID) String() string {
	return xid.ID(id).String()
}

type Name string

type Names []Name

type Tag struct {
	ID     ID
	Name   Name
	Colour string
}

type Tags []*Tag

func (t Tags) Names() []Name {
	return dt.Map(t, func(t *Tag) Name {
		return t.Name
	})
}

func Map(in *ent.Tag) *Tag {
	return &Tag{
		ID:     ID(in.ID),
		Name:   Name(in.Name),
		Colour: deriveTagColour(in.Name),
	}
}

func deriveTagColour(s string) string {
	hash := dt.Reduce([]byte(s), func(r uint16, b byte) uint16 {
		s := uint16(b) * 42
		x := (r + 1) * s % 360
		return x
	}, 69)

	hue := float64(hash)

	c2 := csscolorparser.FromOklch(0.7226, 0.12, hue, 1.0)

	return c2.HexString()
}
