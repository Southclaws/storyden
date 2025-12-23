package redis_search

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/redis/rueidis"
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

type RedisSearcher struct {
	client    rueidis.Client
	indexName string
	hydrator  *hydrate.Hydrator
}

type Document struct {
	Kind        string
	Name        string
	Slug        string
	Description string
	Content     string
	CreatedAt   int64
	AuthorID    string
	CategoryID  string
	Tags        []string
}

type SearchResult struct {
	Total int
	Hits  []SearchHit
}

type SearchHit struct {
	ID          xid.ID
	Kind        string
	Name        string
	Slug        string
	Description string
	CreatedAt   int64
	Score       float64
}

func New(ctx context.Context, cfg config.Config, redisClient rueidis.Client, hydrator *hydrate.Hydrator) (*RedisSearcher, error) {
	if cfg.SearchProvider != "redis" {
		return nil, nil
	}

	if redisClient == nil {
		return nil, fault.New("REDIS_URL is required when SEARCH_PROVIDER is set to 'redis'")
	}

	indexName := cfg.RedisSearchIndexName
	if indexName == "" {
		indexName = "storyden"
	}

	rs := &RedisSearcher{
		client:    redisClient,
		indexName: indexName,
		hydrator:  hydrator,
	}

	if err := rs.ensureIndex(ctx); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return rs, nil
}

func (c *RedisSearcher) key(item datagraph.ItemRef) string {
	return fmt.Sprintf("%s%s",
		c.prefix(item.GetKind()),
		item.GetID().String())
}

// idFromKey gets the ID from "storyden:datagraph:thread:d3lal65o2dtv213s95o0"
func (c *RedisSearcher) idFromKey(key string) (xid.ID, error) {
	parts := strings.Split(key, ":")
	if len(parts) < 4 {
		return xid.NilID(), errors.New("invalid key format")
	}

	idStr := parts[3]
	id, err := xid.FromString(idStr)
	if err != nil {
		return xid.NilID(), fault.Wrap(err, fmsg.With("failed to parse xid from key"))
	}
	return id, nil
}

func (c *RedisSearcher) prefix(k datagraph.Kind) string {
	return fmt.Sprintf("%s:datagraph:%s:", c.indexName, k.String())
}

func (c *RedisSearcher) prefixes() []string {
	return []string{
		c.prefix(datagraph.KindPost),
		c.prefix(datagraph.KindThread),
		c.prefix(datagraph.KindReply),
		c.prefix(datagraph.KindNode),
		c.prefix(datagraph.KindCollection),
		c.prefix(datagraph.KindProfile),
		c.prefix(datagraph.KindEvent),
	}
}

func (c *RedisSearcher) createIndex(ctx context.Context) error {
	p := c.prefixes()
	cmd := c.client.B().FtCreate().
		Index(c.indexName).
		OnHash().
		Prefix(int64(len(p))).
		Prefix(p...).
		Stopwords(0). // Disable stopwords: Storyden is not just for English.
		Schema().     // This has a very minor memory impact, but it's fine.
		// NOTE: We also disable stemming, again, anglocentric. Can fix later.
		FieldName("kind").Tag().
		FieldName("name").Text().Nostem().
		FieldName("slug").Text().Nostem().
		FieldName("description").Text().Nostem().
		FieldName("content").Text().Nostem().
		FieldName("created_at").Numeric().Sortable().
		FieldName("author_id").Tag().
		FieldName("category_id").Tag().
		FieldName("tags").Tag().Separator(",").
		Build()

	err := c.client.Do(ctx, cmd).Error()
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create search index"))
	}

	return nil
}

func (s *RedisSearcher) Search(ctx context.Context, q string, p pagination.Parameters, opts searcher.Options) (*pagination.Result[datagraph.Item], error) {
	escapedQuery := s.buildQuery(q, opts)

	offset := (p.PageOneIndexed() - 1) * p.Size()
	limit := p.Size()

	cmd := s.client.B().FtSearch().
		Index(s.indexName).
		Query(escapedQuery).
		Withscores().
		Limit().OffsetNum(int64(offset), int64(limit)).
		Build()

	total, docs, err := s.client.Do(ctx, cmd).AsFtSearch()
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to search redis index"))
	}

	hits := make([]SearchHit, 0, len(docs))
	for _, doc := range docs {
		hit := SearchHit{}

		hit.ID, err = s.idFromKey(doc.Key)
		if err != nil {
			// TODO: Log error
			continue
		}
		hit.Score = doc.Score
		if kind, ok := doc.Doc["kind"]; ok {
			hit.Kind = kind
		}
		if name, ok := doc.Doc["name"]; ok {
			hit.Name = name
		}
		if slug, ok := doc.Doc["slug"]; ok {
			hit.Slug = slug
		}
		if desc, ok := doc.Doc["description"]; ok {
			hit.Description = desc
		}
		if created, ok := doc.Doc["created_at"]; ok {
			if ts, err := strconv.ParseInt(created, 10, 64); err == nil {
				hit.CreatedAt = ts
			}
		}

		hits = append(hits, hit)
	}

	refs := make([]*datagraph.Ref, 0, len(hits))
	for _, hit := range hits {

		kind, err := datagraph.NewKind(hit.Kind)
		if err != nil {
			continue
		}
		refs = append(refs, &datagraph.Ref{
			ID:        hit.ID,
			Kind:      kind,
			Relevance: hit.Score,
		})
	}

	items, err := s.hydrator.Hydrate(ctx, refs...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	totalPages := int(total) / p.Size()
	if int(total)%p.Size() > 0 {
		totalPages++
	}

	nextPage := opt.NewEmpty[int]()
	if p.PageOneIndexed() < totalPages {
		nextPage = opt.New(p.PageOneIndexed() + 1)
	}

	return &pagination.Result[datagraph.Item]{
		Size:        p.Size(),
		Results:     int(total),
		TotalPages:  totalPages,
		CurrentPage: p.PageOneIndexed(),
		NextPage:    nextPage,
		Items:       items,
	}, nil
}

