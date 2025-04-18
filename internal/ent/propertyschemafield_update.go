// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/Southclaws/storyden/internal/ent/predicate"
	"github.com/Southclaws/storyden/internal/ent/property"
	"github.com/Southclaws/storyden/internal/ent/propertyschema"
	"github.com/Southclaws/storyden/internal/ent/propertyschemafield"
	"github.com/rs/xid"
)

// PropertySchemaFieldUpdate is the builder for updating PropertySchemaField entities.
type PropertySchemaFieldUpdate struct {
	config
	hooks     []Hook
	mutation  *PropertySchemaFieldMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the PropertySchemaFieldUpdate builder.
func (psfu *PropertySchemaFieldUpdate) Where(ps ...predicate.PropertySchemaField) *PropertySchemaFieldUpdate {
	psfu.mutation.Where(ps...)
	return psfu
}

// SetName sets the "name" field.
func (psfu *PropertySchemaFieldUpdate) SetName(s string) *PropertySchemaFieldUpdate {
	psfu.mutation.SetName(s)
	return psfu
}

// SetNillableName sets the "name" field if the given value is not nil.
func (psfu *PropertySchemaFieldUpdate) SetNillableName(s *string) *PropertySchemaFieldUpdate {
	if s != nil {
		psfu.SetName(*s)
	}
	return psfu
}

// SetType sets the "type" field.
func (psfu *PropertySchemaFieldUpdate) SetType(s string) *PropertySchemaFieldUpdate {
	psfu.mutation.SetType(s)
	return psfu
}

// SetNillableType sets the "type" field if the given value is not nil.
func (psfu *PropertySchemaFieldUpdate) SetNillableType(s *string) *PropertySchemaFieldUpdate {
	if s != nil {
		psfu.SetType(*s)
	}
	return psfu
}

// SetSort sets the "sort" field.
func (psfu *PropertySchemaFieldUpdate) SetSort(s string) *PropertySchemaFieldUpdate {
	psfu.mutation.SetSort(s)
	return psfu
}

// SetNillableSort sets the "sort" field if the given value is not nil.
func (psfu *PropertySchemaFieldUpdate) SetNillableSort(s *string) *PropertySchemaFieldUpdate {
	if s != nil {
		psfu.SetSort(*s)
	}
	return psfu
}

// SetSchemaID sets the "schema_id" field.
func (psfu *PropertySchemaFieldUpdate) SetSchemaID(x xid.ID) *PropertySchemaFieldUpdate {
	psfu.mutation.SetSchemaID(x)
	return psfu
}

// SetNillableSchemaID sets the "schema_id" field if the given value is not nil.
func (psfu *PropertySchemaFieldUpdate) SetNillableSchemaID(x *xid.ID) *PropertySchemaFieldUpdate {
	if x != nil {
		psfu.SetSchemaID(*x)
	}
	return psfu
}

// SetSchema sets the "schema" edge to the PropertySchema entity.
func (psfu *PropertySchemaFieldUpdate) SetSchema(p *PropertySchema) *PropertySchemaFieldUpdate {
	return psfu.SetSchemaID(p.ID)
}

// AddPropertyIDs adds the "properties" edge to the Property entity by IDs.
func (psfu *PropertySchemaFieldUpdate) AddPropertyIDs(ids ...xid.ID) *PropertySchemaFieldUpdate {
	psfu.mutation.AddPropertyIDs(ids...)
	return psfu
}

// AddProperties adds the "properties" edges to the Property entity.
func (psfu *PropertySchemaFieldUpdate) AddProperties(p ...*Property) *PropertySchemaFieldUpdate {
	ids := make([]xid.ID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return psfu.AddPropertyIDs(ids...)
}

// Mutation returns the PropertySchemaFieldMutation object of the builder.
func (psfu *PropertySchemaFieldUpdate) Mutation() *PropertySchemaFieldMutation {
	return psfu.mutation
}

// ClearSchema clears the "schema" edge to the PropertySchema entity.
func (psfu *PropertySchemaFieldUpdate) ClearSchema() *PropertySchemaFieldUpdate {
	psfu.mutation.ClearSchema()
	return psfu
}

// ClearProperties clears all "properties" edges to the Property entity.
func (psfu *PropertySchemaFieldUpdate) ClearProperties() *PropertySchemaFieldUpdate {
	psfu.mutation.ClearProperties()
	return psfu
}

// RemovePropertyIDs removes the "properties" edge to Property entities by IDs.
func (psfu *PropertySchemaFieldUpdate) RemovePropertyIDs(ids ...xid.ID) *PropertySchemaFieldUpdate {
	psfu.mutation.RemovePropertyIDs(ids...)
	return psfu
}

// RemoveProperties removes "properties" edges to Property entities.
func (psfu *PropertySchemaFieldUpdate) RemoveProperties(p ...*Property) *PropertySchemaFieldUpdate {
	ids := make([]xid.ID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return psfu.RemovePropertyIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (psfu *PropertySchemaFieldUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, psfu.sqlSave, psfu.mutation, psfu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (psfu *PropertySchemaFieldUpdate) SaveX(ctx context.Context) int {
	affected, err := psfu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (psfu *PropertySchemaFieldUpdate) Exec(ctx context.Context) error {
	_, err := psfu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (psfu *PropertySchemaFieldUpdate) ExecX(ctx context.Context) {
	if err := psfu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (psfu *PropertySchemaFieldUpdate) check() error {
	if psfu.mutation.SchemaCleared() && len(psfu.mutation.SchemaIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "PropertySchemaField.schema"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (psfu *PropertySchemaFieldUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *PropertySchemaFieldUpdate {
	psfu.modifiers = append(psfu.modifiers, modifiers...)
	return psfu
}

func (psfu *PropertySchemaFieldUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := psfu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(propertyschemafield.Table, propertyschemafield.Columns, sqlgraph.NewFieldSpec(propertyschemafield.FieldID, field.TypeString))
	if ps := psfu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := psfu.mutation.Name(); ok {
		_spec.SetField(propertyschemafield.FieldName, field.TypeString, value)
	}
	if value, ok := psfu.mutation.GetType(); ok {
		_spec.SetField(propertyschemafield.FieldType, field.TypeString, value)
	}
	if value, ok := psfu.mutation.Sort(); ok {
		_spec.SetField(propertyschemafield.FieldSort, field.TypeString, value)
	}
	if psfu.mutation.SchemaCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   propertyschemafield.SchemaTable,
			Columns: []string{propertyschemafield.SchemaColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(propertyschema.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := psfu.mutation.SchemaIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   propertyschemafield.SchemaTable,
			Columns: []string{propertyschemafield.SchemaColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(propertyschema.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if psfu.mutation.PropertiesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   propertyschemafield.PropertiesTable,
			Columns: []string{propertyschemafield.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(property.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := psfu.mutation.RemovedPropertiesIDs(); len(nodes) > 0 && !psfu.mutation.PropertiesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   propertyschemafield.PropertiesTable,
			Columns: []string{propertyschemafield.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(property.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := psfu.mutation.PropertiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   propertyschemafield.PropertiesTable,
			Columns: []string{propertyschemafield.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(property.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(psfu.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, psfu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{propertyschemafield.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	psfu.mutation.done = true
	return n, nil
}

// PropertySchemaFieldUpdateOne is the builder for updating a single PropertySchemaField entity.
type PropertySchemaFieldUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *PropertySchemaFieldMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetName sets the "name" field.
func (psfuo *PropertySchemaFieldUpdateOne) SetName(s string) *PropertySchemaFieldUpdateOne {
	psfuo.mutation.SetName(s)
	return psfuo
}

// SetNillableName sets the "name" field if the given value is not nil.
func (psfuo *PropertySchemaFieldUpdateOne) SetNillableName(s *string) *PropertySchemaFieldUpdateOne {
	if s != nil {
		psfuo.SetName(*s)
	}
	return psfuo
}

// SetType sets the "type" field.
func (psfuo *PropertySchemaFieldUpdateOne) SetType(s string) *PropertySchemaFieldUpdateOne {
	psfuo.mutation.SetType(s)
	return psfuo
}

// SetNillableType sets the "type" field if the given value is not nil.
func (psfuo *PropertySchemaFieldUpdateOne) SetNillableType(s *string) *PropertySchemaFieldUpdateOne {
	if s != nil {
		psfuo.SetType(*s)
	}
	return psfuo
}

// SetSort sets the "sort" field.
func (psfuo *PropertySchemaFieldUpdateOne) SetSort(s string) *PropertySchemaFieldUpdateOne {
	psfuo.mutation.SetSort(s)
	return psfuo
}

// SetNillableSort sets the "sort" field if the given value is not nil.
func (psfuo *PropertySchemaFieldUpdateOne) SetNillableSort(s *string) *PropertySchemaFieldUpdateOne {
	if s != nil {
		psfuo.SetSort(*s)
	}
	return psfuo
}

// SetSchemaID sets the "schema_id" field.
func (psfuo *PropertySchemaFieldUpdateOne) SetSchemaID(x xid.ID) *PropertySchemaFieldUpdateOne {
	psfuo.mutation.SetSchemaID(x)
	return psfuo
}

// SetNillableSchemaID sets the "schema_id" field if the given value is not nil.
func (psfuo *PropertySchemaFieldUpdateOne) SetNillableSchemaID(x *xid.ID) *PropertySchemaFieldUpdateOne {
	if x != nil {
		psfuo.SetSchemaID(*x)
	}
	return psfuo
}

// SetSchema sets the "schema" edge to the PropertySchema entity.
func (psfuo *PropertySchemaFieldUpdateOne) SetSchema(p *PropertySchema) *PropertySchemaFieldUpdateOne {
	return psfuo.SetSchemaID(p.ID)
}

// AddPropertyIDs adds the "properties" edge to the Property entity by IDs.
func (psfuo *PropertySchemaFieldUpdateOne) AddPropertyIDs(ids ...xid.ID) *PropertySchemaFieldUpdateOne {
	psfuo.mutation.AddPropertyIDs(ids...)
	return psfuo
}

// AddProperties adds the "properties" edges to the Property entity.
func (psfuo *PropertySchemaFieldUpdateOne) AddProperties(p ...*Property) *PropertySchemaFieldUpdateOne {
	ids := make([]xid.ID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return psfuo.AddPropertyIDs(ids...)
}

// Mutation returns the PropertySchemaFieldMutation object of the builder.
func (psfuo *PropertySchemaFieldUpdateOne) Mutation() *PropertySchemaFieldMutation {
	return psfuo.mutation
}

// ClearSchema clears the "schema" edge to the PropertySchema entity.
func (psfuo *PropertySchemaFieldUpdateOne) ClearSchema() *PropertySchemaFieldUpdateOne {
	psfuo.mutation.ClearSchema()
	return psfuo
}

// ClearProperties clears all "properties" edges to the Property entity.
func (psfuo *PropertySchemaFieldUpdateOne) ClearProperties() *PropertySchemaFieldUpdateOne {
	psfuo.mutation.ClearProperties()
	return psfuo
}

// RemovePropertyIDs removes the "properties" edge to Property entities by IDs.
func (psfuo *PropertySchemaFieldUpdateOne) RemovePropertyIDs(ids ...xid.ID) *PropertySchemaFieldUpdateOne {
	psfuo.mutation.RemovePropertyIDs(ids...)
	return psfuo
}

// RemoveProperties removes "properties" edges to Property entities.
func (psfuo *PropertySchemaFieldUpdateOne) RemoveProperties(p ...*Property) *PropertySchemaFieldUpdateOne {
	ids := make([]xid.ID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return psfuo.RemovePropertyIDs(ids...)
}

// Where appends a list predicates to the PropertySchemaFieldUpdate builder.
func (psfuo *PropertySchemaFieldUpdateOne) Where(ps ...predicate.PropertySchemaField) *PropertySchemaFieldUpdateOne {
	psfuo.mutation.Where(ps...)
	return psfuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (psfuo *PropertySchemaFieldUpdateOne) Select(field string, fields ...string) *PropertySchemaFieldUpdateOne {
	psfuo.fields = append([]string{field}, fields...)
	return psfuo
}

// Save executes the query and returns the updated PropertySchemaField entity.
func (psfuo *PropertySchemaFieldUpdateOne) Save(ctx context.Context) (*PropertySchemaField, error) {
	return withHooks(ctx, psfuo.sqlSave, psfuo.mutation, psfuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (psfuo *PropertySchemaFieldUpdateOne) SaveX(ctx context.Context) *PropertySchemaField {
	node, err := psfuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (psfuo *PropertySchemaFieldUpdateOne) Exec(ctx context.Context) error {
	_, err := psfuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (psfuo *PropertySchemaFieldUpdateOne) ExecX(ctx context.Context) {
	if err := psfuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (psfuo *PropertySchemaFieldUpdateOne) check() error {
	if psfuo.mutation.SchemaCleared() && len(psfuo.mutation.SchemaIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "PropertySchemaField.schema"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (psfuo *PropertySchemaFieldUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *PropertySchemaFieldUpdateOne {
	psfuo.modifiers = append(psfuo.modifiers, modifiers...)
	return psfuo
}

func (psfuo *PropertySchemaFieldUpdateOne) sqlSave(ctx context.Context) (_node *PropertySchemaField, err error) {
	if err := psfuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(propertyschemafield.Table, propertyschemafield.Columns, sqlgraph.NewFieldSpec(propertyschemafield.FieldID, field.TypeString))
	id, ok := psfuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "PropertySchemaField.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := psfuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, propertyschemafield.FieldID)
		for _, f := range fields {
			if !propertyschemafield.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != propertyschemafield.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := psfuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := psfuo.mutation.Name(); ok {
		_spec.SetField(propertyschemafield.FieldName, field.TypeString, value)
	}
	if value, ok := psfuo.mutation.GetType(); ok {
		_spec.SetField(propertyschemafield.FieldType, field.TypeString, value)
	}
	if value, ok := psfuo.mutation.Sort(); ok {
		_spec.SetField(propertyschemafield.FieldSort, field.TypeString, value)
	}
	if psfuo.mutation.SchemaCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   propertyschemafield.SchemaTable,
			Columns: []string{propertyschemafield.SchemaColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(propertyschema.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := psfuo.mutation.SchemaIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   propertyschemafield.SchemaTable,
			Columns: []string{propertyschemafield.SchemaColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(propertyschema.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if psfuo.mutation.PropertiesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   propertyschemafield.PropertiesTable,
			Columns: []string{propertyschemafield.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(property.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := psfuo.mutation.RemovedPropertiesIDs(); len(nodes) > 0 && !psfuo.mutation.PropertiesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   propertyschemafield.PropertiesTable,
			Columns: []string{propertyschemafield.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(property.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := psfuo.mutation.PropertiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   propertyschemafield.PropertiesTable,
			Columns: []string{propertyschemafield.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(property.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(psfuo.modifiers...)
	_node = &PropertySchemaField{config: psfuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, psfuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{propertyschemafield.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	psfuo.mutation.done = true
	return _node, nil
}
