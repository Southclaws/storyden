package datagraph

import (
	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/internal/ent"
)

func ClusterFromModel(c *ent.Cluster) (*Cluster, error) {
	accEdge, err := c.Edges.OwnerOrErr()
	if err != nil {
		return nil, fault.Wrap(err)
	}

	pro, err := profile.FromModel(accEdge)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	clusters, err := dt.MapErr(c.Edges.Clusters, ClusterFromModel)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	visibility, err := post.NewVisibility(c.Visibility.String())
	if err != nil {
		return nil, fault.Wrap(err)
	}

	assets := dt.Map(c.Edges.Assets, asset.FromModel)

	return &Cluster{
		ID:          ClusterID(c.ID),
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
		Name:        c.Name,
		Slug:        c.Slug,
		Assets:      assets,
		Links:       dt.Map(c.Edges.Links, LinkFromModel),
		Description: c.Description,
		Content:     opt.NewPtr(c.Content),
		Owner:       *pro,
		Clusters:    clusters,
		Visibility:  visibility,
		Properties:  c.Properties,
	}, nil
}
