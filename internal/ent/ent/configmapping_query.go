// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent/configmapping"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent/predicate"
)

// ConfigMappingQuery is the builder for querying ConfigMapping entities.
type ConfigMappingQuery struct {
	config
	limit      *int
	offset     *int
	unique     *bool
	order      []OrderFunc
	fields     []string
	predicates []predicate.ConfigMapping
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the ConfigMappingQuery builder.
func (cmq *ConfigMappingQuery) Where(ps ...predicate.ConfigMapping) *ConfigMappingQuery {
	cmq.predicates = append(cmq.predicates, ps...)
	return cmq
}

// Limit adds a limit step to the query.
func (cmq *ConfigMappingQuery) Limit(limit int) *ConfigMappingQuery {
	cmq.limit = &limit
	return cmq
}

// Offset adds an offset step to the query.
func (cmq *ConfigMappingQuery) Offset(offset int) *ConfigMappingQuery {
	cmq.offset = &offset
	return cmq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (cmq *ConfigMappingQuery) Unique(unique bool) *ConfigMappingQuery {
	cmq.unique = &unique
	return cmq
}

// Order adds an order step to the query.
func (cmq *ConfigMappingQuery) Order(o ...OrderFunc) *ConfigMappingQuery {
	cmq.order = append(cmq.order, o...)
	return cmq
}

// First returns the first ConfigMapping entity from the query.
// Returns a *NotFoundError when no ConfigMapping was found.
func (cmq *ConfigMappingQuery) First(ctx context.Context) (*ConfigMapping, error) {
	nodes, err := cmq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{configmapping.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (cmq *ConfigMappingQuery) FirstX(ctx context.Context) *ConfigMapping {
	node, err := cmq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first ConfigMapping ID from the query.
// Returns a *NotFoundError when no ConfigMapping ID was found.
func (cmq *ConfigMappingQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = cmq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{configmapping.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (cmq *ConfigMappingQuery) FirstIDX(ctx context.Context) int {
	id, err := cmq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single ConfigMapping entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one ConfigMapping entity is found.
// Returns a *NotFoundError when no ConfigMapping entities are found.
func (cmq *ConfigMappingQuery) Only(ctx context.Context) (*ConfigMapping, error) {
	nodes, err := cmq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{configmapping.Label}
	default:
		return nil, &NotSingularError{configmapping.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (cmq *ConfigMappingQuery) OnlyX(ctx context.Context) *ConfigMapping {
	node, err := cmq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only ConfigMapping ID in the query.
// Returns a *NotSingularError when more than one ConfigMapping ID is found.
// Returns a *NotFoundError when no entities are found.
func (cmq *ConfigMappingQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = cmq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{configmapping.Label}
	default:
		err = &NotSingularError{configmapping.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (cmq *ConfigMappingQuery) OnlyIDX(ctx context.Context) int {
	id, err := cmq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of ConfigMappings.
func (cmq *ConfigMappingQuery) All(ctx context.Context) ([]*ConfigMapping, error) {
	if err := cmq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return cmq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (cmq *ConfigMappingQuery) AllX(ctx context.Context) []*ConfigMapping {
	nodes, err := cmq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of ConfigMapping IDs.
func (cmq *ConfigMappingQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := cmq.Select(configmapping.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (cmq *ConfigMappingQuery) IDsX(ctx context.Context) []int {
	ids, err := cmq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (cmq *ConfigMappingQuery) Count(ctx context.Context) (int, error) {
	if err := cmq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return cmq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (cmq *ConfigMappingQuery) CountX(ctx context.Context) int {
	count, err := cmq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (cmq *ConfigMappingQuery) Exist(ctx context.Context) (bool, error) {
	if err := cmq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return cmq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (cmq *ConfigMappingQuery) ExistX(ctx context.Context) bool {
	exist, err := cmq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the ConfigMappingQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (cmq *ConfigMappingQuery) Clone() *ConfigMappingQuery {
	if cmq == nil {
		return nil
	}
	return &ConfigMappingQuery{
		config:     cmq.config,
		limit:      cmq.limit,
		offset:     cmq.offset,
		order:      append([]OrderFunc{}, cmq.order...),
		predicates: append([]predicate.ConfigMapping{}, cmq.predicates...),
		// clone intermediate query.
		sql:    cmq.sql.Clone(),
		path:   cmq.path,
		unique: cmq.unique,
	}
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		ClientID int32 `json:"client_id,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.ConfigMapping.Query().
//		GroupBy(configmapping.FieldClientID).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (cmq *ConfigMappingQuery) GroupBy(field string, fields ...string) *ConfigMappingGroupBy {
	grbuild := &ConfigMappingGroupBy{config: cmq.config}
	grbuild.fields = append([]string{field}, fields...)
	grbuild.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := cmq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return cmq.sqlQuery(ctx), nil
	}
	grbuild.label = configmapping.Label
	grbuild.flds, grbuild.scan = &grbuild.fields, grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		ClientID int32 `json:"client_id,omitempty"`
//	}
//
//	client.ConfigMapping.Query().
//		Select(configmapping.FieldClientID).
//		Scan(ctx, &v)
func (cmq *ConfigMappingQuery) Select(fields ...string) *ConfigMappingSelect {
	cmq.fields = append(cmq.fields, fields...)
	selbuild := &ConfigMappingSelect{ConfigMappingQuery: cmq}
	selbuild.label = configmapping.Label
	selbuild.flds, selbuild.scan = &cmq.fields, selbuild.Scan
	return selbuild
}

func (cmq *ConfigMappingQuery) prepareQuery(ctx context.Context) error {
	for _, f := range cmq.fields {
		if !configmapping.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if cmq.path != nil {
		prev, err := cmq.path(ctx)
		if err != nil {
			return err
		}
		cmq.sql = prev
	}
	return nil
}

func (cmq *ConfigMappingQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*ConfigMapping, error) {
	var (
		nodes = []*ConfigMapping{}
		_spec = cmq.querySpec()
	)
	_spec.ScanValues = func(columns []string) ([]interface{}, error) {
		return (*ConfigMapping).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []interface{}) error {
		node := &ConfigMapping{config: cmq.config}
		nodes = append(nodes, node)
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, cmq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (cmq *ConfigMappingQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := cmq.querySpec()
	_spec.Node.Columns = cmq.fields
	if len(cmq.fields) > 0 {
		_spec.Unique = cmq.unique != nil && *cmq.unique
	}
	return sqlgraph.CountNodes(ctx, cmq.driver, _spec)
}

func (cmq *ConfigMappingQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := cmq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %w", err)
	}
	return n > 0, nil
}

func (cmq *ConfigMappingQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   configmapping.Table,
			Columns: configmapping.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: configmapping.FieldID,
			},
		},
		From:   cmq.sql,
		Unique: true,
	}
	if unique := cmq.unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := cmq.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, configmapping.FieldID)
		for i := range fields {
			if fields[i] != configmapping.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := cmq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := cmq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := cmq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := cmq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (cmq *ConfigMappingQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(cmq.driver.Dialect())
	t1 := builder.Table(configmapping.Table)
	columns := cmq.fields
	if len(columns) == 0 {
		columns = configmapping.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if cmq.sql != nil {
		selector = cmq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if cmq.unique != nil && *cmq.unique {
		selector.Distinct()
	}
	for _, p := range cmq.predicates {
		p(selector)
	}
	for _, p := range cmq.order {
		p(selector)
	}
	if offset := cmq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := cmq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ConfigMappingGroupBy is the group-by builder for ConfigMapping entities.
type ConfigMappingGroupBy struct {
	config
	selector
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (cmgb *ConfigMappingGroupBy) Aggregate(fns ...AggregateFunc) *ConfigMappingGroupBy {
	cmgb.fns = append(cmgb.fns, fns...)
	return cmgb
}

// Scan applies the group-by query and scans the result into the given value.
func (cmgb *ConfigMappingGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := cmgb.path(ctx)
	if err != nil {
		return err
	}
	cmgb.sql = query
	return cmgb.sqlScan(ctx, v)
}

func (cmgb *ConfigMappingGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	for _, f := range cmgb.fields {
		if !configmapping.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for group-by", f)}
		}
	}
	selector := cmgb.sqlQuery()
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := cmgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (cmgb *ConfigMappingGroupBy) sqlQuery() *sql.Selector {
	selector := cmgb.sql.Select()
	aggregation := make([]string, 0, len(cmgb.fns))
	for _, fn := range cmgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	// If no columns were selected in a custom aggregation function, the default
	// selection is the fields used for "group-by", and the aggregation functions.
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(cmgb.fields)+len(cmgb.fns))
		for _, f := range cmgb.fields {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	return selector.GroupBy(selector.Columns(cmgb.fields...)...)
}

// ConfigMappingSelect is the builder for selecting fields of ConfigMapping entities.
type ConfigMappingSelect struct {
	*ConfigMappingQuery
	selector
	// intermediate query (i.e. traversal path).
	sql *sql.Selector
}

// Scan applies the selector query and scans the result into the given value.
func (cms *ConfigMappingSelect) Scan(ctx context.Context, v interface{}) error {
	if err := cms.prepareQuery(ctx); err != nil {
		return err
	}
	cms.sql = cms.ConfigMappingQuery.sqlQuery(ctx)
	return cms.sqlScan(ctx, v)
}

func (cms *ConfigMappingSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := cms.sql.Query()
	if err := cms.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}