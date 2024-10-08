// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/Southclaws/storyden/internal/ent/accountfollow"
	"github.com/Southclaws/storyden/internal/ent/predicate"
)

// AccountFollowDelete is the builder for deleting a AccountFollow entity.
type AccountFollowDelete struct {
	config
	hooks    []Hook
	mutation *AccountFollowMutation
}

// Where appends a list predicates to the AccountFollowDelete builder.
func (afd *AccountFollowDelete) Where(ps ...predicate.AccountFollow) *AccountFollowDelete {
	afd.mutation.Where(ps...)
	return afd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (afd *AccountFollowDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, afd.sqlExec, afd.mutation, afd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (afd *AccountFollowDelete) ExecX(ctx context.Context) int {
	n, err := afd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (afd *AccountFollowDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(accountfollow.Table, sqlgraph.NewFieldSpec(accountfollow.FieldID, field.TypeString))
	if ps := afd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, afd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	afd.mutation.done = true
	return affected, err
}

// AccountFollowDeleteOne is the builder for deleting a single AccountFollow entity.
type AccountFollowDeleteOne struct {
	afd *AccountFollowDelete
}

// Where appends a list predicates to the AccountFollowDelete builder.
func (afdo *AccountFollowDeleteOne) Where(ps ...predicate.AccountFollow) *AccountFollowDeleteOne {
	afdo.afd.mutation.Where(ps...)
	return afdo
}

// Exec executes the deletion query.
func (afdo *AccountFollowDeleteOne) Exec(ctx context.Context) error {
	n, err := afdo.afd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{accountfollow.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (afdo *AccountFollowDeleteOne) ExecX(ctx context.Context) {
	if err := afdo.Exec(ctx); err != nil {
		panic(err)
	}
}
