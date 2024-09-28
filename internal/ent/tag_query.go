// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"database/sql/driver"
	"fmt"
	"math"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/Southclaws/storyden/internal/ent/account"
	"github.com/Southclaws/storyden/internal/ent/node"
	"github.com/Southclaws/storyden/internal/ent/post"
	"github.com/Southclaws/storyden/internal/ent/predicate"
	"github.com/Southclaws/storyden/internal/ent/tag"
	"github.com/rs/xid"
)

// TagQuery is the builder for querying Tag entities.
type TagQuery struct {
	config
	ctx          *QueryContext
	order        []tag.OrderOption
	inters       []Interceptor
	predicates   []predicate.Tag
	withPosts    *PostQuery
	withNodes    *NodeQuery
	withAccounts *AccountQuery
	modifiers    []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the TagQuery builder.
func (tq *TagQuery) Where(ps ...predicate.Tag) *TagQuery {
	tq.predicates = append(tq.predicates, ps...)
	return tq
}

// Limit the number of records to be returned by this query.
func (tq *TagQuery) Limit(limit int) *TagQuery {
	tq.ctx.Limit = &limit
	return tq
}

// Offset to start from.
func (tq *TagQuery) Offset(offset int) *TagQuery {
	tq.ctx.Offset = &offset
	return tq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (tq *TagQuery) Unique(unique bool) *TagQuery {
	tq.ctx.Unique = &unique
	return tq
}

// Order specifies how the records should be ordered.
func (tq *TagQuery) Order(o ...tag.OrderOption) *TagQuery {
	tq.order = append(tq.order, o...)
	return tq
}

// QueryPosts chains the current query on the "posts" edge.
func (tq *TagQuery) QueryPosts() *PostQuery {
	query := (&PostClient{config: tq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := tq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := tq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(tag.Table, tag.FieldID, selector),
			sqlgraph.To(post.Table, post.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, tag.PostsTable, tag.PostsPrimaryKey...),
		)
		fromU = sqlgraph.SetNeighbors(tq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryNodes chains the current query on the "nodes" edge.
func (tq *TagQuery) QueryNodes() *NodeQuery {
	query := (&NodeClient{config: tq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := tq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := tq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(tag.Table, tag.FieldID, selector),
			sqlgraph.To(node.Table, node.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, tag.NodesTable, tag.NodesPrimaryKey...),
		)
		fromU = sqlgraph.SetNeighbors(tq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryAccounts chains the current query on the "accounts" edge.
func (tq *TagQuery) QueryAccounts() *AccountQuery {
	query := (&AccountClient{config: tq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := tq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := tq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(tag.Table, tag.FieldID, selector),
			sqlgraph.To(account.Table, account.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, tag.AccountsTable, tag.AccountsPrimaryKey...),
		)
		fromU = sqlgraph.SetNeighbors(tq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Tag entity from the query.
// Returns a *NotFoundError when no Tag was found.
func (tq *TagQuery) First(ctx context.Context) (*Tag, error) {
	nodes, err := tq.Limit(1).All(setContextOp(ctx, tq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{tag.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (tq *TagQuery) FirstX(ctx context.Context) *Tag {
	node, err := tq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Tag ID from the query.
// Returns a *NotFoundError when no Tag ID was found.
func (tq *TagQuery) FirstID(ctx context.Context) (id xid.ID, err error) {
	var ids []xid.ID
	if ids, err = tq.Limit(1).IDs(setContextOp(ctx, tq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{tag.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (tq *TagQuery) FirstIDX(ctx context.Context) xid.ID {
	id, err := tq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Tag entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Tag entity is found.
// Returns a *NotFoundError when no Tag entities are found.
func (tq *TagQuery) Only(ctx context.Context) (*Tag, error) {
	nodes, err := tq.Limit(2).All(setContextOp(ctx, tq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{tag.Label}
	default:
		return nil, &NotSingularError{tag.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (tq *TagQuery) OnlyX(ctx context.Context) *Tag {
	node, err := tq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Tag ID in the query.
// Returns a *NotSingularError when more than one Tag ID is found.
// Returns a *NotFoundError when no entities are found.
func (tq *TagQuery) OnlyID(ctx context.Context) (id xid.ID, err error) {
	var ids []xid.ID
	if ids, err = tq.Limit(2).IDs(setContextOp(ctx, tq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{tag.Label}
	default:
		err = &NotSingularError{tag.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (tq *TagQuery) OnlyIDX(ctx context.Context) xid.ID {
	id, err := tq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Tags.
func (tq *TagQuery) All(ctx context.Context) ([]*Tag, error) {
	ctx = setContextOp(ctx, tq.ctx, ent.OpQueryAll)
	if err := tq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*Tag, *TagQuery]()
	return withInterceptors[[]*Tag](ctx, tq, qr, tq.inters)
}

// AllX is like All, but panics if an error occurs.
func (tq *TagQuery) AllX(ctx context.Context) []*Tag {
	nodes, err := tq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Tag IDs.
func (tq *TagQuery) IDs(ctx context.Context) (ids []xid.ID, err error) {
	if tq.ctx.Unique == nil && tq.path != nil {
		tq.Unique(true)
	}
	ctx = setContextOp(ctx, tq.ctx, ent.OpQueryIDs)
	if err = tq.Select(tag.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (tq *TagQuery) IDsX(ctx context.Context) []xid.ID {
	ids, err := tq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (tq *TagQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, tq.ctx, ent.OpQueryCount)
	if err := tq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, tq, querierCount[*TagQuery](), tq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (tq *TagQuery) CountX(ctx context.Context) int {
	count, err := tq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (tq *TagQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, tq.ctx, ent.OpQueryExist)
	switch _, err := tq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (tq *TagQuery) ExistX(ctx context.Context) bool {
	exist, err := tq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the TagQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (tq *TagQuery) Clone() *TagQuery {
	if tq == nil {
		return nil
	}
	return &TagQuery{
		config:       tq.config,
		ctx:          tq.ctx.Clone(),
		order:        append([]tag.OrderOption{}, tq.order...),
		inters:       append([]Interceptor{}, tq.inters...),
		predicates:   append([]predicate.Tag{}, tq.predicates...),
		withPosts:    tq.withPosts.Clone(),
		withNodes:    tq.withNodes.Clone(),
		withAccounts: tq.withAccounts.Clone(),
		// clone intermediate query.
		sql:       tq.sql.Clone(),
		path:      tq.path,
		modifiers: append([]func(*sql.Selector){}, tq.modifiers...),
	}
}

// WithPosts tells the query-builder to eager-load the nodes that are connected to
// the "posts" edge. The optional arguments are used to configure the query builder of the edge.
func (tq *TagQuery) WithPosts(opts ...func(*PostQuery)) *TagQuery {
	query := (&PostClient{config: tq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	tq.withPosts = query
	return tq
}

// WithNodes tells the query-builder to eager-load the nodes that are connected to
// the "nodes" edge. The optional arguments are used to configure the query builder of the edge.
func (tq *TagQuery) WithNodes(opts ...func(*NodeQuery)) *TagQuery {
	query := (&NodeClient{config: tq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	tq.withNodes = query
	return tq
}

// WithAccounts tells the query-builder to eager-load the nodes that are connected to
// the "accounts" edge. The optional arguments are used to configure the query builder of the edge.
func (tq *TagQuery) WithAccounts(opts ...func(*AccountQuery)) *TagQuery {
	query := (&AccountClient{config: tq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	tq.withAccounts = query
	return tq
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
//	client.Tag.Query().
//		GroupBy(tag.FieldCreatedAt).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (tq *TagQuery) GroupBy(field string, fields ...string) *TagGroupBy {
	tq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &TagGroupBy{build: tq}
	grbuild.flds = &tq.ctx.Fields
	grbuild.label = tag.Label
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
//	client.Tag.Query().
//		Select(tag.FieldCreatedAt).
//		Scan(ctx, &v)
func (tq *TagQuery) Select(fields ...string) *TagSelect {
	tq.ctx.Fields = append(tq.ctx.Fields, fields...)
	sbuild := &TagSelect{TagQuery: tq}
	sbuild.label = tag.Label
	sbuild.flds, sbuild.scan = &tq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a TagSelect configured with the given aggregations.
func (tq *TagQuery) Aggregate(fns ...AggregateFunc) *TagSelect {
	return tq.Select().Aggregate(fns...)
}

func (tq *TagQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range tq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, tq); err != nil {
				return err
			}
		}
	}
	for _, f := range tq.ctx.Fields {
		if !tag.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if tq.path != nil {
		prev, err := tq.path(ctx)
		if err != nil {
			return err
		}
		tq.sql = prev
	}
	return nil
}

func (tq *TagQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Tag, error) {
	var (
		nodes       = []*Tag{}
		_spec       = tq.querySpec()
		loadedTypes = [3]bool{
			tq.withPosts != nil,
			tq.withNodes != nil,
			tq.withAccounts != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Tag).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Tag{config: tq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if len(tq.modifiers) > 0 {
		_spec.Modifiers = tq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, tq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := tq.withPosts; query != nil {
		if err := tq.loadPosts(ctx, query, nodes,
			func(n *Tag) { n.Edges.Posts = []*Post{} },
			func(n *Tag, e *Post) { n.Edges.Posts = append(n.Edges.Posts, e) }); err != nil {
			return nil, err
		}
	}
	if query := tq.withNodes; query != nil {
		if err := tq.loadNodes(ctx, query, nodes,
			func(n *Tag) { n.Edges.Nodes = []*Node{} },
			func(n *Tag, e *Node) { n.Edges.Nodes = append(n.Edges.Nodes, e) }); err != nil {
			return nil, err
		}
	}
	if query := tq.withAccounts; query != nil {
		if err := tq.loadAccounts(ctx, query, nodes,
			func(n *Tag) { n.Edges.Accounts = []*Account{} },
			func(n *Tag, e *Account) { n.Edges.Accounts = append(n.Edges.Accounts, e) }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (tq *TagQuery) loadPosts(ctx context.Context, query *PostQuery, nodes []*Tag, init func(*Tag), assign func(*Tag, *Post)) error {
	edgeIDs := make([]driver.Value, len(nodes))
	byID := make(map[xid.ID]*Tag)
	nids := make(map[xid.ID]map[*Tag]struct{})
	for i, node := range nodes {
		edgeIDs[i] = node.ID
		byID[node.ID] = node
		if init != nil {
			init(node)
		}
	}
	query.Where(func(s *sql.Selector) {
		joinT := sql.Table(tag.PostsTable)
		s.Join(joinT).On(s.C(post.FieldID), joinT.C(tag.PostsPrimaryKey[1]))
		s.Where(sql.InValues(joinT.C(tag.PostsPrimaryKey[0]), edgeIDs...))
		columns := s.SelectedColumns()
		s.Select(joinT.C(tag.PostsPrimaryKey[0]))
		s.AppendSelect(columns...)
		s.SetDistinct(false)
	})
	if err := query.prepareQuery(ctx); err != nil {
		return err
	}
	qr := QuerierFunc(func(ctx context.Context, q Query) (Value, error) {
		return query.sqlAll(ctx, func(_ context.Context, spec *sqlgraph.QuerySpec) {
			assign := spec.Assign
			values := spec.ScanValues
			spec.ScanValues = func(columns []string) ([]any, error) {
				values, err := values(columns[1:])
				if err != nil {
					return nil, err
				}
				return append([]any{new(xid.ID)}, values...), nil
			}
			spec.Assign = func(columns []string, values []any) error {
				outValue := *values[0].(*xid.ID)
				inValue := *values[1].(*xid.ID)
				if nids[inValue] == nil {
					nids[inValue] = map[*Tag]struct{}{byID[outValue]: {}}
					return assign(columns[1:], values[1:])
				}
				nids[inValue][byID[outValue]] = struct{}{}
				return nil
			}
		})
	})
	neighbors, err := withInterceptors[[]*Post](ctx, query, qr, query.inters)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected "posts" node returned %v`, n.ID)
		}
		for kn := range nodes {
			assign(kn, n)
		}
	}
	return nil
}
func (tq *TagQuery) loadNodes(ctx context.Context, query *NodeQuery, nodes []*Tag, init func(*Tag), assign func(*Tag, *Node)) error {
	edgeIDs := make([]driver.Value, len(nodes))
	byID := make(map[xid.ID]*Tag)
	nids := make(map[xid.ID]map[*Tag]struct{})
	for i, node := range nodes {
		edgeIDs[i] = node.ID
		byID[node.ID] = node
		if init != nil {
			init(node)
		}
	}
	query.Where(func(s *sql.Selector) {
		joinT := sql.Table(tag.NodesTable)
		s.Join(joinT).On(s.C(node.FieldID), joinT.C(tag.NodesPrimaryKey[1]))
		s.Where(sql.InValues(joinT.C(tag.NodesPrimaryKey[0]), edgeIDs...))
		columns := s.SelectedColumns()
		s.Select(joinT.C(tag.NodesPrimaryKey[0]))
		s.AppendSelect(columns...)
		s.SetDistinct(false)
	})
	if err := query.prepareQuery(ctx); err != nil {
		return err
	}
	qr := QuerierFunc(func(ctx context.Context, q Query) (Value, error) {
		return query.sqlAll(ctx, func(_ context.Context, spec *sqlgraph.QuerySpec) {
			assign := spec.Assign
			values := spec.ScanValues
			spec.ScanValues = func(columns []string) ([]any, error) {
				values, err := values(columns[1:])
				if err != nil {
					return nil, err
				}
				return append([]any{new(xid.ID)}, values...), nil
			}
			spec.Assign = func(columns []string, values []any) error {
				outValue := *values[0].(*xid.ID)
				inValue := *values[1].(*xid.ID)
				if nids[inValue] == nil {
					nids[inValue] = map[*Tag]struct{}{byID[outValue]: {}}
					return assign(columns[1:], values[1:])
				}
				nids[inValue][byID[outValue]] = struct{}{}
				return nil
			}
		})
	})
	neighbors, err := withInterceptors[[]*Node](ctx, query, qr, query.inters)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected "nodes" node returned %v`, n.ID)
		}
		for kn := range nodes {
			assign(kn, n)
		}
	}
	return nil
}
func (tq *TagQuery) loadAccounts(ctx context.Context, query *AccountQuery, nodes []*Tag, init func(*Tag), assign func(*Tag, *Account)) error {
	edgeIDs := make([]driver.Value, len(nodes))
	byID := make(map[xid.ID]*Tag)
	nids := make(map[xid.ID]map[*Tag]struct{})
	for i, node := range nodes {
		edgeIDs[i] = node.ID
		byID[node.ID] = node
		if init != nil {
			init(node)
		}
	}
	query.Where(func(s *sql.Selector) {
		joinT := sql.Table(tag.AccountsTable)
		s.Join(joinT).On(s.C(account.FieldID), joinT.C(tag.AccountsPrimaryKey[0]))
		s.Where(sql.InValues(joinT.C(tag.AccountsPrimaryKey[1]), edgeIDs...))
		columns := s.SelectedColumns()
		s.Select(joinT.C(tag.AccountsPrimaryKey[1]))
		s.AppendSelect(columns...)
		s.SetDistinct(false)
	})
	if err := query.prepareQuery(ctx); err != nil {
		return err
	}
	qr := QuerierFunc(func(ctx context.Context, q Query) (Value, error) {
		return query.sqlAll(ctx, func(_ context.Context, spec *sqlgraph.QuerySpec) {
			assign := spec.Assign
			values := spec.ScanValues
			spec.ScanValues = func(columns []string) ([]any, error) {
				values, err := values(columns[1:])
				if err != nil {
					return nil, err
				}
				return append([]any{new(xid.ID)}, values...), nil
			}
			spec.Assign = func(columns []string, values []any) error {
				outValue := *values[0].(*xid.ID)
				inValue := *values[1].(*xid.ID)
				if nids[inValue] == nil {
					nids[inValue] = map[*Tag]struct{}{byID[outValue]: {}}
					return assign(columns[1:], values[1:])
				}
				nids[inValue][byID[outValue]] = struct{}{}
				return nil
			}
		})
	})
	neighbors, err := withInterceptors[[]*Account](ctx, query, qr, query.inters)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected "accounts" node returned %v`, n.ID)
		}
		for kn := range nodes {
			assign(kn, n)
		}
	}
	return nil
}

func (tq *TagQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := tq.querySpec()
	if len(tq.modifiers) > 0 {
		_spec.Modifiers = tq.modifiers
	}
	_spec.Node.Columns = tq.ctx.Fields
	if len(tq.ctx.Fields) > 0 {
		_spec.Unique = tq.ctx.Unique != nil && *tq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, tq.driver, _spec)
}

func (tq *TagQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(tag.Table, tag.Columns, sqlgraph.NewFieldSpec(tag.FieldID, field.TypeString))
	_spec.From = tq.sql
	if unique := tq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if tq.path != nil {
		_spec.Unique = true
	}
	if fields := tq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, tag.FieldID)
		for i := range fields {
			if fields[i] != tag.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := tq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := tq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := tq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := tq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (tq *TagQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(tq.driver.Dialect())
	t1 := builder.Table(tag.Table)
	columns := tq.ctx.Fields
	if len(columns) == 0 {
		columns = tag.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if tq.sql != nil {
		selector = tq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if tq.ctx.Unique != nil && *tq.ctx.Unique {
		selector.Distinct()
	}
	for _, m := range tq.modifiers {
		m(selector)
	}
	for _, p := range tq.predicates {
		p(selector)
	}
	for _, p := range tq.order {
		p(selector)
	}
	if offset := tq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := tq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// Modify adds a query modifier for attaching custom logic to queries.
func (tq *TagQuery) Modify(modifiers ...func(s *sql.Selector)) *TagSelect {
	tq.modifiers = append(tq.modifiers, modifiers...)
	return tq.Select()
}

// TagGroupBy is the group-by builder for Tag entities.
type TagGroupBy struct {
	selector
	build *TagQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (tgb *TagGroupBy) Aggregate(fns ...AggregateFunc) *TagGroupBy {
	tgb.fns = append(tgb.fns, fns...)
	return tgb
}

// Scan applies the selector query and scans the result into the given value.
func (tgb *TagGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, tgb.build.ctx, ent.OpQueryGroupBy)
	if err := tgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*TagQuery, *TagGroupBy](ctx, tgb.build, tgb, tgb.build.inters, v)
}

func (tgb *TagGroupBy) sqlScan(ctx context.Context, root *TagQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(tgb.fns))
	for _, fn := range tgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*tgb.flds)+len(tgb.fns))
		for _, f := range *tgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*tgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := tgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// TagSelect is the builder for selecting fields of Tag entities.
type TagSelect struct {
	*TagQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (ts *TagSelect) Aggregate(fns ...AggregateFunc) *TagSelect {
	ts.fns = append(ts.fns, fns...)
	return ts
}

// Scan applies the selector query and scans the result into the given value.
func (ts *TagSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ts.ctx, ent.OpQuerySelect)
	if err := ts.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*TagQuery, *TagSelect](ctx, ts.TagQuery, ts, ts.inters, v)
}

func (ts *TagSelect) sqlScan(ctx context.Context, root *TagQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(ts.fns))
	for _, fn := range ts.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*ts.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ts.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (ts *TagSelect) Modify(modifiers ...func(s *sql.Selector)) *TagSelect {
	ts.modifiers = append(ts.modifiers, modifiers...)
	return ts
}
