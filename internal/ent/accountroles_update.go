// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/Southclaws/storyden/internal/ent/account"
	"github.com/Southclaws/storyden/internal/ent/accountroles"
	"github.com/Southclaws/storyden/internal/ent/predicate"
	"github.com/Southclaws/storyden/internal/ent/role"
	"github.com/rs/xid"
)

// AccountRolesUpdate is the builder for updating AccountRoles entities.
type AccountRolesUpdate struct {
	config
	hooks     []Hook
	mutation  *AccountRolesMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the AccountRolesUpdate builder.
func (aru *AccountRolesUpdate) Where(ps ...predicate.AccountRoles) *AccountRolesUpdate {
	aru.mutation.Where(ps...)
	return aru
}

// SetAccountID sets the "account_id" field.
func (aru *AccountRolesUpdate) SetAccountID(x xid.ID) *AccountRolesUpdate {
	aru.mutation.SetAccountID(x)
	return aru
}

// SetNillableAccountID sets the "account_id" field if the given value is not nil.
func (aru *AccountRolesUpdate) SetNillableAccountID(x *xid.ID) *AccountRolesUpdate {
	if x != nil {
		aru.SetAccountID(*x)
	}
	return aru
}

// SetRoleID sets the "role_id" field.
func (aru *AccountRolesUpdate) SetRoleID(x xid.ID) *AccountRolesUpdate {
	aru.mutation.SetRoleID(x)
	return aru
}

// SetNillableRoleID sets the "role_id" field if the given value is not nil.
func (aru *AccountRolesUpdate) SetNillableRoleID(x *xid.ID) *AccountRolesUpdate {
	if x != nil {
		aru.SetRoleID(*x)
	}
	return aru
}

// SetBadge sets the "badge" field.
func (aru *AccountRolesUpdate) SetBadge(b bool) *AccountRolesUpdate {
	aru.mutation.SetBadge(b)
	return aru
}

// SetNillableBadge sets the "badge" field if the given value is not nil.
func (aru *AccountRolesUpdate) SetNillableBadge(b *bool) *AccountRolesUpdate {
	if b != nil {
		aru.SetBadge(*b)
	}
	return aru
}

// ClearBadge clears the value of the "badge" field.
func (aru *AccountRolesUpdate) ClearBadge() *AccountRolesUpdate {
	aru.mutation.ClearBadge()
	return aru
}

// SetAccount sets the "account" edge to the Account entity.
func (aru *AccountRolesUpdate) SetAccount(a *Account) *AccountRolesUpdate {
	return aru.SetAccountID(a.ID)
}

// SetRole sets the "role" edge to the Role entity.
func (aru *AccountRolesUpdate) SetRole(r *Role) *AccountRolesUpdate {
	return aru.SetRoleID(r.ID)
}

// Mutation returns the AccountRolesMutation object of the builder.
func (aru *AccountRolesUpdate) Mutation() *AccountRolesMutation {
	return aru.mutation
}

// ClearAccount clears the "account" edge to the Account entity.
func (aru *AccountRolesUpdate) ClearAccount() *AccountRolesUpdate {
	aru.mutation.ClearAccount()
	return aru
}

// ClearRole clears the "role" edge to the Role entity.
func (aru *AccountRolesUpdate) ClearRole() *AccountRolesUpdate {
	aru.mutation.ClearRole()
	return aru
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (aru *AccountRolesUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, aru.sqlSave, aru.mutation, aru.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (aru *AccountRolesUpdate) SaveX(ctx context.Context) int {
	affected, err := aru.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (aru *AccountRolesUpdate) Exec(ctx context.Context) error {
	_, err := aru.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (aru *AccountRolesUpdate) ExecX(ctx context.Context) {
	if err := aru.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (aru *AccountRolesUpdate) check() error {
	if aru.mutation.AccountCleared() && len(aru.mutation.AccountIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "AccountRoles.account"`)
	}
	if aru.mutation.RoleCleared() && len(aru.mutation.RoleIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "AccountRoles.role"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (aru *AccountRolesUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *AccountRolesUpdate {
	aru.modifiers = append(aru.modifiers, modifiers...)
	return aru
}

func (aru *AccountRolesUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := aru.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(accountroles.Table, accountroles.Columns, sqlgraph.NewFieldSpec(accountroles.FieldID, field.TypeString))
	if ps := aru.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := aru.mutation.Badge(); ok {
		_spec.SetField(accountroles.FieldBadge, field.TypeBool, value)
	}
	if aru.mutation.BadgeCleared() {
		_spec.ClearField(accountroles.FieldBadge, field.TypeBool)
	}
	if aru.mutation.AccountCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   accountroles.AccountTable,
			Columns: []string{accountroles.AccountColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(account.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := aru.mutation.AccountIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   accountroles.AccountTable,
			Columns: []string{accountroles.AccountColumn},
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
	if aru.mutation.RoleCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   accountroles.RoleTable,
			Columns: []string{accountroles.RoleColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(role.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := aru.mutation.RoleIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   accountroles.RoleTable,
			Columns: []string{accountroles.RoleColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(role.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(aru.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, aru.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{accountroles.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	aru.mutation.done = true
	return n, nil
}

// AccountRolesUpdateOne is the builder for updating a single AccountRoles entity.
type AccountRolesUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *AccountRolesMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetAccountID sets the "account_id" field.
func (aruo *AccountRolesUpdateOne) SetAccountID(x xid.ID) *AccountRolesUpdateOne {
	aruo.mutation.SetAccountID(x)
	return aruo
}

// SetNillableAccountID sets the "account_id" field if the given value is not nil.
func (aruo *AccountRolesUpdateOne) SetNillableAccountID(x *xid.ID) *AccountRolesUpdateOne {
	if x != nil {
		aruo.SetAccountID(*x)
	}
	return aruo
}

// SetRoleID sets the "role_id" field.
func (aruo *AccountRolesUpdateOne) SetRoleID(x xid.ID) *AccountRolesUpdateOne {
	aruo.mutation.SetRoleID(x)
	return aruo
}

// SetNillableRoleID sets the "role_id" field if the given value is not nil.
func (aruo *AccountRolesUpdateOne) SetNillableRoleID(x *xid.ID) *AccountRolesUpdateOne {
	if x != nil {
		aruo.SetRoleID(*x)
	}
	return aruo
}

// SetBadge sets the "badge" field.
func (aruo *AccountRolesUpdateOne) SetBadge(b bool) *AccountRolesUpdateOne {
	aruo.mutation.SetBadge(b)
	return aruo
}

// SetNillableBadge sets the "badge" field if the given value is not nil.
func (aruo *AccountRolesUpdateOne) SetNillableBadge(b *bool) *AccountRolesUpdateOne {
	if b != nil {
		aruo.SetBadge(*b)
	}
	return aruo
}

// ClearBadge clears the value of the "badge" field.
func (aruo *AccountRolesUpdateOne) ClearBadge() *AccountRolesUpdateOne {
	aruo.mutation.ClearBadge()
	return aruo
}

// SetAccount sets the "account" edge to the Account entity.
func (aruo *AccountRolesUpdateOne) SetAccount(a *Account) *AccountRolesUpdateOne {
	return aruo.SetAccountID(a.ID)
}

// SetRole sets the "role" edge to the Role entity.
func (aruo *AccountRolesUpdateOne) SetRole(r *Role) *AccountRolesUpdateOne {
	return aruo.SetRoleID(r.ID)
}

// Mutation returns the AccountRolesMutation object of the builder.
func (aruo *AccountRolesUpdateOne) Mutation() *AccountRolesMutation {
	return aruo.mutation
}

// ClearAccount clears the "account" edge to the Account entity.
func (aruo *AccountRolesUpdateOne) ClearAccount() *AccountRolesUpdateOne {
	aruo.mutation.ClearAccount()
	return aruo
}

// ClearRole clears the "role" edge to the Role entity.
func (aruo *AccountRolesUpdateOne) ClearRole() *AccountRolesUpdateOne {
	aruo.mutation.ClearRole()
	return aruo
}

// Where appends a list predicates to the AccountRolesUpdate builder.
func (aruo *AccountRolesUpdateOne) Where(ps ...predicate.AccountRoles) *AccountRolesUpdateOne {
	aruo.mutation.Where(ps...)
	return aruo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (aruo *AccountRolesUpdateOne) Select(field string, fields ...string) *AccountRolesUpdateOne {
	aruo.fields = append([]string{field}, fields...)
	return aruo
}

// Save executes the query and returns the updated AccountRoles entity.
func (aruo *AccountRolesUpdateOne) Save(ctx context.Context) (*AccountRoles, error) {
	return withHooks(ctx, aruo.sqlSave, aruo.mutation, aruo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (aruo *AccountRolesUpdateOne) SaveX(ctx context.Context) *AccountRoles {
	node, err := aruo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (aruo *AccountRolesUpdateOne) Exec(ctx context.Context) error {
	_, err := aruo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (aruo *AccountRolesUpdateOne) ExecX(ctx context.Context) {
	if err := aruo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (aruo *AccountRolesUpdateOne) check() error {
	if aruo.mutation.AccountCleared() && len(aruo.mutation.AccountIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "AccountRoles.account"`)
	}
	if aruo.mutation.RoleCleared() && len(aruo.mutation.RoleIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "AccountRoles.role"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (aruo *AccountRolesUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *AccountRolesUpdateOne {
	aruo.modifiers = append(aruo.modifiers, modifiers...)
	return aruo
}

func (aruo *AccountRolesUpdateOne) sqlSave(ctx context.Context) (_node *AccountRoles, err error) {
	if err := aruo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(accountroles.Table, accountroles.Columns, sqlgraph.NewFieldSpec(accountroles.FieldID, field.TypeString))
	id, ok := aruo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "AccountRoles.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := aruo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, accountroles.FieldID)
		for _, f := range fields {
			if !accountroles.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != accountroles.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := aruo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := aruo.mutation.Badge(); ok {
		_spec.SetField(accountroles.FieldBadge, field.TypeBool, value)
	}
	if aruo.mutation.BadgeCleared() {
		_spec.ClearField(accountroles.FieldBadge, field.TypeBool)
	}
	if aruo.mutation.AccountCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   accountroles.AccountTable,
			Columns: []string{accountroles.AccountColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(account.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := aruo.mutation.AccountIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   accountroles.AccountTable,
			Columns: []string{accountroles.AccountColumn},
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
	if aruo.mutation.RoleCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   accountroles.RoleTable,
			Columns: []string{accountroles.RoleColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(role.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := aruo.mutation.RoleIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   accountroles.RoleTable,
			Columns: []string{accountroles.RoleColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(role.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(aruo.modifiers...)
	_node = &AccountRoles{config: aruo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, aruo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{accountroles.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	aruo.mutation.done = true
	return _node, nil
}
