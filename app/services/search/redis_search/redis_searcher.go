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
	ID          string
	Kind        string
	Name        string
	Slug        string
	Description string
	Content     string
	CreatedAt   int64
}

type SearchResult struct {
	Total int
	Hits  []SearchHit
}

type SearchHit struct {
	ID          string
	Kind        string
	Name        string
	Slug        string
	Description string
	CreatedAt   int64
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
		Schema().
		FieldName("id").Text().Sortable().
		FieldName("kind").Tag().Sortable().
		FieldName("name").Text().Sortable().
		FieldName("slug").Text().Sortable().
		FieldName("description").Text().
		FieldName("content").Text().
		FieldName("created_at").Numeric().Sortable().
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

		if id, ok := doc.Doc["id"]; ok {
			hit.ID = id
		}
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
		id, err := xid.FromString(hit.ID)
		if err != nil {
			continue
		}
		kind, err := datagraph.NewKind(hit.Kind)
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
	if s.client == nil {
		return nil, searcher.ErrFastMatchesUnavailable
	}

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
		match, ok := s.matchFromDoc(doc.Doc)
		if !ok {
			continue
		}
		matches = append(matches, match)
	}

	return matches, nil
}

func (s *RedisSearcher) Index(ctx context.Context, item datagraph.Item) error {
	doc := Document{
		ID:          item.GetID().String(),
		Kind:        item.GetKind().String(),
		Name:        item.GetName(),
		Slug:        item.GetSlug(),
		Description: item.GetDesc(),
		Content:     item.GetContent().HTML(),
		CreatedAt:   item.GetCreated().Unix(),
	}

	key := s.key(item)

	cmd := s.client.B().Hset().
		Key(key).
		FieldValue().
		FieldValue("id", doc.ID).
		FieldValue("kind", doc.Kind).
		FieldValue("name", doc.Name).
		FieldValue("slug", doc.Slug).
		FieldValue("description", doc.Description).
		FieldValue("content", doc.Content).
		FieldValue("created_at", strconv.FormatInt(doc.CreatedAt, 10)).
		Build()

	err := s.client.Do(ctx, cmd).Error()
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx), fmsg.With(fmt.Sprintf("failed to index document in redis %s", item.GetID())))
	}

	return nil
}

func (s *RedisSearcher) buildQuery(q string, opts searcher.Options) string {
	escapedQuery := escapeRedisSearch(q)

	if kinds, ok := opts.Kinds.Get(); ok && len(kinds) > 0 {
		kindStrs := make([]string, len(kinds))
		for i, k := range kinds {
			kindStrs[i] = k.String()
		}
		kindFilter := fmt.Sprintf("@kind:{%s}", strings.Join(kindStrs, "|"))
		return fmt.Sprintf("(%s) %s", escapedQuery, kindFilter)
	}

	return escapedQuery
}

func (s *RedisSearcher) buildPrefixQuery(q string, opts searcher.Options) string {
	words := strings.Fields(q)
	if len(words) == 0 {
		return "*"
	}

	prefixTerms := make([]string, len(words))
	for i, word := range words {
		escaped := escapeRedisSearch(word)
		prefixTerms[i] = escaped + "*"
	}

	nameQuery := fmt.Sprintf("@name:(%s)", strings.Join(prefixTerms, " "))

	if kinds, ok := opts.Kinds.Get(); ok && len(kinds) > 0 {
		kindStrs := make([]string, len(kinds))
		for i, k := range kinds {
			kindStrs[i] = k.String()
		}
		kindFilter := fmt.Sprintf("@kind:{%s}", strings.Join(kindStrs, "|"))
		return fmt.Sprintf("%s %s", nameQuery, kindFilter)
	}

	return nameQuery
}

func (s *RedisSearcher) matchFromDoc(fields map[string]string) (datagraph.Match, bool) {
	id, err := xid.FromString(fields["id"])
	if err != nil {
		return datagraph.Match{}, false
	}

	kind, err := datagraph.NewKind(fields["kind"])
	if err != nil {
		return datagraph.Match{}, false
	}

	return datagraph.Match{
		ID:          id,
		Kind:        kind,
		Slug:        fields["slug"],
		Name:        fields["name"],
		Description: fields["description"],
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
			fmt.Println(re.Error())
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
