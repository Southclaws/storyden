// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/Southclaws/storyden/internal/ent/account"
	"github.com/Southclaws/storyden/internal/ent/asset"
	"github.com/Southclaws/storyden/internal/ent/collection"
	"github.com/Southclaws/storyden/internal/ent/link"
	"github.com/Southclaws/storyden/internal/ent/node"
	"github.com/Southclaws/storyden/internal/ent/tag"
	"github.com/rs/xid"
)

// NodeCreate is the builder for creating a Node entity.
type NodeCreate struct {
	config
	mutation *NodeMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCreatedAt sets the "created_at" field.
func (nc *NodeCreate) SetCreatedAt(t time.Time) *NodeCreate {
	nc.mutation.SetCreatedAt(t)
	return nc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (nc *NodeCreate) SetNillableCreatedAt(t *time.Time) *NodeCreate {
	if t != nil {
		nc.SetCreatedAt(*t)
	}
	return nc
}

// SetUpdatedAt sets the "updated_at" field.
func (nc *NodeCreate) SetUpdatedAt(t time.Time) *NodeCreate {
	nc.mutation.SetUpdatedAt(t)
	return nc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (nc *NodeCreate) SetNillableUpdatedAt(t *time.Time) *NodeCreate {
	if t != nil {
		nc.SetUpdatedAt(*t)
	}
	return nc
}

// SetDeletedAt sets the "deleted_at" field.
func (nc *NodeCreate) SetDeletedAt(t time.Time) *NodeCreate {
	nc.mutation.SetDeletedAt(t)
	return nc
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (nc *NodeCreate) SetNillableDeletedAt(t *time.Time) *NodeCreate {
	if t != nil {
		nc.SetDeletedAt(*t)
	}
	return nc
}

// SetName sets the "name" field.
func (nc *NodeCreate) SetName(s string) *NodeCreate {
	nc.mutation.SetName(s)
	return nc
}

// SetSlug sets the "slug" field.
func (nc *NodeCreate) SetSlug(s string) *NodeCreate {
	nc.mutation.SetSlug(s)
	return nc
}

// SetDescription sets the "description" field.
func (nc *NodeCreate) SetDescription(s string) *NodeCreate {
	nc.mutation.SetDescription(s)
	return nc
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (nc *NodeCreate) SetNillableDescription(s *string) *NodeCreate {
	if s != nil {
		nc.SetDescription(*s)
	}
	return nc
}

// SetContent sets the "content" field.
func (nc *NodeCreate) SetContent(s string) *NodeCreate {
	nc.mutation.SetContent(s)
	return nc
}

// SetNillableContent sets the "content" field if the given value is not nil.
func (nc *NodeCreate) SetNillableContent(s *string) *NodeCreate {
	if s != nil {
		nc.SetContent(*s)
	}
	return nc
}

// SetParentNodeID sets the "parent_node_id" field.
func (nc *NodeCreate) SetParentNodeID(x xid.ID) *NodeCreate {
	nc.mutation.SetParentNodeID(x)
	return nc
}

// SetNillableParentNodeID sets the "parent_node_id" field if the given value is not nil.
func (nc *NodeCreate) SetNillableParentNodeID(x *xid.ID) *NodeCreate {
	if x != nil {
		nc.SetParentNodeID(*x)
	}
	return nc
}

// SetAccountID sets the "account_id" field.
func (nc *NodeCreate) SetAccountID(x xid.ID) *NodeCreate {
	nc.mutation.SetAccountID(x)
	return nc
}

// SetVisibility sets the "visibility" field.
func (nc *NodeCreate) SetVisibility(n node.Visibility) *NodeCreate {
	nc.mutation.SetVisibility(n)
	return nc
}

// SetNillableVisibility sets the "visibility" field if the given value is not nil.
func (nc *NodeCreate) SetNillableVisibility(n *node.Visibility) *NodeCreate {
	if n != nil {
		nc.SetVisibility(*n)
	}
	return nc
}

// SetProperties sets the "properties" field.
func (nc *NodeCreate) SetProperties(a any) *NodeCreate {
	nc.mutation.SetProperties(a)
	return nc
}

// SetID sets the "id" field.
func (nc *NodeCreate) SetID(x xid.ID) *NodeCreate {
	nc.mutation.SetID(x)
	return nc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (nc *NodeCreate) SetNillableID(x *xid.ID) *NodeCreate {
	if x != nil {
		nc.SetID(*x)
	}
	return nc
}

// SetOwnerID sets the "owner" edge to the Account entity by ID.
func (nc *NodeCreate) SetOwnerID(id xid.ID) *NodeCreate {
	nc.mutation.SetOwnerID(id)
	return nc
}

// SetOwner sets the "owner" edge to the Account entity.
func (nc *NodeCreate) SetOwner(a *Account) *NodeCreate {
	return nc.SetOwnerID(a.ID)
}

// SetParentID sets the "parent" edge to the Node entity by ID.
func (nc *NodeCreate) SetParentID(id xid.ID) *NodeCreate {
	nc.mutation.SetParentID(id)
	return nc
}

// SetNillableParentID sets the "parent" edge to the Node entity by ID if the given value is not nil.
func (nc *NodeCreate) SetNillableParentID(id *xid.ID) *NodeCreate {
	if id != nil {
		nc = nc.SetParentID(*id)
	}
	return nc
}

// SetParent sets the "parent" edge to the Node entity.
func (nc *NodeCreate) SetParent(n *Node) *NodeCreate {
	return nc.SetParentID(n.ID)
}

// AddNodeIDs adds the "nodes" edge to the Node entity by IDs.
func (nc *NodeCreate) AddNodeIDs(ids ...xid.ID) *NodeCreate {
	nc.mutation.AddNodeIDs(ids...)
	return nc
}

// AddNodes adds the "nodes" edges to the Node entity.
func (nc *NodeCreate) AddNodes(n ...*Node) *NodeCreate {
	ids := make([]xid.ID, len(n))
	for i := range n {
		ids[i] = n[i].ID
	}
	return nc.AddNodeIDs(ids...)
}

// AddAssetIDs adds the "assets" edge to the Asset entity by IDs.
func (nc *NodeCreate) AddAssetIDs(ids ...xid.ID) *NodeCreate {
	nc.mutation.AddAssetIDs(ids...)
	return nc
}

// AddAssets adds the "assets" edges to the Asset entity.
func (nc *NodeCreate) AddAssets(a ...*Asset) *NodeCreate {
	ids := make([]xid.ID, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return nc.AddAssetIDs(ids...)
}

// AddTagIDs adds the "tags" edge to the Tag entity by IDs.
func (nc *NodeCreate) AddTagIDs(ids ...xid.ID) *NodeCreate {
	nc.mutation.AddTagIDs(ids...)
	return nc
}

// AddTags adds the "tags" edges to the Tag entity.
func (nc *NodeCreate) AddTags(t ...*Tag) *NodeCreate {
	ids := make([]xid.ID, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return nc.AddTagIDs(ids...)
}

// AddLinkIDs adds the "links" edge to the Link entity by IDs.
func (nc *NodeCreate) AddLinkIDs(ids ...xid.ID) *NodeCreate {
	nc.mutation.AddLinkIDs(ids...)
	return nc
}

// AddLinks adds the "links" edges to the Link entity.
func (nc *NodeCreate) AddLinks(l ...*Link) *NodeCreate {
	ids := make([]xid.ID, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return nc.AddLinkIDs(ids...)
}

// AddCollectionIDs adds the "collections" edge to the Collection entity by IDs.
func (nc *NodeCreate) AddCollectionIDs(ids ...xid.ID) *NodeCreate {
	nc.mutation.AddCollectionIDs(ids...)
	return nc
}

// AddCollections adds the "collections" edges to the Collection entity.
func (nc *NodeCreate) AddCollections(c ...*Collection) *NodeCreate {
	ids := make([]xid.ID, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return nc.AddCollectionIDs(ids...)
}

// Mutation returns the NodeMutation object of the builder.
func (nc *NodeCreate) Mutation() *NodeMutation {
	return nc.mutation
}

// Save creates the Node in the database.
func (nc *NodeCreate) Save(ctx context.Context) (*Node, error) {
	nc.defaults()
	return withHooks(ctx, nc.sqlSave, nc.mutation, nc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (nc *NodeCreate) SaveX(ctx context.Context) *Node {
	v, err := nc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (nc *NodeCreate) Exec(ctx context.Context) error {
	_, err := nc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (nc *NodeCreate) ExecX(ctx context.Context) {
	if err := nc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (nc *NodeCreate) defaults() {
	if _, ok := nc.mutation.CreatedAt(); !ok {
		v := node.DefaultCreatedAt()
		nc.mutation.SetCreatedAt(v)
	}
	if _, ok := nc.mutation.UpdatedAt(); !ok {
		v := node.DefaultUpdatedAt()
		nc.mutation.SetUpdatedAt(v)
	}
	if _, ok := nc.mutation.Visibility(); !ok {
		v := node.DefaultVisibility
		nc.mutation.SetVisibility(v)
	}
	if _, ok := nc.mutation.ID(); !ok {
		v := node.DefaultID()
		nc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (nc *NodeCreate) check() error {
	if _, ok := nc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "Node.created_at"`)}
	}
	if _, ok := nc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "Node.updated_at"`)}
	}
	if _, ok := nc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "Node.name"`)}
	}
	if _, ok := nc.mutation.Slug(); !ok {
		return &ValidationError{Name: "slug", err: errors.New(`ent: missing required field "Node.slug"`)}
	}
	if _, ok := nc.mutation.AccountID(); !ok {
		return &ValidationError{Name: "account_id", err: errors.New(`ent: missing required field "Node.account_id"`)}
	}
	if _, ok := nc.mutation.Visibility(); !ok {
		return &ValidationError{Name: "visibility", err: errors.New(`ent: missing required field "Node.visibility"`)}
	}
	if v, ok := nc.mutation.Visibility(); ok {
		if err := node.VisibilityValidator(v); err != nil {
			return &ValidationError{Name: "visibility", err: fmt.Errorf(`ent: validator failed for field "Node.visibility": %w`, err)}
		}
	}
	if v, ok := nc.mutation.ID(); ok {
		if err := node.IDValidator(v.String()); err != nil {
			return &ValidationError{Name: "id", err: fmt.Errorf(`ent: validator failed for field "Node.id": %w`, err)}
		}
	}
	if _, ok := nc.mutation.OwnerID(); !ok {
		return &ValidationError{Name: "owner", err: errors.New(`ent: missing required edge "Node.owner"`)}
	}
	return nil
}

func (nc *NodeCreate) sqlSave(ctx context.Context) (*Node, error) {
	if err := nc.check(); err != nil {
		return nil, err
	}
	_node, _spec := nc.createSpec()
	if err := sqlgraph.CreateNode(ctx, nc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(*xid.ID); ok {
			_node.ID = *id
		} else if err := _node.ID.Scan(_spec.ID.Value); err != nil {
			return nil, err
		}
	}
	nc.mutation.id = &_node.ID
	nc.mutation.done = true
	return _node, nil
}

func (nc *NodeCreate) createSpec() (*Node, *sqlgraph.CreateSpec) {
	var (
		_node = &Node{config: nc.config}
		_spec = sqlgraph.NewCreateSpec(node.Table, sqlgraph.NewFieldSpec(node.FieldID, field.TypeString))
	)
	_spec.OnConflict = nc.conflict
	if id, ok := nc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := nc.mutation.CreatedAt(); ok {
		_spec.SetField(node.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := nc.mutation.UpdatedAt(); ok {
		_spec.SetField(node.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	if value, ok := nc.mutation.DeletedAt(); ok {
		_spec.SetField(node.FieldDeletedAt, field.TypeTime, value)
		_node.DeletedAt = &value
	}
	if value, ok := nc.mutation.Name(); ok {
		_spec.SetField(node.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := nc.mutation.Slug(); ok {
		_spec.SetField(node.FieldSlug, field.TypeString, value)
		_node.Slug = value
	}
	if value, ok := nc.mutation.Description(); ok {
		_spec.SetField(node.FieldDescription, field.TypeString, value)
		_node.Description = &value
	}
	if value, ok := nc.mutation.Content(); ok {
		_spec.SetField(node.FieldContent, field.TypeString, value)
		_node.Content = &value
	}
	if value, ok := nc.mutation.Visibility(); ok {
		_spec.SetField(node.FieldVisibility, field.TypeEnum, value)
		_node.Visibility = value
	}
	if value, ok := nc.mutation.Properties(); ok {
		_spec.SetField(node.FieldProperties, field.TypeJSON, value)
		_node.Properties = value
	}
	if nodes := nc.mutation.OwnerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   node.OwnerTable,
			Columns: []string{node.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(account.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.AccountID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := nc.mutation.ParentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   node.ParentTable,
			Columns: []string{node.ParentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(node.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.ParentNodeID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := nc.mutation.NodesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   node.NodesTable,
			Columns: []string{node.NodesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(node.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := nc.mutation.AssetsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   node.AssetsTable,
			Columns: node.AssetsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(asset.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := nc.mutation.TagsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   node.TagsTable,
			Columns: node.TagsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tag.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := nc.mutation.LinksIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   node.LinksTable,
			Columns: node.LinksPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(link.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := nc.mutation.CollectionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   node.CollectionsTable,
			Columns: node.CollectionsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(collection.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Node.Create().
//		SetCreatedAt(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.NodeUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
func (nc *NodeCreate) OnConflict(opts ...sql.ConflictOption) *NodeUpsertOne {
	nc.conflict = opts
	return &NodeUpsertOne{
		create: nc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Node.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (nc *NodeCreate) OnConflictColumns(columns ...string) *NodeUpsertOne {
	nc.conflict = append(nc.conflict, sql.ConflictColumns(columns...))
	return &NodeUpsertOne{
		create: nc,
	}
}

type (
	// NodeUpsertOne is the builder for "upsert"-ing
	//  one Node node.
	NodeUpsertOne struct {
		create *NodeCreate
	}

	// NodeUpsert is the "OnConflict" setter.
	NodeUpsert struct {
		*sql.UpdateSet
	}
)

// SetUpdatedAt sets the "updated_at" field.
func (u *NodeUpsert) SetUpdatedAt(v time.Time) *NodeUpsert {
	u.Set(node.FieldUpdatedAt, v)
	return u
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *NodeUpsert) UpdateUpdatedAt() *NodeUpsert {
	u.SetExcluded(node.FieldUpdatedAt)
	return u
}

// SetDeletedAt sets the "deleted_at" field.
func (u *NodeUpsert) SetDeletedAt(v time.Time) *NodeUpsert {
	u.Set(node.FieldDeletedAt, v)
	return u
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *NodeUpsert) UpdateDeletedAt() *NodeUpsert {
	u.SetExcluded(node.FieldDeletedAt)
	return u
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *NodeUpsert) ClearDeletedAt() *NodeUpsert {
	u.SetNull(node.FieldDeletedAt)
	return u
}

// SetName sets the "name" field.
func (u *NodeUpsert) SetName(v string) *NodeUpsert {
	u.Set(node.FieldName, v)
	return u
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *NodeUpsert) UpdateName() *NodeUpsert {
	u.SetExcluded(node.FieldName)
	return u
}

// SetSlug sets the "slug" field.
func (u *NodeUpsert) SetSlug(v string) *NodeUpsert {
	u.Set(node.FieldSlug, v)
	return u
}

// UpdateSlug sets the "slug" field to the value that was provided on create.
func (u *NodeUpsert) UpdateSlug() *NodeUpsert {
	u.SetExcluded(node.FieldSlug)
	return u
}

// SetDescription sets the "description" field.
func (u *NodeUpsert) SetDescription(v string) *NodeUpsert {
	u.Set(node.FieldDescription, v)
	return u
}

// UpdateDescription sets the "description" field to the value that was provided on create.
func (u *NodeUpsert) UpdateDescription() *NodeUpsert {
	u.SetExcluded(node.FieldDescription)
	return u
}

// ClearDescription clears the value of the "description" field.
func (u *NodeUpsert) ClearDescription() *NodeUpsert {
	u.SetNull(node.FieldDescription)
	return u
}

// SetContent sets the "content" field.
func (u *NodeUpsert) SetContent(v string) *NodeUpsert {
	u.Set(node.FieldContent, v)
	return u
}

// UpdateContent sets the "content" field to the value that was provided on create.
func (u *NodeUpsert) UpdateContent() *NodeUpsert {
	u.SetExcluded(node.FieldContent)
	return u
}

// ClearContent clears the value of the "content" field.
func (u *NodeUpsert) ClearContent() *NodeUpsert {
	u.SetNull(node.FieldContent)
	return u
}

// SetParentNodeID sets the "parent_node_id" field.
func (u *NodeUpsert) SetParentNodeID(v xid.ID) *NodeUpsert {
	u.Set(node.FieldParentNodeID, v)
	return u
}

// UpdateParentNodeID sets the "parent_node_id" field to the value that was provided on create.
func (u *NodeUpsert) UpdateParentNodeID() *NodeUpsert {
	u.SetExcluded(node.FieldParentNodeID)
	return u
}

// ClearParentNodeID clears the value of the "parent_node_id" field.
func (u *NodeUpsert) ClearParentNodeID() *NodeUpsert {
	u.SetNull(node.FieldParentNodeID)
	return u
}

// SetAccountID sets the "account_id" field.
func (u *NodeUpsert) SetAccountID(v xid.ID) *NodeUpsert {
	u.Set(node.FieldAccountID, v)
	return u
}

// UpdateAccountID sets the "account_id" field to the value that was provided on create.
func (u *NodeUpsert) UpdateAccountID() *NodeUpsert {
	u.SetExcluded(node.FieldAccountID)
	return u
}

// SetVisibility sets the "visibility" field.
func (u *NodeUpsert) SetVisibility(v node.Visibility) *NodeUpsert {
	u.Set(node.FieldVisibility, v)
	return u
}

// UpdateVisibility sets the "visibility" field to the value that was provided on create.
func (u *NodeUpsert) UpdateVisibility() *NodeUpsert {
	u.SetExcluded(node.FieldVisibility)
	return u
}

// SetProperties sets the "properties" field.
func (u *NodeUpsert) SetProperties(v any) *NodeUpsert {
	u.Set(node.FieldProperties, v)
	return u
}

// UpdateProperties sets the "properties" field to the value that was provided on create.
func (u *NodeUpsert) UpdateProperties() *NodeUpsert {
	u.SetExcluded(node.FieldProperties)
	return u
}

// ClearProperties clears the value of the "properties" field.
func (u *NodeUpsert) ClearProperties() *NodeUpsert {
	u.SetNull(node.FieldProperties)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.Node.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(node.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *NodeUpsertOne) UpdateNewValues() *NodeUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(node.FieldID)
		}
		if _, exists := u.create.mutation.CreatedAt(); exists {
			s.SetIgnore(node.FieldCreatedAt)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Node.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *NodeUpsertOne) Ignore() *NodeUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *NodeUpsertOne) DoNothing() *NodeUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the NodeCreate.OnConflict
// documentation for more info.
func (u *NodeUpsertOne) Update(set func(*NodeUpsert)) *NodeUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&NodeUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdatedAt sets the "updated_at" field.
func (u *NodeUpsertOne) SetUpdatedAt(v time.Time) *NodeUpsertOne {
	return u.Update(func(s *NodeUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *NodeUpsertOne) UpdateUpdatedAt() *NodeUpsertOne {
	return u.Update(func(s *NodeUpsert) {
		s.UpdateUpdatedAt()
	})
}

// SetDeletedAt sets the "deleted_at" field.
func (u *NodeUpsertOne) SetDeletedAt(v time.Time) *NodeUpsertOne {
	return u.Update(func(s *NodeUpsert) {
		s.SetDeletedAt(v)
	})
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *NodeUpsertOne) UpdateDeletedAt() *NodeUpsertOne {
	return u.Update(func(s *NodeUpsert) {
		s.UpdateDeletedAt()
	})
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *NodeUpsertOne) ClearDeletedAt() *NodeUpsertOne {
	return u.Update(func(s *NodeUpsert) {
		s.ClearDeletedAt()
	})
}

// SetName sets the "name" field.
func (u *NodeUpsertOne) SetName(v string) *NodeUpsertOne {
	return u.Update(func(s *NodeUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *NodeUpsertOne) UpdateName() *NodeUpsertOne {
	return u.Update(func(s *NodeUpsert) {
		s.UpdateName()
	})
}

// SetSlug sets the "slug" field.
func (u *NodeUpsertOne) SetSlug(v string) *NodeUpsertOne {
	return u.Update(func(s *NodeUpsert) {
		s.SetSlug(v)
	})
}

// UpdateSlug sets the "slug" field to the value that was provided on create.
func (u *NodeUpsertOne) UpdateSlug() *NodeUpsertOne {
	return u.Update(func(s *NodeUpsert) {
		s.UpdateSlug()
	})
}

// SetDescription sets the "description" field.
func (u *NodeUpsertOne) SetDescription(v string) *NodeUpsertOne {
	return u.Update(func(s *NodeUpsert) {
		s.SetDescription(v)
	})
}

// UpdateDescription sets the "description" field to the value that was provided on create.
func (u *NodeUpsertOne) UpdateDescription() *NodeUpsertOne {
	return u.Update(func(s *NodeUpsert) {
		s.UpdateDescription()
	})
}

// ClearDescription clears the value of the "description" field.
func (u *NodeUpsertOne) ClearDescription() *NodeUpsertOne {
	return u.Update(func(s *NodeUpsert) {
		s.ClearDescription()
	})
}

// SetContent sets the "content" field.
func (u *NodeUpsertOne) SetContent(v string) *NodeUpsertOne {
	return u.Update(func(s *NodeUpsert) {
		s.SetContent(v)
	})
}

// UpdateContent sets the "content" field to the value that was provided on create.
func (u *NodeUpsertOne) UpdateContent() *NodeUpsertOne {
	return u.Update(func(s *NodeUpsert) {
		s.UpdateContent()
	})
}

// ClearContent clears the value of the "content" field.
func (u *NodeUpsertOne) ClearContent() *NodeUpsertOne {
	return u.Update(func(s *NodeUpsert) {
		s.ClearContent()
	})
}

// SetParentNodeID sets the "parent_node_id" field.
func (u *NodeUpsertOne) SetParentNodeID(v xid.ID) *NodeUpsertOne {
	return u.Update(func(s *NodeUpsert) {
		s.SetParentNodeID(v)
	})
}

// UpdateParentNodeID sets the "parent_node_id" field to the value that was provided on create.
func (u *NodeUpsertOne) UpdateParentNodeID() *NodeUpsertOne {
	return u.Update(func(s *NodeUpsert) {
		s.UpdateParentNodeID()
	})
}

// ClearParentNodeID clears the value of the "parent_node_id" field.
func (u *NodeUpsertOne) ClearParentNodeID() *NodeUpsertOne {
	return u.Update(func(s *NodeUpsert) {
		s.ClearParentNodeID()
	})
}

// SetAccountID sets the "account_id" field.
func (u *NodeUpsertOne) SetAccountID(v xid.ID) *NodeUpsertOne {
	return u.Update(func(s *NodeUpsert) {
		s.SetAccountID(v)
	})
}

// UpdateAccountID sets the "account_id" field to the value that was provided on create.
func (u *NodeUpsertOne) UpdateAccountID() *NodeUpsertOne {
	return u.Update(func(s *NodeUpsert) {
		s.UpdateAccountID()
	})
}

// SetVisibility sets the "visibility" field.
func (u *NodeUpsertOne) SetVisibility(v node.Visibility) *NodeUpsertOne {
	return u.Update(func(s *NodeUpsert) {
		s.SetVisibility(v)
	})
}

// UpdateVisibility sets the "visibility" field to the value that was provided on create.
func (u *NodeUpsertOne) UpdateVisibility() *NodeUpsertOne {
	return u.Update(func(s *NodeUpsert) {
		s.UpdateVisibility()
	})
}

// SetProperties sets the "properties" field.
func (u *NodeUpsertOne) SetProperties(v any) *NodeUpsertOne {
	return u.Update(func(s *NodeUpsert) {
		s.SetProperties(v)
	})
}

// UpdateProperties sets the "properties" field to the value that was provided on create.
func (u *NodeUpsertOne) UpdateProperties() *NodeUpsertOne {
	return u.Update(func(s *NodeUpsert) {
		s.UpdateProperties()
	})
}

// ClearProperties clears the value of the "properties" field.
func (u *NodeUpsertOne) ClearProperties() *NodeUpsertOne {
	return u.Update(func(s *NodeUpsert) {
		s.ClearProperties()
	})
}

// Exec executes the query.
func (u *NodeUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for NodeCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *NodeUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *NodeUpsertOne) ID(ctx context.Context) (id xid.ID, err error) {
	if u.create.driver.Dialect() == dialect.MySQL {
		// In case of "ON CONFLICT", there is no way to get back non-numeric ID
		// fields from the database since MySQL does not support the RETURNING clause.
		return id, errors.New("ent: NodeUpsertOne.ID is not supported by MySQL driver. Use NodeUpsertOne.Exec instead")
	}
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *NodeUpsertOne) IDX(ctx context.Context) xid.ID {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// NodeCreateBulk is the builder for creating many Node entities in bulk.
type NodeCreateBulk struct {
	config
	err      error
	builders []*NodeCreate
	conflict []sql.ConflictOption
}

// Save creates the Node entities in the database.
func (ncb *NodeCreateBulk) Save(ctx context.Context) ([]*Node, error) {
	if ncb.err != nil {
		return nil, ncb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(ncb.builders))
	nodes := make([]*Node, len(ncb.builders))
	mutators := make([]Mutator, len(ncb.builders))
	for i := range ncb.builders {
		func(i int, root context.Context) {
			builder := ncb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*NodeMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, ncb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = ncb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ncb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, ncb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ncb *NodeCreateBulk) SaveX(ctx context.Context) []*Node {
	v, err := ncb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ncb *NodeCreateBulk) Exec(ctx context.Context) error {
	_, err := ncb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ncb *NodeCreateBulk) ExecX(ctx context.Context) {
	if err := ncb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Node.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.NodeUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
func (ncb *NodeCreateBulk) OnConflict(opts ...sql.ConflictOption) *NodeUpsertBulk {
	ncb.conflict = opts
	return &NodeUpsertBulk{
		create: ncb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Node.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (ncb *NodeCreateBulk) OnConflictColumns(columns ...string) *NodeUpsertBulk {
	ncb.conflict = append(ncb.conflict, sql.ConflictColumns(columns...))
	return &NodeUpsertBulk{
		create: ncb,
	}
}

// NodeUpsertBulk is the builder for "upsert"-ing
// a bulk of Node nodes.
type NodeUpsertBulk struct {
	create *NodeCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Node.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(node.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *NodeUpsertBulk) UpdateNewValues() *NodeUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(node.FieldID)
			}
			if _, exists := b.mutation.CreatedAt(); exists {
				s.SetIgnore(node.FieldCreatedAt)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Node.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *NodeUpsertBulk) Ignore() *NodeUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *NodeUpsertBulk) DoNothing() *NodeUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the NodeCreateBulk.OnConflict
// documentation for more info.
func (u *NodeUpsertBulk) Update(set func(*NodeUpsert)) *NodeUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&NodeUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdatedAt sets the "updated_at" field.
func (u *NodeUpsertBulk) SetUpdatedAt(v time.Time) *NodeUpsertBulk {
	return u.Update(func(s *NodeUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *NodeUpsertBulk) UpdateUpdatedAt() *NodeUpsertBulk {
	return u.Update(func(s *NodeUpsert) {
		s.UpdateUpdatedAt()
	})
}

// SetDeletedAt sets the "deleted_at" field.
func (u *NodeUpsertBulk) SetDeletedAt(v time.Time) *NodeUpsertBulk {
	return u.Update(func(s *NodeUpsert) {
		s.SetDeletedAt(v)
	})
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *NodeUpsertBulk) UpdateDeletedAt() *NodeUpsertBulk {
	return u.Update(func(s *NodeUpsert) {
		s.UpdateDeletedAt()
	})
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *NodeUpsertBulk) ClearDeletedAt() *NodeUpsertBulk {
	return u.Update(func(s *NodeUpsert) {
		s.ClearDeletedAt()
	})
}

// SetName sets the "name" field.
func (u *NodeUpsertBulk) SetName(v string) *NodeUpsertBulk {
	return u.Update(func(s *NodeUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *NodeUpsertBulk) UpdateName() *NodeUpsertBulk {
	return u.Update(func(s *NodeUpsert) {
		s.UpdateName()
	})
}

// SetSlug sets the "slug" field.
func (u *NodeUpsertBulk) SetSlug(v string) *NodeUpsertBulk {
	return u.Update(func(s *NodeUpsert) {
		s.SetSlug(v)
	})
}

// UpdateSlug sets the "slug" field to the value that was provided on create.
func (u *NodeUpsertBulk) UpdateSlug() *NodeUpsertBulk {
	return u.Update(func(s *NodeUpsert) {
		s.UpdateSlug()
	})
}

// SetDescription sets the "description" field.
func (u *NodeUpsertBulk) SetDescription(v string) *NodeUpsertBulk {
	return u.Update(func(s *NodeUpsert) {
		s.SetDescription(v)
	})
}

// UpdateDescription sets the "description" field to the value that was provided on create.
func (u *NodeUpsertBulk) UpdateDescription() *NodeUpsertBulk {
	return u.Update(func(s *NodeUpsert) {
		s.UpdateDescription()
	})
}

// ClearDescription clears the value of the "description" field.
func (u *NodeUpsertBulk) ClearDescription() *NodeUpsertBulk {
	return u.Update(func(s *NodeUpsert) {
		s.ClearDescription()
	})
}

// SetContent sets the "content" field.
func (u *NodeUpsertBulk) SetContent(v string) *NodeUpsertBulk {
	return u.Update(func(s *NodeUpsert) {
		s.SetContent(v)
	})
}

// UpdateContent sets the "content" field to the value that was provided on create.
func (u *NodeUpsertBulk) UpdateContent() *NodeUpsertBulk {
	return u.Update(func(s *NodeUpsert) {
		s.UpdateContent()
	})
}

// ClearContent clears the value of the "content" field.
func (u *NodeUpsertBulk) ClearContent() *NodeUpsertBulk {
	return u.Update(func(s *NodeUpsert) {
		s.ClearContent()
	})
}

// SetParentNodeID sets the "parent_node_id" field.
func (u *NodeUpsertBulk) SetParentNodeID(v xid.ID) *NodeUpsertBulk {
	return u.Update(func(s *NodeUpsert) {
		s.SetParentNodeID(v)
	})
}

// UpdateParentNodeID sets the "parent_node_id" field to the value that was provided on create.
func (u *NodeUpsertBulk) UpdateParentNodeID() *NodeUpsertBulk {
	return u.Update(func(s *NodeUpsert) {
		s.UpdateParentNodeID()
	})
}

// ClearParentNodeID clears the value of the "parent_node_id" field.
func (u *NodeUpsertBulk) ClearParentNodeID() *NodeUpsertBulk {
	return u.Update(func(s *NodeUpsert) {
		s.ClearParentNodeID()
	})
}

// SetAccountID sets the "account_id" field.
func (u *NodeUpsertBulk) SetAccountID(v xid.ID) *NodeUpsertBulk {
	return u.Update(func(s *NodeUpsert) {
		s.SetAccountID(v)
	})
}

// UpdateAccountID sets the "account_id" field to the value that was provided on create.
func (u *NodeUpsertBulk) UpdateAccountID() *NodeUpsertBulk {
	return u.Update(func(s *NodeUpsert) {
		s.UpdateAccountID()
	})
}

// SetVisibility sets the "visibility" field.
func (u *NodeUpsertBulk) SetVisibility(v node.Visibility) *NodeUpsertBulk {
	return u.Update(func(s *NodeUpsert) {
		s.SetVisibility(v)
	})
}

// UpdateVisibility sets the "visibility" field to the value that was provided on create.
func (u *NodeUpsertBulk) UpdateVisibility() *NodeUpsertBulk {
	return u.Update(func(s *NodeUpsert) {
		s.UpdateVisibility()
	})
}

// SetProperties sets the "properties" field.
func (u *NodeUpsertBulk) SetProperties(v any) *NodeUpsertBulk {
	return u.Update(func(s *NodeUpsert) {
		s.SetProperties(v)
	})
}

// UpdateProperties sets the "properties" field to the value that was provided on create.
func (u *NodeUpsertBulk) UpdateProperties() *NodeUpsertBulk {
	return u.Update(func(s *NodeUpsert) {
		s.UpdateProperties()
	})
}

// ClearProperties clears the value of the "properties" field.
func (u *NodeUpsertBulk) ClearProperties() *NodeUpsertBulk {
	return u.Update(func(s *NodeUpsert) {
		s.ClearProperties()
	})
}

// Exec executes the query.
func (u *NodeUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the NodeCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for NodeCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *NodeUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
