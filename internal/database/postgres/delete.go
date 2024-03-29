package postgres

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/kilianstallz/stage-sync/pkg/models"
	"go.uber.org/zap"
)

// BuildDeleteQuery builds an delete query from the given parameters
func BuildDeleteQuery(tableName string, rows models.Row) string {
	q := goqu.Dialect("postgres").Delete(tableName)
	for _, column := range rows {
		q = q.Where(
			goqu.C(column.Name).Eq(column.Value),
		)
	}
	query, _, err := q.ToSQL()
	if err != nil {
		zap.S().Fatal(zap.Error(err))
	}
	return query + ";"
}
