package chromem_semdexer

import (
	"context"
	"fmt"
	"math"
	"net/url"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/philippgille/chromem-go"
	"github.com/rs/xid"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/datagraph/hydrate"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/services/search/searcher"
	"github.com/Southclaws/storyden/app/services/semdex"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/ai"
)

type chromemRefIndex struct {
	db       *chromem.DB
	c        *chromem.Collection
	hydrator *hydrate.Hydrator
}

func New(cfg config.Config, rh *hydrate.Hydrator, aip ai.Prompter) (semdex.Semdexer, error) {
	db, err := chromem.NewPersistentDB(cfg.SemdexLocalPath, false)
	if err != nil {
		return nil, err
	}

	if _, ok := aip.(*ai.Disabled); ok {
		return nil, fault.New("a language model provider must be enabled for the embedded semdexer to be enabled")
	}

	ef := aip.EmbeddingFunc()

	collection, err := db.GetOrCreateCollection("semdex", nil, ef)
	if err != nil {
		return nil, err
	}

	return &chromemRefIndex{
		db:       db,
		c:        collection,
		hydrator: rh,
	}, nil
}

func (c *chromemRefIndex) Index(ctx context.Context, object datagraph.Item) (int, error) {
	err := c.c.AddDocument(ctx, chromem.Document{
		ID:      object.GetID().String(),
		Content: object.GetContent().Plaintext(),
		Metadata: map[string]string{
			"datagraph_kind": object.GetKind().String(),
		},
	})
	if err != nil {
		return 0, fault.Wrap(err, fctx.With(ctx))
	}

	return 1, nil
}

func (c *chromemRefIndex) Delete(ctx context.Context, object xid.ID) (int, error) {
	err := c.c.Delete(ctx, nil, nil, object.String())
	if err != nil {
		return 0, fault.Wrap(err, fctx.With(ctx))
	}

	return 1, nil
}

