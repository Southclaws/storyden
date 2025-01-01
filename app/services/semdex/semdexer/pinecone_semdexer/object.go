package pinecone_semdexer

import (
	"fmt"
	"net/url"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/services/semdex"
	"github.com/Southclaws/storyden/internal/infrastructure/vector/pinecone"
)

type Object struct {
	ID        xid.ID
	Kind      datagraph.Kind
	Relevance float64
	URL       url.URL
	Content   string
}

type Objects []*Object

func (o *Object) ToChunk() *semdex.Chunk {
	return &semdex.Chunk{
		ID:      o.ID,
		Kind:    o.Kind,
		URL:     o.URL,
		Content: o.Content,
	}
}

func (o *Object) ToRef() *datagraph.Ref {
	return &datagraph.Ref{
		ID:        o.ID,
		Kind:      o.Kind,
		Relevance: o.Relevance,
	}
}

func (o Objects) ToChunks() []*semdex.Chunk {
	chunks := make([]*semdex.Chunk, len(o))
	for i, object := range o {
		chunks[i] = object.ToChunk()
	}
	return chunks
}

func (o Objects) ToRefs() datagraph.RefList {
	refs := make(datagraph.RefList, len(o))
	for i, object := range o {
		refs[i] = object.ToRef()
	}
	return refs
}

func mapObject(v *pinecone.ScoredVector) (*Object, error) {
	meta := v.Vector.Metadata.AsMap()

	idRaw, ok := meta["datagraph_id"]
	if !ok {
		return nil, fault.New("missing datagraph_id in metadata")
	}

	typeRaw, ok := meta["datagraph_type"]
	if !ok {
		return nil, fault.New("missing datagraph_type in metadata")
	}

	contentRaw, ok := meta["content"]
	if !ok {
		return nil, fault.New("missing content in metadata")
	}

	//

	idString, ok := idRaw.(string)
	if !ok {
		return nil, fault.New("datagraph_id in metadata is not a string")
	}

	typeString, ok := typeRaw.(string)
	if !ok {
		return nil, fault.New("datagraph_type in metadata is not a string")
	}

	content, ok := contentRaw.(string)
	if !ok {
		return nil, fault.New("content in metadata is not a string")
	}

	//

	id, err := xid.FromString(idString)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	dk, err := datagraph.NewKind(typeString)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	sdr, err := url.Parse(fmt.Sprintf("%s:%s/%s", datagraph.RefScheme, dk.String(), id.String()))
	if err != nil {
		return nil, err
	}

	return &Object{
		ID:        id,
		Kind:      dk,
		Relevance: float64((v.Score + 1) / 2),
		URL:       *sdr,
		Content:   content,
	}, nil
}

func mapObjects(objects []*pinecone.ScoredVector) (Objects, error) {
	return dt.MapErr(objects, mapObject)
}
