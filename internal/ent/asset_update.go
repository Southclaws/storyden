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
	"github.com/Southclaws/storyden/internal/ent/asset"
	"github.com/Southclaws/storyden/internal/ent/post"
	"github.com/Southclaws/storyden/internal/ent/predicate"
	"github.com/rs/xid"
)

// AssetUpdate is the builder for updating Asset entities.
type AssetUpdate struct {
	config
	hooks     []Hook
	mutation  *AssetMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the AssetUpdate builder.
func (au *AssetUpdate) Where(ps ...predicate.Asset) *AssetUpdate {
	au.mutation.Where(ps...)
	return au
}

// SetUpdatedAt sets the "updated_at" field.
func (au *AssetUpdate) SetUpdatedAt(t time.Time) *AssetUpdate {
	au.mutation.SetUpdatedAt(t)
	return au
}

// SetURL sets the "url" field.
func (au *AssetUpdate) SetURL(s string) *AssetUpdate {
	au.mutation.SetURL(s)
	return au
}

// SetMimetype sets the "mimetype" field.
func (au *AssetUpdate) SetMimetype(s string) *AssetUpdate {
	au.mutation.SetMimetype(s)
	return au
}

// SetWidth sets the "width" field.
func (au *AssetUpdate) SetWidth(i int) *AssetUpdate {
	au.mutation.ResetWidth()
	au.mutation.SetWidth(i)
	return au
}

// AddWidth adds i to the "width" field.
func (au *AssetUpdate) AddWidth(i int) *AssetUpdate {
	au.mutation.AddWidth(i)
	return au
}

// SetHeight sets the "height" field.
func (au *AssetUpdate) SetHeight(i int) *AssetUpdate {
	au.mutation.ResetHeight()
	au.mutation.SetHeight(i)
	return au
}

// AddHeight adds i to the "height" field.
func (au *AssetUpdate) AddHeight(i int) *AssetUpdate {
	au.mutation.AddHeight(i)
	return au
}

// SetAccountID sets the "account_id" field.
func (au *AssetUpdate) SetAccountID(x xid.ID) *AssetUpdate {
	au.mutation.SetAccountID(x)
	return au
}

// AddPostIDs adds the "posts" edge to the Post entity by IDs.
func (au *AssetUpdate) AddPostIDs(ids ...xid.ID) *AssetUpdate {
	au.mutation.AddPostIDs(ids...)
	return au
}

// AddPosts adds the "posts" edges to the Post entity.
func (au *AssetUpdate) AddPosts(p ...*Post) *AssetUpdate {
	ids := make([]xid.ID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return au.AddPostIDs(ids...)
}

// SetOwnerID sets the "owner" edge to the Account entity by ID.
func (au *AssetUpdate) SetOwnerID(id xid.ID) *AssetUpdate {
	au.mutation.SetOwnerID(id)
	return au
}

// SetOwner sets the "owner" edge to the Account entity.
func (au *AssetUpdate) SetOwner(a *Account) *AssetUpdate {
	return au.SetOwnerID(a.ID)
}

// Mutation returns the AssetMutation object of the builder.
func (au *AssetUpdate) Mutation() *AssetMutation {
	return au.mutation
}

// ClearPosts clears all "posts" edges to the Post entity.
func (au *AssetUpdate) ClearPosts() *AssetUpdate {
	au.mutation.ClearPosts()
	return au
}

// RemovePostIDs removes the "posts" edge to Post entities by IDs.
func (au *AssetUpdate) RemovePostIDs(ids ...xid.ID) *AssetUpdate {
	au.mutation.RemovePostIDs(ids...)
	return au
}

// RemovePosts removes "posts" edges to Post entities.
func (au *AssetUpdate) RemovePosts(p ...*Post) *AssetUpdate {
	ids := make([]xid.ID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return au.RemovePostIDs(ids...)
}

// ClearOwner clears the "owner" edge to the Account entity.
func (au *AssetUpdate) ClearOwner() *AssetUpdate {
	au.mutation.ClearOwner()
	return au
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (au *AssetUpdate) Save(ctx context.Context) (int, error) {
	au.defaults()
	return withHooks[int, AssetMutation](ctx, au.sqlSave, au.mutation, au.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (au *AssetUpdate) SaveX(ctx context.Context) int {
	affected, err := au.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (au *AssetUpdate) Exec(ctx context.Context) error {
	_, err := au.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (au *AssetUpdate) ExecX(ctx context.Context) {
	if err := au.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (au *AssetUpdate) defaults() {
	if _, ok := au.mutation.UpdatedAt(); !ok {
		v := asset.UpdateDefaultUpdatedAt()
		au.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (au *AssetUpdate) check() error {
	if _, ok := au.mutation.OwnerID(); au.mutation.OwnerCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "Asset.owner"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (au *AssetUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *AssetUpdate {
	au.modifiers = append(au.modifiers, modifiers...)
	return au
}

func (au *AssetUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := au.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(asset.Table, asset.Columns, sqlgraph.NewFieldSpec(asset.FieldID, field.TypeString))
	if ps := au.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := au.mutation.UpdatedAt(); ok {
		_spec.SetField(asset.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := au.mutation.URL(); ok {
		_spec.SetField(asset.FieldURL, field.TypeString, value)
	}
	if value, ok := au.mutation.Mimetype(); ok {
		_spec.SetField(asset.FieldMimetype, field.TypeString, value)
	}
	if value, ok := au.mutation.Width(); ok {
		_spec.SetField(asset.FieldWidth, field.TypeInt, value)
	}
	if value, ok := au.mutation.AddedWidth(); ok {
		_spec.AddField(asset.FieldWidth, field.TypeInt, value)
	}
	if value, ok := au.mutation.Height(); ok {
		_spec.SetField(asset.FieldHeight, field.TypeInt, value)
	}
	if value, ok := au.mutation.AddedHeight(); ok {
		_spec.AddField(asset.FieldHeight, field.TypeInt, value)
	}
	if au.mutation.PostsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   asset.PostsTable,
			Columns: asset.PostsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(post.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := au.mutation.RemovedPostsIDs(); len(nodes) > 0 && !au.mutation.PostsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   asset.PostsTable,
			Columns: asset.PostsPrimaryKey,
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
	if nodes := au.mutation.PostsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   asset.PostsTable,
			Columns: asset.PostsPrimaryKey,
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
	if au.mutation.OwnerCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   asset.OwnerTable,
			Columns: []string{asset.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(account.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := au.mutation.OwnerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   asset.OwnerTable,
			Columns: []string{asset.OwnerColumn},
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
	_spec.AddModifiers(au.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, au.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{asset.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	au.mutation.done = true
	return n, nil
}

// AssetUpdateOne is the builder for updating a single Asset entity.
type AssetUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *AssetMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetUpdatedAt sets the "updated_at" field.
func (auo *AssetUpdateOne) SetUpdatedAt(t time.Time) *AssetUpdateOne {
	auo.mutation.SetUpdatedAt(t)
	return auo
}

// SetURL sets the "url" field.
func (auo *AssetUpdateOne) SetURL(s string) *AssetUpdateOne {
	auo.mutation.SetURL(s)
	return auo
}

// SetMimetype sets the "mimetype" field.
func (auo *AssetUpdateOne) SetMimetype(s string) *AssetUpdateOne {
	auo.mutation.SetMimetype(s)
	return auo
}

// SetWidth sets the "width" field.
func (auo *AssetUpdateOne) SetWidth(i int) *AssetUpdateOne {
	auo.mutation.ResetWidth()
	auo.mutation.SetWidth(i)
	return auo
}

// AddWidth adds i to the "width" field.
func (auo *AssetUpdateOne) AddWidth(i int) *AssetUpdateOne {
	auo.mutation.AddWidth(i)
	return auo
}

// SetHeight sets the "height" field.
func (auo *AssetUpdateOne) SetHeight(i int) *AssetUpdateOne {
	auo.mutation.ResetHeight()
	auo.mutation.SetHeight(i)
	return auo
}

// AddHeight adds i to the "height" field.
func (auo *AssetUpdateOne) AddHeight(i int) *AssetUpdateOne {
	auo.mutation.AddHeight(i)
	return auo
}

// SetAccountID sets the "account_id" field.
func (auo *AssetUpdateOne) SetAccountID(x xid.ID) *AssetUpdateOne {
	auo.mutation.SetAccountID(x)
	return auo
}

// AddPostIDs adds the "posts" edge to the Post entity by IDs.
func (auo *AssetUpdateOne) AddPostIDs(ids ...xid.ID) *AssetUpdateOne {
	auo.mutation.AddPostIDs(ids...)
	return auo
}

// AddPosts adds the "posts" edges to the Post entity.
func (auo *AssetUpdateOne) AddPosts(p ...*Post) *AssetUpdateOne {
	ids := make([]xid.ID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return auo.AddPostIDs(ids...)
}

// SetOwnerID sets the "owner" edge to the Account entity by ID.
func (auo *AssetUpdateOne) SetOwnerID(id xid.ID) *AssetUpdateOne {
	auo.mutation.SetOwnerID(id)
	return auo
}

// SetOwner sets the "owner" edge to the Account entity.
func (auo *AssetUpdateOne) SetOwner(a *Account) *AssetUpdateOne {
	return auo.SetOwnerID(a.ID)
}

// Mutation returns the AssetMutation object of the builder.
func (auo *AssetUpdateOne) Mutation() *AssetMutation {
	return auo.mutation
}

// ClearPosts clears all "posts" edges to the Post entity.
func (auo *AssetUpdateOne) ClearPosts() *AssetUpdateOne {
	auo.mutation.ClearPosts()
	return auo
}

// RemovePostIDs removes the "posts" edge to Post entities by IDs.
func (auo *AssetUpdateOne) RemovePostIDs(ids ...xid.ID) *AssetUpdateOne {
	auo.mutation.RemovePostIDs(ids...)
	return auo
}

// RemovePosts removes "posts" edges to Post entities.
func (auo *AssetUpdateOne) RemovePosts(p ...*Post) *AssetUpdateOne {
	ids := make([]xid.ID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return auo.RemovePostIDs(ids...)
}

// ClearOwner clears the "owner" edge to the Account entity.
func (auo *AssetUpdateOne) ClearOwner() *AssetUpdateOne {
	auo.mutation.ClearOwner()
	return auo
}

// Where appends a list predicates to the AssetUpdate builder.
func (auo *AssetUpdateOne) Where(ps ...predicate.Asset) *AssetUpdateOne {
	auo.mutation.Where(ps...)
	return auo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (auo *AssetUpdateOne) Select(field string, fields ...string) *AssetUpdateOne {
	auo.fields = append([]string{field}, fields...)
	return auo
}

// Save executes the query and returns the updated Asset entity.
func (auo *AssetUpdateOne) Save(ctx context.Context) (*Asset, error) {
	auo.defaults()
	return withHooks[*Asset, AssetMutation](ctx, auo.sqlSave, auo.mutation, auo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (auo *AssetUpdateOne) SaveX(ctx context.Context) *Asset {
	node, err := auo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (auo *AssetUpdateOne) Exec(ctx context.Context) error {
	_, err := auo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (auo *AssetUpdateOne) ExecX(ctx context.Context) {
	if err := auo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (auo *AssetUpdateOne) defaults() {
	if _, ok := auo.mutation.UpdatedAt(); !ok {
		v := asset.UpdateDefaultUpdatedAt()
		auo.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (auo *AssetUpdateOne) check() error {
	if _, ok := auo.mutation.OwnerID(); auo.mutation.OwnerCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "Asset.owner"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (auo *AssetUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *AssetUpdateOne {
	auo.modifiers = append(auo.modifiers, modifiers...)
	return auo
}

func (auo *AssetUpdateOne) sqlSave(ctx context.Context) (_node *Asset, err error) {
	if err := auo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(asset.Table, asset.Columns, sqlgraph.NewFieldSpec(asset.FieldID, field.TypeString))
	id, ok := auo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Asset.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := auo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, asset.FieldID)
		for _, f := range fields {
			if !asset.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != asset.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := auo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := auo.mutation.UpdatedAt(); ok {
		_spec.SetField(asset.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := auo.mutation.URL(); ok {
		_spec.SetField(asset.FieldURL, field.TypeString, value)
	}
	if value, ok := auo.mutation.Mimetype(); ok {
		_spec.SetField(asset.FieldMimetype, field.TypeString, value)
	}
	if value, ok := auo.mutation.Width(); ok {
		_spec.SetField(asset.FieldWidth, field.TypeInt, value)
	}
	if value, ok := auo.mutation.AddedWidth(); ok {
		_spec.AddField(asset.FieldWidth, field.TypeInt, value)
	}
	if value, ok := auo.mutation.Height(); ok {
		_spec.SetField(asset.FieldHeight, field.TypeInt, value)
	}
	if value, ok := auo.mutation.AddedHeight(); ok {
		_spec.AddField(asset.FieldHeight, field.TypeInt, value)
	}
	if auo.mutation.PostsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   asset.PostsTable,
			Columns: asset.PostsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(post.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := auo.mutation.RemovedPostsIDs(); len(nodes) > 0 && !auo.mutation.PostsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   asset.PostsTable,
			Columns: asset.PostsPrimaryKey,
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
	if nodes := auo.mutation.PostsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   asset.PostsTable,
			Columns: asset.PostsPrimaryKey,
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
	if auo.mutation.OwnerCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   asset.OwnerTable,
			Columns: []string{asset.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(account.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := auo.mutation.OwnerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   asset.OwnerTable,
			Columns: []string{asset.OwnerColumn},
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
	_spec.AddModifiers(auo.modifiers...)
	_node = &Asset{config: auo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, auo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{asset.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	auo.mutation.done = true
	return _node, nil
}
