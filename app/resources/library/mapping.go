package library

import (
	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/role/held"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/collection/collection_item_status"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/link/link_ref"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/internal/ent"
)

func MapNode(isRoot bool, ps *PropertySchemaTable, roleHydratorFn func(accID xid.ID) (held.Roles, error)) func(c *ent.Node) (*Node, error) {
	profileMapper := profile.RefMapper(roleHydratorFn)

	return func(c *ent.Node) (*Node, error) {
		accEdge, err := c.Edges.OwnerOrErr()
		if err != nil {
			return nil, fault.Wrap(err)
		}

		pro, err := profileMapper(accEdge)
		if err != nil {
			return nil, fault.Wrap(err)
		}

		parent, err := opt.MapErr(opt.NewPtr(c.Edges.Parent), func(c ent.Node) (Node, error) {
			p, err := MapNode(false, ps, roleHydratorFn)(&c)
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

		nodes, err := dt.MapErr(c.Edges.Nodes, MapNode(false, ps, roleHydratorFn))
		if err != nil {
			return nil, fault.Wrap(err)
		}

		visibility, err := visibility.NewVisibility(c.Visibility.String())
		if err != nil {
			return nil, fault.Wrap(err)
		}

		assets := dt.Map(c.Edges.Assets, asset.Map)

		richContent, err := opt.MapErr(opt.NewPtr(c.Content), datagraph.NewRichText)
		if err != nil {
			return nil, err
		}

		primaryImage := opt.Map(opt.NewPtr(c.Edges.PrimaryImage), func(e ent.Asset) asset.Asset {
			return *asset.Map(&e)
		})

		// This edge is optional anyway, so if not loaded, nothing bad happens.
		link := opt.Map(opt.NewPtr(c.Edges.Link), func(in ent.Link) link_ref.LinkRef {
			return *link_ref.Map(&in)
		})

		metadata := c.Metadata
		if metadata == nil {
			metadata = make(map[string]any)
		}

		n := &Node{
			Mark:          NewMark(c.ID, c.Slug),
			CreatedAt:     c.CreatedAt,
			UpdatedAt:     c.UpdatedAt,
			IndexedAt:     opt.NewPtr(c.IndexedAt),
			Name:          c.Name,
			Assets:        assets,
			WebLink:       link,
			Content:       richContent,
			Description:   opt.NewPtr(c.Description),
			PrimaryImage:  primaryImage,
			Owner:         *pro,
			Parent:        parent,
			HideChildTree: c.HideChildTree,
			Tags:          tags,
			Collections:   collection_item_status.Status{
				// NOTE: Members cannot yet add nodes to collections.
			},
			Nodes:      nodes,
			Visibility: visibility,
			SortKey:    c.Sort,
			Metadata:   metadata,
		}

		if ps != nil {
			// Sibling properties may contain values, so we pass in the edge.
			n.Properties = opt.NewPtr(ps.BuildPropertyTable(c.Edges.Properties, isRoot))

			if isRoot {
				// Child properties don't contain values, only the property schemas.
				n.ChildProperties = opt.NewPtr(ps.ChildSchemas())
			}
		}

		return n, nil
	}
}

func ItemRef(c *ent.Node) (datagraph.Item, error) {
	content, err := opt.MapErr(opt.NewPtr(c.Content), datagraph.NewRichText)
	if err != nil {
		return nil, err
	}

	return &Node{
		Mark:        NewMark(c.ID, c.Slug),
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
		IndexedAt:   opt.NewPtr(c.IndexedAt),
		Name:        c.Name,
		Content:     content,
		Description: opt.NewPtr(c.Description),
		Visibility:  visibility.NewVisibilityFromEnt(c.Visibility),
		SortKey:     c.Sort,
		Metadata:    c.Metadata,
	}, nil
}
