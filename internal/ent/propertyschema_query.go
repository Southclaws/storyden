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
	"github.com/Southclaws/storyden/internal/ent/node"
	"github.com/Southclaws/storyden/internal/ent/predicate"
	"github.com/Southclaws/storyden/internal/ent/propertyschema"
	"github.com/Southclaws/storyden/internal/ent/propertyschemafield"
	"github.com/rs/xid"
)

// PropertySchemaQuery is the builder for querying PropertySchema entities.
type PropertySchemaQuery struct {
	config
	ctx        *QueryContext
	order      []propertyschema.OrderOption
	inters     []Interceptor
	predicates []predicate.PropertySchema
	withNode   *NodeQuery
	withFields *PropertySchemaFieldQuery
	modifiers  []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the PropertySchemaQuery builder.
func (psq *PropertySchemaQuery) Where(ps ...predicate.PropertySchema) *PropertySchemaQuery {
	psq.predicates = append(psq.predicates, ps...)
	return psq
}

// Limit the number of records to be returned by this query.
func (psq *PropertySchemaQuery) Limit(limit int) *PropertySchemaQuery {
	psq.ctx.Limit = &limit
	return psq
}

// Offset to start from.
func (psq *PropertySchemaQuery) Offset(offset int) *PropertySchemaQuery {
	psq.ctx.Offset = &offset
	return psq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (psq *PropertySchemaQuery) Unique(unique bool) *PropertySchemaQuery {
	psq.ctx.Unique = &unique
	return psq
}

// Order specifies how the records should be ordered.
func (psq *PropertySchemaQuery) Order(o ...propertyschema.OrderOption) *PropertySchemaQuery {
	psq.order = append(psq.order, o...)
	return psq
}

// QueryNode chains the current query on the "node" edge.
func (psq *PropertySchemaQuery) QueryNode() *NodeQuery {
	query := (&NodeClient{config: psq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := psq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := psq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(propertyschema.Table, propertyschema.FieldID, selector),
			sqlgraph.To(node.Table, node.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, propertyschema.NodeTable, propertyschema.NodeColumn),
		)
		fromU = sqlgraph.SetNeighbors(psq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryFields chains the current query on the "fields" edge.
func (psq *PropertySchemaQuery) QueryFields() *PropertySchemaFieldQuery {
	query := (&PropertySchemaFieldClient{config: psq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := psq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := psq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(propertyschema.Table, propertyschema.FieldID, selector),
			sqlgraph.To(propertyschemafield.Table, propertyschemafield.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, propertyschema.FieldsTable, propertyschema.FieldsColumn),
		)
		fromU = sqlgraph.SetNeighbors(psq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first PropertySchema entity from the query.
// Returns a *NotFoundError when no PropertySchema was found.
func (psq *PropertySchemaQuery) First(ctx context.Context) (*PropertySchema, error) {
	nodes, err := psq.Limit(1).All(setContextOp(ctx, psq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{propertyschema.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (psq *PropertySchemaQuery) FirstX(ctx context.Context) *PropertySchema {
	node, err := psq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first PropertySchema ID from the query.
// Returns a *NotFoundError when no PropertySchema ID was found.
func (psq *PropertySchemaQuery) FirstID(ctx context.Context) (id xid.ID, err error) {
	var ids []xid.ID
	if ids, err = psq.Limit(1).IDs(setContextOp(ctx, psq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{propertyschema.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (psq *PropertySchemaQuery) FirstIDX(ctx context.Context) xid.ID {
	id, err := psq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single PropertySchema entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one PropertySchema entity is found.
// Returns a *NotFoundError when no PropertySchema entities are found.
func (psq *PropertySchemaQuery) Only(ctx context.Context) (*PropertySchema, error) {
	nodes, err := psq.Limit(2).All(setContextOp(ctx, psq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{propertyschema.Label}
	default:
		return nil, &NotSingularError{propertyschema.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (psq *PropertySchemaQuery) OnlyX(ctx context.Context) *PropertySchema {
	node, err := psq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only PropertySchema ID in the query.
// Returns a *NotSingularError when more than one PropertySchema ID is found.
// Returns a *NotFoundError when no entities are found.
func (psq *PropertySchemaQuery) OnlyID(ctx context.Context) (id xid.ID, err error) {
	var ids []xid.ID
	if ids, err = psq.Limit(2).IDs(setContextOp(ctx, psq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{propertyschema.Label}
	default:
		err = &NotSingularError{propertyschema.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (psq *PropertySchemaQuery) OnlyIDX(ctx context.Context) xid.ID {
	id, err := psq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of PropertySchemas.
func (psq *PropertySchemaQuery) All(ctx context.Context) ([]*PropertySchema, error) {
	ctx = setContextOp(ctx, psq.ctx, ent.OpQueryAll)
	if err := psq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*PropertySchema, *PropertySchemaQuery]()
	return withInterceptors[[]*PropertySchema](ctx, psq, qr, psq.inters)
}

// AllX is like All, but panics if an error occurs.
func (psq *PropertySchemaQuery) AllX(ctx context.Context) []*PropertySchema {
	nodes, err := psq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of PropertySchema IDs.
func (psq *PropertySchemaQuery) IDs(ctx context.Context) (ids []xid.ID, err error) {
	if psq.ctx.Unique == nil && psq.path != nil {
		psq.Unique(true)
	}
	ctx = setContextOp(ctx, psq.ctx, ent.OpQueryIDs)
	if err = psq.Select(propertyschema.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (psq *PropertySchemaQuery) IDsX(ctx context.Context) []xid.ID {
	ids, err := psq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (psq *PropertySchemaQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, psq.ctx, ent.OpQueryCount)
	if err := psq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, psq, querierCount[*PropertySchemaQuery](), psq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (psq *PropertySchemaQuery) CountX(ctx context.Context) int {
	count, err := psq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (psq *PropertySchemaQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, psq.ctx, ent.OpQueryExist)
	switch _, err := psq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (psq *PropertySchemaQuery) ExistX(ctx context.Context) bool {
	exist, err := psq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the PropertySchemaQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (psq *PropertySchemaQuery) Clone() *PropertySchemaQuery {
	if psq == nil {
		return nil
	}
	return &PropertySchemaQuery{
		config:     psq.config,
		ctx:        psq.ctx.Clone(),
		order:      append([]propertyschema.OrderOption{}, psq.order...),
		inters:     append([]Interceptor{}, psq.inters...),
		predicates: append([]predicate.PropertySchema{}, psq.predicates...),
		withNode:   psq.withNode.Clone(),
		withFields: psq.withFields.Clone(),
		// clone intermediate query.
		sql:       psq.sql.Clone(),
		path:      psq.path,
		modifiers: append([]func(*sql.Selector){}, psq.modifiers...),
	}
}

// WithNode tells the query-builder to eager-load the nodes that are connected to
// the "node" edge. The optional arguments are used to configure the query builder of the edge.
func (psq *PropertySchemaQuery) WithNode(opts ...func(*NodeQuery)) *PropertySchemaQuery {
	query := (&NodeClient{config: psq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	psq.withNode = query
	return psq
}

// WithFields tells the query-builder to eager-load the nodes that are connected to
// the "fields" edge. The optional arguments are used to configure the query builder of the edge.
func (psq *PropertySchemaQuery) WithFields(opts ...func(*PropertySchemaFieldQuery)) *PropertySchemaQuery {
	query := (&PropertySchemaFieldClient{config: psq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	psq.withFields = query
	return psq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
func (psq *PropertySchemaQuery) GroupBy(field string, fields ...string) *PropertySchemaGroupBy {
	psq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &PropertySchemaGroupBy{build: psq}
	grbuild.flds = &psq.ctx.Fields
	grbuild.label = propertyschema.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
func (psq *PropertySchemaQuery) Select(fields ...string) *PropertySchemaSelect {
	psq.ctx.Fields = append(psq.ctx.Fields, fields...)
	sbuild := &PropertySchemaSelect{PropertySchemaQuery: psq}
	sbuild.label = propertyschema.Label
	sbuild.flds, sbuild.scan = &psq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a PropertySchemaSelect configured with the given aggregations.
func (psq *PropertySchemaQuery) Aggregate(fns ...AggregateFunc) *PropertySchemaSelect {
	return psq.Select().Aggregate(fns...)
}

func (psq *PropertySchemaQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range psq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, psq); err != nil {
				return err
			}
		}
	}
	for _, f := range psq.ctx.Fields {
		if !propertyschema.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if psq.path != nil {
		prev, err := psq.path(ctx)
		if err != nil {
			return err
		}
		psq.sql = prev
	}
	return nil
}

func (psq *PropertySchemaQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*PropertySchema, error) {
	var (
		nodes       = []*PropertySchema{}
		_spec       = psq.querySpec()
		loadedTypes = [2]bool{
			psq.withNode != nil,
			psq.withFields != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*PropertySchema).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &PropertySchema{config: psq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if len(psq.modifiers) > 0 {
		_spec.Modifiers = psq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, psq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := psq.withNode; query != nil {
		if err := psq.loadNode(ctx, query, nodes,
			func(n *PropertySchema) { n.Edges.Node = []*Node{} },
			func(n *PropertySchema, e *Node) { n.Edges.Node = append(n.Edges.Node, e) }); err != nil {
			return nil, err
		}
	}
	if query := psq.withFields; query != nil {
		if err := psq.loadFields(ctx, query, nodes,
			func(n *PropertySchema) { n.Edges.Fields = []*PropertySchemaField{} },
			func(n *PropertySchema, e *PropertySchemaField) { n.Edges.Fields = append(n.Edges.Fields, e) }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (psq *PropertySchemaQuery) loadNode(ctx context.Context, query *NodeQuery, nodes []*PropertySchema, init func(*PropertySchema), assign func(*PropertySchema, *Node)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[xid.ID]*PropertySchema)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(node.FieldPropertySchemaID)
	}
	query.Where(predicate.Node(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(propertyschema.NodeColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.PropertySchemaID
		if fk == nil {
			return fmt.Errorf(`foreign-key "property_schema_id" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "property_schema_id" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}
func (psq *PropertySchemaQuery) loadFields(ctx context.Context, query *PropertySchemaFieldQuery, nodes []*PropertySchema, init func(*PropertySchema), assign func(*PropertySchema, *PropertySchemaField)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[xid.ID]*PropertySchema)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(propertyschemafield.FieldSchemaID)
	}
	query.Where(predicate.PropertySchemaField(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(propertyschema.FieldsColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.SchemaID
		node, ok := nodeids[fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "schema_id" returned %v for node %v`, fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}

func (psq *PropertySchemaQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := psq.querySpec()
	if len(psq.modifiers) > 0 {
		_spec.Modifiers = psq.modifiers
	}
	_spec.Node.Columns = psq.ctx.Fields
	if len(psq.ctx.Fields) > 0 {
		_spec.Unique = psq.ctx.Unique != nil && *psq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, psq.driver, _spec)
}

func (psq *PropertySchemaQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(propertyschema.Table, propertyschema.Columns, sqlgraph.NewFieldSpec(propertyschema.FieldID, field.TypeString))
	_spec.From = psq.sql
	if unique := psq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if psq.path != nil {
		_spec.Unique = true
	}
	if fields := psq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, propertyschema.FieldID)
		for i := range fields {
			if fields[i] != propertyschema.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := psq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := psq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := psq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := psq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (psq *PropertySchemaQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(psq.driver.Dialect())
	t1 := builder.Table(propertyschema.Table)
	columns := psq.ctx.Fields
	if len(columns) == 0 {
		columns = propertyschema.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if psq.sql != nil {
		selector = psq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if psq.ctx.Unique != nil && *psq.ctx.Unique {
		selector.Distinct()
	}
	for _, m := range psq.modifiers {
		m(selector)
	}
	for _, p := range psq.predicates {
		p(selector)
	}
	for _, p := range psq.order {
		p(selector)
	}
	if offset := psq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := psq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// Modify adds a query modifier for attaching custom logic to queries.
func (psq *PropertySchemaQuery) Modify(modifiers ...func(s *sql.Selector)) *PropertySchemaSelect {
	psq.modifiers = append(psq.modifiers, modifiers...)
	return psq.Select()
}

// PropertySchemaGroupBy is the group-by builder for PropertySchema entities.
type PropertySchemaGroupBy struct {
	selector
	build *PropertySchemaQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (psgb *PropertySchemaGroupBy) Aggregate(fns ...AggregateFunc) *PropertySchemaGroupBy {
	psgb.fns = append(psgb.fns, fns...)
	return psgb
}

// Scan applies the selector query and scans the result into the given value.
func (psgb *PropertySchemaGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, psgb.build.ctx, ent.OpQueryGroupBy)
	if err := psgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*PropertySchemaQuery, *PropertySchemaGroupBy](ctx, psgb.build, psgb, psgb.build.inters, v)
}

func (psgb *PropertySchemaGroupBy) sqlScan(ctx context.Context, root *PropertySchemaQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(psgb.fns))
	for _, fn := range psgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*psgb.flds)+len(psgb.fns))
		for _, f := range *psgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*psgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := psgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// PropertySchemaSelect is the builder for selecting fields of PropertySchema entities.
type PropertySchemaSelect struct {
	*PropertySchemaQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (pss *PropertySchemaSelect) Aggregate(fns ...AggregateFunc) *PropertySchemaSelect {
	pss.fns = append(pss.fns, fns...)
	return pss
}

// Scan applies the selector query and scans the result into the given value.
func (pss *PropertySchemaSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, pss.ctx, ent.OpQuerySelect)
	if err := pss.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*PropertySchemaQuery, *PropertySchemaSelect](ctx, pss.PropertySchemaQuery, pss, pss.inters, v)
}

func (pss *PropertySchemaSelect) sqlScan(ctx context.Context, root *PropertySchemaQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(pss.fns))
	for _, fn := range pss.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*pss.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := pss.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (pss *PropertySchemaSelect) Modify(modifiers ...func(s *sql.Selector)) *PropertySchemaSelect {
	pss.modifiers = append(pss.modifiers, modifiers...)
	return pss
}