func (s *RedisSearcher) MatchFast(ctx context.Context, q string, limit int, opts searcher.Options) (datagraph.MatchList, error) {
	cmd := s.client.B().FtSearch().
		Index(s.indexName).
		Query(s.buildPrefixQuery(q, opts)).
		Limit().OffsetNum(0, int64(limit)).
		Build()

	_, docs, err := s.client.Do(ctx, cmd).AsFtSearch()
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to search redis index"))
	}

	matches := make(datagraph.MatchList, 0, len(docs))
	for _, doc := range docs {
		match, ok := s.matchFromDoc(doc)
		if !ok {
			continue
		}
		matches = append(matches, match)
	}

	return matches, nil
}

func (s *RedisSearcher) buildDocument(item datagraph.Item) Document {
	doc := Document{
		Kind:        item.GetKind().String(),
		Name:        item.GetName(),
		Slug:        item.GetSlug(),
		Description: item.GetDesc(),
		Content:     item.GetContent().Plaintext(),
		CreatedAt:   item.GetCreated().Unix(),
	}

	if v, ok := item.(datagraph.WithAuthor); ok {
		doc.AuthorID = v.GetAuthor().String()
	}

	if v, ok := item.(datagraph.WithCategory); ok {
		doc.CategoryID = v.GetCategory().String()
	}

	if v, ok := item.(datagraph.WithTagNames); ok {
		doc.Tags = v.GetTags()
	}

	return doc
}

func (s *RedisSearcher) Index(ctx context.Context, item datagraph.Item) error {
	doc := s.buildDocument(item)

	key := s.key(item)

	builder := s.client.B().Hset().
		Key(key).
		FieldValue().
		FieldValue("kind", doc.Kind).
		FieldValue("name", doc.Name).
		FieldValue("slug", doc.Slug).
		FieldValue("description", doc.Description).
		FieldValue("content", doc.Content).
		FieldValue("created_at", strconv.FormatInt(doc.CreatedAt, 10))

	if doc.AuthorID != "" {
		builder = builder.FieldValue("author_id", doc.AuthorID)
	}

	if doc.CategoryID != "" {
		builder = builder.FieldValue("category_id", doc.CategoryID)
	}

	if len(doc.Tags) > 0 {
		builder = builder.FieldValue("tags", strings.Join(doc.Tags, ","))
	}

	cmd := builder.Build()

	err := s.client.Do(ctx, cmd).Error()
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx), fmsg.With(fmt.Sprintf("failed to index document in redis %s", item.GetID())))
	}

	return nil
}

func (s *RedisSearcher) buildQuery(q string, opts searcher.Options) string {
	escapedQuery := escapeRedisSearch(q)
	filters := []string{}

	if kinds, ok := opts.Kinds.Get(); ok && len(kinds) > 0 {
		kindStrs := make([]string, len(kinds))
		for i, k := range kinds {
			kindStrs[i] = k.String()
		}
		filters = append(filters, fmt.Sprintf("@kind:{%s}", strings.Join(kindStrs, "|")))
	}

	if authors, ok := opts.Authors.Get(); ok && len(authors) > 0 {
		authorStrs := make([]string, len(authors))
		for i, a := range authors {
			authorStrs[i] = a.String()
		}
		filters = append(filters, fmt.Sprintf("@author_id:{%s}", strings.Join(authorStrs, "|")))
	}

	if categories, ok := opts.Categories.Get(); ok && len(categories) > 0 {
		categoryStrs := make([]string, len(categories))
		for i, c := range categories {
			categoryStrs[i] = c.String()
		}
		filters = append(filters, fmt.Sprintf("@category_id:{%s}", strings.Join(categoryStrs, "|")))
	}

	if tags, ok := opts.Tags.Get(); ok && len(tags) > 0 {
		for _, t := range tags {
			filters = append(filters, fmt.Sprintf("@tags:{%s}", escapeRedisSearch(t.String())))
		}
	}

	if len(filters) > 0 {
		return fmt.Sprintf("(%s) %s", escapedQuery, strings.Join(filters, " "))
	}

	return escapedQuery
}

