package database

import (
	"fmt"
	"github.com/kilianstallz/stage-sync/models"
	"github.com/kilianstallz/stage-sync/pkg/config"
	"strconv"
	"strings"
)

//// BuildConnectionString builds a postgres connection string in url form from the config
//func BuildConnectionString(credentials config.ConfigDB) string {
//	userPart := ""
//	if credentials.User != "" {
//		userPart = fmt.Sprintf("%s", credentials.User)
//	}
//	if credentials.Password != "" {
//		userPart = fmt.Sprintf("%s:%s", userPart, credentials.Password)
//	}
//	if userPart != "" {
//		userPart = fmt.Sprintf("%s@", userPart)
//	}
//	return fmt.Sprintf("postgres://%s%s:%d/%s?sslmode=disable", userPart, credentials.Host, credentials.Port, credentials.Database)
//}

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
		} else {
			query += fmt.Sprint(column.Value)
		}
	}
	query += ")"

	return query
}

// BuildDeleteQuery builds an delete query from the given parameters
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
		} else {
			query += fmt.Sprintf("%q", column.Name) + " = " + fmt.Sprint(column.Value)
		}
	}
	return query
}

// BuildUpdateQuery builds an update query from the given parameters witch only updates the changed columns of the row
func BuildUpdateQuery(tableName string, changedColumn []string, originalRow models.Row, updatedRow models.Row) string {
	query := fmt.Sprintf("UPDATE %q SET ", tableName)
	for i, column := range updatedRow {
		if !contains(changedColumn, column.Name) {
			continue
		}
		if strings.HasPrefix(column.Type, "int") {
			query += fmt.Sprintf("%q = %d", column.Name, column.Value)
		} else if column.Type == "string" {
			query += fmt.Sprintf("%q = '%s'", column.Name, column.Value)
		} else if strings.HasPrefix(column.Type, "float") {
			query += fmt.Sprintf("%q = %f", column.Name, column.Value)
		} else {
			query += fmt.Sprintf("%q", column.Name) + " = " + fmt.Sprint(column.Value)
		}
		if i < len(updatedRow)-1 {
			query += ", "
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
		} else {
			query += fmt.Sprintf("%q", column.Name) + " = " + fmt.Sprint(column.Value)
		}
	}
	return query
}

// contains checks if a string is in a given slice
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
