// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/Southclaws/storyden/internal/ent/accountroles"
	"github.com/Southclaws/storyden/internal/ent/predicate"
)

// AccountRolesDelete is the builder for deleting a AccountRoles entity.
type AccountRolesDelete struct {
	config
	hooks    []Hook
	mutation *AccountRolesMutation
}

// Where appends a list predicates to the AccountRolesDelete builder.
func (ard *AccountRolesDelete) Where(ps ...predicate.AccountRoles) *AccountRolesDelete {
	ard.mutation.Where(ps...)
	return ard
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (ard *AccountRolesDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, ard.sqlExec, ard.mutation, ard.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (ard *AccountRolesDelete) ExecX(ctx context.Context) int {
	n, err := ard.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (ard *AccountRolesDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(accountroles.Table, sqlgraph.NewFieldSpec(accountroles.FieldID, field.TypeString))
	if ps := ard.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, ard.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	ard.mutation.done = true
	return affected, err
}

// AccountRolesDeleteOne is the builder for deleting a single AccountRoles entity.
type AccountRolesDeleteOne struct {
	ard *AccountRolesDelete
}

// Where appends a list predicates to the AccountRolesDelete builder.
func (ardo *AccountRolesDeleteOne) Where(ps ...predicate.AccountRoles) *AccountRolesDeleteOne {
	ardo.ard.mutation.Where(ps...)
	return ardo
}

// Exec executes the deletion query.
func (ardo *AccountRolesDeleteOne) Exec(ctx context.Context) error {
	n, err := ardo.ard.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{accountroles.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (ardo *AccountRolesDeleteOne) ExecX(ctx context.Context) {
	if err := ardo.Exec(ctx); err != nil {
		panic(err)
	}
}
