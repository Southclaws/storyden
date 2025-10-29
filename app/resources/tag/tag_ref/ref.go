package tag_ref

import (
	"github.com/Southclaws/dt"
	"github.com/mazznoer/csscolorparser"
	"github.com/rs/xid"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/internal/ent"
)

type ID xid.ID

func (id ID) String() string {
	return xid.ID(id).String()
}

type Name struct {
	s string
}

func NewName(s string) Name {
	return Name{s: mark.Slugify(s)}
}

func (n Name) String() string {
	return n.s
}

type Names []Name

func (n Names) Strings() []string {
	return dt.Map(n, func(n Name) string {
		return n.String()
	})
}

type Tag struct {
	ID        ID
	Name      Name
	Colour    string
	ItemCount int
}

type Tags []*Tag

func (a Tags) Len() int           { return len(a) }
func (a Tags) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Tags) Less(i, j int) bool { return a[i].ItemCount > a[j].ItemCount }

func (t Tags) Names() []Name {
	return dt.Map(t, func(t *Tag) Name {
		return t.Name
	})
}

func Map(counts TagItemsResults) func(in *ent.Tag) *Tag {
	return func(in *ent.Tag) *Tag {
		return &Tag{
			ID:        ID(in.ID),
			Name:      NewName(in.Name),
			Colour:    deriveTagColour(in.Name),
			ItemCount: counts.Get(in.ID),
		}
	}
}

type TagItemsResult struct {
	TagID xid.ID `db:"tag_id"`
	Count int    `db:"items"`
}

type TagItemsResults []TagItemsResult

func (t TagItemsResults) Get(id xid.ID) int {
	table := lo.KeyBy(t, func(r TagItemsResult) xid.ID { return r.TagID })
	c, ok := table[id]
	if !ok {
		return 0
	}
	return c.Count
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
