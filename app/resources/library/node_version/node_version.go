package node_version

import (
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/schema"
)

type VersionID xid.ID

func (i VersionID) String() string { return xid.ID(i).String() }

func VersionIDFromString(s string) (VersionID, error) {
	id, err := xid.FromString(s)
	if err != nil {
		return VersionID(xid.NilID()), err
	}
	return VersionID(id), nil
}

type PropertySnapshot struct {
	ID    opt.Optional[xid.ID]
	Name  string
	Type  opt.Optional[library.PropertyType]
	Value string
	Sort  opt.Optional[string]
}

type NodeVersion struct {
	ID        VersionID
	CreatedAt time.Time
	UpdatedAt time.Time

	NodeID   library.NodeID
	Author   profile.Ref
	Status   VersionStatus
	Previous opt.Optional[VersionReference]

	Name        string
	Slug        string
	Description opt.Optional[string]
	Content     opt.Optional[datagraph.Content]

	PropertiesSnapshot opt.Optional[[]PropertySnapshot]

	Metadata map[string]any
}

type VersionReference struct {
	ID        VersionID
	CreatedAt time.Time
	UpdatedAt time.Time
	Author    profile.Ref
	Status    VersionStatus
}

type NodeVersionWithNode struct {
	NodeVersion
	Node *library.Node
}

func Map(v *ent.NodeVersion) (*NodeVersion, error) {
	authorEdge, err := v.Edges.AuthorOrErr()
	if err != nil {
		return nil, fault.Wrap(err)
	}

	author, err := profile.MapRef(authorEdge)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	content, err := opt.MapErr(opt.NewPtr(v.Content), func(raw string) (datagraph.Content, error) {
		content, err := datagraph.NewRichTextWithBlocks(raw)
		if err != nil {
			return datagraph.Content{}, err
		}
		return content.Content, nil
	})
	if err != nil {
		return nil, fault.Wrap(err)
	}

	status, err := NewVersionStatus(string(v.Status))
	if err != nil {
		return nil, fault.Wrap(err)
	}

	metadata := v.Metadata
	if metadata == nil {
		metadata = make(map[string]any)
	}

	propertiesSnapshot := opt.NewEmpty[[]PropertySnapshot]()
	if v.PropertiesSnapshot.Set {
		properties, err := mapPropertySnapshots(v.PropertiesSnapshot.Properties)
		if err != nil {
			return nil, fault.Wrap(err)
		}

		propertiesSnapshot = opt.New(properties)
	}

	return &NodeVersion{
		ID:                 VersionID(v.ID),
		CreatedAt:          v.CreatedAt,
		UpdatedAt:          v.UpdatedAt,
		NodeID:             library.NodeID(v.NodeID),
		Author:             *author,
		Status:             status,
		Name:               v.Name,
		Slug:               v.Slug,
		Description:        opt.NewPtr(v.Description),
		Content:            content,
		Metadata:           metadata,
		PropertiesSnapshot: propertiesSnapshot,
	}, nil
}

func MapReference(v *ent.NodeVersion) (*VersionReference, error) {
	authorEdge, err := v.Edges.AuthorOrErr()
	if err != nil {
		return nil, fault.Wrap(err)
	}

	author, err := profile.MapRef(authorEdge)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	status, err := NewVersionStatus(string(v.Status))
	if err != nil {
		return nil, fault.Wrap(err)
	}

	return &VersionReference{
		ID:        VersionID(v.ID),
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
		Author:    *author,
		Status:    status,
	}, nil
}

func mapPropertySnapshots(in []schema.PropertySnapshotEntry) ([]PropertySnapshot, error) {
	return dt.MapErr(in, func(e schema.PropertySnapshotEntry) (PropertySnapshot, error) {
		id, err := opt.MapErr(opt.NewSafe(e.ID, e.ID != ""), func(s string) (xid.ID, error) {
			id, err := xid.FromString(s)
			if err != nil {
				return xid.NilID(), fault.Wrap(err)
			}
			return id, nil
		})
		if err != nil {
			return PropertySnapshot{}, fault.Wrap(err)
		}

		typ, err := opt.MapErr(opt.NewSafe(e.Type, e.Type != ""), func(t string) (library.PropertyType, error) {
			return library.NewPropertyType(t)
		})
		if err != nil {
			return PropertySnapshot{}, fault.Wrap(err)
		}

		return PropertySnapshot{
			ID:    id,
			Name:  e.Name,
			Type:  typ,
			Value: e.Value,
			Sort:  opt.NewSafe(e.Sort, e.Sort != ""),
		}, nil
	})
}

func UnmapPropertySnapshots(in []PropertySnapshot) schema.PropertySnapshot {
	entries := dt.Map(in, func(p PropertySnapshot) schema.PropertySnapshotEntry {
		var id string
		if v, ok := p.ID.Get(); ok {
			id = v.String()
		}

		var typ string
		if v, ok := p.Type.Get(); ok {
			typ = v.String()
		}

		var sort string
		if v, ok := p.Sort.Get(); ok {
			sort = v
		}

		return schema.PropertySnapshotEntry{
			ID:    id,
			Name:  p.Name,
			Type:  typ,
			Value: p.Value,
			Sort:  sort,
		}
	})

	return schema.PropertySnapshot{
		Set:        true,
		Properties: entries,
	}
}
