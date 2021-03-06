// Code generated by entc, DO NOT EDIT.

package model

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/Southclaws/storyden/backend/internal/infrastructure/db/model/post"
	"github.com/Southclaws/storyden/backend/internal/infrastructure/db/model/react"
	"github.com/Southclaws/storyden/backend/internal/infrastructure/db/model/user"
	"github.com/google/uuid"
)

// ReactCreate is the builder for creating a React entity.
type ReactCreate struct {
	config
	mutation *ReactMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetEmoji sets the "emoji" field.
func (rc *ReactCreate) SetEmoji(s string) *ReactCreate {
	rc.mutation.SetEmoji(s)
	return rc
}

// SetCreatedAt sets the "createdAt" field.
func (rc *ReactCreate) SetCreatedAt(t time.Time) *ReactCreate {
	rc.mutation.SetCreatedAt(t)
	return rc
}

// SetNillableCreatedAt sets the "createdAt" field if the given value is not nil.
func (rc *ReactCreate) SetNillableCreatedAt(t *time.Time) *ReactCreate {
	if t != nil {
		rc.SetCreatedAt(*t)
	}
	return rc
}

// SetID sets the "id" field.
func (rc *ReactCreate) SetID(u uuid.UUID) *ReactCreate {
	rc.mutation.SetID(u)
	return rc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (rc *ReactCreate) SetNillableID(u *uuid.UUID) *ReactCreate {
	if u != nil {
		rc.SetID(*u)
	}
	return rc
}

// SetUserID sets the "user" edge to the User entity by ID.
func (rc *ReactCreate) SetUserID(id uuid.UUID) *ReactCreate {
	rc.mutation.SetUserID(id)
	return rc
}

// SetNillableUserID sets the "user" edge to the User entity by ID if the given value is not nil.
func (rc *ReactCreate) SetNillableUserID(id *uuid.UUID) *ReactCreate {
	if id != nil {
		rc = rc.SetUserID(*id)
	}
	return rc
}

// SetUser sets the "user" edge to the User entity.
func (rc *ReactCreate) SetUser(u *User) *ReactCreate {
	return rc.SetUserID(u.ID)
}

// SetPostID sets the "Post" edge to the Post entity by ID.
func (rc *ReactCreate) SetPostID(id uuid.UUID) *ReactCreate {
	rc.mutation.SetPostID(id)
	return rc
}

// SetNillablePostID sets the "Post" edge to the Post entity by ID if the given value is not nil.
func (rc *ReactCreate) SetNillablePostID(id *uuid.UUID) *ReactCreate {
	if id != nil {
		rc = rc.SetPostID(*id)
	}
	return rc
}

// SetPost sets the "Post" edge to the Post entity.
func (rc *ReactCreate) SetPost(p *Post) *ReactCreate {
	return rc.SetPostID(p.ID)
}

// Mutation returns the ReactMutation object of the builder.
func (rc *ReactCreate) Mutation() *ReactMutation {
	return rc.mutation
}

// Save creates the React in the database.
func (rc *ReactCreate) Save(ctx context.Context) (*React, error) {
	var (
		err  error
		node *React
	)
	rc.defaults()
	if len(rc.hooks) == 0 {
		if err = rc.check(); err != nil {
			return nil, err
		}
		node, err = rc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ReactMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = rc.check(); err != nil {
				return nil, err
			}
			rc.mutation = mutation
			if node, err = rc.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(rc.hooks) - 1; i >= 0; i-- {
			if rc.hooks[i] == nil {
				return nil, fmt.Errorf("model: uninitialized hook (forgotten import model/runtime?)")
			}
			mut = rc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, rc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (rc *ReactCreate) SaveX(ctx context.Context) *React {
	v, err := rc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (rc *ReactCreate) Exec(ctx context.Context) error {
	_, err := rc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (rc *ReactCreate) ExecX(ctx context.Context) {
	if err := rc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (rc *ReactCreate) defaults() {
	if _, ok := rc.mutation.CreatedAt(); !ok {
		v := react.DefaultCreatedAt()
		rc.mutation.SetCreatedAt(v)
	}
	if _, ok := rc.mutation.ID(); !ok {
		v := react.DefaultID()
		rc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (rc *ReactCreate) check() error {
	if _, ok := rc.mutation.Emoji(); !ok {
		return &ValidationError{Name: "emoji", err: errors.New(`model: missing required field "React.emoji"`)}
	}
	if _, ok := rc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "createdAt", err: errors.New(`model: missing required field "React.createdAt"`)}
	}
	return nil
}

func (rc *ReactCreate) sqlSave(ctx context.Context) (*React, error) {
	_node, _spec := rc.createSpec()
	if err := sqlgraph.CreateNode(ctx, rc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(*uuid.UUID); ok {
			_node.ID = *id
		} else if err := _node.ID.Scan(_spec.ID.Value); err != nil {
			return nil, err
		}
	}
	return _node, nil
}

func (rc *ReactCreate) createSpec() (*React, *sqlgraph.CreateSpec) {
	var (
		_node = &React{config: rc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: react.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: react.FieldID,
			},
		}
	)
	_spec.OnConflict = rc.conflict
	if id, ok := rc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := rc.mutation.Emoji(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: react.FieldEmoji,
		})
		_node.Emoji = value
	}
	if value, ok := rc.mutation.CreatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: react.FieldCreatedAt,
		})
		_node.CreatedAt = value
	}
	if nodes := rc.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   react.UserTable,
			Columns: []string{react.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: user.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.react_user = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := rc.mutation.PostIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   react.PostTable,
			Columns: []string{react.PostColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: post.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.react_post = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.React.Create().
//		SetEmoji(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.ReactUpsert) {
//			SetEmoji(v+v).
//		}).
//		Exec(ctx)
//
func (rc *ReactCreate) OnConflict(opts ...sql.ConflictOption) *ReactUpsertOne {
	rc.conflict = opts
	return &ReactUpsertOne{
		create: rc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.React.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
//
func (rc *ReactCreate) OnConflictColumns(columns ...string) *ReactUpsertOne {
	rc.conflict = append(rc.conflict, sql.ConflictColumns(columns...))
	return &ReactUpsertOne{
		create: rc,
	}
}

type (
	// ReactUpsertOne is the builder for "upsert"-ing
	//  one React node.
	ReactUpsertOne struct {
		create *ReactCreate
	}

	// ReactUpsert is the "OnConflict" setter.
	ReactUpsert struct {
		*sql.UpdateSet
	}
)

// SetEmoji sets the "emoji" field.
func (u *ReactUpsert) SetEmoji(v string) *ReactUpsert {
	u.Set(react.FieldEmoji, v)
	return u
}

// UpdateEmoji sets the "emoji" field to the value that was provided on create.
func (u *ReactUpsert) UpdateEmoji() *ReactUpsert {
	u.SetExcluded(react.FieldEmoji)
	return u
}

// SetCreatedAt sets the "createdAt" field.
func (u *ReactUpsert) SetCreatedAt(v time.Time) *ReactUpsert {
	u.Set(react.FieldCreatedAt, v)
	return u
}

// UpdateCreatedAt sets the "createdAt" field to the value that was provided on create.
func (u *ReactUpsert) UpdateCreatedAt() *ReactUpsert {
	u.SetExcluded(react.FieldCreatedAt)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.React.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(react.FieldID)
//			}),
//		).
//		Exec(ctx)
//
func (u *ReactUpsertOne) UpdateNewValues() *ReactUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(react.FieldID)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//  client.React.Create().
//      OnConflict(sql.ResolveWithIgnore()).
//      Exec(ctx)
//
func (u *ReactUpsertOne) Ignore() *ReactUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *ReactUpsertOne) DoNothing() *ReactUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the ReactCreate.OnConflict
// documentation for more info.
func (u *ReactUpsertOne) Update(set func(*ReactUpsert)) *ReactUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&ReactUpsert{UpdateSet: update})
	}))
	return u
}

// SetEmoji sets the "emoji" field.
func (u *ReactUpsertOne) SetEmoji(v string) *ReactUpsertOne {
	return u.Update(func(s *ReactUpsert) {
		s.SetEmoji(v)
	})
}

// UpdateEmoji sets the "emoji" field to the value that was provided on create.
func (u *ReactUpsertOne) UpdateEmoji() *ReactUpsertOne {
	return u.Update(func(s *ReactUpsert) {
		s.UpdateEmoji()
	})
}

// SetCreatedAt sets the "createdAt" field.
func (u *ReactUpsertOne) SetCreatedAt(v time.Time) *ReactUpsertOne {
	return u.Update(func(s *ReactUpsert) {
		s.SetCreatedAt(v)
	})
}

// UpdateCreatedAt sets the "createdAt" field to the value that was provided on create.
func (u *ReactUpsertOne) UpdateCreatedAt() *ReactUpsertOne {
	return u.Update(func(s *ReactUpsert) {
		s.UpdateCreatedAt()
	})
}

// Exec executes the query.
func (u *ReactUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("model: missing options for ReactCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *ReactUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *ReactUpsertOne) ID(ctx context.Context) (id uuid.UUID, err error) {
	if u.create.driver.Dialect() == dialect.MySQL {
		// In case of "ON CONFLICT", there is no way to get back non-numeric ID
		// fields from the database since MySQL does not support the RETURNING clause.
		return id, errors.New("model: ReactUpsertOne.ID is not supported by MySQL driver. Use ReactUpsertOne.Exec instead")
	}
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *ReactUpsertOne) IDX(ctx context.Context) uuid.UUID {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// ReactCreateBulk is the builder for creating many React entities in bulk.
type ReactCreateBulk struct {
	config
	builders []*ReactCreate
	conflict []sql.ConflictOption
}

// Save creates the React entities in the database.
func (rcb *ReactCreateBulk) Save(ctx context.Context) ([]*React, error) {
	specs := make([]*sqlgraph.CreateSpec, len(rcb.builders))
	nodes := make([]*React, len(rcb.builders))
	mutators := make([]Mutator, len(rcb.builders))
	for i := range rcb.builders {
		func(i int, root context.Context) {
			builder := rcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ReactMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, rcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = rcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, rcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{err.Error(), err}
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
		if _, err := mutators[0].Mutate(ctx, rcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (rcb *ReactCreateBulk) SaveX(ctx context.Context) []*React {
	v, err := rcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (rcb *ReactCreateBulk) Exec(ctx context.Context) error {
	_, err := rcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (rcb *ReactCreateBulk) ExecX(ctx context.Context) {
	if err := rcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.React.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.ReactUpsert) {
//			SetEmoji(v+v).
//		}).
//		Exec(ctx)
//
func (rcb *ReactCreateBulk) OnConflict(opts ...sql.ConflictOption) *ReactUpsertBulk {
	rcb.conflict = opts
	return &ReactUpsertBulk{
		create: rcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.React.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
//
func (rcb *ReactCreateBulk) OnConflictColumns(columns ...string) *ReactUpsertBulk {
	rcb.conflict = append(rcb.conflict, sql.ConflictColumns(columns...))
	return &ReactUpsertBulk{
		create: rcb,
	}
}

// ReactUpsertBulk is the builder for "upsert"-ing
// a bulk of React nodes.
type ReactUpsertBulk struct {
	create *ReactCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.React.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(react.FieldID)
//			}),
//		).
//		Exec(ctx)
//
func (u *ReactUpsertBulk) UpdateNewValues() *ReactUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(react.FieldID)
				return
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.React.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
//
func (u *ReactUpsertBulk) Ignore() *ReactUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *ReactUpsertBulk) DoNothing() *ReactUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the ReactCreateBulk.OnConflict
// documentation for more info.
func (u *ReactUpsertBulk) Update(set func(*ReactUpsert)) *ReactUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&ReactUpsert{UpdateSet: update})
	}))
	return u
}

// SetEmoji sets the "emoji" field.
func (u *ReactUpsertBulk) SetEmoji(v string) *ReactUpsertBulk {
	return u.Update(func(s *ReactUpsert) {
		s.SetEmoji(v)
	})
}

// UpdateEmoji sets the "emoji" field to the value that was provided on create.
func (u *ReactUpsertBulk) UpdateEmoji() *ReactUpsertBulk {
	return u.Update(func(s *ReactUpsert) {
		s.UpdateEmoji()
	})
}

// SetCreatedAt sets the "createdAt" field.
func (u *ReactUpsertBulk) SetCreatedAt(v time.Time) *ReactUpsertBulk {
	return u.Update(func(s *ReactUpsert) {
		s.SetCreatedAt(v)
	})
}

// UpdateCreatedAt sets the "createdAt" field to the value that was provided on create.
func (u *ReactUpsertBulk) UpdateCreatedAt() *ReactUpsertBulk {
	return u.Update(func(s *ReactUpsert) {
		s.UpdateCreatedAt()
	})
}

// Exec executes the query.
func (u *ReactUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("model: OnConflict was set for builder %d. Set it on the ReactCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("model: missing options for ReactCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *ReactUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