func (c *chromemRefIndex) Search(ctx context.Context, q string, p pagination.Parameters, opts searcher.Options) (*pagination.Result[datagraph.Item], error) {
	nr := min(c.c.Count(), p.Size())
	if nr == 0 {
		res := pagination.NewPageResult[datagraph.Item](p, 0, nil)
		return &res, nil
	}

	rs, err := c.c.Query(ctx, q, nr, nil, nil)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	filtered := lo.Filter(rs, func(r chromem.Result, _ int) bool {
		return r.Similarity > 0.2
	})

	list, err := mapResults(filtered)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	hyd, err := c.hydrator.Hydrate(ctx, list...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	results := pagination.NewPageResult(p, len(rs), hyd)

	return &results, nil
}

func (c *chromemRefIndex) SearchRefs(ctx context.Context, q string, p pagination.Parameters, opts searcher.Options) (*pagination.Result[*datagraph.Ref], error) {
	nr := min(c.c.Count(), p.Size())
	if nr == 0 {
		res := pagination.NewPageResult[*datagraph.Ref](p, 0, nil)
		return &res, nil
	}

	rs, err := c.c.Query(ctx, q, nr, nil, nil)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	filtered := lo.Filter(rs, func(r chromem.Result, _ int) bool {
		return r.Similarity > 0.2
	})

	list, err := mapResults(filtered)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	results := pagination.NewPageResult(p, len(rs), list)

	return &results, nil
}

func (c *chromemRefIndex) SearchChunks(ctx context.Context, q string, p pagination.Parameters, opts searcher.Options) ([]*semdex.Chunk, error) {
	nr := min(c.c.Count(), p.Size())
	if nr == 0 {
		return []*semdex.Chunk{}, nil
	}

	rs, err := c.c.Query(ctx, q, nr, nil, nil)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	filtered := lo.Filter(rs, func(r chromem.Result, _ int) bool {
		return r.Similarity > 0.2
	})

	list, err := mapChunks(filtered)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return list, nil
}

func (c *chromemRefIndex) RecommendRefs(ctx context.Context, object datagraph.Item) (datagraph.RefList, error) {
	doc, err := c.c.GetByID(ctx, object.GetID().String())
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	nr := min(c.c.Count(), 10)

	rs, err := c.c.QueryEmbedding(ctx, doc.Embedding, nr, nil, nil)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	list, err := mapResults(rs)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return list, nil
}

func (c *chromemRefIndex) Recommend(ctx context.Context, object datagraph.Item) (datagraph.ItemList, error) {
	doc, err := c.c.GetByID(ctx, object.GetID().String())
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	nr := min(c.c.Count(), 10)

	rs, err := c.c.QueryEmbedding(ctx, doc.Embedding, nr, nil, nil)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	list, err := mapResults(rs)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	items, err := c.hydrator.Hydrate(ctx, list...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return items, nil
}

func (c *chromemRefIndex) ScoreRelevance(ctx context.Context, object datagraph.Item, ids ...xid.ID) (map[xid.ID]float64, error) {
	src, err := c.c.GetByID(ctx, object.GetID().String())
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	srcCoord := tof64(src.Embedding)

	cluster, err := dt.MapErr(ids, func(id xid.ID) (*chromem.Document, error) {
		doc, err := c.c.GetByID(ctx, id.String())
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		return &doc, nil
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result := dt.Reduce(cluster, func(acc map[xid.ID]float64, curr *chromem.Document) map[xid.ID]float64 {
		cosine, err := cosine(srcCoord, tof64(curr.Embedding))
		if err != nil {
			return acc
		}

		id, err := xid.FromString(curr.ID)
		if err != nil {
			return acc
		}

		acc[id] = cosine
		return acc
	}, map[xid.ID]float64{})

	return result, nil
}

func mapChunks(rs []chromem.Result) ([]*semdex.Chunk, error) {
	return dt.MapErr(rs, mapChunk)
}

func mapChunk(r chromem.Result) (*semdex.Chunk, error) {
	id, err := xid.FromString(r.ID)
	if err != nil {
		return nil, err
	}

	dk, ok := r.Metadata["datagraph_kind"]
	if !ok {
		return nil, fault.New("missing datagraph_kind metadata")
	}

	k, err := datagraph.NewKind(dk)
	if err != nil {
		return nil, err
	}

	sdr, err := url.Parse(fmt.Sprintf("%s:%s/%s", datagraph.RefScheme, k, id.String()))
	if err != nil {
		return nil, err
	}

	return &semdex.Chunk{
		ID:      id,
		Kind:    k,
		URL:     *sdr,
		Content: r.Content,
	}, nil
}

func mapResults(rs []chromem.Result) (datagraph.RefList, error) {
	return dt.MapErr(rs, mapResult)
}

func mapResult(r chromem.Result) (*datagraph.Ref, error) {
	id, err := xid.FromString(r.ID)
	if err != nil {
		return nil, err
	}

	dk, ok := r.Metadata["datagraph_kind"]
	if !ok {
		return nil, fault.New("missing datagraph_kind metadata")
	}

	k, err := datagraph.NewKind(dk)
	if err != nil {
		return nil, err
	}

	return &datagraph.Ref{
		ID:        id,
		Kind:      k,
		Relevance: float64(r.Similarity),
	}, nil
}

func mapDoc(r chromem.Document) (*datagraph.Ref, error) {
	id, err := xid.FromString(r.ID)
	if err != nil {
		return nil, err
	}

	dk, ok := r.Metadata["datagraph_kind"]
	if !ok {
		return nil, fault.New("missing datagraph_kind metadata")
	}

	k, err := datagraph.NewKind(dk)
	if err != nil {
		return nil, err
	}

	return &datagraph.Ref{
		ID:   id,
		Kind: k,
	}, nil
}

func (c *chromemRefIndex) GetMany(ctx context.Context, limit uint, ids ...xid.ID) (datagraph.RefList, error) {
	refs := datagraph.RefList{}

	for _, id := range ids {
		r, err := c.c.GetByID(ctx, id.String())
		if err != nil {
			continue
		}

		ref, err := mapDoc(r)
		if err != nil {
			continue
		}

		refs = append(refs, ref)
	}

	return refs, nil
}

func cosine(a []float64, b []float64) (cosine float64, err error) {
	c := 0
	la := len(a)
	lb := len(b)

	if la > lb {
		c = la
	} else {
		c = lb
	}

	sum := 0.0
	s1 := 0.0
	s2 := 0.0

	for k := 0; k < c; k++ {
		if k >= la {
			s2 += math.Pow(b[k], 2)
			continue
		}

		if k >= lb {
			s1 += math.Pow(a[k], 2)
			continue
		}

		sum += a[k] * b[k]
		s1 += math.Pow(a[k], 2)
		s2 += math.Pow(b[k], 2)
	}

	if s1 == 0 || s2 == 0 {
		return 0.0, fault.New("failed to compute cosine similarity on zero length vectors")
	}

	return sum / (math.Sqrt(s1) * math.Sqrt(s2)), nil
}

func tof64(a []float32) []float64 {
	b := make([]float64, len(a))
	for i, v := range a {
		b[i] = float64(v)
	}
	return b
}
