package postgres

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/kilianstallz/stage-sync/internal/database/utils"
	"github.com/kilianstallz/stage-sync/pkg/models"
	"go.uber.org/zap"
)

// BuildUpdateQuery builds an update query from the given parameters witch only updates the changed columns of the row
func BuildUpdateQuery(tableName string, changedColumn []string, originalRow models.Row, updatedRow models.Row) string {
	query := goqu.Dialect("postgres").Update(tableName)
	for _, column := range updatedRow {
		if !utils.ArrayContains(changedColumn, column.Name) {
			continue
		}
		query = query.Set(goqu.Ex{column.Name: column.Value})
		// if strings.HasPrefix(column.Type, "int") {
		//	query = query.Set(goqu.Ex{column.Name: column.Value})
		// } else if column.Type == "string" {
		//	query = query.Set(goqu.Ex{column.Name: column.Value})
		// } else if strings.HasPrefix(column.Type, "float") {
		//	query = query.Set(goqu.Ex{column.Name: column.Value})
		// } else if column.Type == "NULL" {
		//	query = query.Set(goqu.Ex{column.Name: column.Value})
		// } else {
		//	query = query.Set(goqu.Ex{column.Name: column.Value})
		// }
	}
	for _, column := range originalRow {
		// if strings.HasPrefix(column.Type, "int") {
		//	query = query.Where(goqu.Ex{column.Name: column.Value})
		// } else if column.Type == "string" {
		//	query = query.Where(goqu.Ex{column.Name: column.Value})
		// } else if strings.HasPrefix(column.Type, "float") {
		//	query = query.Where(goqu.Ex{column.Name: column.Value})
		// } else if column.Type == "NULL" {
		//	query = query.Where(goqu.Ex{column.Name: column.Value})
		// } else {
		//	query = query.Where(goqu.Ex{column.Name: column.Value})
		// }
		if !column.IsPK {
			continue
		}
		query = query.Where(goqu.Ex{column.Name: column.Value})
	}
	q, _, err := query.ToSQL()
	if err != nil {
		zap.S().Fatal(zap.Error(err))
	}
	return q + ";"
}
