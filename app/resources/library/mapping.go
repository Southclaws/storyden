package library

import (
	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"

	"github.com/Southclaws/storyden/app/resources/link/link_ref"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/internal/ent"
)

func NodeFromModel(c *ent.Node) (*Node, error) {
	accEdge, err := c.Edges.OwnerOrErr()
	if err != nil {
		return nil, fault.Wrap(err)
	}

	pro, err := profile.ProfileFromModel(accEdge)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	parent, err := opt.MapErr(opt.NewPtr(c.Edges.Parent), func(c ent.Node) (Node, error) {
		p, err := NodeFromModel(&c)
		if err != nil {
			return Node{}, err
		}
		return *p, nil
	})
	if err != nil {
		return nil, err
	}

	tagsEdge := c.Edges.Tags
	tags := dt.Map(tagsEdge, tag_ref.Map(nil))

	nodes, err := dt.MapErr(c.Edges.Nodes, NodeFromModel)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	visibility, err := visibility.NewVisibility(c.Visibility.String())
	if err != nil {
		return nil, fault.Wrap(err)
	}

	assets := dt.Map(c.Edges.Assets, asset.FromModel)

	richContent, err := opt.MapErr(opt.NewPtr(c.Content), datagraph.NewRichText)
	if err != nil {
		return nil, err
	}

	primaryImage := opt.Map(opt.NewPtr(c.Edges.PrimaryImage), func(e ent.Asset) asset.Asset {
		return *asset.FromModel(&e)
	})

	// This edge is optional anyway, so if not loaded, nothing bad happens.
	link := opt.Map(opt.NewPtr(c.Edges.Link), func(in ent.Link) link_ref.LinkRef {
		return *link_ref.Map(&in)
	})

	return &Node{
		Mark:         NewMark(c.ID, c.Slug),
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
		Name:         c.Name,
		Assets:       assets,
		WebLink:      link,
		Content:      richContent,
		Description:  opt.NewPtr(c.Description),
		PrimaryImage: primaryImage,
		Owner:        *pro,
		Parent:       parent,
		Tags:         tags,
		Nodes:        nodes,
		Visibility:   visibility,
		Metadata:     c.Metadata,
	}, nil
}
