// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/Southclaws/storyden/internal/ent/node"
	"github.com/Southclaws/storyden/internal/ent/predicate"
	"github.com/Southclaws/storyden/internal/ent/propertyschema"
	"github.com/Southclaws/storyden/internal/ent/propertyschemafield"
	"github.com/rs/xid"
)

// PropertySchemaUpdate is the builder for updating PropertySchema entities.
type PropertySchemaUpdate struct {
	config
	hooks     []Hook
	mutation  *PropertySchemaMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the PropertySchemaUpdate builder.
func (psu *PropertySchemaUpdate) Where(ps ...predicate.PropertySchema) *PropertySchemaUpdate {
	psu.mutation.Where(ps...)
	return psu
}

// AddNodeIDs adds the "node" edge to the Node entity by IDs.
func (psu *PropertySchemaUpdate) AddNodeIDs(ids ...xid.ID) *PropertySchemaUpdate {
	psu.mutation.AddNodeIDs(ids...)
	return psu
}

// AddNode adds the "node" edges to the Node entity.
func (psu *PropertySchemaUpdate) AddNode(n ...*Node) *PropertySchemaUpdate {
	ids := make([]xid.ID, len(n))
	for i := range n {
		ids[i] = n[i].ID
	}
	return psu.AddNodeIDs(ids...)
}

// AddFieldIDs adds the "fields" edge to the PropertySchemaField entity by IDs.
func (psu *PropertySchemaUpdate) AddFieldIDs(ids ...xid.ID) *PropertySchemaUpdate {
	psu.mutation.AddFieldIDs(ids...)
	return psu
}

// AddFields adds the "fields" edges to the PropertySchemaField entity.
func (psu *PropertySchemaUpdate) AddFields(p ...*PropertySchemaField) *PropertySchemaUpdate {
	ids := make([]xid.ID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return psu.AddFieldIDs(ids...)
}

// Mutation returns the PropertySchemaMutation object of the builder.
func (psu *PropertySchemaUpdate) Mutation() *PropertySchemaMutation {
	return psu.mutation
}

// ClearNode clears all "node" edges to the Node entity.
func (psu *PropertySchemaUpdate) ClearNode() *PropertySchemaUpdate {
	psu.mutation.ClearNode()
	return psu
}

// RemoveNodeIDs removes the "node" edge to Node entities by IDs.
func (psu *PropertySchemaUpdate) RemoveNodeIDs(ids ...xid.ID) *PropertySchemaUpdate {
	psu.mutation.RemoveNodeIDs(ids...)
	return psu
}

// RemoveNode removes "node" edges to Node entities.
func (psu *PropertySchemaUpdate) RemoveNode(n ...*Node) *PropertySchemaUpdate {
	ids := make([]xid.ID, len(n))
	for i := range n {
		ids[i] = n[i].ID
	}
	return psu.RemoveNodeIDs(ids...)
}

// ClearFields clears all "fields" edges to the PropertySchemaField entity.
func (psu *PropertySchemaUpdate) ClearFields() *PropertySchemaUpdate {
	psu.mutation.ClearFields()
	return psu
}

// RemoveFieldIDs removes the "fields" edge to PropertySchemaField entities by IDs.
func (psu *PropertySchemaUpdate) RemoveFieldIDs(ids ...xid.ID) *PropertySchemaUpdate {
	psu.mutation.RemoveFieldIDs(ids...)
	return psu
}

// RemoveFields removes "fields" edges to PropertySchemaField entities.
func (psu *PropertySchemaUpdate) RemoveFields(p ...*PropertySchemaField) *PropertySchemaUpdate {
	ids := make([]xid.ID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return psu.RemoveFieldIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (psu *PropertySchemaUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, psu.sqlSave, psu.mutation, psu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (psu *PropertySchemaUpdate) SaveX(ctx context.Context) int {
	affected, err := psu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (psu *PropertySchemaUpdate) Exec(ctx context.Context) error {
	_, err := psu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (psu *PropertySchemaUpdate) ExecX(ctx context.Context) {
	if err := psu.Exec(ctx); err != nil {
		panic(err)
	}
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (psu *PropertySchemaUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *PropertySchemaUpdate {
	psu.modifiers = append(psu.modifiers, modifiers...)
	return psu
}

func (psu *PropertySchemaUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(propertyschema.Table, propertyschema.Columns, sqlgraph.NewFieldSpec(propertyschema.FieldID, field.TypeString))
	if ps := psu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if psu.mutation.NodeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   propertyschema.NodeTable,
			Columns: []string{propertyschema.NodeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(node.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := psu.mutation.RemovedNodeIDs(); len(nodes) > 0 && !psu.mutation.NodeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   propertyschema.NodeTable,
			Columns: []string{propertyschema.NodeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(node.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := psu.mutation.NodeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   propertyschema.NodeTable,
			Columns: []string{propertyschema.NodeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(node.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if psu.mutation.FieldsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   propertyschema.FieldsTable,
			Columns: []string{propertyschema.FieldsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(propertyschemafield.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := psu.mutation.RemovedFieldsIDs(); len(nodes) > 0 && !psu.mutation.FieldsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   propertyschema.FieldsTable,
			Columns: []string{propertyschema.FieldsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(propertyschemafield.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := psu.mutation.FieldsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   propertyschema.FieldsTable,
			Columns: []string{propertyschema.FieldsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(propertyschemafield.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(psu.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, psu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{propertyschema.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	psu.mutation.done = true
	return n, nil
}

// PropertySchemaUpdateOne is the builder for updating a single PropertySchema entity.
type PropertySchemaUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *PropertySchemaMutation
	modifiers []func(*sql.UpdateBuilder)
}

// AddNodeIDs adds the "node" edge to the Node entity by IDs.
func (psuo *PropertySchemaUpdateOne) AddNodeIDs(ids ...xid.ID) *PropertySchemaUpdateOne {
	psuo.mutation.AddNodeIDs(ids...)
	return psuo
}

// AddNode adds the "node" edges to the Node entity.
func (psuo *PropertySchemaUpdateOne) AddNode(n ...*Node) *PropertySchemaUpdateOne {
	ids := make([]xid.ID, len(n))
	for i := range n {
		ids[i] = n[i].ID
	}
	return psuo.AddNodeIDs(ids...)
}

// AddFieldIDs adds the "fields" edge to the PropertySchemaField entity by IDs.
func (psuo *PropertySchemaUpdateOne) AddFieldIDs(ids ...xid.ID) *PropertySchemaUpdateOne {
	psuo.mutation.AddFieldIDs(ids...)
	return psuo
}

// AddFields adds the "fields" edges to the PropertySchemaField entity.
func (psuo *PropertySchemaUpdateOne) AddFields(p ...*PropertySchemaField) *PropertySchemaUpdateOne {
	ids := make([]xid.ID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return psuo.AddFieldIDs(ids...)
}

// Mutation returns the PropertySchemaMutation object of the builder.
func (psuo *PropertySchemaUpdateOne) Mutation() *PropertySchemaMutation {
	return psuo.mutation
}

// ClearNode clears all "node" edges to the Node entity.
func (psuo *PropertySchemaUpdateOne) ClearNode() *PropertySchemaUpdateOne {
	psuo.mutation.ClearNode()
	return psuo
}

// RemoveNodeIDs removes the "node" edge to Node entities by IDs.
func (psuo *PropertySchemaUpdateOne) RemoveNodeIDs(ids ...xid.ID) *PropertySchemaUpdateOne {
	psuo.mutation.RemoveNodeIDs(ids...)
	return psuo
}

// RemoveNode removes "node" edges to Node entities.
func (psuo *PropertySchemaUpdateOne) RemoveNode(n ...*Node) *PropertySchemaUpdateOne {
	ids := make([]xid.ID, len(n))
	for i := range n {
		ids[i] = n[i].ID
	}
	return psuo.RemoveNodeIDs(ids...)
}

// ClearFields clears all "fields" edges to the PropertySchemaField entity.
func (psuo *PropertySchemaUpdateOne) ClearFields() *PropertySchemaUpdateOne {
	psuo.mutation.ClearFields()
	return psuo
}

// RemoveFieldIDs removes the "fields" edge to PropertySchemaField entities by IDs.
func (psuo *PropertySchemaUpdateOne) RemoveFieldIDs(ids ...xid.ID) *PropertySchemaUpdateOne {
	psuo.mutation.RemoveFieldIDs(ids...)
	return psuo
}

// RemoveFields removes "fields" edges to PropertySchemaField entities.
func (psuo *PropertySchemaUpdateOne) RemoveFields(p ...*PropertySchemaField) *PropertySchemaUpdateOne {
	ids := make([]xid.ID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return psuo.RemoveFieldIDs(ids...)
}

// Where appends a list predicates to the PropertySchemaUpdate builder.
func (psuo *PropertySchemaUpdateOne) Where(ps ...predicate.PropertySchema) *PropertySchemaUpdateOne {
	psuo.mutation.Where(ps...)
	return psuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (psuo *PropertySchemaUpdateOne) Select(field string, fields ...string) *PropertySchemaUpdateOne {
	psuo.fields = append([]string{field}, fields...)
	return psuo
}

// Save executes the query and returns the updated PropertySchema entity.
func (psuo *PropertySchemaUpdateOne) Save(ctx context.Context) (*PropertySchema, error) {
	return withHooks(ctx, psuo.sqlSave, psuo.mutation, psuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (psuo *PropertySchemaUpdateOne) SaveX(ctx context.Context) *PropertySchema {
	node, err := psuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (psuo *PropertySchemaUpdateOne) Exec(ctx context.Context) error {
	_, err := psuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (psuo *PropertySchemaUpdateOne) ExecX(ctx context.Context) {
	if err := psuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (psuo *PropertySchemaUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *PropertySchemaUpdateOne {
	psuo.modifiers = append(psuo.modifiers, modifiers...)
	return psuo
}

func (psuo *PropertySchemaUpdateOne) sqlSave(ctx context.Context) (_node *PropertySchema, err error) {
	_spec := sqlgraph.NewUpdateSpec(propertyschema.Table, propertyschema.Columns, sqlgraph.NewFieldSpec(propertyschema.FieldID, field.TypeString))
	id, ok := psuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "PropertySchema.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := psuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, propertyschema.FieldID)
		for _, f := range fields {
			if !propertyschema.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != propertyschema.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := psuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if psuo.mutation.NodeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   propertyschema.NodeTable,
			Columns: []string{propertyschema.NodeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(node.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := psuo.mutation.RemovedNodeIDs(); len(nodes) > 0 && !psuo.mutation.NodeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   propertyschema.NodeTable,
			Columns: []string{propertyschema.NodeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(node.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := psuo.mutation.NodeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   propertyschema.NodeTable,
			Columns: []string{propertyschema.NodeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(node.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if psuo.mutation.FieldsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   propertyschema.FieldsTable,
			Columns: []string{propertyschema.FieldsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(propertyschemafield.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := psuo.mutation.RemovedFieldsIDs(); len(nodes) > 0 && !psuo.mutation.FieldsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   propertyschema.FieldsTable,
			Columns: []string{propertyschema.FieldsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(propertyschemafield.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := psuo.mutation.FieldsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   propertyschema.FieldsTable,
			Columns: []string{propertyschema.FieldsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(propertyschemafield.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(psuo.modifiers...)
	_node = &PropertySchema{config: psuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, psuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{propertyschema.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	psuo.mutation.done = true
	return _node, nil
}
