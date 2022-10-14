package builder

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"stage-sync/models"
)

// BuildInsertQuery builds an insert query from the given parameters
func BuildInsertQuery(tableName string, row models.Row) string {
	var colMap = make(map[string]interface{})
	lo.ForEach[models.Column](row, func(column models.Column, _ int) {
		colMap[column.Name] = column.Value
	})
	q := goqu.Dialect("postgres").Insert(tableName).Rows(colMap)

	query, _, err := q.ToSQL()
	if err != nil {
		zap.S().Fatal(zap.Error(err))
	}
	return query + ";"
}
