package weaviate_semdexer

import (
	"github.com/Southclaws/fault"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/datagraph"
)

func mapToNodeReference(v WeaviateObject) (*datagraph.Ref, error) {
	id, err := xid.FromString(v.DatagraphID)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	dk, err := datagraph.NewKind(v.DatagraphType)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	return &datagraph.Ref{
		ID:        id,
		Kind:      dk,
		Relevance: min(max(1-v.Additional.Distance, 0), 1),
	}, nil
}
