package database

import (
	"database/sql"
	"fmt"
	"reflect"
	"stage-sync-cli/config"
	"stage-sync-cli/internal/database/builder"
	"stage-sync-cli/models"
	"time"
)

func QueryTables(config *config.Config, sourceDB *sql.DB) []models.Table {
	var tables []models.Table
	for _, table := range config.Tables {

		// Get all the data from the source database table
		q := builder.BuildSelectQuery(table)
		dbRows, err := sourceDB.Query(q)
		if err != nil {
			panic(err)
		}

		// create table struct
		tableStruct := models.Table{
			Name:        table.Name,
			Rows:        []models.Row{},
			PrimaryKeys: table.PrimaryKeys,
		}

		// iterate over the rows
		for dbRows.Next() {
			// rows := make([]models.Column, 0, len(columns))
			cols, _ := dbRows.ColumnTypes()
			pointerValues := make([]interface{}, len(cols))
			object := map[string]interface{}{}
			for i, column := range cols {
				typeName := column.DatabaseTypeName()
				isNullable := true
				switch typeName {
				case "INT":
					if isNullable {
						pointerValues[i] = new(sql.NullInt64)
						object[column.Name()] = new(sql.NullInt64)
					} else {
						pointerValues[i] = new(int)
						object[column.Name()] = new(int)
					}
				case "VARCHAR":
					if isNullable {
						pointerValues[i] = new(sql.NullString)
						object[column.Name()] = new(sql.NullString)
					} else {
						pointerValues[i] = new(string)
						object[column.Name()] = new(string)
					}
				case "DATETIME":
					if isNullable {
						pointerValues[i] = new(sql.NullTime)
						object[column.Name()] = new(sql.NullTime)
					} else {
						pointerValues[i] = new(time.Time)
						object[column.Name()] = new(time.Time)
					}
				case "NUMERIC":
					if isNullable {
						pointerValues[i] = new(sql.NullFloat64)
						object[column.Name()] = new(sql.NullFloat64)
					} else {
						pointerValues[i] = new(float64)
						object[column.Name()] = new(float64)
					}
				case "TEXT":
					if isNullable {
						pointerValues[i] = new(sql.NullString)
						object[column.Name()] = new(sql.NullString)
					} else {
						pointerValues[i] = new(string)
						object[column.Name()] = new(string)
					}
				case "BOOL":
					if isNullable {
						pointerValues[i] = new(sql.NullBool)
						object[column.Name()] = new(sql.NullBool)
					} else {
						pointerValues[i] = new(bool)
						object[column.Name()] = new(bool)
					}
				case "DATE":
					if isNullable {
						pointerValues[i] = new(sql.NullTime)
						object[column.Name()] = new(sql.NullTime)
					} else {
						pointerValues[i] = new(time.Time)
						object[column.Name()] = new(time.Time)
					}
				case "INT4":
					if isNullable {
						pointerValues[i] = new(sql.NullInt64)
						object[column.Name()] = new(sql.NullInt64)
					} else {
						pointerValues[i] = new(int)
						object[column.Name()] = new(int)
					}
				case "DECIMAL":
					if isNullable {
						pointerValues[i] = new(sql.NullFloat64)
						object[column.Name()] = new(sql.NullFloat64)
					} else {
						pointerValues[i] = new(float64)
						object[column.Name()] = new(float64)
					}
				case "TIMESTAMP":
					if isNullable {
						pointerValues[i] = new(sql.NullTime)
						object[column.Name()] = new(sql.NullTime)
					} else {
						pointerValues[i] = new(time.Time)
						object[column.Name()] = new(time.Time)
					}
				default:
					panic(fmt.Sprintf("Unsupported type: %s", typeName))
				}
			}
			err := dbRows.Scan(pointerValues...)
			if err != nil {
				panic(err)
			}
			// pointer values is an array of pointers to the values
			// make an array of actual values
			values := make([]interface{}, len(pointerValues))
			for i, pointerValue := range pointerValues {
				values[i] = reflect.ValueOf(pointerValue).Elem().Interface()
			}

			row := make([]models.Column, 0, len(table.Columns))
			for i, column := range table.Columns {
				val := ConvertDbValue(values[i])
				typ := reflectType(val)
				row = append(row, models.Column{
					Name:  column,
					Value: val,
					Type:  typ,
				})
			}
			tableStruct.Rows = append(tableStruct.Rows, row)
		}

		tables = append(tables, tableStruct)

		dbRows.Close()
	}

	return tables
}

func reflectType(v interface{}) string {
	if v == nil {
		return "NULL"
	}
	return reflect.TypeOf(v).String()
}
