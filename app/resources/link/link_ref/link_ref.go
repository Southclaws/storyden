package link_ref

import (
	"sort"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/internal/ent"
)

type ID = xid.ID

type LinkRef struct {
	ID        ID
	CreatedAt time.Time
	UpdatedAt time.Time

	URL          string
	Slug         string
	Domain       string
	Title        opt.Optional[string]
	Description  opt.Optional[string]
	FaviconImage opt.Optional[asset.Asset]
	PrimaryImage opt.Optional[asset.Asset]
}

type LinkRefs []*LinkRef

func (l LinkRefs) Latest() opt.Optional[*LinkRef] {
	if len(l) == 0 {
		return opt.NewEmpty[*LinkRef]()
	}

	return opt.New(l[0])
}

func (a LinkRefs) Len() int           { return len(a) }
func (a LinkRefs) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a LinkRefs) Less(i, j int) bool { return xid.ID(a[i].ID).String() > xid.ID(a[j].ID).String() }

func LinksFromModel(in []*ent.Link) []*LinkRef {
	list := dt.Map(in, Map)
	sort.Sort(LinkRefs(list))
	return list
}

func Map(in *ent.Link) *LinkRef {
	favicon := opt.NewPtrMap(in.Edges.FaviconImage, func(a ent.Asset) asset.Asset {
		return *asset.Map(&a)
	})

	primary := opt.NewPtrMap(in.Edges.PrimaryImage, func(a ent.Asset) asset.Asset {
		return *asset.Map(&a)
	})

	return &LinkRef{
		ID:           ID(in.ID),
		URL:          in.URL,
		Slug:         in.Slug,
		Domain:       in.Domain,
		Title:        opt.New(in.Title),
		Description:  opt.New(in.Description),
		FaviconImage: favicon,
		PrimaryImage: primary,
	}
}
