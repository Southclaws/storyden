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
	"github.com/Southclaws/storyden/internal/ent/cluster"
	"github.com/Southclaws/storyden/internal/ent/item"
	"github.com/Southclaws/storyden/internal/ent/tag"
	"github.com/rs/xid"
)

// ClusterCreate is the builder for creating a Cluster entity.
type ClusterCreate struct {
	config
	mutation *ClusterMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCreatedAt sets the "created_at" field.
func (cc *ClusterCreate) SetCreatedAt(t time.Time) *ClusterCreate {
	cc.mutation.SetCreatedAt(t)
	return cc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (cc *ClusterCreate) SetNillableCreatedAt(t *time.Time) *ClusterCreate {
	if t != nil {
		cc.SetCreatedAt(*t)
	}
	return cc
}

// SetUpdatedAt sets the "updated_at" field.
func (cc *ClusterCreate) SetUpdatedAt(t time.Time) *ClusterCreate {
	cc.mutation.SetUpdatedAt(t)
	return cc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (cc *ClusterCreate) SetNillableUpdatedAt(t *time.Time) *ClusterCreate {
	if t != nil {
		cc.SetUpdatedAt(*t)
	}
	return cc
}

// SetDeletedAt sets the "deleted_at" field.
func (cc *ClusterCreate) SetDeletedAt(t time.Time) *ClusterCreate {
	cc.mutation.SetDeletedAt(t)
	return cc
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (cc *ClusterCreate) SetNillableDeletedAt(t *time.Time) *ClusterCreate {
	if t != nil {
		cc.SetDeletedAt(*t)
	}
	return cc
}

// SetName sets the "name" field.
func (cc *ClusterCreate) SetName(s string) *ClusterCreate {
	cc.mutation.SetName(s)
	return cc
}

// SetSlug sets the "slug" field.
func (cc *ClusterCreate) SetSlug(s string) *ClusterCreate {
	cc.mutation.SetSlug(s)
	return cc
}

// SetImageURL sets the "image_url" field.
func (cc *ClusterCreate) SetImageURL(s string) *ClusterCreate {
	cc.mutation.SetImageURL(s)
	return cc
}

// SetNillableImageURL sets the "image_url" field if the given value is not nil.
func (cc *ClusterCreate) SetNillableImageURL(s *string) *ClusterCreate {
	if s != nil {
		cc.SetImageURL(*s)
	}
	return cc
}

// SetDescription sets the "description" field.
func (cc *ClusterCreate) SetDescription(s string) *ClusterCreate {
	cc.mutation.SetDescription(s)
	return cc
}

// SetContent sets the "content" field.
func (cc *ClusterCreate) SetContent(s string) *ClusterCreate {
	cc.mutation.SetContent(s)
	return cc
}

// SetNillableContent sets the "content" field if the given value is not nil.
func (cc *ClusterCreate) SetNillableContent(s *string) *ClusterCreate {
	if s != nil {
		cc.SetContent(*s)
	}
	return cc
}

// SetParentClusterID sets the "parent_cluster_id" field.
func (cc *ClusterCreate) SetParentClusterID(x xid.ID) *ClusterCreate {
	cc.mutation.SetParentClusterID(x)
	return cc
}

// SetNillableParentClusterID sets the "parent_cluster_id" field if the given value is not nil.
func (cc *ClusterCreate) SetNillableParentClusterID(x *xid.ID) *ClusterCreate {
	if x != nil {
		cc.SetParentClusterID(*x)
	}
	return cc
}

// SetAccountID sets the "account_id" field.
func (cc *ClusterCreate) SetAccountID(x xid.ID) *ClusterCreate {
	cc.mutation.SetAccountID(x)
	return cc
}

// SetProperties sets the "properties" field.
func (cc *ClusterCreate) SetProperties(a any) *ClusterCreate {
	cc.mutation.SetProperties(a)
	return cc
}

// SetID sets the "id" field.
func (cc *ClusterCreate) SetID(x xid.ID) *ClusterCreate {
	cc.mutation.SetID(x)
	return cc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (cc *ClusterCreate) SetNillableID(x *xid.ID) *ClusterCreate {
	if x != nil {
		cc.SetID(*x)
	}
	return cc
}

// SetOwnerID sets the "owner" edge to the Account entity by ID.
func (cc *ClusterCreate) SetOwnerID(id xid.ID) *ClusterCreate {
	cc.mutation.SetOwnerID(id)
	return cc
}

// SetOwner sets the "owner" edge to the Account entity.
func (cc *ClusterCreate) SetOwner(a *Account) *ClusterCreate {
	return cc.SetOwnerID(a.ID)
}

// SetParentID sets the "parent" edge to the Cluster entity by ID.
func (cc *ClusterCreate) SetParentID(id xid.ID) *ClusterCreate {
	cc.mutation.SetParentID(id)
	return cc
}

// SetNillableParentID sets the "parent" edge to the Cluster entity by ID if the given value is not nil.
func (cc *ClusterCreate) SetNillableParentID(id *xid.ID) *ClusterCreate {
	if id != nil {
		cc = cc.SetParentID(*id)
	}
	return cc
}

// SetParent sets the "parent" edge to the Cluster entity.
func (cc *ClusterCreate) SetParent(c *Cluster) *ClusterCreate {
	return cc.SetParentID(c.ID)
}

// AddClusterIDs adds the "clusters" edge to the Cluster entity by IDs.
func (cc *ClusterCreate) AddClusterIDs(ids ...xid.ID) *ClusterCreate {
	cc.mutation.AddClusterIDs(ids...)
	return cc
}

// AddClusters adds the "clusters" edges to the Cluster entity.
func (cc *ClusterCreate) AddClusters(c ...*Cluster) *ClusterCreate {
	ids := make([]xid.ID, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return cc.AddClusterIDs(ids...)
}

// AddItemIDs adds the "items" edge to the Item entity by IDs.
func (cc *ClusterCreate) AddItemIDs(ids ...xid.ID) *ClusterCreate {
	cc.mutation.AddItemIDs(ids...)
	return cc
}

// AddItems adds the "items" edges to the Item entity.
func (cc *ClusterCreate) AddItems(i ...*Item) *ClusterCreate {
	ids := make([]xid.ID, len(i))
	for j := range i {
		ids[j] = i[j].ID
	}
	return cc.AddItemIDs(ids...)
}

// AddAssetIDs adds the "assets" edge to the Asset entity by IDs.
func (cc *ClusterCreate) AddAssetIDs(ids ...string) *ClusterCreate {
	cc.mutation.AddAssetIDs(ids...)
	return cc
}

// AddAssets adds the "assets" edges to the Asset entity.
func (cc *ClusterCreate) AddAssets(a ...*Asset) *ClusterCreate {
	ids := make([]string, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return cc.AddAssetIDs(ids...)
}

// AddTagIDs adds the "tags" edge to the Tag entity by IDs.
func (cc *ClusterCreate) AddTagIDs(ids ...xid.ID) *ClusterCreate {
	cc.mutation.AddTagIDs(ids...)
	return cc
}

// AddTags adds the "tags" edges to the Tag entity.
func (cc *ClusterCreate) AddTags(t ...*Tag) *ClusterCreate {
	ids := make([]xid.ID, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return cc.AddTagIDs(ids...)
}

// Mutation returns the ClusterMutation object of the builder.
func (cc *ClusterCreate) Mutation() *ClusterMutation {
	return cc.mutation
}

// Save creates the Cluster in the database.
func (cc *ClusterCreate) Save(ctx context.Context) (*Cluster, error) {
	cc.defaults()
	return withHooks(ctx, cc.sqlSave, cc.mutation, cc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (cc *ClusterCreate) SaveX(ctx context.Context) *Cluster {
	v, err := cc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (cc *ClusterCreate) Exec(ctx context.Context) error {
	_, err := cc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cc *ClusterCreate) ExecX(ctx context.Context) {
	if err := cc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (cc *ClusterCreate) defaults() {
	if _, ok := cc.mutation.CreatedAt(); !ok {
		v := cluster.DefaultCreatedAt()
		cc.mutation.SetCreatedAt(v)
	}
	if _, ok := cc.mutation.UpdatedAt(); !ok {
		v := cluster.DefaultUpdatedAt()
		cc.mutation.SetUpdatedAt(v)
	}
	if _, ok := cc.mutation.ID(); !ok {
		v := cluster.DefaultID()
		cc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (cc *ClusterCreate) check() error {
	if _, ok := cc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "Cluster.created_at"`)}
	}
	if _, ok := cc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "Cluster.updated_at"`)}
	}
	if _, ok := cc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "Cluster.name"`)}
	}
	if _, ok := cc.mutation.Slug(); !ok {
		return &ValidationError{Name: "slug", err: errors.New(`ent: missing required field "Cluster.slug"`)}
	}
	if _, ok := cc.mutation.Description(); !ok {
		return &ValidationError{Name: "description", err: errors.New(`ent: missing required field "Cluster.description"`)}
	}
	if _, ok := cc.mutation.AccountID(); !ok {
		return &ValidationError{Name: "account_id", err: errors.New(`ent: missing required field "Cluster.account_id"`)}
	}
	if v, ok := cc.mutation.ID(); ok {
		if err := cluster.IDValidator(v.String()); err != nil {
			return &ValidationError{Name: "id", err: fmt.Errorf(`ent: validator failed for field "Cluster.id": %w`, err)}
		}
	}
	if _, ok := cc.mutation.OwnerID(); !ok {
		return &ValidationError{Name: "owner", err: errors.New(`ent: missing required edge "Cluster.owner"`)}
	}
	return nil
}

func (cc *ClusterCreate) sqlSave(ctx context.Context) (*Cluster, error) {
	if err := cc.check(); err != nil {
		return nil, err
	}
	_node, _spec := cc.createSpec()
	if err := sqlgraph.CreateNode(ctx, cc.driver, _spec); err != nil {
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
	cc.mutation.id = &_node.ID
	cc.mutation.done = true
	return _node, nil
}

func (cc *ClusterCreate) createSpec() (*Cluster, *sqlgraph.CreateSpec) {
	var (
		_node = &Cluster{config: cc.config}
		_spec = sqlgraph.NewCreateSpec(cluster.Table, sqlgraph.NewFieldSpec(cluster.FieldID, field.TypeString))
	)
	_spec.OnConflict = cc.conflict
	if id, ok := cc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := cc.mutation.CreatedAt(); ok {
		_spec.SetField(cluster.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := cc.mutation.UpdatedAt(); ok {
		_spec.SetField(cluster.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	if value, ok := cc.mutation.DeletedAt(); ok {
		_spec.SetField(cluster.FieldDeletedAt, field.TypeTime, value)
		_node.DeletedAt = &value
	}
	if value, ok := cc.mutation.Name(); ok {
		_spec.SetField(cluster.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := cc.mutation.Slug(); ok {
		_spec.SetField(cluster.FieldSlug, field.TypeString, value)
		_node.Slug = value
	}
	if value, ok := cc.mutation.ImageURL(); ok {
		_spec.SetField(cluster.FieldImageURL, field.TypeString, value)
		_node.ImageURL = &value
	}
	if value, ok := cc.mutation.Description(); ok {
		_spec.SetField(cluster.FieldDescription, field.TypeString, value)
		_node.Description = value
	}
	if value, ok := cc.mutation.Content(); ok {
		_spec.SetField(cluster.FieldContent, field.TypeString, value)
		_node.Content = &value
	}
	if value, ok := cc.mutation.Properties(); ok {
		_spec.SetField(cluster.FieldProperties, field.TypeJSON, value)
		_node.Properties = value
	}
	if nodes := cc.mutation.OwnerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   cluster.OwnerTable,
			Columns: []string{cluster.OwnerColumn},
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
	if nodes := cc.mutation.ParentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   cluster.ParentTable,
			Columns: []string{cluster.ParentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(cluster.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.ParentClusterID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := cc.mutation.ClustersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   cluster.ClustersTable,
			Columns: []string{cluster.ClustersColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(cluster.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := cc.mutation.ItemsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   cluster.ItemsTable,
			Columns: cluster.ItemsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(item.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := cc.mutation.AssetsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   cluster.AssetsTable,
			Columns: cluster.AssetsPrimaryKey,
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
	if nodes := cc.mutation.TagsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   cluster.TagsTable,
			Columns: cluster.TagsPrimaryKey,
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
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Cluster.Create().
//		SetCreatedAt(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.ClusterUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
func (cc *ClusterCreate) OnConflict(opts ...sql.ConflictOption) *ClusterUpsertOne {
	cc.conflict = opts
	return &ClusterUpsertOne{
		create: cc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Cluster.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (cc *ClusterCreate) OnConflictColumns(columns ...string) *ClusterUpsertOne {
	cc.conflict = append(cc.conflict, sql.ConflictColumns(columns...))
	return &ClusterUpsertOne{
		create: cc,
	}
}

type (
	// ClusterUpsertOne is the builder for "upsert"-ing
	//  one Cluster node.
	ClusterUpsertOne struct {
		create *ClusterCreate
	}

	// ClusterUpsert is the "OnConflict" setter.
	ClusterUpsert struct {
		*sql.UpdateSet
	}
)

// SetUpdatedAt sets the "updated_at" field.
func (u *ClusterUpsert) SetUpdatedAt(v time.Time) *ClusterUpsert {
	u.Set(cluster.FieldUpdatedAt, v)
	return u
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *ClusterUpsert) UpdateUpdatedAt() *ClusterUpsert {
	u.SetExcluded(cluster.FieldUpdatedAt)
	return u
}

// SetDeletedAt sets the "deleted_at" field.
func (u *ClusterUpsert) SetDeletedAt(v time.Time) *ClusterUpsert {
	u.Set(cluster.FieldDeletedAt, v)
	return u
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *ClusterUpsert) UpdateDeletedAt() *ClusterUpsert {
	u.SetExcluded(cluster.FieldDeletedAt)
	return u
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *ClusterUpsert) ClearDeletedAt() *ClusterUpsert {
	u.SetNull(cluster.FieldDeletedAt)
	return u
}

// SetName sets the "name" field.
func (u *ClusterUpsert) SetName(v string) *ClusterUpsert {
	u.Set(cluster.FieldName, v)
	return u
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *ClusterUpsert) UpdateName() *ClusterUpsert {
	u.SetExcluded(cluster.FieldName)
	return u
}

// SetSlug sets the "slug" field.
func (u *ClusterUpsert) SetSlug(v string) *ClusterUpsert {
	u.Set(cluster.FieldSlug, v)
	return u
}

// UpdateSlug sets the "slug" field to the value that was provided on create.
func (u *ClusterUpsert) UpdateSlug() *ClusterUpsert {
	u.SetExcluded(cluster.FieldSlug)
	return u
}

// SetImageURL sets the "image_url" field.
func (u *ClusterUpsert) SetImageURL(v string) *ClusterUpsert {
	u.Set(cluster.FieldImageURL, v)
	return u
}

// UpdateImageURL sets the "image_url" field to the value that was provided on create.
func (u *ClusterUpsert) UpdateImageURL() *ClusterUpsert {
	u.SetExcluded(cluster.FieldImageURL)
	return u
}

// ClearImageURL clears the value of the "image_url" field.
func (u *ClusterUpsert) ClearImageURL() *ClusterUpsert {
	u.SetNull(cluster.FieldImageURL)
	return u
}

// SetDescription sets the "description" field.
func (u *ClusterUpsert) SetDescription(v string) *ClusterUpsert {
	u.Set(cluster.FieldDescription, v)
	return u
}

// UpdateDescription sets the "description" field to the value that was provided on create.
func (u *ClusterUpsert) UpdateDescription() *ClusterUpsert {
	u.SetExcluded(cluster.FieldDescription)
	return u
}

// SetContent sets the "content" field.
func (u *ClusterUpsert) SetContent(v string) *ClusterUpsert {
	u.Set(cluster.FieldContent, v)
	return u
}

// UpdateContent sets the "content" field to the value that was provided on create.
func (u *ClusterUpsert) UpdateContent() *ClusterUpsert {
	u.SetExcluded(cluster.FieldContent)
	return u
}

// ClearContent clears the value of the "content" field.
func (u *ClusterUpsert) ClearContent() *ClusterUpsert {
	u.SetNull(cluster.FieldContent)
	return u
}

// SetParentClusterID sets the "parent_cluster_id" field.
func (u *ClusterUpsert) SetParentClusterID(v xid.ID) *ClusterUpsert {
	u.Set(cluster.FieldParentClusterID, v)
	return u
}

// UpdateParentClusterID sets the "parent_cluster_id" field to the value that was provided on create.
func (u *ClusterUpsert) UpdateParentClusterID() *ClusterUpsert {
	u.SetExcluded(cluster.FieldParentClusterID)
	return u
}

// ClearParentClusterID clears the value of the "parent_cluster_id" field.
func (u *ClusterUpsert) ClearParentClusterID() *ClusterUpsert {
	u.SetNull(cluster.FieldParentClusterID)
	return u
}

// SetAccountID sets the "account_id" field.
func (u *ClusterUpsert) SetAccountID(v xid.ID) *ClusterUpsert {
	u.Set(cluster.FieldAccountID, v)
	return u
}

// UpdateAccountID sets the "account_id" field to the value that was provided on create.
func (u *ClusterUpsert) UpdateAccountID() *ClusterUpsert {
	u.SetExcluded(cluster.FieldAccountID)
	return u
}

// SetProperties sets the "properties" field.
func (u *ClusterUpsert) SetProperties(v any) *ClusterUpsert {
	u.Set(cluster.FieldProperties, v)
	return u
}

// UpdateProperties sets the "properties" field to the value that was provided on create.
func (u *ClusterUpsert) UpdateProperties() *ClusterUpsert {
	u.SetExcluded(cluster.FieldProperties)
	return u
}

// ClearProperties clears the value of the "properties" field.
func (u *ClusterUpsert) ClearProperties() *ClusterUpsert {
	u.SetNull(cluster.FieldProperties)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.Cluster.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(cluster.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *ClusterUpsertOne) UpdateNewValues() *ClusterUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(cluster.FieldID)
		}
		if _, exists := u.create.mutation.CreatedAt(); exists {
			s.SetIgnore(cluster.FieldCreatedAt)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Cluster.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *ClusterUpsertOne) Ignore() *ClusterUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *ClusterUpsertOne) DoNothing() *ClusterUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the ClusterCreate.OnConflict
// documentation for more info.
func (u *ClusterUpsertOne) Update(set func(*ClusterUpsert)) *ClusterUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&ClusterUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdatedAt sets the "updated_at" field.
func (u *ClusterUpsertOne) SetUpdatedAt(v time.Time) *ClusterUpsertOne {
	return u.Update(func(s *ClusterUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *ClusterUpsertOne) UpdateUpdatedAt() *ClusterUpsertOne {
	return u.Update(func(s *ClusterUpsert) {
		s.UpdateUpdatedAt()
	})
}

// SetDeletedAt sets the "deleted_at" field.
func (u *ClusterUpsertOne) SetDeletedAt(v time.Time) *ClusterUpsertOne {
	return u.Update(func(s *ClusterUpsert) {
		s.SetDeletedAt(v)
	})
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *ClusterUpsertOne) UpdateDeletedAt() *ClusterUpsertOne {
	return u.Update(func(s *ClusterUpsert) {
		s.UpdateDeletedAt()
	})
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *ClusterUpsertOne) ClearDeletedAt() *ClusterUpsertOne {
	return u.Update(func(s *ClusterUpsert) {
		s.ClearDeletedAt()
	})
}

// SetName sets the "name" field.
func (u *ClusterUpsertOne) SetName(v string) *ClusterUpsertOne {
	return u.Update(func(s *ClusterUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *ClusterUpsertOne) UpdateName() *ClusterUpsertOne {
	return u.Update(func(s *ClusterUpsert) {
		s.UpdateName()
	})
}

// SetSlug sets the "slug" field.
func (u *ClusterUpsertOne) SetSlug(v string) *ClusterUpsertOne {
	return u.Update(func(s *ClusterUpsert) {
		s.SetSlug(v)
	})
}

// UpdateSlug sets the "slug" field to the value that was provided on create.
func (u *ClusterUpsertOne) UpdateSlug() *ClusterUpsertOne {
	return u.Update(func(s *ClusterUpsert) {
		s.UpdateSlug()
	})
}

// SetImageURL sets the "image_url" field.
func (u *ClusterUpsertOne) SetImageURL(v string) *ClusterUpsertOne {
	return u.Update(func(s *ClusterUpsert) {
		s.SetImageURL(v)
	})
}

// UpdateImageURL sets the "image_url" field to the value that was provided on create.
func (u *ClusterUpsertOne) UpdateImageURL() *ClusterUpsertOne {
	return u.Update(func(s *ClusterUpsert) {
		s.UpdateImageURL()
	})
}

// ClearImageURL clears the value of the "image_url" field.
func (u *ClusterUpsertOne) ClearImageURL() *ClusterUpsertOne {
	return u.Update(func(s *ClusterUpsert) {
		s.ClearImageURL()
	})
}

// SetDescription sets the "description" field.
func (u *ClusterUpsertOne) SetDescription(v string) *ClusterUpsertOne {
	return u.Update(func(s *ClusterUpsert) {
		s.SetDescription(v)
	})
}

// UpdateDescription sets the "description" field to the value that was provided on create.
func (u *ClusterUpsertOne) UpdateDescription() *ClusterUpsertOne {
	return u.Update(func(s *ClusterUpsert) {
		s.UpdateDescription()
	})
}

// SetContent sets the "content" field.
func (u *ClusterUpsertOne) SetContent(v string) *ClusterUpsertOne {
	return u.Update(func(s *ClusterUpsert) {
		s.SetContent(v)
	})
}

// UpdateContent sets the "content" field to the value that was provided on create.
func (u *ClusterUpsertOne) UpdateContent() *ClusterUpsertOne {
	return u.Update(func(s *ClusterUpsert) {
		s.UpdateContent()
	})
}

// ClearContent clears the value of the "content" field.
func (u *ClusterUpsertOne) ClearContent() *ClusterUpsertOne {
	return u.Update(func(s *ClusterUpsert) {
		s.ClearContent()
	})
}

// SetParentClusterID sets the "parent_cluster_id" field.
func (u *ClusterUpsertOne) SetParentClusterID(v xid.ID) *ClusterUpsertOne {
	return u.Update(func(s *ClusterUpsert) {
		s.SetParentClusterID(v)
	})
}

// UpdateParentClusterID sets the "parent_cluster_id" field to the value that was provided on create.
func (u *ClusterUpsertOne) UpdateParentClusterID() *ClusterUpsertOne {
	return u.Update(func(s *ClusterUpsert) {
		s.UpdateParentClusterID()
	})
}

// ClearParentClusterID clears the value of the "parent_cluster_id" field.
func (u *ClusterUpsertOne) ClearParentClusterID() *ClusterUpsertOne {
	return u.Update(func(s *ClusterUpsert) {
		s.ClearParentClusterID()
	})
}

// SetAccountID sets the "account_id" field.
func (u *ClusterUpsertOne) SetAccountID(v xid.ID) *ClusterUpsertOne {
	return u.Update(func(s *ClusterUpsert) {
		s.SetAccountID(v)
	})
}

// UpdateAccountID sets the "account_id" field to the value that was provided on create.
func (u *ClusterUpsertOne) UpdateAccountID() *ClusterUpsertOne {
	return u.Update(func(s *ClusterUpsert) {
		s.UpdateAccountID()
	})
}

// SetProperties sets the "properties" field.
func (u *ClusterUpsertOne) SetProperties(v any) *ClusterUpsertOne {
	return u.Update(func(s *ClusterUpsert) {
		s.SetProperties(v)
	})
}

// UpdateProperties sets the "properties" field to the value that was provided on create.
func (u *ClusterUpsertOne) UpdateProperties() *ClusterUpsertOne {
	return u.Update(func(s *ClusterUpsert) {
		s.UpdateProperties()
	})
}

// ClearProperties clears the value of the "properties" field.
func (u *ClusterUpsertOne) ClearProperties() *ClusterUpsertOne {
	return u.Update(func(s *ClusterUpsert) {
		s.ClearProperties()
	})
}

// Exec executes the query.
func (u *ClusterUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for ClusterCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *ClusterUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *ClusterUpsertOne) ID(ctx context.Context) (id xid.ID, err error) {
	if u.create.driver.Dialect() == dialect.MySQL {
		// In case of "ON CONFLICT", there is no way to get back non-numeric ID
		// fields from the database since MySQL does not support the RETURNING clause.
		return id, errors.New("ent: ClusterUpsertOne.ID is not supported by MySQL driver. Use ClusterUpsertOne.Exec instead")
	}
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *ClusterUpsertOne) IDX(ctx context.Context) xid.ID {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// ClusterCreateBulk is the builder for creating many Cluster entities in bulk.
type ClusterCreateBulk struct {
	config
	err      error
	builders []*ClusterCreate
	conflict []sql.ConflictOption
}

// Save creates the Cluster entities in the database.
func (ccb *ClusterCreateBulk) Save(ctx context.Context) ([]*Cluster, error) {
	if ccb.err != nil {
		return nil, ccb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(ccb.builders))
	nodes := make([]*Cluster, len(ccb.builders))
	mutators := make([]Mutator, len(ccb.builders))
	for i := range ccb.builders {
		func(i int, root context.Context) {
			builder := ccb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ClusterMutation)
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
					_, err = mutators[i+1].Mutate(root, ccb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = ccb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ccb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, ccb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ccb *ClusterCreateBulk) SaveX(ctx context.Context) []*Cluster {
	v, err := ccb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ccb *ClusterCreateBulk) Exec(ctx context.Context) error {
	_, err := ccb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ccb *ClusterCreateBulk) ExecX(ctx context.Context) {
	if err := ccb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Cluster.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.ClusterUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
func (ccb *ClusterCreateBulk) OnConflict(opts ...sql.ConflictOption) *ClusterUpsertBulk {
	ccb.conflict = opts
	return &ClusterUpsertBulk{
		create: ccb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Cluster.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (ccb *ClusterCreateBulk) OnConflictColumns(columns ...string) *ClusterUpsertBulk {
	ccb.conflict = append(ccb.conflict, sql.ConflictColumns(columns...))
	return &ClusterUpsertBulk{
		create: ccb,
	}
}

// ClusterUpsertBulk is the builder for "upsert"-ing
// a bulk of Cluster nodes.
type ClusterUpsertBulk struct {
	create *ClusterCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Cluster.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(cluster.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *ClusterUpsertBulk) UpdateNewValues() *ClusterUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(cluster.FieldID)
			}
			if _, exists := b.mutation.CreatedAt(); exists {
				s.SetIgnore(cluster.FieldCreatedAt)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Cluster.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *ClusterUpsertBulk) Ignore() *ClusterUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *ClusterUpsertBulk) DoNothing() *ClusterUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the ClusterCreateBulk.OnConflict
// documentation for more info.
func (u *ClusterUpsertBulk) Update(set func(*ClusterUpsert)) *ClusterUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&ClusterUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdatedAt sets the "updated_at" field.
func (u *ClusterUpsertBulk) SetUpdatedAt(v time.Time) *ClusterUpsertBulk {
	return u.Update(func(s *ClusterUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *ClusterUpsertBulk) UpdateUpdatedAt() *ClusterUpsertBulk {
	return u.Update(func(s *ClusterUpsert) {
		s.UpdateUpdatedAt()
	})
}

// SetDeletedAt sets the "deleted_at" field.
func (u *ClusterUpsertBulk) SetDeletedAt(v time.Time) *ClusterUpsertBulk {
	return u.Update(func(s *ClusterUpsert) {
		s.SetDeletedAt(v)
	})
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *ClusterUpsertBulk) UpdateDeletedAt() *ClusterUpsertBulk {
	return u.Update(func(s *ClusterUpsert) {
		s.UpdateDeletedAt()
	})
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (u *ClusterUpsertBulk) ClearDeletedAt() *ClusterUpsertBulk {
	return u.Update(func(s *ClusterUpsert) {
		s.ClearDeletedAt()
	})
}

// SetName sets the "name" field.
func (u *ClusterUpsertBulk) SetName(v string) *ClusterUpsertBulk {
	return u.Update(func(s *ClusterUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *ClusterUpsertBulk) UpdateName() *ClusterUpsertBulk {
	return u.Update(func(s *ClusterUpsert) {
		s.UpdateName()
	})
}

// SetSlug sets the "slug" field.
func (u *ClusterUpsertBulk) SetSlug(v string) *ClusterUpsertBulk {
	return u.Update(func(s *ClusterUpsert) {
		s.SetSlug(v)
	})
}

// UpdateSlug sets the "slug" field to the value that was provided on create.
func (u *ClusterUpsertBulk) UpdateSlug() *ClusterUpsertBulk {
	return u.Update(func(s *ClusterUpsert) {
		s.UpdateSlug()
	})
}

// SetImageURL sets the "image_url" field.
func (u *ClusterUpsertBulk) SetImageURL(v string) *ClusterUpsertBulk {
	return u.Update(func(s *ClusterUpsert) {
		s.SetImageURL(v)
	})
}

// UpdateImageURL sets the "image_url" field to the value that was provided on create.
func (u *ClusterUpsertBulk) UpdateImageURL() *ClusterUpsertBulk {
	return u.Update(func(s *ClusterUpsert) {
		s.UpdateImageURL()
	})
}

// ClearImageURL clears the value of the "image_url" field.
func (u *ClusterUpsertBulk) ClearImageURL() *ClusterUpsertBulk {
	return u.Update(func(s *ClusterUpsert) {
		s.ClearImageURL()
	})
}

// SetDescription sets the "description" field.
func (u *ClusterUpsertBulk) SetDescription(v string) *ClusterUpsertBulk {
	return u.Update(func(s *ClusterUpsert) {
		s.SetDescription(v)
	})
}

// UpdateDescription sets the "description" field to the value that was provided on create.
func (u *ClusterUpsertBulk) UpdateDescription() *ClusterUpsertBulk {
	return u.Update(func(s *ClusterUpsert) {
		s.UpdateDescription()
	})
}

// SetContent sets the "content" field.
func (u *ClusterUpsertBulk) SetContent(v string) *ClusterUpsertBulk {
	return u.Update(func(s *ClusterUpsert) {
		s.SetContent(v)
	})
}

// UpdateContent sets the "content" field to the value that was provided on create.
func (u *ClusterUpsertBulk) UpdateContent() *ClusterUpsertBulk {
	return u.Update(func(s *ClusterUpsert) {
		s.UpdateContent()
	})
}

// ClearContent clears the value of the "content" field.
func (u *ClusterUpsertBulk) ClearContent() *ClusterUpsertBulk {
	return u.Update(func(s *ClusterUpsert) {
		s.ClearContent()
	})
}

// SetParentClusterID sets the "parent_cluster_id" field.
func (u *ClusterUpsertBulk) SetParentClusterID(v xid.ID) *ClusterUpsertBulk {
	return u.Update(func(s *ClusterUpsert) {
		s.SetParentClusterID(v)
	})
}

// UpdateParentClusterID sets the "parent_cluster_id" field to the value that was provided on create.
func (u *ClusterUpsertBulk) UpdateParentClusterID() *ClusterUpsertBulk {
	return u.Update(func(s *ClusterUpsert) {
		s.UpdateParentClusterID()
	})
}

// ClearParentClusterID clears the value of the "parent_cluster_id" field.
func (u *ClusterUpsertBulk) ClearParentClusterID() *ClusterUpsertBulk {
	return u.Update(func(s *ClusterUpsert) {
		s.ClearParentClusterID()
	})
}

// SetAccountID sets the "account_id" field.
func (u *ClusterUpsertBulk) SetAccountID(v xid.ID) *ClusterUpsertBulk {
	return u.Update(func(s *ClusterUpsert) {
		s.SetAccountID(v)
	})
}

// UpdateAccountID sets the "account_id" field to the value that was provided on create.
func (u *ClusterUpsertBulk) UpdateAccountID() *ClusterUpsertBulk {
	return u.Update(func(s *ClusterUpsert) {
		s.UpdateAccountID()
	})
}

// SetProperties sets the "properties" field.
func (u *ClusterUpsertBulk) SetProperties(v any) *ClusterUpsertBulk {
	return u.Update(func(s *ClusterUpsert) {
		s.SetProperties(v)
	})
}

// UpdateProperties sets the "properties" field to the value that was provided on create.
func (u *ClusterUpsertBulk) UpdateProperties() *ClusterUpsertBulk {
	return u.Update(func(s *ClusterUpsert) {
		s.UpdateProperties()
	})
}

// ClearProperties clears the value of the "properties" field.
func (u *ClusterUpsertBulk) ClearProperties() *ClusterUpsertBulk {
	return u.Update(func(s *ClusterUpsert) {
		s.ClearProperties()
	})
}

// Exec executes the query.
func (u *ClusterUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the ClusterCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for ClusterCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *ClusterUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
