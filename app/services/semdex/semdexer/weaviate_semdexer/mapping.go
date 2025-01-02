package weaviate_semdexer

import (
	"encoding/json"
	"hash/fnv"
	"strconv"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/google/uuid"
	"github.com/rs/xid"
	"github.com/weaviate/weaviate/entities/models"

	"github.com/Southclaws/storyden/app/resources/datagraph"
)

type WeaviateObject struct {
	DatagraphID   string             `json:"datagraph_id"`
	DatagraphType string             `json:"datagraph_type"`
	Name          string             `json:"name"`
	Content       string             `json:"content"`
	Additional    WeaviateAdditional `json:"_additional"`
}

type WeaviateAdditional struct {
	Distance float64 `json:"distance"`
	Summary  []struct {
		Property string `json:"property"`
		Result   string `json:"result"`
	} `json:"summary"`
	Generate struct {
		SingleResult string `json:"singleResult"`
		Error        string `json:"error"`
	} `json:"generate"`
	Score        string `json:"score"`
	ExplainScore string `json:"explainScore"`
}

type WeaviateContent map[string][]WeaviateObject

type WeaviateResponse struct {
	Get     WeaviateContent
	Explore WeaviateContent
}

func mapToNodeReference(v WeaviateObject) (*datagraph.Ref, error) {
	id, err := xid.FromString(v.DatagraphID)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	dk, err := datagraph.NewKind(v.DatagraphType)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	var relevance float64
	if v.Additional.Distance > 0 {
		// Distances are inverse to "scores" (complexionary or relevance)
		relevance = min(max(1-v.Additional.Distance, 0), 1)
	} else if v.Additional.Score != "" {
		relevance, err = strconv.ParseFloat(v.Additional.Score, 64)
		if err != nil {
			return nil, fault.Wrap(err)
		}
	}

	return &datagraph.Ref{
		ID:        id,
		Kind:      dk,
		Relevance: relevance,
	}, nil
}

func mapResponseObjects(raw map[string]models.JSONObject) (*WeaviateResponse, error) {
	j, err := json.Marshal(raw)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	parsed := WeaviateResponse{}
	err = json.Unmarshal(j, &parsed)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	return &parsed, nil
}

func generateChunkID(id xid.ID, chunk string) uuid.UUID {
	// We don't currently support sharing chunks across content nodes, so append
	// the object's ID to the chunk's hash, to ensure it's unique to the object.
	payload := []byte(append(id.Bytes(), chunk...))

	return uuid.NewHash(fnv.New128(), uuid.NameSpaceOID, payload, 4)
}

func chunkIDsFor(id xid.ID) func(chunk string) uuid.UUID {
	return func(chunk string) uuid.UUID {
		// We don't currently support sharing chunks across content nodes, so append
		// the object's ID to the chunk's hash, to ensure it's unique to the object.
		payload := []byte(append(id.Bytes(), chunk...))

		return uuid.NewHash(fnv.New128(), uuid.NameSpaceOID, payload, 4)
	}
}

func chunkIDsForItem(object datagraph.Item) []uuid.UUID {
	return dt.Map(object.GetContent().Split(), chunkIDsFor(object.GetID()))
}
