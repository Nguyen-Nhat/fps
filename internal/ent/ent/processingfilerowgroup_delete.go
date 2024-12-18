// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent/predicate"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent/processingfilerowgroup"
)

// ProcessingFileRowGroupDelete is the builder for deleting a ProcessingFileRowGroup entity.
type ProcessingFileRowGroupDelete struct {
	config
	hooks    []Hook
	mutation *ProcessingFileRowGroupMutation
}

// Where appends a list predicates to the ProcessingFileRowGroupDelete builder.
func (pfrgd *ProcessingFileRowGroupDelete) Where(ps ...predicate.ProcessingFileRowGroup) *ProcessingFileRowGroupDelete {
	pfrgd.mutation.Where(ps...)
	return pfrgd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (pfrgd *ProcessingFileRowGroupDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(pfrgd.hooks) == 0 {
		affected, err = pfrgd.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ProcessingFileRowGroupMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			pfrgd.mutation = mutation
			affected, err = pfrgd.sqlExec(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(pfrgd.hooks) - 1; i >= 0; i-- {
			if pfrgd.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = pfrgd.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, pfrgd.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// ExecX is like Exec, but panics if an error occurs.
func (pfrgd *ProcessingFileRowGroupDelete) ExecX(ctx context.Context) int {
	n, err := pfrgd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (pfrgd *ProcessingFileRowGroupDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: processingfilerowgroup.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: processingfilerowgroup.FieldID,
			},
		},
	}
	if ps := pfrgd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, pfrgd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	return affected, err
}

// ProcessingFileRowGroupDeleteOne is the builder for deleting a single ProcessingFileRowGroup entity.
type ProcessingFileRowGroupDeleteOne struct {
	pfrgd *ProcessingFileRowGroupDelete
}

// Exec executes the deletion query.
func (pfrgdo *ProcessingFileRowGroupDeleteOne) Exec(ctx context.Context) error {
	n, err := pfrgdo.pfrgd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{processingfilerowgroup.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (pfrgdo *ProcessingFileRowGroupDeleteOne) ExecX(ctx context.Context) {
	pfrgdo.pfrgd.ExecX(ctx)
}
