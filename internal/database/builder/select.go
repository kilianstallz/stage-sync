package builder

import (
	"fmt"
	"stage-sync-cli/config"
	"strconv"
)

// BuildSelectQuery build a select query from the given parameters
func BuildSelectQuery(table config.ConfigTable) string {
	query := "SELECT "
	for i, column := range table.Columns {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("%q", column)
	}

	query += " FROM " + fmt.Sprintf("%q", table.Name)

	// where
	if len(table.OnlyWhere) > 0 {
		query += " WHERE "
		for i, column := range table.OnlyWhere {
			if i > 0 {
				query += " AND "
			}
			switch column.Type {
			case "string":
				query += fmt.Sprintf("%q = '%s'", column.Name, column.Value)
				break
			case "int":
				// parse value(string) to int
				val, err := strconv.Atoi(column.Value)
				if err != nil {
					break
				}
				query += fmt.Sprintf("%q = %d", column.Name, val)
				break
			case "float":
				// parse value(string) to float
				val, err := strconv.ParseFloat(column.Value, 64)
				if err != nil {
					break
				}
				query += fmt.Sprintf("%q = %f", column.Name, val)
				break
			}
		}
	}
	return query
}
