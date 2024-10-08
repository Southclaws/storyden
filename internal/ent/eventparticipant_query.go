// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/Southclaws/storyden/internal/ent/account"
	"github.com/Southclaws/storyden/internal/ent/event"
	"github.com/Southclaws/storyden/internal/ent/eventparticipant"
	"github.com/Southclaws/storyden/internal/ent/predicate"
	"github.com/rs/xid"
)

// EventParticipantQuery is the builder for querying EventParticipant entities.
type EventParticipantQuery struct {
	config
	ctx         *QueryContext
	order       []eventparticipant.OrderOption
	inters      []Interceptor
	predicates  []predicate.EventParticipant
	withAccount *AccountQuery
	withEvent   *EventQuery
	modifiers   []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the EventParticipantQuery builder.
func (epq *EventParticipantQuery) Where(ps ...predicate.EventParticipant) *EventParticipantQuery {
	epq.predicates = append(epq.predicates, ps...)
	return epq
}

// Limit the number of records to be returned by this query.
func (epq *EventParticipantQuery) Limit(limit int) *EventParticipantQuery {
	epq.ctx.Limit = &limit
	return epq
}

// Offset to start from.
func (epq *EventParticipantQuery) Offset(offset int) *EventParticipantQuery {
	epq.ctx.Offset = &offset
	return epq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (epq *EventParticipantQuery) Unique(unique bool) *EventParticipantQuery {
	epq.ctx.Unique = &unique
	return epq
}

// Order specifies how the records should be ordered.
func (epq *EventParticipantQuery) Order(o ...eventparticipant.OrderOption) *EventParticipantQuery {
	epq.order = append(epq.order, o...)
	return epq
}

// QueryAccount chains the current query on the "account" edge.
func (epq *EventParticipantQuery) QueryAccount() *AccountQuery {
	query := (&AccountClient{config: epq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := epq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := epq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(eventparticipant.Table, eventparticipant.FieldID, selector),
			sqlgraph.To(account.Table, account.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, eventparticipant.AccountTable, eventparticipant.AccountColumn),
		)
		fromU = sqlgraph.SetNeighbors(epq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryEvent chains the current query on the "event" edge.
func (epq *EventParticipantQuery) QueryEvent() *EventQuery {
	query := (&EventClient{config: epq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := epq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := epq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(eventparticipant.Table, eventparticipant.FieldID, selector),
			sqlgraph.To(event.Table, event.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, eventparticipant.EventTable, eventparticipant.EventColumn),
		)
		fromU = sqlgraph.SetNeighbors(epq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first EventParticipant entity from the query.
// Returns a *NotFoundError when no EventParticipant was found.
func (epq *EventParticipantQuery) First(ctx context.Context) (*EventParticipant, error) {
	nodes, err := epq.Limit(1).All(setContextOp(ctx, epq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{eventparticipant.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (epq *EventParticipantQuery) FirstX(ctx context.Context) *EventParticipant {
	node, err := epq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first EventParticipant ID from the query.
// Returns a *NotFoundError when no EventParticipant ID was found.
func (epq *EventParticipantQuery) FirstID(ctx context.Context) (id xid.ID, err error) {
	var ids []xid.ID
	if ids, err = epq.Limit(1).IDs(setContextOp(ctx, epq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{eventparticipant.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (epq *EventParticipantQuery) FirstIDX(ctx context.Context) xid.ID {
	id, err := epq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single EventParticipant entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one EventParticipant entity is found.
// Returns a *NotFoundError when no EventParticipant entities are found.
func (epq *EventParticipantQuery) Only(ctx context.Context) (*EventParticipant, error) {
	nodes, err := epq.Limit(2).All(setContextOp(ctx, epq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{eventparticipant.Label}
	default:
		return nil, &NotSingularError{eventparticipant.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (epq *EventParticipantQuery) OnlyX(ctx context.Context) *EventParticipant {
	node, err := epq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only EventParticipant ID in the query.
// Returns a *NotSingularError when more than one EventParticipant ID is found.
// Returns a *NotFoundError when no entities are found.
func (epq *EventParticipantQuery) OnlyID(ctx context.Context) (id xid.ID, err error) {
	var ids []xid.ID
	if ids, err = epq.Limit(2).IDs(setContextOp(ctx, epq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{eventparticipant.Label}
	default:
		err = &NotSingularError{eventparticipant.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (epq *EventParticipantQuery) OnlyIDX(ctx context.Context) xid.ID {
	id, err := epq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of EventParticipants.
func (epq *EventParticipantQuery) All(ctx context.Context) ([]*EventParticipant, error) {
	ctx = setContextOp(ctx, epq.ctx, ent.OpQueryAll)
	if err := epq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*EventParticipant, *EventParticipantQuery]()
	return withInterceptors[[]*EventParticipant](ctx, epq, qr, epq.inters)
}

// AllX is like All, but panics if an error occurs.
func (epq *EventParticipantQuery) AllX(ctx context.Context) []*EventParticipant {
	nodes, err := epq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of EventParticipant IDs.
func (epq *EventParticipantQuery) IDs(ctx context.Context) (ids []xid.ID, err error) {
	if epq.ctx.Unique == nil && epq.path != nil {
		epq.Unique(true)
	}
	ctx = setContextOp(ctx, epq.ctx, ent.OpQueryIDs)
	if err = epq.Select(eventparticipant.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (epq *EventParticipantQuery) IDsX(ctx context.Context) []xid.ID {
	ids, err := epq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (epq *EventParticipantQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, epq.ctx, ent.OpQueryCount)
	if err := epq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, epq, querierCount[*EventParticipantQuery](), epq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (epq *EventParticipantQuery) CountX(ctx context.Context) int {
	count, err := epq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (epq *EventParticipantQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, epq.ctx, ent.OpQueryExist)
	switch _, err := epq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (epq *EventParticipantQuery) ExistX(ctx context.Context) bool {
	exist, err := epq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the EventParticipantQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (epq *EventParticipantQuery) Clone() *EventParticipantQuery {
	if epq == nil {
		return nil
	}
	return &EventParticipantQuery{
		config:      epq.config,
		ctx:         epq.ctx.Clone(),
		order:       append([]eventparticipant.OrderOption{}, epq.order...),
		inters:      append([]Interceptor{}, epq.inters...),
		predicates:  append([]predicate.EventParticipant{}, epq.predicates...),
		withAccount: epq.withAccount.Clone(),
		withEvent:   epq.withEvent.Clone(),
		// clone intermediate query.
		sql:       epq.sql.Clone(),
		path:      epq.path,
		modifiers: append([]func(*sql.Selector){}, epq.modifiers...),
	}
}

// WithAccount tells the query-builder to eager-load the nodes that are connected to
// the "account" edge. The optional arguments are used to configure the query builder of the edge.
func (epq *EventParticipantQuery) WithAccount(opts ...func(*AccountQuery)) *EventParticipantQuery {
	query := (&AccountClient{config: epq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	epq.withAccount = query
	return epq
}

// WithEvent tells the query-builder to eager-load the nodes that are connected to
// the "event" edge. The optional arguments are used to configure the query builder of the edge.
func (epq *EventParticipantQuery) WithEvent(opts ...func(*EventQuery)) *EventParticipantQuery {
	query := (&EventClient{config: epq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	epq.withEvent = query
	return epq
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
//	client.EventParticipant.Query().
//		GroupBy(eventparticipant.FieldCreatedAt).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (epq *EventParticipantQuery) GroupBy(field string, fields ...string) *EventParticipantGroupBy {
	epq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &EventParticipantGroupBy{build: epq}
	grbuild.flds = &epq.ctx.Fields
	grbuild.label = eventparticipant.Label
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
//	client.EventParticipant.Query().
//		Select(eventparticipant.FieldCreatedAt).
//		Scan(ctx, &v)
func (epq *EventParticipantQuery) Select(fields ...string) *EventParticipantSelect {
	epq.ctx.Fields = append(epq.ctx.Fields, fields...)
	sbuild := &EventParticipantSelect{EventParticipantQuery: epq}
	sbuild.label = eventparticipant.Label
	sbuild.flds, sbuild.scan = &epq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a EventParticipantSelect configured with the given aggregations.
func (epq *EventParticipantQuery) Aggregate(fns ...AggregateFunc) *EventParticipantSelect {
	return epq.Select().Aggregate(fns...)
}

func (epq *EventParticipantQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range epq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, epq); err != nil {
				return err
			}
		}
	}
	for _, f := range epq.ctx.Fields {
		if !eventparticipant.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if epq.path != nil {
		prev, err := epq.path(ctx)
		if err != nil {
			return err
		}
		epq.sql = prev
	}
	return nil
}

func (epq *EventParticipantQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*EventParticipant, error) {
	var (
		nodes       = []*EventParticipant{}
		_spec       = epq.querySpec()
		loadedTypes = [2]bool{
			epq.withAccount != nil,
			epq.withEvent != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*EventParticipant).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &EventParticipant{config: epq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if len(epq.modifiers) > 0 {
		_spec.Modifiers = epq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, epq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := epq.withAccount; query != nil {
		if err := epq.loadAccount(ctx, query, nodes, nil,
			func(n *EventParticipant, e *Account) { n.Edges.Account = e }); err != nil {
			return nil, err
		}
	}
	if query := epq.withEvent; query != nil {
		if err := epq.loadEvent(ctx, query, nodes, nil,
			func(n *EventParticipant, e *Event) { n.Edges.Event = e }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (epq *EventParticipantQuery) loadAccount(ctx context.Context, query *AccountQuery, nodes []*EventParticipant, init func(*EventParticipant), assign func(*EventParticipant, *Account)) error {
	ids := make([]xid.ID, 0, len(nodes))
	nodeids := make(map[xid.ID][]*EventParticipant)
	for i := range nodes {
		fk := nodes[i].AccountID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(account.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "account_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (epq *EventParticipantQuery) loadEvent(ctx context.Context, query *EventQuery, nodes []*EventParticipant, init func(*EventParticipant), assign func(*EventParticipant, *Event)) error {
	ids := make([]xid.ID, 0, len(nodes))
	nodeids := make(map[xid.ID][]*EventParticipant)
	for i := range nodes {
		fk := nodes[i].EventID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(event.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "event_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}

func (epq *EventParticipantQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := epq.querySpec()
	if len(epq.modifiers) > 0 {
		_spec.Modifiers = epq.modifiers
	}
	_spec.Node.Columns = epq.ctx.Fields
	if len(epq.ctx.Fields) > 0 {
		_spec.Unique = epq.ctx.Unique != nil && *epq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, epq.driver, _spec)
}

func (epq *EventParticipantQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(eventparticipant.Table, eventparticipant.Columns, sqlgraph.NewFieldSpec(eventparticipant.FieldID, field.TypeString))
	_spec.From = epq.sql
	if unique := epq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if epq.path != nil {
		_spec.Unique = true
	}
	if fields := epq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, eventparticipant.FieldID)
		for i := range fields {
			if fields[i] != eventparticipant.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
		if epq.withAccount != nil {
			_spec.Node.AddColumnOnce(eventparticipant.FieldAccountID)
		}
		if epq.withEvent != nil {
			_spec.Node.AddColumnOnce(eventparticipant.FieldEventID)
		}
	}
	if ps := epq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := epq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := epq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := epq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (epq *EventParticipantQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(epq.driver.Dialect())
	t1 := builder.Table(eventparticipant.Table)
	columns := epq.ctx.Fields
	if len(columns) == 0 {
		columns = eventparticipant.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if epq.sql != nil {
		selector = epq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if epq.ctx.Unique != nil && *epq.ctx.Unique {
		selector.Distinct()
	}
	for _, m := range epq.modifiers {
		m(selector)
	}
	for _, p := range epq.predicates {
		p(selector)
	}
	for _, p := range epq.order {
		p(selector)
	}
	if offset := epq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := epq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// Modify adds a query modifier for attaching custom logic to queries.
func (epq *EventParticipantQuery) Modify(modifiers ...func(s *sql.Selector)) *EventParticipantSelect {
	epq.modifiers = append(epq.modifiers, modifiers...)
	return epq.Select()
}

// EventParticipantGroupBy is the group-by builder for EventParticipant entities.
type EventParticipantGroupBy struct {
	selector
	build *EventParticipantQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (epgb *EventParticipantGroupBy) Aggregate(fns ...AggregateFunc) *EventParticipantGroupBy {
	epgb.fns = append(epgb.fns, fns...)
	return epgb
}

// Scan applies the selector query and scans the result into the given value.
func (epgb *EventParticipantGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, epgb.build.ctx, ent.OpQueryGroupBy)
	if err := epgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*EventParticipantQuery, *EventParticipantGroupBy](ctx, epgb.build, epgb, epgb.build.inters, v)
}

func (epgb *EventParticipantGroupBy) sqlScan(ctx context.Context, root *EventParticipantQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(epgb.fns))
	for _, fn := range epgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*epgb.flds)+len(epgb.fns))
		for _, f := range *epgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*epgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := epgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// EventParticipantSelect is the builder for selecting fields of EventParticipant entities.
type EventParticipantSelect struct {
	*EventParticipantQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (eps *EventParticipantSelect) Aggregate(fns ...AggregateFunc) *EventParticipantSelect {
	eps.fns = append(eps.fns, fns...)
	return eps
}

// Scan applies the selector query and scans the result into the given value.
func (eps *EventParticipantSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, eps.ctx, ent.OpQuerySelect)
	if err := eps.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*EventParticipantQuery, *EventParticipantSelect](ctx, eps.EventParticipantQuery, eps, eps.inters, v)
}

func (eps *EventParticipantSelect) sqlScan(ctx context.Context, root *EventParticipantQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(eps.fns))
	for _, fn := range eps.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*eps.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := eps.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (eps *EventParticipantSelect) Modify(modifiers ...func(s *sql.Selector)) *EventParticipantSelect {
	eps.modifiers = append(eps.modifiers, modifiers...)
	return eps
}
