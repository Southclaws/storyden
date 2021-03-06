// Code generated by entc, DO NOT EDIT.

package model

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/Southclaws/storyden/backend/internal/infrastructure/db/model/post"
	"github.com/Southclaws/storyden/backend/internal/infrastructure/db/model/predicate"
	"github.com/Southclaws/storyden/backend/internal/infrastructure/db/model/react"
	"github.com/Southclaws/storyden/backend/internal/infrastructure/db/model/user"
	"github.com/google/uuid"
)

// ReactUpdate is the builder for updating React entities.
type ReactUpdate struct {
	config
	hooks    []Hook
	mutation *ReactMutation
}

// Where appends a list predicates to the ReactUpdate builder.
func (ru *ReactUpdate) Where(ps ...predicate.React) *ReactUpdate {
	ru.mutation.Where(ps...)
	return ru
}

// SetEmoji sets the "emoji" field.
func (ru *ReactUpdate) SetEmoji(s string) *ReactUpdate {
	ru.mutation.SetEmoji(s)
	return ru
}

// SetCreatedAt sets the "createdAt" field.
func (ru *ReactUpdate) SetCreatedAt(t time.Time) *ReactUpdate {
	ru.mutation.SetCreatedAt(t)
	return ru
}

// SetNillableCreatedAt sets the "createdAt" field if the given value is not nil.
func (ru *ReactUpdate) SetNillableCreatedAt(t *time.Time) *ReactUpdate {
	if t != nil {
		ru.SetCreatedAt(*t)
	}
	return ru
}

// SetUserID sets the "user" edge to the User entity by ID.
func (ru *ReactUpdate) SetUserID(id uuid.UUID) *ReactUpdate {
	ru.mutation.SetUserID(id)
	return ru
}

// SetNillableUserID sets the "user" edge to the User entity by ID if the given value is not nil.
func (ru *ReactUpdate) SetNillableUserID(id *uuid.UUID) *ReactUpdate {
	if id != nil {
		ru = ru.SetUserID(*id)
	}
	return ru
}

// SetUser sets the "user" edge to the User entity.
func (ru *ReactUpdate) SetUser(u *User) *ReactUpdate {
	return ru.SetUserID(u.ID)
}

// SetPostID sets the "Post" edge to the Post entity by ID.
func (ru *ReactUpdate) SetPostID(id uuid.UUID) *ReactUpdate {
	ru.mutation.SetPostID(id)
	return ru
}

// SetNillablePostID sets the "Post" edge to the Post entity by ID if the given value is not nil.
func (ru *ReactUpdate) SetNillablePostID(id *uuid.UUID) *ReactUpdate {
	if id != nil {
		ru = ru.SetPostID(*id)
	}
	return ru
}

// SetPost sets the "Post" edge to the Post entity.
func (ru *ReactUpdate) SetPost(p *Post) *ReactUpdate {
	return ru.SetPostID(p.ID)
}

// Mutation returns the ReactMutation object of the builder.
func (ru *ReactUpdate) Mutation() *ReactMutation {
	return ru.mutation
}

// ClearUser clears the "user" edge to the User entity.
func (ru *ReactUpdate) ClearUser() *ReactUpdate {
	ru.mutation.ClearUser()
	return ru
}

// ClearPost clears the "Post" edge to the Post entity.
func (ru *ReactUpdate) ClearPost() *ReactUpdate {
	ru.mutation.ClearPost()
	return ru
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (ru *ReactUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(ru.hooks) == 0 {
		affected, err = ru.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ReactMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			ru.mutation = mutation
			affected, err = ru.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(ru.hooks) - 1; i >= 0; i-- {
			if ru.hooks[i] == nil {
				return 0, fmt.Errorf("model: uninitialized hook (forgotten import model/runtime?)")
			}
			mut = ru.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, ru.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
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

func (ru *ReactUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   react.Table,
			Columns: react.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: react.FieldID,
			},
		},
	}
	if ps := ru.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ru.mutation.Emoji(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: react.FieldEmoji,
		})
	}
	if value, ok := ru.mutation.CreatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: react.FieldCreatedAt,
		})
	}
	if ru.mutation.UserCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ru.mutation.UserIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if ru.mutation.PostCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ru.mutation.PostIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, ru.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{react.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// ReactUpdateOne is the builder for updating a single React entity.
type ReactUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *ReactMutation
}

