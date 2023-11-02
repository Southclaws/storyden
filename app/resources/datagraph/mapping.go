package datagraph

import (
	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/link"
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

	items, err := dt.MapErr(c.Edges.Items, ItemFromModel)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	clusters, err := dt.MapErr(c.Edges.Clusters, ClusterFromModel)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	return &Cluster{
		ID:          ClusterID(c.ID),
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
		Name:        c.Name,
		Slug:        c.Slug,
		ImageURL:    opt.NewPtr(c.ImageURL),
		Links:       dt.Map(c.Edges.Links, link.Map),
		Description: c.Description,
		Content:     opt.NewPtr(c.Content),
		Owner:       *pro,
		Items:       items,
		Clusters:    clusters,
		Properties:  c.Properties,
	}, nil
}

func ItemFromModel(c *ent.Item) (*Item, error) {
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

	return &Item{
		ID:          ItemID(c.ID),
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
		Name:        c.Name,
		Slug:        c.Slug,
		ImageURL:    opt.NewPtr(c.ImageURL),
		Links:       dt.Map(c.Edges.Links, link.Map),
		Description: c.Description,
		Content:     opt.NewPtr(c.Content),
		Owner:       *pro,
		In:          clusters,
		Properties:  c.Properties,
	}, nil
}
