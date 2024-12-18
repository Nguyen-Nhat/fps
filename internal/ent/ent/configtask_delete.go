// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent/configtask"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent/predicate"
)

// ConfigTaskDelete is the builder for deleting a ConfigTask entity.
type ConfigTaskDelete struct {
	config
	hooks    []Hook
	mutation *ConfigTaskMutation
}

// Where appends a list predicates to the ConfigTaskDelete builder.
func (ctd *ConfigTaskDelete) Where(ps ...predicate.ConfigTask) *ConfigTaskDelete {
	ctd.mutation.Where(ps...)
	return ctd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (ctd *ConfigTaskDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(ctd.hooks) == 0 {
		affected, err = ctd.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ConfigTaskMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			ctd.mutation = mutation
			affected, err = ctd.sqlExec(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(ctd.hooks) - 1; i >= 0; i-- {
			if ctd.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = ctd.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, ctd.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// ExecX is like Exec, but panics if an error occurs.
func (ctd *ConfigTaskDelete) ExecX(ctx context.Context) int {
	n, err := ctd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (ctd *ConfigTaskDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: configtask.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: configtask.FieldID,
			},
		},
	}
	if ps := ctd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, ctd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	return affected, err
}

// ConfigTaskDeleteOne is the builder for deleting a single ConfigTask entity.
type ConfigTaskDeleteOne struct {
	ctd *ConfigTaskDelete
}

// Exec executes the deletion query.
func (ctdo *ConfigTaskDeleteOne) Exec(ctx context.Context) error {
	n, err := ctdo.ctd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{configtask.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (ctdo *ConfigTaskDeleteOne) ExecX(ctx context.Context) {
	ctdo.ctd.ExecX(ctx)
}