// SetEmoji sets the "emoji" field.
func (ruo *ReactUpdateOne) SetEmoji(s string) *ReactUpdateOne {
	ruo.mutation.SetEmoji(s)
	return ruo
}

// SetCreatedAt sets the "createdAt" field.
func (ruo *ReactUpdateOne) SetCreatedAt(t time.Time) *ReactUpdateOne {
	ruo.mutation.SetCreatedAt(t)
	return ruo
}

// SetNillableCreatedAt sets the "createdAt" field if the given value is not nil.
func (ruo *ReactUpdateOne) SetNillableCreatedAt(t *time.Time) *ReactUpdateOne {
	if t != nil {
		ruo.SetCreatedAt(*t)
	}
	return ruo
}

// SetUserID sets the "user" edge to the User entity by ID.
func (ruo *ReactUpdateOne) SetUserID(id uuid.UUID) *ReactUpdateOne {
	ruo.mutation.SetUserID(id)
	return ruo
}

// SetNillableUserID sets the "user" edge to the User entity by ID if the given value is not nil.
func (ruo *ReactUpdateOne) SetNillableUserID(id *uuid.UUID) *ReactUpdateOne {
	if id != nil {
		ruo = ruo.SetUserID(*id)
	}
	return ruo
}

// SetUser sets the "user" edge to the User entity.
func (ruo *ReactUpdateOne) SetUser(u *User) *ReactUpdateOne {
	return ruo.SetUserID(u.ID)
}

// SetPostID sets the "Post" edge to the Post entity by ID.
func (ruo *ReactUpdateOne) SetPostID(id uuid.UUID) *ReactUpdateOne {
	ruo.mutation.SetPostID(id)
	return ruo
}

// SetNillablePostID sets the "Post" edge to the Post entity by ID if the given value is not nil.
func (ruo *ReactUpdateOne) SetNillablePostID(id *uuid.UUID) *ReactUpdateOne {
	if id != nil {
		ruo = ruo.SetPostID(*id)
	}
	return ruo
}

// SetPost sets the "Post" edge to the Post entity.
func (ruo *ReactUpdateOne) SetPost(p *Post) *ReactUpdateOne {
	return ruo.SetPostID(p.ID)
}

// Mutation returns the ReactMutation object of the builder.
func (ruo *ReactUpdateOne) Mutation() *ReactMutation {
	return ruo.mutation
}

// ClearUser clears the "user" edge to the User entity.
func (ruo *ReactUpdateOne) ClearUser() *ReactUpdateOne {
	ruo.mutation.ClearUser()
	return ruo
}

// ClearPost clears the "Post" edge to the Post entity.
func (ruo *ReactUpdateOne) ClearPost() *ReactUpdateOne {
	ruo.mutation.ClearPost()
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
	var (
		err  error
		node *React
	)
	if len(ruo.hooks) == 0 {
		node, err = ruo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ReactMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			ruo.mutation = mutation
			node, err = ruo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(ruo.hooks) - 1; i >= 0; i-- {
			if ruo.hooks[i] == nil {
				return nil, fmt.Errorf("model: uninitialized hook (forgotten import model/runtime?)")
			}
			mut = ruo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, ruo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
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

func (ruo *ReactUpdateOne) sqlSave(ctx context.Context) (_node *React, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   react.Table,
			Columns: react.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: react.FieldID,
			},
		},
	}
	id, ok := ruo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`model: missing "React.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := ruo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, react.FieldID)
		for _, f := range fields {
			if !react.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("model: invalid field %q for query", f)}
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
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: react.FieldEmoji,
		})
	}
	if value, ok := ruo.mutation.CreatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: react.FieldCreatedAt,
		})
	}
	if ruo.mutation.UserCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ruo.mutation.UserIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if ruo.mutation.PostCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ruo.mutation.PostIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &React{config: ruo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, ruo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{react.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}
