package builder

import (
	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
	"stage-sync/config"
)

// BuildSelectQuery build a select query from the given parameters
func BuildSelectQuery(table config.ConfigTable) string {
	query := goqu.From(table.Name)
	var selects []interface{}
	for _, column := range table.Columns {
		selects = append(selects, column)
	}
	query = query.Select(selects...)
	// where
	if len(table.OnlyWhere) > 0 {
		whereMap := make(map[string][]interface{})
		for _, column := range table.OnlyWhere {
			if whereMap[column.Name] == nil {
				whereMap[column.Name] = []interface{}{column.Value}
			} else {
				whereMap[column.Name] = append(whereMap[column.Name], column.Value)
			}
		}
		for key, value := range whereMap {
			query = query.Where(goqu.I(key).In(value))
		}
	}
	q, _, err := query.ToSQL()
	if err != nil {
		zap.S().Fatal(err)
	}
	return q + ";"
}
