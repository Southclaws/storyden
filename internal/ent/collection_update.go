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
	"github.com/Southclaws/storyden/internal/ent/collection"
	"github.com/Southclaws/storyden/internal/ent/post"
	"github.com/Southclaws/storyden/internal/ent/predicate"
	"github.com/rs/xid"
)

// CollectionUpdate is the builder for updating Collection entities.
type CollectionUpdate struct {
	config
	hooks     []Hook
	mutation  *CollectionMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the CollectionUpdate builder.
func (cu *CollectionUpdate) Where(ps ...predicate.Collection) *CollectionUpdate {
	cu.mutation.Where(ps...)
	return cu
}

// SetUpdatedAt sets the "updated_at" field.
func (cu *CollectionUpdate) SetUpdatedAt(t time.Time) *CollectionUpdate {
	cu.mutation.SetUpdatedAt(t)
	return cu
}

// SetName sets the "name" field.
func (cu *CollectionUpdate) SetName(s string) *CollectionUpdate {
	cu.mutation.SetName(s)
	return cu
}

// SetDescription sets the "description" field.
func (cu *CollectionUpdate) SetDescription(s string) *CollectionUpdate {
	cu.mutation.SetDescription(s)
	return cu
}

// SetOwnerID sets the "owner" edge to the Account entity by ID.
func (cu *CollectionUpdate) SetOwnerID(id xid.ID) *CollectionUpdate {
	cu.mutation.SetOwnerID(id)
	return cu
}

// SetNillableOwnerID sets the "owner" edge to the Account entity by ID if the given value is not nil.
func (cu *CollectionUpdate) SetNillableOwnerID(id *xid.ID) *CollectionUpdate {
	if id != nil {
		cu = cu.SetOwnerID(*id)
	}
	return cu
}

// SetOwner sets the "owner" edge to the Account entity.
func (cu *CollectionUpdate) SetOwner(a *Account) *CollectionUpdate {
	return cu.SetOwnerID(a.ID)
}

// AddPostIDs adds the "posts" edge to the Post entity by IDs.
func (cu *CollectionUpdate) AddPostIDs(ids ...xid.ID) *CollectionUpdate {
	cu.mutation.AddPostIDs(ids...)
	return cu
}

// AddPosts adds the "posts" edges to the Post entity.
func (cu *CollectionUpdate) AddPosts(p ...*Post) *CollectionUpdate {
	ids := make([]xid.ID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return cu.AddPostIDs(ids...)
}

// Mutation returns the CollectionMutation object of the builder.
func (cu *CollectionUpdate) Mutation() *CollectionMutation {
	return cu.mutation
}

// ClearOwner clears the "owner" edge to the Account entity.
func (cu *CollectionUpdate) ClearOwner() *CollectionUpdate {
	cu.mutation.ClearOwner()
	return cu
}

// ClearPosts clears all "posts" edges to the Post entity.
func (cu *CollectionUpdate) ClearPosts() *CollectionUpdate {
	cu.mutation.ClearPosts()
	return cu
}

// RemovePostIDs removes the "posts" edge to Post entities by IDs.
func (cu *CollectionUpdate) RemovePostIDs(ids ...xid.ID) *CollectionUpdate {
	cu.mutation.RemovePostIDs(ids...)
	return cu
}

// RemovePosts removes "posts" edges to Post entities.
func (cu *CollectionUpdate) RemovePosts(p ...*Post) *CollectionUpdate {
	ids := make([]xid.ID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return cu.RemovePostIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (cu *CollectionUpdate) Save(ctx context.Context) (int, error) {
	cu.defaults()
	return withHooks(ctx, cu.sqlSave, cu.mutation, cu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (cu *CollectionUpdate) SaveX(ctx context.Context) int {
	affected, err := cu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (cu *CollectionUpdate) Exec(ctx context.Context) error {
	_, err := cu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cu *CollectionUpdate) ExecX(ctx context.Context) {
	if err := cu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (cu *CollectionUpdate) defaults() {
	if _, ok := cu.mutation.UpdatedAt(); !ok {
		v := collection.UpdateDefaultUpdatedAt()
		cu.mutation.SetUpdatedAt(v)
	}
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (cu *CollectionUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *CollectionUpdate {
	cu.modifiers = append(cu.modifiers, modifiers...)
	return cu
}

func (cu *CollectionUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(collection.Table, collection.Columns, sqlgraph.NewFieldSpec(collection.FieldID, field.TypeString))
	if ps := cu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := cu.mutation.UpdatedAt(); ok {
		_spec.SetField(collection.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := cu.mutation.Name(); ok {
		_spec.SetField(collection.FieldName, field.TypeString, value)
	}
	if value, ok := cu.mutation.Description(); ok {
		_spec.SetField(collection.FieldDescription, field.TypeString, value)
	}
	if cu.mutation.OwnerCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   collection.OwnerTable,
			Columns: []string{collection.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(account.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cu.mutation.OwnerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   collection.OwnerTable,
			Columns: []string{collection.OwnerColumn},
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
	if cu.mutation.PostsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   collection.PostsTable,
			Columns: collection.PostsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(post.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cu.mutation.RemovedPostsIDs(); len(nodes) > 0 && !cu.mutation.PostsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   collection.PostsTable,
			Columns: collection.PostsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(post.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cu.mutation.PostsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   collection.PostsTable,
			Columns: collection.PostsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(post.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(cu.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, cu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{collection.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	cu.mutation.done = true
	return n, nil
}

// CollectionUpdateOne is the builder for updating a single Collection entity.
type CollectionUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *CollectionMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetUpdatedAt sets the "updated_at" field.
func (cuo *CollectionUpdateOne) SetUpdatedAt(t time.Time) *CollectionUpdateOne {
	cuo.mutation.SetUpdatedAt(t)
	return cuo
}

// SetName sets the "name" field.
func (cuo *CollectionUpdateOne) SetName(s string) *CollectionUpdateOne {
	cuo.mutation.SetName(s)
	return cuo
}

// SetDescription sets the "description" field.
func (cuo *CollectionUpdateOne) SetDescription(s string) *CollectionUpdateOne {
	cuo.mutation.SetDescription(s)
	return cuo
}

// SetOwnerID sets the "owner" edge to the Account entity by ID.
func (cuo *CollectionUpdateOne) SetOwnerID(id xid.ID) *CollectionUpdateOne {
	cuo.mutation.SetOwnerID(id)
	return cuo
}

// SetNillableOwnerID sets the "owner" edge to the Account entity by ID if the given value is not nil.
func (cuo *CollectionUpdateOne) SetNillableOwnerID(id *xid.ID) *CollectionUpdateOne {
	if id != nil {
		cuo = cuo.SetOwnerID(*id)
	}
	return cuo
}

// SetOwner sets the "owner" edge to the Account entity.
func (cuo *CollectionUpdateOne) SetOwner(a *Account) *CollectionUpdateOne {
	return cuo.SetOwnerID(a.ID)
}

// AddPostIDs adds the "posts" edge to the Post entity by IDs.
func (cuo *CollectionUpdateOne) AddPostIDs(ids ...xid.ID) *CollectionUpdateOne {
	cuo.mutation.AddPostIDs(ids...)
	return cuo
}

// AddPosts adds the "posts" edges to the Post entity.
func (cuo *CollectionUpdateOne) AddPosts(p ...*Post) *CollectionUpdateOne {
	ids := make([]xid.ID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return cuo.AddPostIDs(ids...)
}

// Mutation returns the CollectionMutation object of the builder.
func (cuo *CollectionUpdateOne) Mutation() *CollectionMutation {
	return cuo.mutation
}

// ClearOwner clears the "owner" edge to the Account entity.
func (cuo *CollectionUpdateOne) ClearOwner() *CollectionUpdateOne {
	cuo.mutation.ClearOwner()
	return cuo
}

// ClearPosts clears all "posts" edges to the Post entity.
func (cuo *CollectionUpdateOne) ClearPosts() *CollectionUpdateOne {
	cuo.mutation.ClearPosts()
	return cuo
}

// RemovePostIDs removes the "posts" edge to Post entities by IDs.
func (cuo *CollectionUpdateOne) RemovePostIDs(ids ...xid.ID) *CollectionUpdateOne {
	cuo.mutation.RemovePostIDs(ids...)
	return cuo
}

// RemovePosts removes "posts" edges to Post entities.
func (cuo *CollectionUpdateOne) RemovePosts(p ...*Post) *CollectionUpdateOne {
	ids := make([]xid.ID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return cuo.RemovePostIDs(ids...)
}

// Where appends a list predicates to the CollectionUpdate builder.
func (cuo *CollectionUpdateOne) Where(ps ...predicate.Collection) *CollectionUpdateOne {
	cuo.mutation.Where(ps...)
	return cuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (cuo *CollectionUpdateOne) Select(field string, fields ...string) *CollectionUpdateOne {
	cuo.fields = append([]string{field}, fields...)
	return cuo
}

// Save executes the query and returns the updated Collection entity.
func (cuo *CollectionUpdateOne) Save(ctx context.Context) (*Collection, error) {
	cuo.defaults()
	return withHooks(ctx, cuo.sqlSave, cuo.mutation, cuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (cuo *CollectionUpdateOne) SaveX(ctx context.Context) *Collection {
	node, err := cuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (cuo *CollectionUpdateOne) Exec(ctx context.Context) error {
	_, err := cuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cuo *CollectionUpdateOne) ExecX(ctx context.Context) {
	if err := cuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (cuo *CollectionUpdateOne) defaults() {
	if _, ok := cuo.mutation.UpdatedAt(); !ok {
		v := collection.UpdateDefaultUpdatedAt()
		cuo.mutation.SetUpdatedAt(v)
	}
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (cuo *CollectionUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *CollectionUpdateOne {
	cuo.modifiers = append(cuo.modifiers, modifiers...)
	return cuo
}

func (cuo *CollectionUpdateOne) sqlSave(ctx context.Context) (_node *Collection, err error) {
	_spec := sqlgraph.NewUpdateSpec(collection.Table, collection.Columns, sqlgraph.NewFieldSpec(collection.FieldID, field.TypeString))
	id, ok := cuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Collection.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := cuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, collection.FieldID)
		for _, f := range fields {
			if !collection.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != collection.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := cuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := cuo.mutation.UpdatedAt(); ok {
		_spec.SetField(collection.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := cuo.mutation.Name(); ok {
		_spec.SetField(collection.FieldName, field.TypeString, value)
	}
	if value, ok := cuo.mutation.Description(); ok {
		_spec.SetField(collection.FieldDescription, field.TypeString, value)
	}
	if cuo.mutation.OwnerCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   collection.OwnerTable,
			Columns: []string{collection.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(account.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cuo.mutation.OwnerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   collection.OwnerTable,
			Columns: []string{collection.OwnerColumn},
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
	if cuo.mutation.PostsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   collection.PostsTable,
			Columns: collection.PostsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(post.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cuo.mutation.RemovedPostsIDs(); len(nodes) > 0 && !cuo.mutation.PostsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   collection.PostsTable,
			Columns: collection.PostsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(post.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cuo.mutation.PostsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   collection.PostsTable,
			Columns: collection.PostsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(post.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(cuo.modifiers...)
	_node = &Collection{config: cuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, cuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{collection.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	cuo.mutation.done = true
	return _node, nil
}
