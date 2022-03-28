package builder

import (
	"fmt"
	"stage-sync-cli/internal/database/utils"
	"stage-sync-cli/models"
	"strings"
)

// BuildUpdateQuery builds an update query from the given parameters witch only updates the changed columns of the row
func BuildUpdateQuery(tableName string, changedColumn []string, originalRow models.Row, updatedRow models.Row) string {
	query := fmt.Sprintf("UPDATE %q SET ", tableName)
	for _, column := range updatedRow {
		if !utils.ArrayContains(changedColumn, column.Name) {
			continue
		}
		if strings.HasPrefix(column.Type, "int") {
			query += fmt.Sprintf("%q = %d", column.Name, column.Value)
		} else if column.Type == "string" {
			query += fmt.Sprintf("%q = '%s'", column.Name, column.Value)
		} else if strings.HasPrefix(column.Type, "float") {
			query += fmt.Sprintf("%q = %f", column.Name, column.Value)
		} else if column.Type == "NULL" {
			query += fmt.Sprintf("NULL")
		} else {
			query += fmt.Sprintf("%q", column.Name) + " = " + fmt.Sprint(column.Value)
		}
	}
	query += " WHERE "
	for i, column := range originalRow {
		if i > 0 {
			query += " AND "
		}
		if strings.HasPrefix(column.Type, "int") {
			query += fmt.Sprintf("%q = %d", column.Name, column.Value)
		} else if column.Type == "string" {
			query += fmt.Sprintf("%q = '%s'", column.Name, column.Value)
		} else if strings.HasPrefix(column.Type, "float") {
			query += fmt.Sprintf("%q = %f", column.Name, column.Value)
		} else if column.Type == "NULL" {
			query += fmt.Sprintf("%q = NULL", column.Name)
		} else {
			query += fmt.Sprintf("%q", column.Name) + " = " + fmt.Sprint(column.Value)
		}
	}
	return query
}
