package bleve_search

import (
	"context"
	"fmt"
	"strings"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis"
	"github.com/blevesearch/bleve/v2/analysis/token/lowercase"
	"github.com/blevesearch/bleve/v2/analysis/tokenizer/unicode"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/blevesearch/bleve/v2/registry"
	"github.com/blevesearch/bleve/v2/search"
	"github.com/blevesearch/bleve/v2/search/query"
	"github.com/rs/xid"

	"github.com/Southclaws/dt"
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
	ID          string `json:"id"`
	Kind        string `json:"kind"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Content     string `json:"content"`
	CreatedAt   int64  `json:"created_at"`
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
	searchQuery := s.buildSearchQuery(q, opts)

	req := bleve.NewSearchRequestOptions(searchQuery, p.Size(), (p.PageOneIndexed()-1)*p.Size(), false)
	req.Fields = []string{"id", "kind", "name", "slug", "description", "created_at"}
	req.SortBy([]string{"-_score"})

	result, err := s.client.Search(req)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return s.processResults(ctx, result, p)
}

func (s *BleveSearcher) MatchFast(ctx context.Context, q string, limit int, opts searcher.Options) (datagraph.MatchList, error) {
	matchQuery := s.buildMatchQuery(q, opts)

	req := bleve.NewSearchRequestOptions(matchQuery, limit, 0, false)
	req.Fields = []string{"id", "kind", "name", "slug", "description", "created_at"}
	req.SortBy([]string{"_score"})

	result, err := s.client.Search(req)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return dt.Reduce(result.Hits, func(acc datagraph.MatchList, hit *search.DocumentMatch) datagraph.MatchList {
		match, ok := s.matchFromHit(hit)
		if !ok {
			return acc
		}
		return append(acc, match)
	}, datagraph.MatchList{}), nil
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

func (s *BleveSearcher) buildSearchQuery(q string, opts searcher.Options) query.Query {
	textQuery := bleve.NewMatchQuery(q)

	if kinds, ok := opts.Kinds.Get(); ok && len(kinds) > 0 {
		kindQueries := make([]query.Query, 0, len(kinds))
		for _, k := range kinds {
			kq := bleve.NewMatchQuery(k.String())
			kq.SetField("kind")
			kindQueries = append(kindQueries, kq)
		}
		kindQuery := bleve.NewDisjunctionQuery(kindQueries...)
		return bleve.NewConjunctionQuery(textQuery, kindQuery)
	}

	return textQuery
}

func (s *BleveSearcher) buildMatchQuery(q string, opts searcher.Options) query.Query {
	lowercaseQ := strings.ToLower(q)
	tokens := strings.Fields(lowercaseQ)

	textQuery := bleve.NewBooleanQuery()

	for _, tok := range tokens {
		pq := bleve.NewPrefixQuery(tok)
		pq.SetField("name")
		textQuery.AddShould(pq)
	}

	if kinds, ok := opts.Kinds.Get(); ok && len(kinds) > 0 {
		kindQueries := make([]query.Query, 0, len(kinds))
		for _, k := range kinds {
			kq := bleve.NewMatchQuery(k.String())
			kq.SetField("kind")
			kindQueries = append(kindQueries, kq)
		}
		kindQuery := bleve.NewDisjunctionQuery(kindQueries...)
		return bleve.NewConjunctionQuery(textQuery, kindQuery)
	}

	return textQuery
}

func (s *BleveSearcher) matchFromHit(hit *search.DocumentMatch) (datagraph.Match, bool) {
	id, err := xid.FromString(hit.ID)
	if err != nil {
		return datagraph.Match{}, false
	}

	kindStr, ok := hit.Fields["kind"].(string)
	if !ok {
		return datagraph.Match{}, false
	}

	kind, err := datagraph.NewKind(kindStr)
	if err != nil {
		return datagraph.Match{}, false
	}

	return datagraph.Match{
		ID:          id,
		Kind:        kind,
		Slug:        valueAsString(hit.Fields, "slug"),
		Name:        valueAsString(hit.Fields, "name"),
		Description: valueAsString(hit.Fields, "description"),
	}, true
}

func valueAsString(fields map[string]any, key string) string {
	if v, ok := fields[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}

	return ""
}

func (s *BleveSearcher) Index(ctx context.Context, item datagraph.Item) error {
	doc := Document{
		ID:          item.GetID().String(),
		Kind:        item.GetKind().String(),
		Name:        item.GetName(),
		Slug:        item.GetSlug(),
		Description: item.GetDesc(),
		Content:     item.GetContent().Plaintext(), // We index plaintext only.
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
	err := registry.RegisterAnalyzer("intl", InternationalAnalyser)
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to register international analyzer"))
	}

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

	indexMapping.DefaultAnalyzer = "intl"

	docMapping := bleve.NewDocumentMapping()

	idFieldMapping := bleve.NewTextFieldMapping()
	idFieldMapping.Store = true
	idFieldMapping.Index = false
	idFieldMapping.Analyzer = "keyword"
	docMapping.AddFieldMappingsAt("id", idFieldMapping)

	kindFieldMapping := bleve.NewTextFieldMapping()
	kindFieldMapping.Store = true
	kindFieldMapping.Index = true
	kindFieldMapping.Analyzer = "keyword"
	docMapping.AddFieldMappingsAt("kind", kindFieldMapping)

	nameFieldMapping := bleve.NewTextFieldMapping()
	nameFieldMapping.Store = true
	nameFieldMapping.Index = true
	nameFieldMapping.Analyzer = "intl"
	docMapping.AddFieldMappingsAt("name", nameFieldMapping)

	slugFieldMapping := bleve.NewTextFieldMapping()
	slugFieldMapping.Store = true
	slugFieldMapping.Index = true
	slugFieldMapping.Analyzer = "intl"
	docMapping.AddFieldMappingsAt("slug", slugFieldMapping)

	descFieldMapping := bleve.NewTextFieldMapping()
	descFieldMapping.Store = true
	descFieldMapping.Index = true
	descFieldMapping.Analyzer = "intl"
	docMapping.AddFieldMappingsAt("description", descFieldMapping)

	contentFieldMapping := bleve.NewTextFieldMapping()
	contentFieldMapping.Store = false
	contentFieldMapping.Index = true
	contentFieldMapping.Analyzer = "intl"
	docMapping.AddFieldMappingsAt("content", contentFieldMapping)

	createdAtFieldMapping := bleve.NewNumericFieldMapping()
	createdAtFieldMapping.Store = true
	createdAtFieldMapping.Index = true
	docMapping.AddFieldMappingsAt("created_at", createdAtFieldMapping)

	indexMapping.DefaultMapping = docMapping

	return indexMapping
}

// InternationalAnalyser is a copy of the "standard" analyser but without the
// use of English stopwords. This makes it suitable for any language.
func InternationalAnalyser(config map[string]interface{}, cache *registry.Cache) (analysis.Analyzer, error) {
	tokenizer, err := cache.TokenizerNamed(unicode.Name)
	if err != nil {
		return nil, err
	}
	toLowerFilter, err := cache.TokenFilterNamed(lowercase.Name)
	if err != nil {
		return nil, err
	}

	rv := analysis.DefaultAnalyzer{
		Tokenizer: tokenizer,
		TokenFilters: []analysis.TokenFilter{
			toLowerFilter,
		},
	}
	return &rv, nil
}
