// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent/predicate"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent/processingfilerowgroup"
)

// ProcessingFileRowGroupQuery is the builder for querying ProcessingFileRowGroup entities.
type ProcessingFileRowGroupQuery struct {
	config
	limit      *int
	offset     *int
	unique     *bool
	order      []OrderFunc
	fields     []string
	predicates []predicate.ProcessingFileRowGroup
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the ProcessingFileRowGroupQuery builder.
func (pfrgq *ProcessingFileRowGroupQuery) Where(ps ...predicate.ProcessingFileRowGroup) *ProcessingFileRowGroupQuery {
	pfrgq.predicates = append(pfrgq.predicates, ps...)
	return pfrgq
}

// Limit adds a limit step to the query.
func (pfrgq *ProcessingFileRowGroupQuery) Limit(limit int) *ProcessingFileRowGroupQuery {
	pfrgq.limit = &limit
	return pfrgq
}

// Offset adds an offset step to the query.
func (pfrgq *ProcessingFileRowGroupQuery) Offset(offset int) *ProcessingFileRowGroupQuery {
	pfrgq.offset = &offset
	return pfrgq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (pfrgq *ProcessingFileRowGroupQuery) Unique(unique bool) *ProcessingFileRowGroupQuery {
	pfrgq.unique = &unique
	return pfrgq
}

// Order adds an order step to the query.
func (pfrgq *ProcessingFileRowGroupQuery) Order(o ...OrderFunc) *ProcessingFileRowGroupQuery {
	pfrgq.order = append(pfrgq.order, o...)
	return pfrgq
}

// First returns the first ProcessingFileRowGroup entity from the query.
// Returns a *NotFoundError when no ProcessingFileRowGroup was found.
func (pfrgq *ProcessingFileRowGroupQuery) First(ctx context.Context) (*ProcessingFileRowGroup, error) {
	nodes, err := pfrgq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{processingfilerowgroup.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (pfrgq *ProcessingFileRowGroupQuery) FirstX(ctx context.Context) *ProcessingFileRowGroup {
	node, err := pfrgq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first ProcessingFileRowGroup ID from the query.
// Returns a *NotFoundError when no ProcessingFileRowGroup ID was found.
func (pfrgq *ProcessingFileRowGroupQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = pfrgq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{processingfilerowgroup.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (pfrgq *ProcessingFileRowGroupQuery) FirstIDX(ctx context.Context) int {
	id, err := pfrgq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single ProcessingFileRowGroup entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one ProcessingFileRowGroup entity is found.
// Returns a *NotFoundError when no ProcessingFileRowGroup entities are found.
func (pfrgq *ProcessingFileRowGroupQuery) Only(ctx context.Context) (*ProcessingFileRowGroup, error) {
	nodes, err := pfrgq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{processingfilerowgroup.Label}
	default:
		return nil, &NotSingularError{processingfilerowgroup.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (pfrgq *ProcessingFileRowGroupQuery) OnlyX(ctx context.Context) *ProcessingFileRowGroup {
	node, err := pfrgq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only ProcessingFileRowGroup ID in the query.
// Returns a *NotSingularError when more than one ProcessingFileRowGroup ID is found.
// Returns a *NotFoundError when no entities are found.
func (pfrgq *ProcessingFileRowGroupQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = pfrgq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{processingfilerowgroup.Label}
	default:
		err = &NotSingularError{processingfilerowgroup.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (pfrgq *ProcessingFileRowGroupQuery) OnlyIDX(ctx context.Context) int {
	id, err := pfrgq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of ProcessingFileRowGroups.
func (pfrgq *ProcessingFileRowGroupQuery) All(ctx context.Context) ([]*ProcessingFileRowGroup, error) {
	if err := pfrgq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return pfrgq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (pfrgq *ProcessingFileRowGroupQuery) AllX(ctx context.Context) []*ProcessingFileRowGroup {
	nodes, err := pfrgq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of ProcessingFileRowGroup IDs.
func (pfrgq *ProcessingFileRowGroupQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := pfrgq.Select(processingfilerowgroup.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (pfrgq *ProcessingFileRowGroupQuery) IDsX(ctx context.Context) []int {
	ids, err := pfrgq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (pfrgq *ProcessingFileRowGroupQuery) Count(ctx context.Context) (int, error) {
	if err := pfrgq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return pfrgq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (pfrgq *ProcessingFileRowGroupQuery) CountX(ctx context.Context) int {
	count, err := pfrgq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (pfrgq *ProcessingFileRowGroupQuery) Exist(ctx context.Context) (bool, error) {
	if err := pfrgq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return pfrgq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (pfrgq *ProcessingFileRowGroupQuery) ExistX(ctx context.Context) bool {
	exist, err := pfrgq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the ProcessingFileRowGroupQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (pfrgq *ProcessingFileRowGroupQuery) Clone() *ProcessingFileRowGroupQuery {
	if pfrgq == nil {
		return nil
	}
	return &ProcessingFileRowGroupQuery{
		config:     pfrgq.config,
		limit:      pfrgq.limit,
		offset:     pfrgq.offset,
		order:      append([]OrderFunc{}, pfrgq.order...),
		predicates: append([]predicate.ProcessingFileRowGroup{}, pfrgq.predicates...),
		// clone intermediate query.
		sql:    pfrgq.sql.Clone(),
		path:   pfrgq.path,
		unique: pfrgq.unique,
	}
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		FileID int64 `json:"file_id,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.ProcessingFileRowGroup.Query().
//		GroupBy(processingfilerowgroup.FieldFileID).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (pfrgq *ProcessingFileRowGroupQuery) GroupBy(field string, fields ...string) *ProcessingFileRowGroupGroupBy {
	grbuild := &ProcessingFileRowGroupGroupBy{config: pfrgq.config}
	grbuild.fields = append([]string{field}, fields...)
	grbuild.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := pfrgq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return pfrgq.sqlQuery(ctx), nil
	}
	grbuild.label = processingfilerowgroup.Label
	grbuild.flds, grbuild.scan = &grbuild.fields, grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		FileID int64 `json:"file_id,omitempty"`
//	}
//
//	client.ProcessingFileRowGroup.Query().
//		Select(processingfilerowgroup.FieldFileID).
//		Scan(ctx, &v)
func (pfrgq *ProcessingFileRowGroupQuery) Select(fields ...string) *ProcessingFileRowGroupSelect {
	pfrgq.fields = append(pfrgq.fields, fields...)
	selbuild := &ProcessingFileRowGroupSelect{ProcessingFileRowGroupQuery: pfrgq}
	selbuild.label = processingfilerowgroup.Label
	selbuild.flds, selbuild.scan = &pfrgq.fields, selbuild.Scan
	return selbuild
}

func (pfrgq *ProcessingFileRowGroupQuery) prepareQuery(ctx context.Context) error {
	for _, f := range pfrgq.fields {
		if !processingfilerowgroup.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if pfrgq.path != nil {
		prev, err := pfrgq.path(ctx)
		if err != nil {
			return err
		}
		pfrgq.sql = prev
	}
	return nil
}

func (pfrgq *ProcessingFileRowGroupQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*ProcessingFileRowGroup, error) {
	var (
		nodes = []*ProcessingFileRowGroup{}
		_spec = pfrgq.querySpec()
	)
	_spec.ScanValues = func(columns []string) ([]interface{}, error) {
		return (*ProcessingFileRowGroup).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []interface{}) error {
		node := &ProcessingFileRowGroup{config: pfrgq.config}
		nodes = append(nodes, node)
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, pfrgq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (pfrgq *ProcessingFileRowGroupQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := pfrgq.querySpec()
	_spec.Node.Columns = pfrgq.fields
	if len(pfrgq.fields) > 0 {
		_spec.Unique = pfrgq.unique != nil && *pfrgq.unique
	}
	return sqlgraph.CountNodes(ctx, pfrgq.driver, _spec)
}

func (pfrgq *ProcessingFileRowGroupQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := pfrgq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %w", err)
	}
	return n > 0, nil
}

func (pfrgq *ProcessingFileRowGroupQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   processingfilerowgroup.Table,
			Columns: processingfilerowgroup.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: processingfilerowgroup.FieldID,
			},
		},
		From:   pfrgq.sql,
		Unique: true,
	}
	if unique := pfrgq.unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := pfrgq.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, processingfilerowgroup.FieldID)
		for i := range fields {
			if fields[i] != processingfilerowgroup.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := pfrgq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := pfrgq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := pfrgq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := pfrgq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (pfrgq *ProcessingFileRowGroupQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(pfrgq.driver.Dialect())
	t1 := builder.Table(processingfilerowgroup.Table)
	columns := pfrgq.fields
	if len(columns) == 0 {
		columns = processingfilerowgroup.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if pfrgq.sql != nil {
		selector = pfrgq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if pfrgq.unique != nil && *pfrgq.unique {
		selector.Distinct()
	}
	for _, p := range pfrgq.predicates {
		p(selector)
	}
	for _, p := range pfrgq.order {
		p(selector)
	}
	if offset := pfrgq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := pfrgq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ProcessingFileRowGroupGroupBy is the group-by builder for ProcessingFileRowGroup entities.
type ProcessingFileRowGroupGroupBy struct {
	config
	selector
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (pfrggb *ProcessingFileRowGroupGroupBy) Aggregate(fns ...AggregateFunc) *ProcessingFileRowGroupGroupBy {
	pfrggb.fns = append(pfrggb.fns, fns...)
	return pfrggb
}

// Scan applies the group-by query and scans the result into the given value.
func (pfrggb *ProcessingFileRowGroupGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := pfrggb.path(ctx)
	if err != nil {
		return err
	}
	pfrggb.sql = query
	return pfrggb.sqlScan(ctx, v)
}

func (pfrggb *ProcessingFileRowGroupGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	for _, f := range pfrggb.fields {
		if !processingfilerowgroup.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for group-by", f)}
		}
	}
	selector := pfrggb.sqlQuery()
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := pfrggb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (pfrggb *ProcessingFileRowGroupGroupBy) sqlQuery() *sql.Selector {
	selector := pfrggb.sql.Select()
	aggregation := make([]string, 0, len(pfrggb.fns))
	for _, fn := range pfrggb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	// If no columns were selected in a custom aggregation function, the default
	// selection is the fields used for "group-by", and the aggregation functions.
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(pfrggb.fields)+len(pfrggb.fns))
		for _, f := range pfrggb.fields {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	return selector.GroupBy(selector.Columns(pfrggb.fields...)...)
}

// ProcessingFileRowGroupSelect is the builder for selecting fields of ProcessingFileRowGroup entities.
type ProcessingFileRowGroupSelect struct {
	*ProcessingFileRowGroupQuery
	selector
	// intermediate query (i.e. traversal path).
	sql *sql.Selector
}

// Scan applies the selector query and scans the result into the given value.
func (pfrgs *ProcessingFileRowGroupSelect) Scan(ctx context.Context, v interface{}) error {
	if err := pfrgs.prepareQuery(ctx); err != nil {
		return err
	}
	pfrgs.sql = pfrgs.ProcessingFileRowGroupQuery.sqlQuery(ctx)
	return pfrgs.sqlScan(ctx, v)
}

func (pfrgs *ProcessingFileRowGroupSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := pfrgs.sql.Query()
	if err := pfrgs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
