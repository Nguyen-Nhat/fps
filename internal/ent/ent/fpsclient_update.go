// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent/fpsclient"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent/predicate"
)

// FpsClientUpdate is the builder for updating FpsClient entities.
type FpsClientUpdate struct {
	config
	hooks    []Hook
	mutation *FpsClientMutation
}

// Where appends a list predicates to the FpsClientUpdate builder.
func (fcu *FpsClientUpdate) Where(ps ...predicate.FpsClient) *FpsClientUpdate {
	fcu.mutation.Where(ps...)
	return fcu
}

// SetClientID sets the "client_id" field.
func (fcu *FpsClientUpdate) SetClientID(i int32) *FpsClientUpdate {
	fcu.mutation.ResetClientID()
	fcu.mutation.SetClientID(i)
	return fcu
}

// AddClientID adds i to the "client_id" field.
func (fcu *FpsClientUpdate) AddClientID(i int32) *FpsClientUpdate {
	fcu.mutation.AddClientID(i)
	return fcu
}

// SetName sets the "name" field.
func (fcu *FpsClientUpdate) SetName(s string) *FpsClientUpdate {
	fcu.mutation.SetName(s)
	return fcu
}

// SetDescription sets the "description" field.
func (fcu *FpsClientUpdate) SetDescription(s string) *FpsClientUpdate {
	fcu.mutation.SetDescription(s)
	return fcu
}

// SetSampleFileURL sets the "sample_file_url" field.
func (fcu *FpsClientUpdate) SetSampleFileURL(s string) *FpsClientUpdate {
	fcu.mutation.SetSampleFileURL(s)
	return fcu
}

// SetNillableSampleFileURL sets the "sample_file_url" field if the given value is not nil.
func (fcu *FpsClientUpdate) SetNillableSampleFileURL(s *string) *FpsClientUpdate {
	if s != nil {
		fcu.SetSampleFileURL(*s)
	}
	return fcu
}

// SetImportFileTemplateURL sets the "import_file_template_url" field.
func (fcu *FpsClientUpdate) SetImportFileTemplateURL(s string) *FpsClientUpdate {
	fcu.mutation.SetImportFileTemplateURL(s)
	return fcu
}

// SetNillableImportFileTemplateURL sets the "import_file_template_url" field if the given value is not nil.
func (fcu *FpsClientUpdate) SetNillableImportFileTemplateURL(s *string) *FpsClientUpdate {
	if s != nil {
		fcu.SetImportFileTemplateURL(*s)
	}
	return fcu
}

// ClearImportFileTemplateURL clears the value of the "import_file_template_url" field.
func (fcu *FpsClientUpdate) ClearImportFileTemplateURL() *FpsClientUpdate {
	fcu.mutation.ClearImportFileTemplateURL()
	return fcu
}

