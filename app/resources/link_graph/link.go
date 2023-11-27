package link_graph

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/reply"
	"github.com/Southclaws/storyden/app/resources/thread"
	"github.com/Southclaws/storyden/internal/ent"
)

type ID = xid.ID

type WithRefs struct {
	ID          ID
	URL         string
	Slug        string
	Domain      string
	Title       opt.Optional[string]
	Description opt.Optional[string]
	Assets      []*asset.Asset
	Threads     []*thread.Thread
	Replies     []*reply.Reply
	Clusters    []*datagraph.Cluster
	Items       []*datagraph.Item
}

func (l *WithRefs) AssetIDs() []asset.AssetID {
	return dt.Map(l.Assets, func(a *asset.Asset) asset.AssetID { return a.ID })
}

type Repository interface {
	Get(ctx context.Context, slug string) (*WithRefs, error)
}

func Map(in *ent.Link) (*WithRefs, error) {
	postEdge, err := in.Edges.PostsOrErr()
	if err != nil {
		return nil, fault.Wrap(err)
	}

	clusterEdge, err := in.Edges.ClustersOrErr()
	if err != nil {
		return nil, fault.Wrap(err)
	}

	itemEdge, err := in.Edges.ItemsOrErr()
	if err != nil {
		return nil, fault.Wrap(err)
	}

	// Mapping

	threads, err := dt.MapErr(dt.Filter(postEdge, func(p *ent.Post) bool { return p.First }), thread.FromModel)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	replies, err := dt.MapErr(dt.Filter(postEdge, func(p *ent.Post) bool { return !p.First }), reply.FromModel)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	clusters, err := dt.MapErr(clusterEdge, datagraph.ClusterFromModel)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	items, err := dt.MapErr(itemEdge, datagraph.ItemFromModel)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	return &WithRefs{
		ID:          in.ID,
		URL:         in.URL,
		Slug:        in.Slug,
		Domain:      in.Domain,
		Title:       opt.New(in.Title),
		Description: opt.New(in.Description),
		Assets:      dt.Map(in.Edges.Assets, asset.FromModel),
		Threads:     threads,
		Replies:     replies,
		Clusters:    clusters,
		Items:       items,
	}, nil
}
