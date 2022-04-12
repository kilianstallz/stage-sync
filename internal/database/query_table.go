package database

import (
	"database/sql"
	"fmt"
	numeric "github.com/jackc/pgtype/ext/shopspring-numeric"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"reflect"
	"stage-sync/config"
	"stage-sync/internal/database/builder"
	"stage-sync/models"
	"time"
)

func QueryTables(config *config.Config, sourceDB *sql.DB) []models.Table {
	var tables []models.Table
	for _, table := range config.Tables {
		tables = append(tables, buildTable(table, sourceDB))
	}

	return tables
}

func buildTable(table config.ConfigTable, sourceDB *sql.DB) models.Table {
	// Get all the data from the source database table
	q := builder.BuildSelectQuery(table)
	dbRows, err := sourceDB.Query(q)
	if err != nil {
		zap.S().Fatal("Error querying source database", zap.Error(err))
	}
	defer dbRows.Close()

	// create table struct
	tableStruct := models.Table{
		Name:        table.Name,
		Rows:        []models.Row{},
		PrimaryKeys: table.PrimaryKeys,
		NoDelete:    table.NoDelete,
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
				pointerValues[i] = new(numeric.Numeric)
				object[column.Name()] = new(numeric.Numeric)
			case "DECIMAL":
				pointerValues[i] = new(decimal.NullDecimal)
				object[column.Name()] = new(decimal.NullDecimal)
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

	return tableStruct
}

func reflectType(v interface{}) string {
	if v == nil {
		return "NULL"
	}
	return reflect.TypeOf(v).String()
}
