package builder

import (
	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
	"stage-sync-cli/models"
)

// BuildInsertQuery builds an insert query from the given parameters
func BuildInsertQuery(tableName string, row models.Row) string {
	var colMap = make(map[string]interface{})
	for _, column := range row {
		colMap[column.Name] = column.Value
	}
	q := goqu.Dialect("postgres").Insert(tableName).Rows(colMap)

	query, _, err := q.ToSQL()
	if err != nil {
		zap.S().Fatal(zap.Error(err))
	}
	return query + ";"
}
