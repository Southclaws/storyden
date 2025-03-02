// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/Southclaws/storyden/internal/ent/predicate"
	"github.com/Southclaws/storyden/internal/ent/propertyschema"
)

// PropertySchemaDelete is the builder for deleting a PropertySchema entity.
type PropertySchemaDelete struct {
	config
	hooks    []Hook
	mutation *PropertySchemaMutation
}

// Where appends a list predicates to the PropertySchemaDelete builder.
func (psd *PropertySchemaDelete) Where(ps ...predicate.PropertySchema) *PropertySchemaDelete {
	psd.mutation.Where(ps...)
	return psd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (psd *PropertySchemaDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, psd.sqlExec, psd.mutation, psd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (psd *PropertySchemaDelete) ExecX(ctx context.Context) int {
	n, err := psd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (psd *PropertySchemaDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(propertyschema.Table, sqlgraph.NewFieldSpec(propertyschema.FieldID, field.TypeString))
	if ps := psd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, psd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	psd.mutation.done = true
	return affected, err
}

// PropertySchemaDeleteOne is the builder for deleting a single PropertySchema entity.
type PropertySchemaDeleteOne struct {
	psd *PropertySchemaDelete
}

// Where appends a list predicates to the PropertySchemaDelete builder.
func (psdo *PropertySchemaDeleteOne) Where(ps ...predicate.PropertySchema) *PropertySchemaDeleteOne {
	psdo.psd.mutation.Where(ps...)
	return psdo
}

// Exec executes the deletion query.
func (psdo *PropertySchemaDeleteOne) Exec(ctx context.Context) error {
	n, err := psdo.psd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{propertyschema.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (psdo *PropertySchemaDeleteOne) ExecX(ctx context.Context) {
	if err := psdo.Exec(ctx); err != nil {
		panic(err)
	}
}
