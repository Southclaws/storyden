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
	"github.com/Southclaws/storyden/internal/ent/accountfollow"
	"github.com/rs/xid"
)

// AccountFollowCreate is the builder for creating a AccountFollow entity.
type AccountFollowCreate struct {
	config
	mutation *AccountFollowMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCreatedAt sets the "created_at" field.
func (afc *AccountFollowCreate) SetCreatedAt(t time.Time) *AccountFollowCreate {
	afc.mutation.SetCreatedAt(t)
	return afc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (afc *AccountFollowCreate) SetNillableCreatedAt(t *time.Time) *AccountFollowCreate {
	if t != nil {
		afc.SetCreatedAt(*t)
	}
	return afc
}

// SetFollowerAccountID sets the "follower_account_id" field.
func (afc *AccountFollowCreate) SetFollowerAccountID(x xid.ID) *AccountFollowCreate {
	afc.mutation.SetFollowerAccountID(x)
	return afc
}

// SetFollowingAccountID sets the "following_account_id" field.
func (afc *AccountFollowCreate) SetFollowingAccountID(x xid.ID) *AccountFollowCreate {
	afc.mutation.SetFollowingAccountID(x)
	return afc
}

// SetID sets the "id" field.
func (afc *AccountFollowCreate) SetID(x xid.ID) *AccountFollowCreate {
	afc.mutation.SetID(x)
	return afc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (afc *AccountFollowCreate) SetNillableID(x *xid.ID) *AccountFollowCreate {
	if x != nil {
		afc.SetID(*x)
	}
	return afc
}

// SetFollowerID sets the "follower" edge to the Account entity by ID.
func (afc *AccountFollowCreate) SetFollowerID(id xid.ID) *AccountFollowCreate {
	afc.mutation.SetFollowerID(id)
	return afc
}

// SetFollower sets the "follower" edge to the Account entity.
func (afc *AccountFollowCreate) SetFollower(a *Account) *AccountFollowCreate {
	return afc.SetFollowerID(a.ID)
}

// SetFollowingID sets the "following" edge to the Account entity by ID.
func (afc *AccountFollowCreate) SetFollowingID(id xid.ID) *AccountFollowCreate {
	afc.mutation.SetFollowingID(id)
	return afc
}

// SetFollowing sets the "following" edge to the Account entity.
func (afc *AccountFollowCreate) SetFollowing(a *Account) *AccountFollowCreate {
	return afc.SetFollowingID(a.ID)
}

// Mutation returns the AccountFollowMutation object of the builder.
func (afc *AccountFollowCreate) Mutation() *AccountFollowMutation {
	return afc.mutation
}

// Save creates the AccountFollow in the database.
func (afc *AccountFollowCreate) Save(ctx context.Context) (*AccountFollow, error) {
	afc.defaults()
	return withHooks(ctx, afc.sqlSave, afc.mutation, afc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (afc *AccountFollowCreate) SaveX(ctx context.Context) *AccountFollow {
	v, err := afc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (afc *AccountFollowCreate) Exec(ctx context.Context) error {
	_, err := afc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (afc *AccountFollowCreate) ExecX(ctx context.Context) {
	if err := afc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (afc *AccountFollowCreate) defaults() {
	if _, ok := afc.mutation.CreatedAt(); !ok {
		v := accountfollow.DefaultCreatedAt()
		afc.mutation.SetCreatedAt(v)
	}
	if _, ok := afc.mutation.ID(); !ok {
		v := accountfollow.DefaultID()
		afc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (afc *AccountFollowCreate) check() error {
	if _, ok := afc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "AccountFollow.created_at"`)}
	}
	if _, ok := afc.mutation.FollowerAccountID(); !ok {
		return &ValidationError{Name: "follower_account_id", err: errors.New(`ent: missing required field "AccountFollow.follower_account_id"`)}
	}
	if _, ok := afc.mutation.FollowingAccountID(); !ok {
		return &ValidationError{Name: "following_account_id", err: errors.New(`ent: missing required field "AccountFollow.following_account_id"`)}
	}
	if v, ok := afc.mutation.ID(); ok {
		if err := accountfollow.IDValidator(v.String()); err != nil {
			return &ValidationError{Name: "id", err: fmt.Errorf(`ent: validator failed for field "AccountFollow.id": %w`, err)}
		}
	}
	if len(afc.mutation.FollowerIDs()) == 0 {
		return &ValidationError{Name: "follower", err: errors.New(`ent: missing required edge "AccountFollow.follower"`)}
	}
	if len(afc.mutation.FollowingIDs()) == 0 {
		return &ValidationError{Name: "following", err: errors.New(`ent: missing required edge "AccountFollow.following"`)}
	}
	return nil
}

func (afc *AccountFollowCreate) sqlSave(ctx context.Context) (*AccountFollow, error) {
	if err := afc.check(); err != nil {
		return nil, err
	}
	_node, _spec := afc.createSpec()
	if err := sqlgraph.CreateNode(ctx, afc.driver, _spec); err != nil {
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
	afc.mutation.id = &_node.ID
	afc.mutation.done = true
	return _node, nil
}

func (afc *AccountFollowCreate) createSpec() (*AccountFollow, *sqlgraph.CreateSpec) {
	var (
		_node = &AccountFollow{config: afc.config}
		_spec = sqlgraph.NewCreateSpec(accountfollow.Table, sqlgraph.NewFieldSpec(accountfollow.FieldID, field.TypeString))
	)
	_spec.OnConflict = afc.conflict
	if id, ok := afc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := afc.mutation.CreatedAt(); ok {
		_spec.SetField(accountfollow.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if nodes := afc.mutation.FollowerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   accountfollow.FollowerTable,
			Columns: []string{accountfollow.FollowerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(account.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.FollowerAccountID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := afc.mutation.FollowingIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   accountfollow.FollowingTable,
			Columns: []string{accountfollow.FollowingColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(account.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.FollowingAccountID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.AccountFollow.Create().
//		SetCreatedAt(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.AccountFollowUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
func (afc *AccountFollowCreate) OnConflict(opts ...sql.ConflictOption) *AccountFollowUpsertOne {
	afc.conflict = opts
	return &AccountFollowUpsertOne{
		create: afc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.AccountFollow.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (afc *AccountFollowCreate) OnConflictColumns(columns ...string) *AccountFollowUpsertOne {
	afc.conflict = append(afc.conflict, sql.ConflictColumns(columns...))
	return &AccountFollowUpsertOne{
		create: afc,
	}
}

type (
	// AccountFollowUpsertOne is the builder for "upsert"-ing
	//  one AccountFollow node.
	AccountFollowUpsertOne struct {
		create *AccountFollowCreate
	}

	// AccountFollowUpsert is the "OnConflict" setter.
	AccountFollowUpsert struct {
		*sql.UpdateSet
	}
)

// SetFollowerAccountID sets the "follower_account_id" field.
func (u *AccountFollowUpsert) SetFollowerAccountID(v xid.ID) *AccountFollowUpsert {
	u.Set(accountfollow.FieldFollowerAccountID, v)
	return u
}

// UpdateFollowerAccountID sets the "follower_account_id" field to the value that was provided on create.
func (u *AccountFollowUpsert) UpdateFollowerAccountID() *AccountFollowUpsert {
	u.SetExcluded(accountfollow.FieldFollowerAccountID)
	return u
}

// SetFollowingAccountID sets the "following_account_id" field.
func (u *AccountFollowUpsert) SetFollowingAccountID(v xid.ID) *AccountFollowUpsert {
	u.Set(accountfollow.FieldFollowingAccountID, v)
	return u
}

// UpdateFollowingAccountID sets the "following_account_id" field to the value that was provided on create.
func (u *AccountFollowUpsert) UpdateFollowingAccountID() *AccountFollowUpsert {
	u.SetExcluded(accountfollow.FieldFollowingAccountID)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.AccountFollow.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(accountfollow.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *AccountFollowUpsertOne) UpdateNewValues() *AccountFollowUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(accountfollow.FieldID)
		}
		if _, exists := u.create.mutation.CreatedAt(); exists {
			s.SetIgnore(accountfollow.FieldCreatedAt)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.AccountFollow.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *AccountFollowUpsertOne) Ignore() *AccountFollowUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *AccountFollowUpsertOne) DoNothing() *AccountFollowUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the AccountFollowCreate.OnConflict
// documentation for more info.
func (u *AccountFollowUpsertOne) Update(set func(*AccountFollowUpsert)) *AccountFollowUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&AccountFollowUpsert{UpdateSet: update})
	}))
	return u
}

// SetFollowerAccountID sets the "follower_account_id" field.
func (u *AccountFollowUpsertOne) SetFollowerAccountID(v xid.ID) *AccountFollowUpsertOne {
	return u.Update(func(s *AccountFollowUpsert) {
		s.SetFollowerAccountID(v)
	})
}

// UpdateFollowerAccountID sets the "follower_account_id" field to the value that was provided on create.
func (u *AccountFollowUpsertOne) UpdateFollowerAccountID() *AccountFollowUpsertOne {
	return u.Update(func(s *AccountFollowUpsert) {
		s.UpdateFollowerAccountID()
	})
}

// SetFollowingAccountID sets the "following_account_id" field.
func (u *AccountFollowUpsertOne) SetFollowingAccountID(v xid.ID) *AccountFollowUpsertOne {
	return u.Update(func(s *AccountFollowUpsert) {
		s.SetFollowingAccountID(v)
	})
}

// UpdateFollowingAccountID sets the "following_account_id" field to the value that was provided on create.
func (u *AccountFollowUpsertOne) UpdateFollowingAccountID() *AccountFollowUpsertOne {
	return u.Update(func(s *AccountFollowUpsert) {
		s.UpdateFollowingAccountID()
	})
}

// Exec executes the query.
func (u *AccountFollowUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for AccountFollowCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *AccountFollowUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *AccountFollowUpsertOne) ID(ctx context.Context) (id xid.ID, err error) {
	if u.create.driver.Dialect() == dialect.MySQL {
		// In case of "ON CONFLICT", there is no way to get back non-numeric ID
		// fields from the database since MySQL does not support the RETURNING clause.
		return id, errors.New("ent: AccountFollowUpsertOne.ID is not supported by MySQL driver. Use AccountFollowUpsertOne.Exec instead")
	}
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *AccountFollowUpsertOne) IDX(ctx context.Context) xid.ID {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// AccountFollowCreateBulk is the builder for creating many AccountFollow entities in bulk.
type AccountFollowCreateBulk struct {
	config
	err      error
	builders []*AccountFollowCreate
	conflict []sql.ConflictOption
}

// Save creates the AccountFollow entities in the database.
func (afcb *AccountFollowCreateBulk) Save(ctx context.Context) ([]*AccountFollow, error) {
	if afcb.err != nil {
		return nil, afcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(afcb.builders))
	nodes := make([]*AccountFollow, len(afcb.builders))
	mutators := make([]Mutator, len(afcb.builders))
	for i := range afcb.builders {
		func(i int, root context.Context) {
			builder := afcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*AccountFollowMutation)
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
					_, err = mutators[i+1].Mutate(root, afcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = afcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, afcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, afcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (afcb *AccountFollowCreateBulk) SaveX(ctx context.Context) []*AccountFollow {
	v, err := afcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (afcb *AccountFollowCreateBulk) Exec(ctx context.Context) error {
	_, err := afcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (afcb *AccountFollowCreateBulk) ExecX(ctx context.Context) {
	if err := afcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.AccountFollow.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.AccountFollowUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
func (afcb *AccountFollowCreateBulk) OnConflict(opts ...sql.ConflictOption) *AccountFollowUpsertBulk {
	afcb.conflict = opts
	return &AccountFollowUpsertBulk{
		create: afcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.AccountFollow.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (afcb *AccountFollowCreateBulk) OnConflictColumns(columns ...string) *AccountFollowUpsertBulk {
	afcb.conflict = append(afcb.conflict, sql.ConflictColumns(columns...))
	return &AccountFollowUpsertBulk{
		create: afcb,
	}
}

// AccountFollowUpsertBulk is the builder for "upsert"-ing
// a bulk of AccountFollow nodes.
type AccountFollowUpsertBulk struct {
	create *AccountFollowCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.AccountFollow.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(accountfollow.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *AccountFollowUpsertBulk) UpdateNewValues() *AccountFollowUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(accountfollow.FieldID)
			}
			if _, exists := b.mutation.CreatedAt(); exists {
				s.SetIgnore(accountfollow.FieldCreatedAt)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.AccountFollow.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *AccountFollowUpsertBulk) Ignore() *AccountFollowUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *AccountFollowUpsertBulk) DoNothing() *AccountFollowUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the AccountFollowCreateBulk.OnConflict
// documentation for more info.
func (u *AccountFollowUpsertBulk) Update(set func(*AccountFollowUpsert)) *AccountFollowUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&AccountFollowUpsert{UpdateSet: update})
	}))
	return u
}

// SetFollowerAccountID sets the "follower_account_id" field.
func (u *AccountFollowUpsertBulk) SetFollowerAccountID(v xid.ID) *AccountFollowUpsertBulk {
	return u.Update(func(s *AccountFollowUpsert) {
		s.SetFollowerAccountID(v)
	})
}

// UpdateFollowerAccountID sets the "follower_account_id" field to the value that was provided on create.
func (u *AccountFollowUpsertBulk) UpdateFollowerAccountID() *AccountFollowUpsertBulk {
	return u.Update(func(s *AccountFollowUpsert) {
		s.UpdateFollowerAccountID()
	})
}

// SetFollowingAccountID sets the "following_account_id" field.
func (u *AccountFollowUpsertBulk) SetFollowingAccountID(v xid.ID) *AccountFollowUpsertBulk {
	return u.Update(func(s *AccountFollowUpsert) {
		s.SetFollowingAccountID(v)
	})
}

// UpdateFollowingAccountID sets the "following_account_id" field to the value that was provided on create.
func (u *AccountFollowUpsertBulk) UpdateFollowingAccountID() *AccountFollowUpsertBulk {
	return u.Update(func(s *AccountFollowUpsert) {
		s.UpdateFollowingAccountID()
	})
}

// Exec executes the query.
func (u *AccountFollowUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the AccountFollowCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for AccountFollowCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *AccountFollowUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
