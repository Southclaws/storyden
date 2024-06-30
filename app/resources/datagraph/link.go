package datagraph

import (
	"sort"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/internal/ent"
)

type LinkID xid.ID

type Link struct {
	ID          LinkID
	URL         string
	Slug        string
	Domain      string
	Title       opt.Optional[string]
	Description opt.Optional[string]
	Assets      []*asset.Asset
}

func (l *Link) AssetIDs() []asset.AssetID {
	return dt.Map(l.Assets, func(a *asset.Asset) asset.AssetID { return a.ID })
}

func NewLink(url, title, description string) Link {
	return Link{
		URL:         url,
		Title:       opt.New(title),
		Description: opt.New(description),
	}
}

func NewLinkOpt(purl, ptitle, pdescription *string) opt.Optional[Link] {
	if purl == nil {
		return opt.NewEmpty[Link]()
	}

	return opt.New(Link{
		URL:         opt.NewPtr(purl).String(),
		Title:       opt.NewPtr(ptitle),
		Description: opt.NewPtr(pdescription),
	})
}

func LinkFromModel(in *ent.Link) *Link {
	return &Link{
		ID:          LinkID(in.ID),
		URL:         in.URL,
		Slug:        in.Slug,
		Domain:      in.Domain,
		Title:       opt.New(in.Title),
		Description: opt.New(in.Description),
		Assets:      dt.Map(in.Edges.Assets, asset.FromModel),
	}
}

type Links []*Link

func (l Links) Latest() opt.Optional[*Link] {
	if len(l) == 0 {
		return opt.NewEmpty[*Link]()
	}

	return opt.New(l[0])
}

func (a Links) Len() int           { return len(a) }
func (a Links) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Links) Less(i, j int) bool { return xid.ID(a[i].ID).String() > xid.ID(a[j].ID).String() }

func LinksFromModel(in []*ent.Link) []*Link {
	list := dt.Map(in, LinkFromModel)
	sort.Sort(Links(list))
	return list
}
