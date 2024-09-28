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
	"github.com/Southclaws/storyden/internal/ent/collection"
	"github.com/Southclaws/storyden/internal/ent/collectionnode"
	"github.com/Southclaws/storyden/internal/ent/node"
	"github.com/rs/xid"
)

// CollectionNodeCreate is the builder for creating a CollectionNode entity.
type CollectionNodeCreate struct {
	config
	mutation *CollectionNodeMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCreatedAt sets the "created_at" field.
func (cnc *CollectionNodeCreate) SetCreatedAt(t time.Time) *CollectionNodeCreate {
	cnc.mutation.SetCreatedAt(t)
	return cnc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (cnc *CollectionNodeCreate) SetNillableCreatedAt(t *time.Time) *CollectionNodeCreate {
	if t != nil {
		cnc.SetCreatedAt(*t)
	}
	return cnc
}

// SetCollectionID sets the "collection_id" field.
func (cnc *CollectionNodeCreate) SetCollectionID(x xid.ID) *CollectionNodeCreate {
	cnc.mutation.SetCollectionID(x)
	return cnc
}

// SetNillableCollectionID sets the "collection_id" field if the given value is not nil.
func (cnc *CollectionNodeCreate) SetNillableCollectionID(x *xid.ID) *CollectionNodeCreate {
	if x != nil {
		cnc.SetCollectionID(*x)
	}
	return cnc
}

// SetNodeID sets the "node_id" field.
func (cnc *CollectionNodeCreate) SetNodeID(x xid.ID) *CollectionNodeCreate {
	cnc.mutation.SetNodeID(x)
	return cnc
}

// SetNillableNodeID sets the "node_id" field if the given value is not nil.
func (cnc *CollectionNodeCreate) SetNillableNodeID(x *xid.ID) *CollectionNodeCreate {
	if x != nil {
		cnc.SetNodeID(*x)
	}
	return cnc
}

// SetMembershipType sets the "membership_type" field.
func (cnc *CollectionNodeCreate) SetMembershipType(s string) *CollectionNodeCreate {
	cnc.mutation.SetMembershipType(s)
	return cnc
}

// SetNillableMembershipType sets the "membership_type" field if the given value is not nil.
func (cnc *CollectionNodeCreate) SetNillableMembershipType(s *string) *CollectionNodeCreate {
	if s != nil {
		cnc.SetMembershipType(*s)
	}
	return cnc
}

// SetCollection sets the "collection" edge to the Collection entity.
func (cnc *CollectionNodeCreate) SetCollection(c *Collection) *CollectionNodeCreate {
	return cnc.SetCollectionID(c.ID)
}

// SetNode sets the "node" edge to the Node entity.
func (cnc *CollectionNodeCreate) SetNode(n *Node) *CollectionNodeCreate {
	return cnc.SetNodeID(n.ID)
}

// Mutation returns the CollectionNodeMutation object of the builder.
func (cnc *CollectionNodeCreate) Mutation() *CollectionNodeMutation {
	return cnc.mutation
}

// Save creates the CollectionNode in the database.
func (cnc *CollectionNodeCreate) Save(ctx context.Context) (*CollectionNode, error) {
	cnc.defaults()
	return withHooks(ctx, cnc.sqlSave, cnc.mutation, cnc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (cnc *CollectionNodeCreate) SaveX(ctx context.Context) *CollectionNode {
	v, err := cnc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (cnc *CollectionNodeCreate) Exec(ctx context.Context) error {
	_, err := cnc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cnc *CollectionNodeCreate) ExecX(ctx context.Context) {
	if err := cnc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (cnc *CollectionNodeCreate) defaults() {
	if _, ok := cnc.mutation.CreatedAt(); !ok {
		v := collectionnode.DefaultCreatedAt()
		cnc.mutation.SetCreatedAt(v)
	}
	if _, ok := cnc.mutation.CollectionID(); !ok {
		v := collectionnode.DefaultCollectionID()
		cnc.mutation.SetCollectionID(v)
	}
	if _, ok := cnc.mutation.NodeID(); !ok {
		v := collectionnode.DefaultNodeID()
		cnc.mutation.SetNodeID(v)
	}
	if _, ok := cnc.mutation.MembershipType(); !ok {
		v := collectionnode.DefaultMembershipType
		cnc.mutation.SetMembershipType(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (cnc *CollectionNodeCreate) check() error {
	if _, ok := cnc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "CollectionNode.created_at"`)}
	}
	if _, ok := cnc.mutation.CollectionID(); !ok {
		return &ValidationError{Name: "collection_id", err: errors.New(`ent: missing required field "CollectionNode.collection_id"`)}
	}
	if v, ok := cnc.mutation.CollectionID(); ok {
		if err := collectionnode.CollectionIDValidator(v.String()); err != nil {
			return &ValidationError{Name: "collection_id", err: fmt.Errorf(`ent: validator failed for field "CollectionNode.collection_id": %w`, err)}
		}
	}
	if _, ok := cnc.mutation.NodeID(); !ok {
		return &ValidationError{Name: "node_id", err: errors.New(`ent: missing required field "CollectionNode.node_id"`)}
	}
	if v, ok := cnc.mutation.NodeID(); ok {
		if err := collectionnode.NodeIDValidator(v.String()); err != nil {
			return &ValidationError{Name: "node_id", err: fmt.Errorf(`ent: validator failed for field "CollectionNode.node_id": %w`, err)}
		}
	}
	if _, ok := cnc.mutation.MembershipType(); !ok {
		return &ValidationError{Name: "membership_type", err: errors.New(`ent: missing required field "CollectionNode.membership_type"`)}
	}
	if len(cnc.mutation.CollectionIDs()) == 0 {
		return &ValidationError{Name: "collection", err: errors.New(`ent: missing required edge "CollectionNode.collection"`)}
	}
	if len(cnc.mutation.NodeIDs()) == 0 {
		return &ValidationError{Name: "node", err: errors.New(`ent: missing required edge "CollectionNode.node"`)}
	}
	return nil
}

func (cnc *CollectionNodeCreate) sqlSave(ctx context.Context) (*CollectionNode, error) {
	if err := cnc.check(); err != nil {
		return nil, err
	}
	_node, _spec := cnc.createSpec()
	if err := sqlgraph.CreateNode(ctx, cnc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	return _node, nil
}

func (cnc *CollectionNodeCreate) createSpec() (*CollectionNode, *sqlgraph.CreateSpec) {
	var (
		_node = &CollectionNode{config: cnc.config}
		_spec = sqlgraph.NewCreateSpec(collectionnode.Table, nil)
	)
	_spec.OnConflict = cnc.conflict
	if value, ok := cnc.mutation.CreatedAt(); ok {
		_spec.SetField(collectionnode.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := cnc.mutation.MembershipType(); ok {
		_spec.SetField(collectionnode.FieldMembershipType, field.TypeString, value)
		_node.MembershipType = value
	}
	if nodes := cnc.mutation.CollectionIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   collectionnode.CollectionTable,
			Columns: []string{collectionnode.CollectionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(collection.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.CollectionID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := cnc.mutation.NodeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   collectionnode.NodeTable,
			Columns: []string{collectionnode.NodeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(node.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.NodeID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.CollectionNode.Create().
//		SetCreatedAt(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.CollectionNodeUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
func (cnc *CollectionNodeCreate) OnConflict(opts ...sql.ConflictOption) *CollectionNodeUpsertOne {
	cnc.conflict = opts
	return &CollectionNodeUpsertOne{
		create: cnc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.CollectionNode.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (cnc *CollectionNodeCreate) OnConflictColumns(columns ...string) *CollectionNodeUpsertOne {
	cnc.conflict = append(cnc.conflict, sql.ConflictColumns(columns...))
	return &CollectionNodeUpsertOne{
		create: cnc,
	}
}

type (
	// CollectionNodeUpsertOne is the builder for "upsert"-ing
	//  one CollectionNode node.
	CollectionNodeUpsertOne struct {
		create *CollectionNodeCreate
	}

	// CollectionNodeUpsert is the "OnConflict" setter.
	CollectionNodeUpsert struct {
		*sql.UpdateSet
	}
)

// SetCollectionID sets the "collection_id" field.
func (u *CollectionNodeUpsert) SetCollectionID(v xid.ID) *CollectionNodeUpsert {
	u.Set(collectionnode.FieldCollectionID, v)
	return u
}

// UpdateCollectionID sets the "collection_id" field to the value that was provided on create.
func (u *CollectionNodeUpsert) UpdateCollectionID() *CollectionNodeUpsert {
	u.SetExcluded(collectionnode.FieldCollectionID)
	return u
}

// SetNodeID sets the "node_id" field.
func (u *CollectionNodeUpsert) SetNodeID(v xid.ID) *CollectionNodeUpsert {
	u.Set(collectionnode.FieldNodeID, v)
	return u
}

// UpdateNodeID sets the "node_id" field to the value that was provided on create.
func (u *CollectionNodeUpsert) UpdateNodeID() *CollectionNodeUpsert {
	u.SetExcluded(collectionnode.FieldNodeID)
	return u
}

// SetMembershipType sets the "membership_type" field.
func (u *CollectionNodeUpsert) SetMembershipType(v string) *CollectionNodeUpsert {
	u.Set(collectionnode.FieldMembershipType, v)
	return u
}

// UpdateMembershipType sets the "membership_type" field to the value that was provided on create.
func (u *CollectionNodeUpsert) UpdateMembershipType() *CollectionNodeUpsert {
	u.SetExcluded(collectionnode.FieldMembershipType)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.CollectionNode.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *CollectionNodeUpsertOne) UpdateNewValues() *CollectionNodeUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.CreatedAt(); exists {
			s.SetIgnore(collectionnode.FieldCreatedAt)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.CollectionNode.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *CollectionNodeUpsertOne) Ignore() *CollectionNodeUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *CollectionNodeUpsertOne) DoNothing() *CollectionNodeUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the CollectionNodeCreate.OnConflict
// documentation for more info.
func (u *CollectionNodeUpsertOne) Update(set func(*CollectionNodeUpsert)) *CollectionNodeUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&CollectionNodeUpsert{UpdateSet: update})
	}))
	return u
}

// SetCollectionID sets the "collection_id" field.
func (u *CollectionNodeUpsertOne) SetCollectionID(v xid.ID) *CollectionNodeUpsertOne {
	return u.Update(func(s *CollectionNodeUpsert) {
		s.SetCollectionID(v)
	})
}

// UpdateCollectionID sets the "collection_id" field to the value that was provided on create.
func (u *CollectionNodeUpsertOne) UpdateCollectionID() *CollectionNodeUpsertOne {
	return u.Update(func(s *CollectionNodeUpsert) {
		s.UpdateCollectionID()
	})
}

// SetNodeID sets the "node_id" field.
func (u *CollectionNodeUpsertOne) SetNodeID(v xid.ID) *CollectionNodeUpsertOne {
	return u.Update(func(s *CollectionNodeUpsert) {
		s.SetNodeID(v)
	})
}

// UpdateNodeID sets the "node_id" field to the value that was provided on create.
func (u *CollectionNodeUpsertOne) UpdateNodeID() *CollectionNodeUpsertOne {
	return u.Update(func(s *CollectionNodeUpsert) {
		s.UpdateNodeID()
	})
}

// SetMembershipType sets the "membership_type" field.
func (u *CollectionNodeUpsertOne) SetMembershipType(v string) *CollectionNodeUpsertOne {
	return u.Update(func(s *CollectionNodeUpsert) {
		s.SetMembershipType(v)
	})
}

// UpdateMembershipType sets the "membership_type" field to the value that was provided on create.
func (u *CollectionNodeUpsertOne) UpdateMembershipType() *CollectionNodeUpsertOne {
	return u.Update(func(s *CollectionNodeUpsert) {
		s.UpdateMembershipType()
	})
}

// Exec executes the query.
func (u *CollectionNodeUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for CollectionNodeCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *CollectionNodeUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// CollectionNodeCreateBulk is the builder for creating many CollectionNode entities in bulk.
type CollectionNodeCreateBulk struct {
	config
	err      error
	builders []*CollectionNodeCreate
	conflict []sql.ConflictOption
}

// Save creates the CollectionNode entities in the database.
func (cncb *CollectionNodeCreateBulk) Save(ctx context.Context) ([]*CollectionNode, error) {
	if cncb.err != nil {
		return nil, cncb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(cncb.builders))
	nodes := make([]*CollectionNode, len(cncb.builders))
	mutators := make([]Mutator, len(cncb.builders))
	for i := range cncb.builders {
		func(i int, root context.Context) {
			builder := cncb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*CollectionNodeMutation)
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
					_, err = mutators[i+1].Mutate(root, cncb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = cncb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, cncb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
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
		if _, err := mutators[0].Mutate(ctx, cncb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (cncb *CollectionNodeCreateBulk) SaveX(ctx context.Context) []*CollectionNode {
	v, err := cncb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (cncb *CollectionNodeCreateBulk) Exec(ctx context.Context) error {
	_, err := cncb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cncb *CollectionNodeCreateBulk) ExecX(ctx context.Context) {
	if err := cncb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.CollectionNode.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.CollectionNodeUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
func (cncb *CollectionNodeCreateBulk) OnConflict(opts ...sql.ConflictOption) *CollectionNodeUpsertBulk {
	cncb.conflict = opts
	return &CollectionNodeUpsertBulk{
		create: cncb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.CollectionNode.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (cncb *CollectionNodeCreateBulk) OnConflictColumns(columns ...string) *CollectionNodeUpsertBulk {
	cncb.conflict = append(cncb.conflict, sql.ConflictColumns(columns...))
	return &CollectionNodeUpsertBulk{
		create: cncb,
	}
}

// CollectionNodeUpsertBulk is the builder for "upsert"-ing
// a bulk of CollectionNode nodes.
type CollectionNodeUpsertBulk struct {
	create *CollectionNodeCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.CollectionNode.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *CollectionNodeUpsertBulk) UpdateNewValues() *CollectionNodeUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.CreatedAt(); exists {
				s.SetIgnore(collectionnode.FieldCreatedAt)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.CollectionNode.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *CollectionNodeUpsertBulk) Ignore() *CollectionNodeUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *CollectionNodeUpsertBulk) DoNothing() *CollectionNodeUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the CollectionNodeCreateBulk.OnConflict
// documentation for more info.
func (u *CollectionNodeUpsertBulk) Update(set func(*CollectionNodeUpsert)) *CollectionNodeUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&CollectionNodeUpsert{UpdateSet: update})
	}))
	return u
}

// SetCollectionID sets the "collection_id" field.
func (u *CollectionNodeUpsertBulk) SetCollectionID(v xid.ID) *CollectionNodeUpsertBulk {
	return u.Update(func(s *CollectionNodeUpsert) {
		s.SetCollectionID(v)
	})
}

// UpdateCollectionID sets the "collection_id" field to the value that was provided on create.
func (u *CollectionNodeUpsertBulk) UpdateCollectionID() *CollectionNodeUpsertBulk {
	return u.Update(func(s *CollectionNodeUpsert) {
		s.UpdateCollectionID()
	})
}

// SetNodeID sets the "node_id" field.
func (u *CollectionNodeUpsertBulk) SetNodeID(v xid.ID) *CollectionNodeUpsertBulk {
	return u.Update(func(s *CollectionNodeUpsert) {
		s.SetNodeID(v)
	})
}

// UpdateNodeID sets the "node_id" field to the value that was provided on create.
func (u *CollectionNodeUpsertBulk) UpdateNodeID() *CollectionNodeUpsertBulk {
	return u.Update(func(s *CollectionNodeUpsert) {
		s.UpdateNodeID()
	})
}

// SetMembershipType sets the "membership_type" field.
func (u *CollectionNodeUpsertBulk) SetMembershipType(v string) *CollectionNodeUpsertBulk {
	return u.Update(func(s *CollectionNodeUpsert) {
		s.SetMembershipType(v)
	})
}

// UpdateMembershipType sets the "membership_type" field to the value that was provided on create.
func (u *CollectionNodeUpsertBulk) UpdateMembershipType() *CollectionNodeUpsertBulk {
	return u.Update(func(s *CollectionNodeUpsert) {
		s.UpdateMembershipType()
	})
}

// Exec executes the query.
func (u *CollectionNodeUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the CollectionNodeCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for CollectionNodeCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *CollectionNodeUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
