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
	"github.com/Southclaws/storyden/internal/ent/post"
	"github.com/Southclaws/storyden/internal/ent/predicate"
	"github.com/Southclaws/storyden/internal/ent/react"
	"github.com/rs/xid"
)

// ReactUpdate is the builder for updating React entities.
type ReactUpdate struct {
	config
	hooks     []Hook
	mutation  *ReactMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the ReactUpdate builder.
func (ru *ReactUpdate) Where(ps ...predicate.React) *ReactUpdate {
	ru.mutation.Where(ps...)
	return ru
}

// SetAccountID sets the "account_id" field.
func (ru *ReactUpdate) SetAccountID(x xid.ID) *ReactUpdate {
	ru.mutation.SetAccountID(x)
	return ru
}

// SetPostID sets the "post_id" field.
func (ru *ReactUpdate) SetPostID(x xid.ID) *ReactUpdate {
	ru.mutation.SetPostID(x)
	return ru
}

// SetEmoji sets the "emoji" field.
func (ru *ReactUpdate) SetEmoji(s string) *ReactUpdate {
	ru.mutation.SetEmoji(s)
	return ru
}

// SetAccount sets the "account" edge to the Account entity.
func (ru *ReactUpdate) SetAccount(a *Account) *ReactUpdate {
	return ru.SetAccountID(a.ID)
}

// SetPost sets the "Post" edge to the Post entity.
func (ru *ReactUpdate) SetPost(p *Post) *ReactUpdate {
	return ru.SetPostID(p.ID)
}

// Mutation returns the ReactMutation object of the builder.
func (ru *ReactUpdate) Mutation() *ReactMutation {
	return ru.mutation
}

// ClearAccount clears the "account" edge to the Account entity.
func (ru *ReactUpdate) ClearAccount() *ReactUpdate {
	ru.mutation.ClearAccount()
	return ru
}

// ClearPost clears the "Post" edge to the Post entity.
func (ru *ReactUpdate) ClearPost() *ReactUpdate {
	ru.mutation.ClearPost()
	return ru
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (ru *ReactUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, ru.sqlSave, ru.mutation, ru.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ru *ReactUpdate) SaveX(ctx context.Context) int {
	affected, err := ru.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (ru *ReactUpdate) Exec(ctx context.Context) error {
	_, err := ru.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ru *ReactUpdate) ExecX(ctx context.Context) {
	if err := ru.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ru *ReactUpdate) check() error {
	if _, ok := ru.mutation.AccountID(); ru.mutation.AccountCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "React.account"`)
	}
	if _, ok := ru.mutation.PostID(); ru.mutation.PostCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "React.Post"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (ru *ReactUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ReactUpdate {
	ru.modifiers = append(ru.modifiers, modifiers...)
	return ru
}

func (ru *ReactUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := ru.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(react.Table, react.Columns, sqlgraph.NewFieldSpec(react.FieldID, field.TypeString))
	if ps := ru.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ru.mutation.Emoji(); ok {
		_spec.SetField(react.FieldEmoji, field.TypeString, value)
	}
	if ru.mutation.AccountCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   react.AccountTable,
			Columns: []string{react.AccountColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(account.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ru.mutation.AccountIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   react.AccountTable,
			Columns: []string{react.AccountColumn},
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
	if ru.mutation.PostCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   react.PostTable,
			Columns: []string{react.PostColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(post.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ru.mutation.PostIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   react.PostTable,
			Columns: []string{react.PostColumn},
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
	_spec.AddModifiers(ru.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, ru.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{react.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	ru.mutation.done = true
	return n, nil
}

// ReactUpdateOne is the builder for updating a single React entity.
type ReactUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *ReactMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetAccountID sets the "account_id" field.
func (ruo *ReactUpdateOne) SetAccountID(x xid.ID) *ReactUpdateOne {
	ruo.mutation.SetAccountID(x)
	return ruo
}

// SetPostID sets the "post_id" field.
func (ruo *ReactUpdateOne) SetPostID(x xid.ID) *ReactUpdateOne {
	ruo.mutation.SetPostID(x)
	return ruo
}

// SetEmoji sets the "emoji" field.
func (ruo *ReactUpdateOne) SetEmoji(s string) *ReactUpdateOne {
	ruo.mutation.SetEmoji(s)
	return ruo
}

// SetAccount sets the "account" edge to the Account entity.
func (ruo *ReactUpdateOne) SetAccount(a *Account) *ReactUpdateOne {
	return ruo.SetAccountID(a.ID)
}

// SetPost sets the "Post" edge to the Post entity.
func (ruo *ReactUpdateOne) SetPost(p *Post) *ReactUpdateOne {
	return ruo.SetPostID(p.ID)
}

// Mutation returns the ReactMutation object of the builder.
func (ruo *ReactUpdateOne) Mutation() *ReactMutation {
	return ruo.mutation
}

// ClearAccount clears the "account" edge to the Account entity.
func (ruo *ReactUpdateOne) ClearAccount() *ReactUpdateOne {
	ruo.mutation.ClearAccount()
	return ruo
}

// ClearPost clears the "Post" edge to the Post entity.
func (ruo *ReactUpdateOne) ClearPost() *ReactUpdateOne {
	ruo.mutation.ClearPost()
	return ruo
}

// Where appends a list predicates to the ReactUpdate builder.
func (ruo *ReactUpdateOne) Where(ps ...predicate.React) *ReactUpdateOne {
	ruo.mutation.Where(ps...)
	return ruo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (ruo *ReactUpdateOne) Select(field string, fields ...string) *ReactUpdateOne {
	ruo.fields = append([]string{field}, fields...)
	return ruo
}

// Save executes the query and returns the updated React entity.
func (ruo *ReactUpdateOne) Save(ctx context.Context) (*React, error) {
	return withHooks(ctx, ruo.sqlSave, ruo.mutation, ruo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ruo *ReactUpdateOne) SaveX(ctx context.Context) *React {
	node, err := ruo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (ruo *ReactUpdateOne) Exec(ctx context.Context) error {
	_, err := ruo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ruo *ReactUpdateOne) ExecX(ctx context.Context) {
	if err := ruo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ruo *ReactUpdateOne) check() error {
	if _, ok := ruo.mutation.AccountID(); ruo.mutation.AccountCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "React.account"`)
	}
	if _, ok := ruo.mutation.PostID(); ruo.mutation.PostCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "React.Post"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (ruo *ReactUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ReactUpdateOne {
	ruo.modifiers = append(ruo.modifiers, modifiers...)
	return ruo
}

func (ruo *ReactUpdateOne) sqlSave(ctx context.Context) (_node *React, err error) {
	if err := ruo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(react.Table, react.Columns, sqlgraph.NewFieldSpec(react.FieldID, field.TypeString))
	id, ok := ruo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "React.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := ruo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, react.FieldID)
		for _, f := range fields {
			if !react.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != react.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := ruo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ruo.mutation.Emoji(); ok {
		_spec.SetField(react.FieldEmoji, field.TypeString, value)
	}
	if ruo.mutation.AccountCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   react.AccountTable,
			Columns: []string{react.AccountColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(account.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ruo.mutation.AccountIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   react.AccountTable,
			Columns: []string{react.AccountColumn},
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
	if ruo.mutation.PostCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   react.PostTable,
			Columns: []string{react.PostColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(post.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ruo.mutation.PostIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   react.PostTable,
			Columns: []string{react.PostColumn},
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
	_spec.AddModifiers(ruo.modifiers...)
	_node = &React{config: ruo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, ruo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{react.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	ruo.mutation.done = true
	return _node, nil
}
