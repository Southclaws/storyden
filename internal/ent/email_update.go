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
	"github.com/Southclaws/storyden/internal/ent/email"
	"github.com/Southclaws/storyden/internal/ent/predicate"
	"github.com/rs/xid"
)

// EmailUpdate is the builder for updating Email entities.
type EmailUpdate struct {
	config
	hooks     []Hook
	mutation  *EmailMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the EmailUpdate builder.
func (eu *EmailUpdate) Where(ps ...predicate.Email) *EmailUpdate {
	eu.mutation.Where(ps...)
	return eu
}

// SetAccountID sets the "account_id" field.
func (eu *EmailUpdate) SetAccountID(x xid.ID) *EmailUpdate {
	eu.mutation.SetAccountID(x)
	return eu
}

// SetNillableAccountID sets the "account_id" field if the given value is not nil.
func (eu *EmailUpdate) SetNillableAccountID(x *xid.ID) *EmailUpdate {
	if x != nil {
		eu.SetAccountID(*x)
	}
	return eu
}

// ClearAccountID clears the value of the "account_id" field.
func (eu *EmailUpdate) ClearAccountID() *EmailUpdate {
	eu.mutation.ClearAccountID()
	return eu
}

// SetVerificationCode sets the "verification_code" field.
func (eu *EmailUpdate) SetVerificationCode(s string) *EmailUpdate {
	eu.mutation.SetVerificationCode(s)
	return eu
}

// SetNillableVerificationCode sets the "verification_code" field if the given value is not nil.
func (eu *EmailUpdate) SetNillableVerificationCode(s *string) *EmailUpdate {
	if s != nil {
		eu.SetVerificationCode(*s)
	}
	return eu
}

// SetVerified sets the "verified" field.
func (eu *EmailUpdate) SetVerified(b bool) *EmailUpdate {
	eu.mutation.SetVerified(b)
	return eu
}

// SetNillableVerified sets the "verified" field if the given value is not nil.
func (eu *EmailUpdate) SetNillableVerified(b *bool) *EmailUpdate {
	if b != nil {
		eu.SetVerified(*b)
	}
	return eu
}

// SetAccount sets the "account" edge to the Account entity.
func (eu *EmailUpdate) SetAccount(a *Account) *EmailUpdate {
	return eu.SetAccountID(a.ID)
}

// Mutation returns the EmailMutation object of the builder.
func (eu *EmailUpdate) Mutation() *EmailMutation {
	return eu.mutation
}

// ClearAccount clears the "account" edge to the Account entity.
func (eu *EmailUpdate) ClearAccount() *EmailUpdate {
	eu.mutation.ClearAccount()
	return eu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (eu *EmailUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, eu.sqlSave, eu.mutation, eu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (eu *EmailUpdate) SaveX(ctx context.Context) int {
	affected, err := eu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (eu *EmailUpdate) Exec(ctx context.Context) error {
	_, err := eu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (eu *EmailUpdate) ExecX(ctx context.Context) {
	if err := eu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (eu *EmailUpdate) check() error {
	if v, ok := eu.mutation.VerificationCode(); ok {
		if err := email.VerificationCodeValidator(v); err != nil {
			return &ValidationError{Name: "verification_code", err: fmt.Errorf(`ent: validator failed for field "Email.verification_code": %w`, err)}
		}
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (eu *EmailUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *EmailUpdate {
	eu.modifiers = append(eu.modifiers, modifiers...)
	return eu
}

func (eu *EmailUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := eu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(email.Table, email.Columns, sqlgraph.NewFieldSpec(email.FieldID, field.TypeString))
	if ps := eu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := eu.mutation.VerificationCode(); ok {
		_spec.SetField(email.FieldVerificationCode, field.TypeString, value)
	}
	if value, ok := eu.mutation.Verified(); ok {
		_spec.SetField(email.FieldVerified, field.TypeBool, value)
	}
	if eu.mutation.AccountCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   email.AccountTable,
			Columns: []string{email.AccountColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(account.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := eu.mutation.AccountIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   email.AccountTable,
			Columns: []string{email.AccountColumn},
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
	_spec.AddModifiers(eu.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, eu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{email.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	eu.mutation.done = true
	return n, nil
}

// EmailUpdateOne is the builder for updating a single Email entity.
type EmailUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *EmailMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetAccountID sets the "account_id" field.
func (euo *EmailUpdateOne) SetAccountID(x xid.ID) *EmailUpdateOne {
	euo.mutation.SetAccountID(x)
	return euo
}

// SetNillableAccountID sets the "account_id" field if the given value is not nil.
func (euo *EmailUpdateOne) SetNillableAccountID(x *xid.ID) *EmailUpdateOne {
	if x != nil {
		euo.SetAccountID(*x)
	}
	return euo
}

// ClearAccountID clears the value of the "account_id" field.
func (euo *EmailUpdateOne) ClearAccountID() *EmailUpdateOne {
	euo.mutation.ClearAccountID()
	return euo
}

// SetVerificationCode sets the "verification_code" field.
func (euo *EmailUpdateOne) SetVerificationCode(s string) *EmailUpdateOne {
	euo.mutation.SetVerificationCode(s)
	return euo
}

// SetNillableVerificationCode sets the "verification_code" field if the given value is not nil.
func (euo *EmailUpdateOne) SetNillableVerificationCode(s *string) *EmailUpdateOne {
	if s != nil {
		euo.SetVerificationCode(*s)
	}
	return euo
}

// SetVerified sets the "verified" field.
func (euo *EmailUpdateOne) SetVerified(b bool) *EmailUpdateOne {
	euo.mutation.SetVerified(b)
	return euo
}

// SetNillableVerified sets the "verified" field if the given value is not nil.
func (euo *EmailUpdateOne) SetNillableVerified(b *bool) *EmailUpdateOne {
	if b != nil {
		euo.SetVerified(*b)
	}
	return euo
}

// SetAccount sets the "account" edge to the Account entity.
func (euo *EmailUpdateOne) SetAccount(a *Account) *EmailUpdateOne {
	return euo.SetAccountID(a.ID)
}

// Mutation returns the EmailMutation object of the builder.
func (euo *EmailUpdateOne) Mutation() *EmailMutation {
	return euo.mutation
}

// ClearAccount clears the "account" edge to the Account entity.
func (euo *EmailUpdateOne) ClearAccount() *EmailUpdateOne {
	euo.mutation.ClearAccount()
	return euo
}

// Where appends a list predicates to the EmailUpdate builder.
func (euo *EmailUpdateOne) Where(ps ...predicate.Email) *EmailUpdateOne {
	euo.mutation.Where(ps...)
	return euo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (euo *EmailUpdateOne) Select(field string, fields ...string) *EmailUpdateOne {
	euo.fields = append([]string{field}, fields...)
	return euo
}

// Save executes the query and returns the updated Email entity.
func (euo *EmailUpdateOne) Save(ctx context.Context) (*Email, error) {
	return withHooks(ctx, euo.sqlSave, euo.mutation, euo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (euo *EmailUpdateOne) SaveX(ctx context.Context) *Email {
	node, err := euo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (euo *EmailUpdateOne) Exec(ctx context.Context) error {
	_, err := euo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (euo *EmailUpdateOne) ExecX(ctx context.Context) {
	if err := euo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (euo *EmailUpdateOne) check() error {
	if v, ok := euo.mutation.VerificationCode(); ok {
		if err := email.VerificationCodeValidator(v); err != nil {
			return &ValidationError{Name: "verification_code", err: fmt.Errorf(`ent: validator failed for field "Email.verification_code": %w`, err)}
		}
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (euo *EmailUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *EmailUpdateOne {
	euo.modifiers = append(euo.modifiers, modifiers...)
	return euo
}

func (euo *EmailUpdateOne) sqlSave(ctx context.Context) (_node *Email, err error) {
	if err := euo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(email.Table, email.Columns, sqlgraph.NewFieldSpec(email.FieldID, field.TypeString))
	id, ok := euo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Email.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := euo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, email.FieldID)
		for _, f := range fields {
			if !email.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != email.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := euo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := euo.mutation.VerificationCode(); ok {
		_spec.SetField(email.FieldVerificationCode, field.TypeString, value)
	}
	if value, ok := euo.mutation.Verified(); ok {
		_spec.SetField(email.FieldVerified, field.TypeBool, value)
	}
	if euo.mutation.AccountCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   email.AccountTable,
			Columns: []string{email.AccountColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(account.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := euo.mutation.AccountIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   email.AccountTable,
			Columns: []string{email.AccountColumn},
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
	_spec.AddModifiers(euo.modifiers...)
	_node = &Email{config: euo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, euo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{email.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	euo.mutation.done = true
	return _node, nil
}
