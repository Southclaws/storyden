// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/Southclaws/storyden/internal/ent/collection"
	"github.com/Southclaws/storyden/internal/ent/collectionpost"
	"github.com/Southclaws/storyden/internal/ent/post"
	"github.com/Southclaws/storyden/internal/ent/predicate"
	"github.com/rs/xid"
)

// CollectionPostQuery is the builder for querying CollectionPost entities.
type CollectionPostQuery struct {
	config
	ctx            *QueryContext
	order          []collectionpost.OrderOption
	inters         []Interceptor
	predicates     []predicate.CollectionPost
	withCollection *CollectionQuery
	withPost       *PostQuery
	modifiers      []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the CollectionPostQuery builder.
func (cpq *CollectionPostQuery) Where(ps ...predicate.CollectionPost) *CollectionPostQuery {
	cpq.predicates = append(cpq.predicates, ps...)
	return cpq
}

// Limit the number of records to be returned by this query.
func (cpq *CollectionPostQuery) Limit(limit int) *CollectionPostQuery {
	cpq.ctx.Limit = &limit
	return cpq
}

// Offset to start from.
func (cpq *CollectionPostQuery) Offset(offset int) *CollectionPostQuery {
	cpq.ctx.Offset = &offset
	return cpq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (cpq *CollectionPostQuery) Unique(unique bool) *CollectionPostQuery {
	cpq.ctx.Unique = &unique
	return cpq
}

// Order specifies how the records should be ordered.
func (cpq *CollectionPostQuery) Order(o ...collectionpost.OrderOption) *CollectionPostQuery {
	cpq.order = append(cpq.order, o...)
	return cpq
}

// QueryCollection chains the current query on the "collection" edge.
func (cpq *CollectionPostQuery) QueryCollection() *CollectionQuery {
	query := (&CollectionClient{config: cpq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := cpq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := cpq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(collectionpost.Table, collectionpost.CollectionColumn, selector),
			sqlgraph.To(collection.Table, collection.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, collectionpost.CollectionTable, collectionpost.CollectionColumn),
		)
		fromU = sqlgraph.SetNeighbors(cpq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryPost chains the current query on the "post" edge.
func (cpq *CollectionPostQuery) QueryPost() *PostQuery {
	query := (&PostClient{config: cpq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := cpq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := cpq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(collectionpost.Table, collectionpost.PostColumn, selector),
			sqlgraph.To(post.Table, post.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, collectionpost.PostTable, collectionpost.PostColumn),
		)
		fromU = sqlgraph.SetNeighbors(cpq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first CollectionPost entity from the query.
// Returns a *NotFoundError when no CollectionPost was found.
func (cpq *CollectionPostQuery) First(ctx context.Context) (*CollectionPost, error) {
	nodes, err := cpq.Limit(1).All(setContextOp(ctx, cpq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{collectionpost.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (cpq *CollectionPostQuery) FirstX(ctx context.Context) *CollectionPost {
	node, err := cpq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// Only returns a single CollectionPost entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one CollectionPost entity is found.
// Returns a *NotFoundError when no CollectionPost entities are found.
func (cpq *CollectionPostQuery) Only(ctx context.Context) (*CollectionPost, error) {
	nodes, err := cpq.Limit(2).All(setContextOp(ctx, cpq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{collectionpost.Label}
	default:
		return nil, &NotSingularError{collectionpost.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (cpq *CollectionPostQuery) OnlyX(ctx context.Context) *CollectionPost {
	node, err := cpq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// All executes the query and returns a list of CollectionPosts.
func (cpq *CollectionPostQuery) All(ctx context.Context) ([]*CollectionPost, error) {
	ctx = setContextOp(ctx, cpq.ctx, "All")
	if err := cpq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*CollectionPost, *CollectionPostQuery]()
	return withInterceptors[[]*CollectionPost](ctx, cpq, qr, cpq.inters)
}

// AllX is like All, but panics if an error occurs.
func (cpq *CollectionPostQuery) AllX(ctx context.Context) []*CollectionPost {
	nodes, err := cpq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// Count returns the count of the given query.
func (cpq *CollectionPostQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, cpq.ctx, "Count")
	if err := cpq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, cpq, querierCount[*CollectionPostQuery](), cpq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (cpq *CollectionPostQuery) CountX(ctx context.Context) int {
	count, err := cpq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (cpq *CollectionPostQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, cpq.ctx, "Exist")
	switch _, err := cpq.First(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (cpq *CollectionPostQuery) ExistX(ctx context.Context) bool {
	exist, err := cpq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the CollectionPostQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (cpq *CollectionPostQuery) Clone() *CollectionPostQuery {
	if cpq == nil {
		return nil
	}
	return &CollectionPostQuery{
		config:         cpq.config,
		ctx:            cpq.ctx.Clone(),
		order:          append([]collectionpost.OrderOption{}, cpq.order...),
		inters:         append([]Interceptor{}, cpq.inters...),
		predicates:     append([]predicate.CollectionPost{}, cpq.predicates...),
		withCollection: cpq.withCollection.Clone(),
		withPost:       cpq.withPost.Clone(),
		// clone intermediate query.
		sql:  cpq.sql.Clone(),
		path: cpq.path,
	}
}

// WithCollection tells the query-builder to eager-load the nodes that are connected to
// the "collection" edge. The optional arguments are used to configure the query builder of the edge.
func (cpq *CollectionPostQuery) WithCollection(opts ...func(*CollectionQuery)) *CollectionPostQuery {
	query := (&CollectionClient{config: cpq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	cpq.withCollection = query
	return cpq
}

// WithPost tells the query-builder to eager-load the nodes that are connected to
// the "post" edge. The optional arguments are used to configure the query builder of the edge.
func (cpq *CollectionPostQuery) WithPost(opts ...func(*PostQuery)) *CollectionPostQuery {
	query := (&PostClient{config: cpq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	cpq.withPost = query
	return cpq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		CreatedAt time.Time `json:"created_at,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.CollectionPost.Query().
//		GroupBy(collectionpost.FieldCreatedAt).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (cpq *CollectionPostQuery) GroupBy(field string, fields ...string) *CollectionPostGroupBy {
	cpq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &CollectionPostGroupBy{build: cpq}
	grbuild.flds = &cpq.ctx.Fields
	grbuild.label = collectionpost.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		CreatedAt time.Time `json:"created_at,omitempty"`
//	}
//
//	client.CollectionPost.Query().
//		Select(collectionpost.FieldCreatedAt).
//		Scan(ctx, &v)
func (cpq *CollectionPostQuery) Select(fields ...string) *CollectionPostSelect {
	cpq.ctx.Fields = append(cpq.ctx.Fields, fields...)
	sbuild := &CollectionPostSelect{CollectionPostQuery: cpq}
	sbuild.label = collectionpost.Label
	sbuild.flds, sbuild.scan = &cpq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a CollectionPostSelect configured with the given aggregations.
func (cpq *CollectionPostQuery) Aggregate(fns ...AggregateFunc) *CollectionPostSelect {
	return cpq.Select().Aggregate(fns...)
}

func (cpq *CollectionPostQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range cpq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, cpq); err != nil {
				return err
			}
		}
	}
	for _, f := range cpq.ctx.Fields {
		if !collectionpost.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if cpq.path != nil {
		prev, err := cpq.path(ctx)
		if err != nil {
			return err
		}
		cpq.sql = prev
	}
	return nil
}

func (cpq *CollectionPostQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*CollectionPost, error) {
	var (
		nodes       = []*CollectionPost{}
		_spec       = cpq.querySpec()
		loadedTypes = [2]bool{
			cpq.withCollection != nil,
			cpq.withPost != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*CollectionPost).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &CollectionPost{config: cpq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if len(cpq.modifiers) > 0 {
		_spec.Modifiers = cpq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, cpq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := cpq.withCollection; query != nil {
		if err := cpq.loadCollection(ctx, query, nodes, nil,
			func(n *CollectionPost, e *Collection) { n.Edges.Collection = e }); err != nil {
			return nil, err
		}
	}
	if query := cpq.withPost; query != nil {
		if err := cpq.loadPost(ctx, query, nodes, nil,
			func(n *CollectionPost, e *Post) { n.Edges.Post = e }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (cpq *CollectionPostQuery) loadCollection(ctx context.Context, query *CollectionQuery, nodes []*CollectionPost, init func(*CollectionPost), assign func(*CollectionPost, *Collection)) error {
	ids := make([]xid.ID, 0, len(nodes))
	nodeids := make(map[xid.ID][]*CollectionPost)
	for i := range nodes {
		fk := nodes[i].CollectionID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(collection.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "collection_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (cpq *CollectionPostQuery) loadPost(ctx context.Context, query *PostQuery, nodes []*CollectionPost, init func(*CollectionPost), assign func(*CollectionPost, *Post)) error {
	ids := make([]xid.ID, 0, len(nodes))
	nodeids := make(map[xid.ID][]*CollectionPost)
	for i := range nodes {
		fk := nodes[i].PostID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(post.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "post_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}

func (cpq *CollectionPostQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := cpq.querySpec()
	if len(cpq.modifiers) > 0 {
		_spec.Modifiers = cpq.modifiers
	}
	_spec.Unique = false
	_spec.Node.Columns = nil
	return sqlgraph.CountNodes(ctx, cpq.driver, _spec)
}

func (cpq *CollectionPostQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(collectionpost.Table, collectionpost.Columns, nil)
	_spec.From = cpq.sql
	if unique := cpq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if cpq.path != nil {
		_spec.Unique = true
	}
	if fields := cpq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		for i := range fields {
			_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
		}
		if cpq.withCollection != nil {
			_spec.Node.AddColumnOnce(collectionpost.FieldCollectionID)
		}
		if cpq.withPost != nil {
			_spec.Node.AddColumnOnce(collectionpost.FieldPostID)
		}
	}
	if ps := cpq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := cpq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := cpq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := cpq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (cpq *CollectionPostQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(cpq.driver.Dialect())
	t1 := builder.Table(collectionpost.Table)
	columns := cpq.ctx.Fields
	if len(columns) == 0 {
		columns = collectionpost.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if cpq.sql != nil {
		selector = cpq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if cpq.ctx.Unique != nil && *cpq.ctx.Unique {
		selector.Distinct()
	}
	for _, m := range cpq.modifiers {
		m(selector)
	}
	for _, p := range cpq.predicates {
		p(selector)
	}
	for _, p := range cpq.order {
		p(selector)
	}
	if offset := cpq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := cpq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// Modify adds a query modifier for attaching custom logic to queries.
func (cpq *CollectionPostQuery) Modify(modifiers ...func(s *sql.Selector)) *CollectionPostSelect {
	cpq.modifiers = append(cpq.modifiers, modifiers...)
	return cpq.Select()
}

// CollectionPostGroupBy is the group-by builder for CollectionPost entities.
type CollectionPostGroupBy struct {
	selector
	build *CollectionPostQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (cpgb *CollectionPostGroupBy) Aggregate(fns ...AggregateFunc) *CollectionPostGroupBy {
	cpgb.fns = append(cpgb.fns, fns...)
	return cpgb
}

// Scan applies the selector query and scans the result into the given value.
func (cpgb *CollectionPostGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, cpgb.build.ctx, "GroupBy")
	if err := cpgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*CollectionPostQuery, *CollectionPostGroupBy](ctx, cpgb.build, cpgb, cpgb.build.inters, v)
}

func (cpgb *CollectionPostGroupBy) sqlScan(ctx context.Context, root *CollectionPostQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(cpgb.fns))
	for _, fn := range cpgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*cpgb.flds)+len(cpgb.fns))
		for _, f := range *cpgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*cpgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := cpgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// CollectionPostSelect is the builder for selecting fields of CollectionPost entities.
type CollectionPostSelect struct {
	*CollectionPostQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (cps *CollectionPostSelect) Aggregate(fns ...AggregateFunc) *CollectionPostSelect {
	cps.fns = append(cps.fns, fns...)
	return cps
}

// Scan applies the selector query and scans the result into the given value.
func (cps *CollectionPostSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, cps.ctx, "Select")
	if err := cps.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*CollectionPostQuery, *CollectionPostSelect](ctx, cps.CollectionPostQuery, cps, cps.inters, v)
}

func (cps *CollectionPostSelect) sqlScan(ctx context.Context, root *CollectionPostQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(cps.fns))
	for _, fn := range cps.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*cps.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := cps.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (cps *CollectionPostSelect) Modify(modifiers ...func(s *sql.Selector)) *CollectionPostSelect {
	cps.modifiers = append(cps.modifiers, modifiers...)
	return cps
}
