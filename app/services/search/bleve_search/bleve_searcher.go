package bleve_search

import (
	"context"
	"fmt"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/blevesearch/bleve/v2/search/query"
	"github.com/rs/xid"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/datagraph/hydrate"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/services/search/searcher"
	"github.com/Southclaws/storyden/internal/config"
)

type Document struct {
	ID          string
	Kind        string
	Name        string
	Slug        string
	Description string
	Content     string
	CreatedAt   int64
}

type BleveSearcher struct {
	client   bleve.Index
	hydrator *hydrate.Hydrator
}

func New(ctx context.Context, cfg config.Config, hydrator *hydrate.Hydrator) (*BleveSearcher, error) {
	if cfg.SearchProvider != "bleve" {
		return nil, nil
	}

	if cfg.BlevePath == "" {
		return nil, fault.New("BLEVE_PATH is required when SEARCH_PROVIDER is set to 'bleve'")
	}

	index, err := openOrCreateIndex(cfg.BlevePath)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to open or create bleve index"))
	}

	return &BleveSearcher{
		client:   index,
		hydrator: hydrator,
	}, nil
}

func (s *BleveSearcher) Search(ctx context.Context, q string, p pagination.Parameters, opts searcher.Options) (*pagination.Result[datagraph.Item], error) {
	if s.client == nil {
		return nil, fault.New("bleve client is not initialized")
	}

	textQuery := bleve.NewMatchQuery(q)

	if kinds, ok := opts.Kinds.Get(); ok && len(kinds) > 0 {
		kindQueries := make([]query.Query, 0, len(kinds))
		for _, k := range kinds {
			kq := bleve.NewMatchQuery(k.String())
			kq.SetField("kind")
			kindQueries = append(kindQueries, kq)
		}
		kindQuery := bleve.NewDisjunctionQuery(kindQueries...)
		combinedQuery := bleve.NewConjunctionQuery(textQuery, kindQuery)

		req := bleve.NewSearchRequestOptions(combinedQuery, p.Size(), (p.PageOneIndexed()-1)*p.Size(), false)
		req.Fields = []string{"id", "kind", "name", "slug", "description", "created_at"}
		req.SortBy([]string{"-_score"})

		result, err := s.client.Search(req)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
		return s.processResults(ctx, result, p)
	}

	req := bleve.NewSearchRequestOptions(textQuery, p.Size(), (p.PageOneIndexed()-1)*p.Size(), false)
	req.Fields = []string{"id", "kind", "name", "slug", "description", "created_at"}
	req.SortBy([]string{"-_score"})

	result, err := s.client.Search(req)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return s.processResults(ctx, result, p)
}

func (s *BleveSearcher) processResults(ctx context.Context, result *bleve.SearchResult, p pagination.Parameters) (*pagination.Result[datagraph.Item], error) {
	refs := make([]*datagraph.Ref, 0, len(result.Hits))
	for _, hit := range result.Hits {
		id, err := xid.FromString(hit.ID)
		if err != nil {
			continue
		}
		kindStr, ok := hit.Fields["kind"].(string)
		if !ok {
			continue
		}
		kind, err := datagraph.NewKind(kindStr)
		if err != nil {
			continue
		}
		refs = append(refs, &datagraph.Ref{
			ID:   id,
			Kind: kind,
		})
	}

	items, err := s.hydrator.Hydrate(ctx, refs...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	totalPages := int(result.Total) / p.Size()
	if int(result.Total)%p.Size() > 0 {
		totalPages++
	}

	nextPage := opt.NewEmpty[int]()
	if p.PageOneIndexed() < totalPages {
		nextPage = opt.New(p.PageOneIndexed() + 1)
	}

	return &pagination.Result[datagraph.Item]{
		Size:        p.Size(),
		Results:     int(result.Total),
		TotalPages:  totalPages,
		CurrentPage: p.PageOneIndexed(),
		NextPage:    nextPage,
		Items:       items,
	}, nil
}

func (s *BleveSearcher) Index(ctx context.Context, item datagraph.Item) error {
	doc := Document{
		ID:          item.GetID().String(),
		Kind:        item.GetKind().String(),
		Name:        item.GetName(),
		Slug:        item.GetSlug(),
		Description: item.GetDesc(),
		Content:     item.GetContent().HTML(),
		CreatedAt:   item.GetCreated().Unix(),
	}

	err := s.client.Index(item.GetID().String(), doc)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx), fmsg.With(fmt.Sprintf("failed to index document in bleve %s", item.GetID())))
	}

	return nil
}

func (s *BleveSearcher) Deindex(ctx context.Context, ir datagraph.ItemRef) error {
	err := s.client.Delete(ir.GetID().String())
	if err != nil {
		return fault.Wrap(err, fmsg.With(fmt.Sprintf("failed to delete document from bleve %s", ir.GetID())))
	}
	return nil
}

func openOrCreateIndex(path string) (bleve.Index, error) {
	index, err := bleve.Open(path)
	if err == bleve.ErrorIndexPathDoesNotExist {
		indexMapping := createIndexMapping()
		index, err = bleve.New(path, indexMapping)
		if err != nil {
			return nil, fault.Wrap(err, fmsg.With("failed to create new bleve index"))
		}
		return index, nil
	}
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to open existing bleve index"))
	}

	return index, nil
}

func createIndexMapping() mapping.IndexMapping {
	indexMapping := bleve.NewIndexMapping()

	docMapping := bleve.NewDocumentMapping()

	idFieldMapping := bleve.NewTextFieldMapping()
	idFieldMapping.Store = true
	idFieldMapping.Index = true
	docMapping.AddFieldMappingsAt("id", idFieldMapping)

	kindFieldMapping := bleve.NewTextFieldMapping()
	kindFieldMapping.Store = true
	kindFieldMapping.Index = true
	docMapping.AddFieldMappingsAt("kind", kindFieldMapping)

	nameFieldMapping := bleve.NewTextFieldMapping()
	nameFieldMapping.Store = true
	nameFieldMapping.Index = true
	nameFieldMapping.Analyzer = "en"
	docMapping.AddFieldMappingsAt("name", nameFieldMapping)

	slugFieldMapping := bleve.NewTextFieldMapping()
	slugFieldMapping.Store = true
	slugFieldMapping.Index = true
	docMapping.AddFieldMappingsAt("slug", slugFieldMapping)

	descFieldMapping := bleve.NewTextFieldMapping()
	descFieldMapping.Store = true
	descFieldMapping.Index = true
	descFieldMapping.Analyzer = "en"
	docMapping.AddFieldMappingsAt("description", descFieldMapping)

	contentFieldMapping := bleve.NewTextFieldMapping()
	contentFieldMapping.Store = false
	contentFieldMapping.Index = true
	contentFieldMapping.Analyzer = "en"
	docMapping.AddFieldMappingsAt("content", contentFieldMapping)

	createdAtFieldMapping := bleve.NewNumericFieldMapping()
	createdAtFieldMapping.Store = true
	createdAtFieldMapping.Index = true
	docMapping.AddFieldMappingsAt("created_at", createdAtFieldMapping)

	indexMapping.DefaultMapping = docMapping

	return indexMapping
}