func (s *RedisSearcher) buildPrefixQuery(q string, opts searcher.Options) string {
	words := strings.Fields(q)
	if len(words) == 0 {
		return "*"
	}

	var terms []string
	for i, w := range words {
		esc := escapeRedisSearch(w)

		isLast := i == len(words)-1

		if isLast {
			// redis ft prefix must be >= 2 chars
			if len([]rune(w)) >= 2 {
				terms = append(terms, esc)
			}
		} else {
			// previous tokens: fixed terms
			terms = append(terms, esc)
		}
	}

	nameQuery := fmt.Sprintf("@name:(%s*)", strings.Join(terms, " "))
	filters := []string{}

	if kinds, ok := opts.Kinds.Get(); ok && len(kinds) > 0 {
		// NOTE: kind = "reply" does not work here (replies have no "name".)
		// TODO: Determine a path forward for this if it's ever useful.
		// Probably not, typeahead is more about resource names. The only
		// downside here is wasteful searching when kinds is empty (all kinds.)
		kindStrs := make([]string, len(kinds))
		for i, k := range kinds {
			kindStrs[i] = k.String()
		}
		filters = append(filters, fmt.Sprintf("@kind:{%s}", strings.Join(kindStrs, "|")))
	}

	if authors, ok := opts.Authors.Get(); ok && len(authors) > 0 {
		authorStrs := make([]string, len(authors))
		for i, a := range authors {
			authorStrs[i] = a.String()
		}
		filters = append(filters, fmt.Sprintf("@author_id:{%s}", strings.Join(authorStrs, "|")))
	}

	if categories, ok := opts.Categories.Get(); ok && len(categories) > 0 {
		categoryStrs := make([]string, len(categories))
		for i, c := range categories {
			categoryStrs[i] = c.String()
		}
		filters = append(filters, fmt.Sprintf("@category_id:{%s}", strings.Join(categoryStrs, "|")))
	}

	if tags, ok := opts.Tags.Get(); ok && len(tags) > 0 {
		for _, t := range tags {
			filters = append(filters, fmt.Sprintf("@tags:{%s}", escapeRedisSearch(t.String())))
		}
	}

	if len(filters) > 0 {
		return fmt.Sprintf("%s %s", nameQuery, strings.Join(filters, " "))
	}

	return nameQuery
}

func (s *RedisSearcher) matchFromDoc(doc rueidis.FtSearchDoc) (datagraph.Match, bool) {
	id, err := s.idFromKey(doc.Key)
	if err != nil {
		return datagraph.Match{}, false
	}

	kind, err := datagraph.NewKind(doc.Doc["kind"])
	if err != nil {
		return datagraph.Match{}, false
	}

	return datagraph.Match{
		ID:          id,
		Kind:        kind,
		Slug:        doc.Doc["slug"],
		Name:        doc.Doc["name"],
		Description: doc.Doc["description"],
	}, true
}

func (s *RedisSearcher) Deindex(ctx context.Context, ir datagraph.ItemRef) error {
	key := s.key(ir)

	cmd := s.client.B().Del().Key(key).Build()
	err := s.client.Do(ctx, cmd).Error()
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx), fmsg.With(fmt.Sprintf("failed to delete document from redis %s", ir.GetID())))
	}

	return nil
}

func (c *RedisSearcher) ensureIndex(ctx context.Context) error {
	cmd := c.client.B().FtInfo().Index(c.indexName).Build()
	err := c.client.Do(ctx, cmd).Error()
	if err != nil {
		re := &rueidis.RedisError{}
		if errors.As(err, &re) {
			// NOTE: Sketchy way to check for index existence, but rueidis
			// doesn't expose proper sentinel errors for some reason.
			if re.Error() == "Unknown index name" {
				return c.createIndex(ctx)
			}
		}

		return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to check if index exists"))
	}

	return nil
}

func escapeRedisSearch(s string) string {
	var out strings.Builder
	out.Grow(len(s) * 2) // worst case

	// List of operator chars from RedisSearch grammar
	special := map[rune]bool{
		'+': true, '-': true, '=': true, '&': true, '|': true,
		'>': true, '<': true, '!': true, '(': true, ')': true,
		'{': true, '}': true, '[': true, ']': true, '^': true,
		'"': true, '~': true, '*': true, '?': true, ':': true,
		'\\': true,
	}

	for _, r := range s {
		if special[r] {
			out.WriteRune('\\')
		}
		out.WriteRune(r)
	}
	return out.String()
}
