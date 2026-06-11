package node_versioning

import (
	"github.com/Southclaws/dt"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_version"
)

func propertyMutationsToSnapshot(in library.PropertyMutationList) []node_version.PropertySnapshot {
	return dt.Map(in, func(p *library.PropertyMutation) node_version.PropertySnapshot {
		return node_version.PropertySnapshot{
			ID:    p.ID,
			Name:  p.Name,
			Type:  p.Type,
			Value: p.Value,
			Sort:  p.Sort,
		}
	})
}

func snapshotToPropertyMutations(in []node_version.PropertySnapshot) library.PropertyMutationList {
	return dt.Map(in, func(p node_version.PropertySnapshot) *library.PropertyMutation {
		return &library.PropertyMutation{
			ID:    p.ID,
			Name:  p.Name,
			Value: p.Value,
			Type:  p.Type,
			Sort:  p.Sort,
		}
	})
}

func propertyTableToSnapshot(in opt.Optional[library.PropertyTable]) []node_version.PropertySnapshot {
	table, ok := in.Get()
	if !ok {
		return []node_version.PropertySnapshot{}
	}

	return dt.Map(table.Properties, func(p *library.Property) node_version.PropertySnapshot {
		return node_version.PropertySnapshot{
			ID:    opt.New(p.Field.ID),
			Name:  p.Field.Name,
			Type:  opt.New(p.Field.Type),
			Value: p.Value.OrZero(),
			Sort:  opt.NewSafe(p.Field.Sort, p.Field.Sort != ""),
		}
	})
}
