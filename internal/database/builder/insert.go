package builder

import (
	"fmt"
	"stage-sync-cli/models"
	"strings"
)

// BuildInsertQuery builds an insert query from the given parameters
func BuildInsertQuery(tableName string, row models.Row) string {
	query := fmt.Sprintf("INSERT INTO %q (", tableName)
	for i, column := range row {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("%q", column.Name)
	}
	query += ") VALUES ("
	for i, column := range row {
		if i > 0 {
			query += ", "
		}
		if strings.HasPrefix(column.Type, "int") || strings.HasPrefix(column.Type, "float") {
			query += fmt.Sprint(column.Value)
		} else if column.Type == "string" {
			query += fmt.Sprintf("'%s'", column.Value)
		} else if column.Type == "NULL" {
			query += fmt.Sprintf("NULL")
		} else {
			query += fmt.Sprint(column.Value)
		}
	}
	query += ")"

	return query
}
