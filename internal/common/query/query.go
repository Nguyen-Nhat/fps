package query

import (
	"context"
	dbsql "database/sql"
)

func RunRawQuery[T any](ctx context.Context, sqlDB *dbsql.DB, rawQuery string, scan func(*dbsql.Rows, *T) error, args ...any) ([]T, error) {
	// 1. Query
	rows, err := sqlDB.QueryContext(ctx, rawQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 2. Convert to DAO
	var list []T
	for rows.Next() {
		var item T
		if errScan := scan(rows, &item); errScan != nil {
			return nil, errScan
		}
		list = append(list, item)
	}

	// 3. Check error
	if errRow := rows.Err(); errRow != nil {
		return nil, errRow
	}

	// 4. Return result
	return list, nil
}
