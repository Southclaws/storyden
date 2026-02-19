package link

import (
	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/role/held"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/link/link_ref"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/internal/ent"
)

type LinkID xid.ID

type Link struct {
	link_ref.LinkRef

	Assets  []*asset.Asset
	Posts   []*post.Post
	Nodes   []*library.Node
	Related datagraph.ItemList
}

func (l *Link) AssetIDs() []asset.AssetID {
	return dt.Map(l.Assets, func(a *asset.Asset) asset.AssetID { return a.ID })
}

func NewLink(url, title, description string) link_ref.LinkRef {
	return link_ref.LinkRef{
		URL:         url,
		Title:       opt.New(title),
		Description: opt.New(description),
	}
}

func NewLinkOpt(purl, ptitle, pdescription *string) opt.Optional[link_ref.LinkRef] {
	if purl == nil {
		return opt.NewEmpty[link_ref.LinkRef]()
	}

	return opt.New(link_ref.LinkRef{
		URL:         opt.NewPtr(purl).String(),
		Title:       opt.NewPtr(ptitle),
		Description: opt.NewPtr(pdescription),
	})
}

func Map(in *ent.Link, roleHydratorFn func(accID xid.ID) (held.Roles, error)) (*Link, error) {
	postEdge, err := in.Edges.PostsOrErr()
	if err != nil {
		return nil, fault.Wrap(err)
	}

	nodeEdge, err := in.Edges.NodesOrErr()
	if err != nil {
		return nil, fault.Wrap(err)
	}

	posts, err := dt.MapErr(postEdge, func(in *ent.Post) (*post.Post, error) {
		return post.Map(in, roleHydratorFn)
	})
	if err != nil {
		return nil, fault.Wrap(err)
	}

	nodes, err := dt.MapErr(nodeEdge, library.MapNode(true, nil, roleHydratorFn))
	if err != nil {
		return nil, fault.Wrap(err)
	}

	return &Link{
		LinkRef: *link_ref.Map(in),
		Assets:  dt.Map(in.Edges.Assets, asset.Map),
		Posts:   posts,
		Nodes:   nodes,
	}, nil
}
