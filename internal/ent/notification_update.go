// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/Southclaws/storyden/internal/ent/account"
	"github.com/Southclaws/storyden/internal/ent/notification"
	"github.com/Southclaws/storyden/internal/ent/predicate"
	"github.com/rs/xid"
)

// NotificationUpdate is the builder for updating Notification entities.
type NotificationUpdate struct {
	config
	hooks     []Hook
	mutation  *NotificationMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the NotificationUpdate builder.
func (nu *NotificationUpdate) Where(ps ...predicate.Notification) *NotificationUpdate {
	nu.mutation.Where(ps...)
	return nu
}

// SetDeletedAt sets the "deleted_at" field.
func (nu *NotificationUpdate) SetDeletedAt(t time.Time) *NotificationUpdate {
	nu.mutation.SetDeletedAt(t)
	return nu
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (nu *NotificationUpdate) SetNillableDeletedAt(t *time.Time) *NotificationUpdate {
	if t != nil {
		nu.SetDeletedAt(*t)
	}
	return nu
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (nu *NotificationUpdate) ClearDeletedAt() *NotificationUpdate {
	nu.mutation.ClearDeletedAt()
	return nu
}

// SetEventType sets the "event_type" field.
func (nu *NotificationUpdate) SetEventType(s string) *NotificationUpdate {
	nu.mutation.SetEventType(s)
	return nu
}

// SetNillableEventType sets the "event_type" field if the given value is not nil.
func (nu *NotificationUpdate) SetNillableEventType(s *string) *NotificationUpdate {
	if s != nil {
		nu.SetEventType(*s)
	}
	return nu
}

// SetDatagraphKind sets the "datagraph_kind" field.
func (nu *NotificationUpdate) SetDatagraphKind(s string) *NotificationUpdate {
	nu.mutation.SetDatagraphKind(s)
	return nu
}

// SetNillableDatagraphKind sets the "datagraph_kind" field if the given value is not nil.
func (nu *NotificationUpdate) SetNillableDatagraphKind(s *string) *NotificationUpdate {
	if s != nil {
		nu.SetDatagraphKind(*s)
	}
	return nu
}

// ClearDatagraphKind clears the value of the "datagraph_kind" field.
func (nu *NotificationUpdate) ClearDatagraphKind() *NotificationUpdate {
	nu.mutation.ClearDatagraphKind()
	return nu
}

// SetDatagraphID sets the "datagraph_id" field.
func (nu *NotificationUpdate) SetDatagraphID(x xid.ID) *NotificationUpdate {
	nu.mutation.SetDatagraphID(x)
	return nu
}

// SetNillableDatagraphID sets the "datagraph_id" field if the given value is not nil.
func (nu *NotificationUpdate) SetNillableDatagraphID(x *xid.ID) *NotificationUpdate {
	if x != nil {
		nu.SetDatagraphID(*x)
	}
	return nu
}

// ClearDatagraphID clears the value of the "datagraph_id" field.
func (nu *NotificationUpdate) ClearDatagraphID() *NotificationUpdate {
	nu.mutation.ClearDatagraphID()
	return nu
}

// SetRead sets the "read" field.
func (nu *NotificationUpdate) SetRead(b bool) *NotificationUpdate {
	nu.mutation.SetRead(b)
	return nu
}

// SetNillableRead sets the "read" field if the given value is not nil.
func (nu *NotificationUpdate) SetNillableRead(b *bool) *NotificationUpdate {
	if b != nil {
		nu.SetRead(*b)
	}
	return nu
}

// SetOwnerAccountID sets the "owner_account_id" field.
func (nu *NotificationUpdate) SetOwnerAccountID(x xid.ID) *NotificationUpdate {
	nu.mutation.SetOwnerAccountID(x)
	return nu
}

// SetNillableOwnerAccountID sets the "owner_account_id" field if the given value is not nil.
func (nu *NotificationUpdate) SetNillableOwnerAccountID(x *xid.ID) *NotificationUpdate {
	if x != nil {
		nu.SetOwnerAccountID(*x)
	}
	return nu
}

// SetSourceAccountID sets the "source_account_id" field.
func (nu *NotificationUpdate) SetSourceAccountID(x xid.ID) *NotificationUpdate {
	nu.mutation.SetSourceAccountID(x)
	return nu
}

// SetNillableSourceAccountID sets the "source_account_id" field if the given value is not nil.
func (nu *NotificationUpdate) SetNillableSourceAccountID(x *xid.ID) *NotificationUpdate {
	if x != nil {
		nu.SetSourceAccountID(*x)
	}
	return nu
}

// ClearSourceAccountID clears the value of the "source_account_id" field.
func (nu *NotificationUpdate) ClearSourceAccountID() *NotificationUpdate {
	nu.mutation.ClearSourceAccountID()
	return nu
}

// SetOwnerID sets the "owner" edge to the Account entity by ID.
func (nu *NotificationUpdate) SetOwnerID(id xid.ID) *NotificationUpdate {
	nu.mutation.SetOwnerID(id)
	return nu
}

// SetOwner sets the "owner" edge to the Account entity.
func (nu *NotificationUpdate) SetOwner(a *Account) *NotificationUpdate {
	return nu.SetOwnerID(a.ID)
}

// SetSourceID sets the "source" edge to the Account entity by ID.
func (nu *NotificationUpdate) SetSourceID(id xid.ID) *NotificationUpdate {
	nu.mutation.SetSourceID(id)
	return nu
}

// SetNillableSourceID sets the "source" edge to the Account entity by ID if the given value is not nil.
func (nu *NotificationUpdate) SetNillableSourceID(id *xid.ID) *NotificationUpdate {
	if id != nil {
		nu = nu.SetSourceID(*id)
	}
	return nu
}

// SetSource sets the "source" edge to the Account entity.
func (nu *NotificationUpdate) SetSource(a *Account) *NotificationUpdate {
	return nu.SetSourceID(a.ID)
}

// Mutation returns the NotificationMutation object of the builder.
func (nu *NotificationUpdate) Mutation() *NotificationMutation {
	return nu.mutation
}

// ClearOwner clears the "owner" edge to the Account entity.
func (nu *NotificationUpdate) ClearOwner() *NotificationUpdate {
	nu.mutation.ClearOwner()
	return nu
}

// ClearSource clears the "source" edge to the Account entity.
func (nu *NotificationUpdate) ClearSource() *NotificationUpdate {
	nu.mutation.ClearSource()
	return nu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (nu *NotificationUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, nu.sqlSave, nu.mutation, nu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (nu *NotificationUpdate) SaveX(ctx context.Context) int {
	affected, err := nu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (nu *NotificationUpdate) Exec(ctx context.Context) error {
	_, err := nu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (nu *NotificationUpdate) ExecX(ctx context.Context) {
	if err := nu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (nu *NotificationUpdate) check() error {
	if nu.mutation.OwnerCleared() && len(nu.mutation.OwnerIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "Notification.owner"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (nu *NotificationUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *NotificationUpdate {
	nu.modifiers = append(nu.modifiers, modifiers...)
	return nu
}

func (nu *NotificationUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := nu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(notification.Table, notification.Columns, sqlgraph.NewFieldSpec(notification.FieldID, field.TypeString))
	if ps := nu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := nu.mutation.DeletedAt(); ok {
		_spec.SetField(notification.FieldDeletedAt, field.TypeTime, value)
	}
	if nu.mutation.DeletedAtCleared() {
		_spec.ClearField(notification.FieldDeletedAt, field.TypeTime)
	}
	if value, ok := nu.mutation.EventType(); ok {
		_spec.SetField(notification.FieldEventType, field.TypeString, value)
	}
	if value, ok := nu.mutation.DatagraphKind(); ok {
		_spec.SetField(notification.FieldDatagraphKind, field.TypeString, value)
	}
	if nu.mutation.DatagraphKindCleared() {
		_spec.ClearField(notification.FieldDatagraphKind, field.TypeString)
	}
	if value, ok := nu.mutation.DatagraphID(); ok {
		_spec.SetField(notification.FieldDatagraphID, field.TypeString, value)
	}
	if nu.mutation.DatagraphIDCleared() {
		_spec.ClearField(notification.FieldDatagraphID, field.TypeString)
	}
	if value, ok := nu.mutation.Read(); ok {
		_spec.SetField(notification.FieldRead, field.TypeBool, value)
	}
	if nu.mutation.OwnerCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   notification.OwnerTable,
			Columns: []string{notification.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(account.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := nu.mutation.OwnerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   notification.OwnerTable,
			Columns: []string{notification.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(account.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nu.mutation.SourceCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   notification.SourceTable,
			Columns: []string{notification.SourceColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(account.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := nu.mutation.SourceIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   notification.SourceTable,
			Columns: []string{notification.SourceColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(account.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(nu.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, nu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{notification.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	nu.mutation.done = true
	return n, nil
}

// NotificationUpdateOne is the builder for updating a single Notification entity.
type NotificationUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *NotificationMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetDeletedAt sets the "deleted_at" field.
func (nuo *NotificationUpdateOne) SetDeletedAt(t time.Time) *NotificationUpdateOne {
	nuo.mutation.SetDeletedAt(t)
	return nuo
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (nuo *NotificationUpdateOne) SetNillableDeletedAt(t *time.Time) *NotificationUpdateOne {
	if t != nil {
		nuo.SetDeletedAt(*t)
	}
	return nuo
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (nuo *NotificationUpdateOne) ClearDeletedAt() *NotificationUpdateOne {
	nuo.mutation.ClearDeletedAt()
	return nuo
}

// SetEventType sets the "event_type" field.
func (nuo *NotificationUpdateOne) SetEventType(s string) *NotificationUpdateOne {
	nuo.mutation.SetEventType(s)
	return nuo
}

// SetNillableEventType sets the "event_type" field if the given value is not nil.
func (nuo *NotificationUpdateOne) SetNillableEventType(s *string) *NotificationUpdateOne {
	if s != nil {
		nuo.SetEventType(*s)
	}
	return nuo
}

// SetDatagraphKind sets the "datagraph_kind" field.
func (nuo *NotificationUpdateOne) SetDatagraphKind(s string) *NotificationUpdateOne {
	nuo.mutation.SetDatagraphKind(s)
	return nuo
}

// SetNillableDatagraphKind sets the "datagraph_kind" field if the given value is not nil.
func (nuo *NotificationUpdateOne) SetNillableDatagraphKind(s *string) *NotificationUpdateOne {
	if s != nil {
		nuo.SetDatagraphKind(*s)
	}
	return nuo
}

// ClearDatagraphKind clears the value of the "datagraph_kind" field.
func (nuo *NotificationUpdateOne) ClearDatagraphKind() *NotificationUpdateOne {
	nuo.mutation.ClearDatagraphKind()
	return nuo
}

// SetDatagraphID sets the "datagraph_id" field.
func (nuo *NotificationUpdateOne) SetDatagraphID(x xid.ID) *NotificationUpdateOne {
	nuo.mutation.SetDatagraphID(x)
	return nuo
}

// SetNillableDatagraphID sets the "datagraph_id" field if the given value is not nil.
func (nuo *NotificationUpdateOne) SetNillableDatagraphID(x *xid.ID) *NotificationUpdateOne {
	if x != nil {
		nuo.SetDatagraphID(*x)
	}
	return nuo
}

// ClearDatagraphID clears the value of the "datagraph_id" field.
func (nuo *NotificationUpdateOne) ClearDatagraphID() *NotificationUpdateOne {
	nuo.mutation.ClearDatagraphID()
	return nuo
}

// SetRead sets the "read" field.
func (nuo *NotificationUpdateOne) SetRead(b bool) *NotificationUpdateOne {
	nuo.mutation.SetRead(b)
	return nuo
}

// SetNillableRead sets the "read" field if the given value is not nil.
func (nuo *NotificationUpdateOne) SetNillableRead(b *bool) *NotificationUpdateOne {
	if b != nil {
		nuo.SetRead(*b)
	}
	return nuo
}

// SetOwnerAccountID sets the "owner_account_id" field.
func (nuo *NotificationUpdateOne) SetOwnerAccountID(x xid.ID) *NotificationUpdateOne {
	nuo.mutation.SetOwnerAccountID(x)
	return nuo
}

// SetNillableOwnerAccountID sets the "owner_account_id" field if the given value is not nil.
func (nuo *NotificationUpdateOne) SetNillableOwnerAccountID(x *xid.ID) *NotificationUpdateOne {
	if x != nil {
		nuo.SetOwnerAccountID(*x)
	}
	return nuo
}

// SetSourceAccountID sets the "source_account_id" field.
func (nuo *NotificationUpdateOne) SetSourceAccountID(x xid.ID) *NotificationUpdateOne {
	nuo.mutation.SetSourceAccountID(x)
	return nuo
}

// SetNillableSourceAccountID sets the "source_account_id" field if the given value is not nil.
func (nuo *NotificationUpdateOne) SetNillableSourceAccountID(x *xid.ID) *NotificationUpdateOne {
	if x != nil {
		nuo.SetSourceAccountID(*x)
	}
	return nuo
}

// ClearSourceAccountID clears the value of the "source_account_id" field.
func (nuo *NotificationUpdateOne) ClearSourceAccountID() *NotificationUpdateOne {
	nuo.mutation.ClearSourceAccountID()
	return nuo
}

// SetOwnerID sets the "owner" edge to the Account entity by ID.
func (nuo *NotificationUpdateOne) SetOwnerID(id xid.ID) *NotificationUpdateOne {
	nuo.mutation.SetOwnerID(id)
	return nuo
}

// SetOwner sets the "owner" edge to the Account entity.
func (nuo *NotificationUpdateOne) SetOwner(a *Account) *NotificationUpdateOne {
	return nuo.SetOwnerID(a.ID)
}

// SetSourceID sets the "source" edge to the Account entity by ID.
func (nuo *NotificationUpdateOne) SetSourceID(id xid.ID) *NotificationUpdateOne {
	nuo.mutation.SetSourceID(id)
	return nuo
}

// SetNillableSourceID sets the "source" edge to the Account entity by ID if the given value is not nil.
func (nuo *NotificationUpdateOne) SetNillableSourceID(id *xid.ID) *NotificationUpdateOne {
	if id != nil {
		nuo = nuo.SetSourceID(*id)
	}
	return nuo
}

// SetSource sets the "source" edge to the Account entity.
func (nuo *NotificationUpdateOne) SetSource(a *Account) *NotificationUpdateOne {
	return nuo.SetSourceID(a.ID)
}

// Mutation returns the NotificationMutation object of the builder.
func (nuo *NotificationUpdateOne) Mutation() *NotificationMutation {
	return nuo.mutation
}

// ClearOwner clears the "owner" edge to the Account entity.
func (nuo *NotificationUpdateOne) ClearOwner() *NotificationUpdateOne {
	nuo.mutation.ClearOwner()
	return nuo
}

// ClearSource clears the "source" edge to the Account entity.
func (nuo *NotificationUpdateOne) ClearSource() *NotificationUpdateOne {
	nuo.mutation.ClearSource()
	return nuo
}

// Where appends a list predicates to the NotificationUpdate builder.
func (nuo *NotificationUpdateOne) Where(ps ...predicate.Notification) *NotificationUpdateOne {
	nuo.mutation.Where(ps...)
	return nuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (nuo *NotificationUpdateOne) Select(field string, fields ...string) *NotificationUpdateOne {
	nuo.fields = append([]string{field}, fields...)
	return nuo
}

// Save executes the query and returns the updated Notification entity.
func (nuo *NotificationUpdateOne) Save(ctx context.Context) (*Notification, error) {
	return withHooks(ctx, nuo.sqlSave, nuo.mutation, nuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (nuo *NotificationUpdateOne) SaveX(ctx context.Context) *Notification {
	node, err := nuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (nuo *NotificationUpdateOne) Exec(ctx context.Context) error {
	_, err := nuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (nuo *NotificationUpdateOne) ExecX(ctx context.Context) {
	if err := nuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (nuo *NotificationUpdateOne) check() error {
	if nuo.mutation.OwnerCleared() && len(nuo.mutation.OwnerIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "Notification.owner"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (nuo *NotificationUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *NotificationUpdateOne {
	nuo.modifiers = append(nuo.modifiers, modifiers...)
	return nuo
}

func (nuo *NotificationUpdateOne) sqlSave(ctx context.Context) (_node *Notification, err error) {
	if err := nuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(notification.Table, notification.Columns, sqlgraph.NewFieldSpec(notification.FieldID, field.TypeString))
	id, ok := nuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Notification.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := nuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, notification.FieldID)
		for _, f := range fields {
			if !notification.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != notification.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := nuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := nuo.mutation.DeletedAt(); ok {
		_spec.SetField(notification.FieldDeletedAt, field.TypeTime, value)
	}
	if nuo.mutation.DeletedAtCleared() {
		_spec.ClearField(notification.FieldDeletedAt, field.TypeTime)
	}
	if value, ok := nuo.mutation.EventType(); ok {
		_spec.SetField(notification.FieldEventType, field.TypeString, value)
	}
	if value, ok := nuo.mutation.DatagraphKind(); ok {
		_spec.SetField(notification.FieldDatagraphKind, field.TypeString, value)
	}
	if nuo.mutation.DatagraphKindCleared() {
		_spec.ClearField(notification.FieldDatagraphKind, field.TypeString)
	}
	if value, ok := nuo.mutation.DatagraphID(); ok {
		_spec.SetField(notification.FieldDatagraphID, field.TypeString, value)
	}
	if nuo.mutation.DatagraphIDCleared() {
		_spec.ClearField(notification.FieldDatagraphID, field.TypeString)
	}
	if value, ok := nuo.mutation.Read(); ok {
		_spec.SetField(notification.FieldRead, field.TypeBool, value)
	}
	if nuo.mutation.OwnerCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   notification.OwnerTable,
			Columns: []string{notification.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(account.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := nuo.mutation.OwnerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   notification.OwnerTable,
			Columns: []string{notification.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(account.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nuo.mutation.SourceCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   notification.SourceTable,
			Columns: []string{notification.SourceColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(account.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := nuo.mutation.SourceIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   notification.SourceTable,
			Columns: []string{notification.SourceColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(account.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(nuo.modifiers...)
	_node = &Notification{config: nuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, nuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{notification.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	nuo.mutation.done = true
	return _node, nil
}
