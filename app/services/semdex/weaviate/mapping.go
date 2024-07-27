package weaviate

import (
	"github.com/Southclaws/fault"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/datagraph"
)

func mapToNodeReference(v WeaviateObject) (*datagraph.NodeReference, error) {
	id, err := xid.FromString(v.DatagraphID)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	dk, err := datagraph.NewKind(v.DatagraphType)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	return &datagraph.NodeReference{
		ID:    id,
		Kind:  dk,
		Name:  v.Name,
		Score: v.Additional.Distance,
	}, nil
}