// SetCreatedAt sets the "created_at" field.
func (fcu *FpsClientUpdate) SetCreatedAt(t time.Time) *FpsClientUpdate {
	fcu.mutation.SetCreatedAt(t)
	return fcu
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (fcu *FpsClientUpdate) SetNillableCreatedAt(t *time.Time) *FpsClientUpdate {
	if t != nil {
		fcu.SetCreatedAt(*t)
	}
	return fcu
}

// SetCreatedBy sets the "created_by" field.
func (fcu *FpsClientUpdate) SetCreatedBy(s string) *FpsClientUpdate {
	fcu.mutation.SetCreatedBy(s)
	return fcu
}

// SetUpdatedAt sets the "updated_at" field.
func (fcu *FpsClientUpdate) SetUpdatedAt(t time.Time) *FpsClientUpdate {
	fcu.mutation.SetUpdatedAt(t)
	return fcu
}

// Mutation returns the FpsClientMutation object of the builder.
func (fcu *FpsClientUpdate) Mutation() *FpsClientMutation {
	return fcu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (fcu *FpsClientUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	fcu.defaults()
	if len(fcu.hooks) == 0 {
		if err = fcu.check(); err != nil {
			return 0, err
		}
		affected, err = fcu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*FpsClientMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = fcu.check(); err != nil {
				return 0, err
			}
			fcu.mutation = mutation
			affected, err = fcu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(fcu.hooks) - 1; i >= 0; i-- {
			if fcu.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = fcu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, fcu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (fcu *FpsClientUpdate) SaveX(ctx context.Context) int {
	affected, err := fcu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (fcu *FpsClientUpdate) Exec(ctx context.Context) error {
	_, err := fcu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fcu *FpsClientUpdate) ExecX(ctx context.Context) {
	if err := fcu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (fcu *FpsClientUpdate) defaults() {
	if _, ok := fcu.mutation.UpdatedAt(); !ok {
		v := fpsclient.UpdateDefaultUpdatedAt()
		fcu.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (fcu *FpsClientUpdate) check() error {
	if v, ok := fcu.mutation.Name(); ok {
		if err := fpsclient.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "FpsClient.name": %w`, err)}
		}
	}
	if v, ok := fcu.mutation.Description(); ok {
		if err := fpsclient.DescriptionValidator(v); err != nil {
			return &ValidationError{Name: "description", err: fmt.Errorf(`ent: validator failed for field "FpsClient.description": %w`, err)}
		}
	}
	if v, ok := fcu.mutation.CreatedBy(); ok {
		if err := fpsclient.CreatedByValidator(v); err != nil {
			return &ValidationError{Name: "created_by", err: fmt.Errorf(`ent: validator failed for field "FpsClient.created_by": %w`, err)}
		}
	}
	return nil
}

func (fcu *FpsClientUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   fpsclient.Table,
			Columns: fpsclient.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: fpsclient.FieldID,
			},
		},
	}
	if ps := fcu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := fcu.mutation.ClientID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt32,
			Value:  value,
			Column: fpsclient.FieldClientID,
		})
	}
	if value, ok := fcu.mutation.AddedClientID(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt32,
			Value:  value,
			Column: fpsclient.FieldClientID,
		})
	}
	if value, ok := fcu.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: fpsclient.FieldName,
		})
	}
	if value, ok := fcu.mutation.Description(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: fpsclient.FieldDescription,
		})
	}
	if value, ok := fcu.mutation.SampleFileURL(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: fpsclient.FieldSampleFileURL,
		})
	}
	if value, ok := fcu.mutation.ImportFileTemplateURL(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: fpsclient.FieldImportFileTemplateURL,
		})
	}
	if fcu.mutation.ImportFileTemplateURLCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: fpsclient.FieldImportFileTemplateURL,
		})
	}
	if value, ok := fcu.mutation.CreatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: fpsclient.FieldCreatedAt,
		})
	}
	if value, ok := fcu.mutation.CreatedBy(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: fpsclient.FieldCreatedBy,
		})
	}
	if value, ok := fcu.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: fpsclient.FieldUpdatedAt,
		})
	}
	if n, err = sqlgraph.UpdateNodes(ctx, fcu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{fpsclient.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	return n, nil
}

// FpsClientUpdateOne is the builder for updating a single FpsClient entity.
type FpsClientUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *FpsClientMutation
}

// SetClientID sets the "client_id" field.
func (fcuo *FpsClientUpdateOne) SetClientID(i int32) *FpsClientUpdateOne {
	fcuo.mutation.ResetClientID()
	fcuo.mutation.SetClientID(i)
	return fcuo
}

// AddClientID adds i to the "client_id" field.
func (fcuo *FpsClientUpdateOne) AddClientID(i int32) *FpsClientUpdateOne {
	fcuo.mutation.AddClientID(i)
	return fcuo
}

// SetName sets the "name" field.
func (fcuo *FpsClientUpdateOne) SetName(s string) *FpsClientUpdateOne {
	fcuo.mutation.SetName(s)
	return fcuo
}

// SetDescription sets the "description" field.
func (fcuo *FpsClientUpdateOne) SetDescription(s string) *FpsClientUpdateOne {
	fcuo.mutation.SetDescription(s)
	return fcuo
}

// SetSampleFileURL sets the "sample_file_url" field.
func (fcuo *FpsClientUpdateOne) SetSampleFileURL(s string) *FpsClientUpdateOne {
	fcuo.mutation.SetSampleFileURL(s)
	return fcuo
}

// SetNillableSampleFileURL sets the "sample_file_url" field if the given value is not nil.
func (fcuo *FpsClientUpdateOne) SetNillableSampleFileURL(s *string) *FpsClientUpdateOne {
	if s != nil {
		fcuo.SetSampleFileURL(*s)
	}
	return fcuo
}

// SetImportFileTemplateURL sets the "import_file_template_url" field.
func (fcuo *FpsClientUpdateOne) SetImportFileTemplateURL(s string) *FpsClientUpdateOne {
	fcuo.mutation.SetImportFileTemplateURL(s)
	return fcuo
}

// SetNillableImportFileTemplateURL sets the "import_file_template_url" field if the given value is not nil.
func (fcuo *FpsClientUpdateOne) SetNillableImportFileTemplateURL(s *string) *FpsClientUpdateOne {
	if s != nil {
		fcuo.SetImportFileTemplateURL(*s)
	}
	return fcuo
}

// ClearImportFileTemplateURL clears the value of the "import_file_template_url" field.
func (fcuo *FpsClientUpdateOne) ClearImportFileTemplateURL() *FpsClientUpdateOne {
	fcuo.mutation.ClearImportFileTemplateURL()
	return fcuo
}

// SetCreatedAt sets the "created_at" field.
func (fcuo *FpsClientUpdateOne) SetCreatedAt(t time.Time) *FpsClientUpdateOne {
	fcuo.mutation.SetCreatedAt(t)
	return fcuo
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (fcuo *FpsClientUpdateOne) SetNillableCreatedAt(t *time.Time) *FpsClientUpdateOne {
	if t != nil {
		fcuo.SetCreatedAt(*t)
	}
	return fcuo
}

// SetCreatedBy sets the "created_by" field.
func (fcuo *FpsClientUpdateOne) SetCreatedBy(s string) *FpsClientUpdateOne {
	fcuo.mutation.SetCreatedBy(s)
	return fcuo
}

// SetUpdatedAt sets the "updated_at" field.
func (fcuo *FpsClientUpdateOne) SetUpdatedAt(t time.Time) *FpsClientUpdateOne {
	fcuo.mutation.SetUpdatedAt(t)
	return fcuo
}

// Mutation returns the FpsClientMutation object of the builder.
func (fcuo *FpsClientUpdateOne) Mutation() *FpsClientMutation {
	return fcuo.mutation
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (fcuo *FpsClientUpdateOne) Select(field string, fields ...string) *FpsClientUpdateOne {
	fcuo.fields = append([]string{field}, fields...)
	return fcuo
}

// Save executes the query and returns the updated FpsClient entity.
func (fcuo *FpsClientUpdateOne) Save(ctx context.Context) (*FpsClient, error) {
	var (
		err  error
		node *FpsClient
	)
	fcuo.defaults()
	if len(fcuo.hooks) == 0 {
		if err = fcuo.check(); err != nil {
			return nil, err
		}
		node, err = fcuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*FpsClientMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = fcuo.check(); err != nil {
				return nil, err
			}
			fcuo.mutation = mutation
			node, err = fcuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(fcuo.hooks) - 1; i >= 0; i-- {
			if fcuo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = fcuo.hooks[i](mut)
		}
		v, err := mut.Mutate(ctx, fcuo.mutation)
		if err != nil {
			return nil, err
		}
		nv, ok := v.(*FpsClient)
		if !ok {
			return nil, fmt.Errorf("unexpected node type %T returned from FpsClientMutation", v)
		}
		node = nv
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (fcuo *FpsClientUpdateOne) SaveX(ctx context.Context) *FpsClient {
	node, err := fcuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (fcuo *FpsClientUpdateOne) Exec(ctx context.Context) error {
	_, err := fcuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fcuo *FpsClientUpdateOne) ExecX(ctx context.Context) {
	if err := fcuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (fcuo *FpsClientUpdateOne) defaults() {
	if _, ok := fcuo.mutation.UpdatedAt(); !ok {
		v := fpsclient.UpdateDefaultUpdatedAt()
		fcuo.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (fcuo *FpsClientUpdateOne) check() error {
	if v, ok := fcuo.mutation.Name(); ok {
		if err := fpsclient.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "FpsClient.name": %w`, err)}
		}
	}
	if v, ok := fcuo.mutation.Description(); ok {
		if err := fpsclient.DescriptionValidator(v); err != nil {
			return &ValidationError{Name: "description", err: fmt.Errorf(`ent: validator failed for field "FpsClient.description": %w`, err)}
		}
	}
	if v, ok := fcuo.mutation.CreatedBy(); ok {
		if err := fpsclient.CreatedByValidator(v); err != nil {
			return &ValidationError{Name: "created_by", err: fmt.Errorf(`ent: validator failed for field "FpsClient.created_by": %w`, err)}
		}
	}
	return nil
}

func (fcuo *FpsClientUpdateOne) sqlSave(ctx context.Context) (_node *FpsClient, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   fpsclient.Table,
			Columns: fpsclient.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: fpsclient.FieldID,
			},
		},
	}
	id, ok := fcuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "FpsClient.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := fcuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, fpsclient.FieldID)
		for _, f := range fields {
			if !fpsclient.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != fpsclient.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := fcuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := fcuo.mutation.ClientID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt32,
			Value:  value,
			Column: fpsclient.FieldClientID,
		})
	}
	if value, ok := fcuo.mutation.AddedClientID(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt32,
			Value:  value,
			Column: fpsclient.FieldClientID,
		})
	}
	if value, ok := fcuo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: fpsclient.FieldName,
		})
	}
	if value, ok := fcuo.mutation.Description(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: fpsclient.FieldDescription,
		})
	}
	if value, ok := fcuo.mutation.SampleFileURL(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: fpsclient.FieldSampleFileURL,
		})
	}
	if value, ok := fcuo.mutation.ImportFileTemplateURL(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: fpsclient.FieldImportFileTemplateURL,
		})
	}
	if fcuo.mutation.ImportFileTemplateURLCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: fpsclient.FieldImportFileTemplateURL,
		})
	}
	if value, ok := fcuo.mutation.CreatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: fpsclient.FieldCreatedAt,
		})
	}
	if value, ok := fcuo.mutation.CreatedBy(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: fpsclient.FieldCreatedBy,
		})
	}
	if value, ok := fcuo.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: fpsclient.FieldUpdatedAt,
		})
	}
	_node = &FpsClient{config: fcuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, fcuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{fpsclient.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	return _node, nil
}
