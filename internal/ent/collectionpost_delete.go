// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/Southclaws/storyden/internal/ent/collectionpost"
	"github.com/Southclaws/storyden/internal/ent/predicate"
)

// CollectionPostDelete is the builder for deleting a CollectionPost entity.
type CollectionPostDelete struct {
	config
	hooks    []Hook
	mutation *CollectionPostMutation
}

// Where appends a list predicates to the CollectionPostDelete builder.
func (cpd *CollectionPostDelete) Where(ps ...predicate.CollectionPost) *CollectionPostDelete {
	cpd.mutation.Where(ps...)
	return cpd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (cpd *CollectionPostDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, cpd.sqlExec, cpd.mutation, cpd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (cpd *CollectionPostDelete) ExecX(ctx context.Context) int {
	n, err := cpd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (cpd *CollectionPostDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(collectionpost.Table, nil)
	if ps := cpd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, cpd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	cpd.mutation.done = true
	return affected, err
}

// CollectionPostDeleteOne is the builder for deleting a single CollectionPost entity.
type CollectionPostDeleteOne struct {
	cpd *CollectionPostDelete
}

// Where appends a list predicates to the CollectionPostDelete builder.
func (cpdo *CollectionPostDeleteOne) Where(ps ...predicate.CollectionPost) *CollectionPostDeleteOne {
	cpdo.cpd.mutation.Where(ps...)
	return cpdo
}

// Exec executes the deletion query.
func (cpdo *CollectionPostDeleteOne) Exec(ctx context.Context) error {
	n, err := cpdo.cpd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{collectionpost.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (cpdo *CollectionPostDeleteOne) ExecX(ctx context.Context) {
	if err := cpdo.Exec(ctx); err != nil {
		panic(err)
	}
}
