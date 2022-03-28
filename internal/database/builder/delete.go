package builder

import (
	"fmt"
	"stage-sync-cli/models"
	"strings"
)

// BuildUpdateQuery builds an update query from the given parameters
func BuildDeleteQuery(tableName string, rows models.Row) string {
	query := fmt.Sprintf("DELETE FROM %q WHERE ", tableName)
	for i, column := range rows {
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
