package link

import (
	"context"
	"sort"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/link"
)

type ID = xid.ID

type Link struct {
	ID          ID
	URL         string
	Title       opt.Optional[string]
	Description opt.Optional[string]
	Assets      []*asset.Asset
}

func (l *Link) AssetIDs() []asset.AssetID {
	return dt.Map(l.Assets, func(a *asset.Asset) asset.AssetID { return a.ID })
}

type Repository interface {
	Store(ctx context.Context, url, title, description string, opts ...Option) (*Link, error)
	Search(ctx context.Context, filters ...Filter) ([]*Link, error)
}

type (
	Option func(*ent.LinkMutation)
	Filter func(*ent.LinkQuery)
)

func WithPosts(ids ...xid.ID) Option {
	return func(lm *ent.LinkMutation) {
		lm.AddPostIDs(ids...)
	}
}

func WithClusters(ids ...xid.ID) Option {
	return func(lm *ent.LinkMutation) {
		lm.AddClusterIDs(ids...)
	}
}

func WithItems(ids ...xid.ID) Option {
	return func(lm *ent.LinkMutation) {
		lm.AddItemIDs(ids...)
	}
}

func WithAssets(ids ...string) Option {
	return func(lm *ent.LinkMutation) {
		lm.AddAssetIDs(ids...)
	}
}

func WithURL(s string) Filter {
	return func(lq *ent.LinkQuery) {
		lq.Where(link.URLContainsFold(s))
	}
}

func WithPage(page, size int) Filter {
	return func(lq *ent.LinkQuery) {
		lq.Limit(size).Offset(page * size)
	}
}

func WithKeyword(s string) Filter {
	return func(lq *ent.LinkQuery) {
		lq.Where(link.Or(
			link.TitleContainsFold(s),
			link.DescriptionContainsFold(s),
			link.URLContainsFold(s),
		))
	}
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

func Map(in *ent.Link) *Link {
	return &Link{
		ID:          in.ID,
		URL:         in.URL,
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
func (a Links) Less(i, j int) bool { return a[i].ID.String() < a[j].ID.String() }

func MapA(in []*ent.Link) []*Link {
	list := dt.Map(in, Map)
	sort.Sort(Links(list))
	return list
}
